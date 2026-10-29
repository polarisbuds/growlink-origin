<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { api } from "$lib/api";
    import { filterItems } from "$lib/searchable-select";
    import UserChipSearchField from "$lib/components/UserChipSearchField.svelte";
    import SearchableSelect from "$lib/components/SearchableSelect.svelte";

    type Target = { kind: string; targetId: number; label: string };
    type Template = { id: number; title?: string };
    type User = { id: number; name: string; email: string };
    type ManagedTeam = { id: number; name: string; members?: User[] };

    function todayISO() {
        const d = new Date();
        return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
    }

    let templates = $state<Template[]>([]);
    let users = $state<User[]>([]);
    let managedTeams = $state<ManagedTeam[]>([]);
    let mentorId = $state(0);
    let targets = $state<Target[]>([]);
    let userTargets = $derived(targets.filter((t) => t.kind === "user"));
    let sharerUserIds = $state<number[]>([]);
    let menteeChips = $derived(
        userTargets.map((t) => ({
            id: t.targetId,
            name: displayNameForTarget(t),
        })),
    );
    let sharerChips = $derived(
        sharerUserIds.map((id) => ({ id, name: userLabel(id) })),
    );
    let submitting = $state(false);
    let error = $state("");

    let form = $state({
        templateId: 0,
        implementedAt: todayISO(),
    });

    let userTargetPickOpen = $state(false);
    let userTargetPickQuery = $state("");
    let managedTeamPickOpen = $state(false);
    let managedTeamPickQuery = $state("");
    let templatePickOpen = $state(false);
    let templatePickQuery = $state("");
    let sharerPickOpen = $state(false);
    let sharerPickQuery = $state("");

    let userTargetPickFiltered = $derived.by(() => {
        const picked = new Set(
            targets.filter((t) => t.kind === "user").map((t) => t.targetId),
        );
        const list = users.filter((u) => !picked.has(u.id));
        return filterItems(
            list,
            userTargetPickQuery,
            (u) => `${u.name} (${u.email})`,
        );
    });
    let managedTeamPickFiltered = $derived(
        filterItems(managedTeams, managedTeamPickQuery, (t) => t.name),
    );
    let templatePickFiltered = $derived(
        filterItems(templates, templatePickQuery, (t) => t.title || ""),
    );
    let sharerPickFiltered = $derived.by(() => {
        const picked = new Set(sharerUserIds);
        const list = users.filter((u) => !picked.has(u.id));
        return filterItems(
            list,
            sharerPickQuery,
            (u) => `${u.name} (${u.email})`,
        );
    });

    function userLabel(id: number) {
        const u = users.find((x) => x.id === id);
        return u ? u.name : `#${id}`;
    }
    function displayNameForTarget(t: Target) {
        const u = users.find((x) => x.id === t.targetId);
        if (u) return u.name;
        const idx = t.label.indexOf(" (");
        if (idx > 0) return t.label.slice(0, idx);
        return t.label || `#${t.targetId}`;
    }
    function targetKey(kind: string, targetId: number) {
        return `${kind}:${targetId}`;
    }
    function addTarget(kind: string, targetId: number, label: string) {
        const key = targetKey(kind, targetId);
        if (targets.some((t) => targetKey(t.kind, t.targetId) === key)) return;
        targets = [...targets, { kind, targetId, label }];
    }
    function removeTarget(kind: string, targetId: number) {
        targets = targets.filter(
            (t) => !(t.kind === kind && t.targetId === targetId),
        );
    }
    function closePicks() {
        templatePickOpen = false;
        sharerPickOpen = false;
        userTargetPickOpen = false;
        managedTeamPickOpen = false;
    }
    function pickUserTarget(u: User) {
        addTarget("user", u.id, `${u.name} (${u.email})`);
        userTargetPickQuery = "";
        userTargetPickOpen = false;
    }
    function pickManagedTeam(team: ManagedTeam) {
        for (const m of team.members || []) {
            addTarget("user", m.id, `${m.name} (${m.email})`);
        }
        managedTeamPickQuery = "";
        managedTeamPickOpen = false;
    }
    function addSharer(u: User) {
        if (!u?.id || sharerUserIds.includes(u.id)) return;
        sharerUserIds = [...sharerUserIds, u.id];
        sharerPickQuery = "";
        sharerPickOpen = false;
    }
    function removeSharer(id: number) {
        sharerUserIds = sharerUserIds.filter((x) => x !== id);
    }

    onMount(async () => {
        const [t, dir] = await Promise.all([
            api<Template[]>("GET", "/api/v1/templates"),
            api<{ users?: User[]; managedTeams?: ManagedTeam[] }>(
                "GET",
                "/api/v1/planning/directory",
            ),
        ]);
        templates = t;
        users = dir.users || [];
        managedTeams = dir.managedTeams || [];
        form.templateId = templates[0]?.id || 0;

        const onMouseDown = (ev: MouseEvent) => {
            const el = ev.target;
            if (!(el instanceof Element)) return;
            if (!el.closest('[data-pick-root="user-target"]'))
                userTargetPickOpen = false;
            if (!el.closest('[data-pick-root="managed-team"]'))
                managedTeamPickOpen = false;
            if (!el.closest('[data-pick-root="template"]'))
                templatePickOpen = false;
            if (!el.closest('[data-pick-root="sharer"]'))
                sharerPickOpen = false;
        };
        window.addEventListener("mousedown", onMouseDown);
        return () => window.removeEventListener("mousedown", onMouseDown);
    });

    async function submit() {
        error = "";
        if (!form.templateId) {
            error = "テンプレートを選択してください。";
            return;
        }
        if (!mentorId) {
            error = "メンターを選択してください。";
            return;
        }
        if (!userTargets.length) {
            error = "メンティを選択してください。";
            return;
        }
        if (!form.implementedAt) {
            error = "実施日を入力してください。";
            return;
        }
        const menteeIds = userTargets.map((t) => t.targetId);
        if (menteeIds.some((id) => id === mentorId)) {
            error = "メンターと同じユーザをメンティに指定できません。";
            return;
        }
        submitting = true;
        try {
            const sharerPayload = sharerUserIds.length
                ? sharerUserIds
                : undefined;
            await api("POST", "/api/v1/meetings", {
                mentorId,
                menteeIds,
                templateId: form.templateId,
                implementedAt: form.implementedAt,
                sharerUserIds: sharerPayload,
            });
            await goto("/planning/meetings");
        } catch (e) {
            error = e instanceof Error ? e.message : "作成に失敗しました。";
        } finally {
            submitting = false;
        }
    }
</script>

<svelte:window onkeydown={(e) => e.key === "Escape" && closePicks()} />

<svelte:head><title>1on1作成 - Growlink</title></svelte:head>

<main class="mx-auto w-full max-w-3xl min-w-0 p-4 md:p-8">
    <h1 class="font-display text-2xl font-semibold">1on1ミーティング作成</h1>
    <p class="mt-1 text-sm text-base-content/60">
        メンター・メンティ・テンプレート・実施日を指定して 1on1 を計画します
    </p>
    <div
        class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm"
    >
        {#if error}
            <div class="alert alert-error text-sm">{error}</div>
        {/if}
        <div class="grid items-start gap-4 md:grid-cols-2">
            <div class="form-control">
                <span class="label-text font-medium">メンター</span>
                <SearchableSelect
                    pickRootId="meeting-new-mentor"
                    items={users}
                    bind:value={mentorId}
                    labelOf={(u) => `${u.name} (${u.email})`}
                    selectedLabel={userLabel}
                    placeholder="名前で検索"
                />
            </div>
            <div>
                <span class="label-text font-medium">メンティ</span>
                <div class="form-control" data-pick-root="user-target">
                    <UserChipSearchField
                        chips={menteeChips}
                        placeholder={menteeChips.length
                            ? "追加で検索"
                            : "名前で検索して追加"}
                        bind:query={userTargetPickQuery}
                        bind:open={userTargetPickOpen}
                        filteredUsers={userTargetPickFiltered}
                        onPick={pickUserTarget}
                        onRemove={(id) => removeTarget("user", id)}
                    />
                </div>
                {#if managedTeams.length > 0}
                    <div
                        class="form-control mt-3"
                        data-pick-root="managed-team"
                    >
                        <span class="label-text font-medium"
                            >管理チームからまとめて追加</span
                        >
                        <div class="relative mt-1">
                            <input
                                type="text"
                                class="input input-bordered w-full placeholder-current"
                                placeholder="チーム名で検索"
                                autocomplete="off"
                                value={managedTeamPickOpen
                                    ? managedTeamPickQuery
                                    : ""}
                                onfocus={() => (managedTeamPickOpen = true)}
                                oninput={(e) => {
                                    managedTeamPickQuery =
                                        e.currentTarget.value;
                                    managedTeamPickOpen = true;
                                }}
                            />
                            {#if managedTeamPickOpen}
                                <ul
                                    class="absolute left-0 top-full z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg"
                                >
                                    {#each managedTeamPickFiltered as team (team.id)}
                                        <li>
                                            <button
                                                type="button"
                                                class="btn btn-ghost btn-sm w-full justify-start font-normal"
                                                onclick={() =>
                                                    pickManagedTeam(team)}
                                                >{team.name}（{team.members
                                                    ?.length ??
                                                    0}人）</button
                                            >
                                        </li>
                                    {/each}
                                </ul>
                            {/if}
                        </div>
                    </div>
                {/if}
            </div>
        </div>
        <div class="form-control" data-pick-root="template">
            <span class="label-text font-medium">テンプレート</span>
            <div class="relative mt-1">
                <input
                    type="text"
                    class="input input-bordered w-full pr-10 placeholder-current"
                    autocomplete="off"
                    placeholder={form.templateId ? "" : "名前で検索して選択"}
                    value={templatePickOpen
                        ? templatePickQuery
                        : form.templateId
                          ? templates.find((t) => t.id === form.templateId)
                                ?.title || ""
                          : ""}
                    onfocus={() => {
                        templatePickOpen = true;
                        if (form.templateId) templatePickQuery = "";
                    }}
                    oninput={(e) => {
                        templatePickQuery = e.currentTarget.value;
                        templatePickOpen = true;
                    }}
                />
                {#if form.templateId}
                    <button
                        type="button"
                        class="btn btn-ghost btn-xs btn-circle absolute right-1 top-1/2 -translate-y-1/2"
                        aria-label="クリア"
                        onclick={() => {
                            form.templateId = 0;
                            templatePickQuery = "";
                            templatePickOpen = false;
                        }}>×</button
                    >
                {/if}
                {#if templatePickOpen}
                    <ul
                        class="absolute left-0 top-full z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg"
                    >
                        {#each templatePickFiltered as t (t.id)}
                            <li>
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-sm w-full justify-start font-normal"
                                    onclick={() => {
                                        form.templateId = t.id;
                                        templatePickQuery = "";
                                        templatePickOpen = false;
                                    }}>{t.title}</button
                                >
                            </li>
                        {/each}
                        {#if templatePickFiltered.length === 0}
                            <li class="px-4 py-3 text-sm text-base-content/50">
                                該当なし
                            </li>
                        {/if}
                    </ul>
                {/if}
            </div>
        </div>
        <div class="form-control" data-pick-root="sharer">
            <span class="label-text font-medium">共有者（任意・複数）</span>
            <UserChipSearchField
                chips={sharerChips}
                placeholder={sharerChips.length
                    ? "追加で検索"
                    : "名前で検索して追加"}
                bind:query={sharerPickQuery}
                bind:open={sharerPickOpen}
                filteredUsers={sharerPickFiltered}
                onPick={addSharer}
                onRemove={removeSharer}
            />
        </div>
        <label class="form-control">
            <span class="label-text font-medium">実施日</span>
            <input
                class="input input-bordered"
                type="date"
                bind:value={form.implementedAt}
            />
        </label>
        <div class="flex flex-wrap gap-2 pt-2">
            <button
                type="button"
                class="btn btn-primary"
                disabled={submitting}
                onclick={submit}>{submitting ? "作成中…" : "作成"}</button
            >
            <a href="/planning/meetings" class="btn btn-ghost">戻る</a>
        </div>
    </div>
</main>

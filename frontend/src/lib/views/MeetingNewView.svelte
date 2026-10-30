<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { api } from "$lib/api";
    import { filterItems } from "$lib/searchable-select";
    import {
        matchPairingIds,
        type MeetingPairingItem,
    } from "$lib/meeting-pair-match";
    import { expandUserPairings } from "$lib/pairing-helpers";
    import {
        buildUserDepartmentMaps,
        buildUserTeamMaps,
        departmentNameForUser,
        teamNameForUser,
    } from "$lib/planning-directory";
    import UserChipSearchField from "$lib/components/UserChipSearchField.svelte";
    import SearchableSelect from "$lib/components/SearchableSelect.svelte";

    type Target = { kind: string; targetId: number; label: string };
    type Template = { id: number; title?: string };
    type User = {
        id: number;
        name: string;
        email: string;
        departmentNames?: string;
    };
    type ManagedTeam = { id: number; name: string; members?: User[] };
    type CreateMode = "planned" | "pair" | "monthly";

    function todayISO() {
        const d = new Date();
        return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
    }
    function currentMonthISO() {
        const d = new Date();
        return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}`;
    }

    let pairings = $state<MeetingPairingItem[]>([]);
    let pairingsLoaded = $state(false);
    let pairingsLoading = $state(false);
    let selectedPairingIds = $state<number[]>([]);
    let pairSearch = $state("");
    let templates = $state<Template[]>([]);
    let users = $state<User[]>([]);
    let managedTeams = $state<ManagedTeam[]>([]);
    let mode = $state<CreateMode>("planned");
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
        baseMonth: currentMonthISO(),
        week: 2,
        weekday: 3,
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

    let deptMaps = $derived(buildUserDepartmentMaps(users));
    let teamMaps = $derived(buildUserTeamMaps(users));
    let pairTableRows = $derived(expandUserPairings(pairings));
    let filteredPairTableRows = $derived.by(() => {
        const q = pairSearch.trim().toLowerCase();
        if (!q) return pairTableRows;
        return pairTableRows.filter((row) => {
            const dept = departmentNameForUser(row.menteeUserId, deptMaps);
            return (
                row.mentorName.toLowerCase().includes(q) ||
                row.menteeName.toLowerCase().includes(q) ||
                dept.toLowerCase().includes(q)
            );
        });
    });
    let displayedPairingIds = $derived([
        ...new Set(filteredPairTableRows.map((r) => r.pairing.id)),
    ]);
    let allPairRowsSelected = $derived(
        displayedPairingIds.length > 0 &&
            displayedPairingIds.every((id) => selectedPairingIds.includes(id)),
    );

    let matchedPairingCount = $derived(
        matchPairingIds(
            pairings,
            mentorId ? [mentorId] : [],
            targets.map((t) => ({ kind: t.kind, targetId: t.targetId })),
        ).length,
    );

    type SelectItem = { id: number };

    function userLabel(id: number) {
        const u = users.find((x) => x.id === id);
        return u ? u.name : `#${id}`;
    }

    function userLabelOf(item: SelectItem) {
        const u = users.find((x) => x.id === item.id);
        return u ? `${u.name} (${u.email})` : `#${item.id}`;
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
    function submitLabel() {
        if (submitting) return "作成中…";
        if (mode === "planned") return "作成";
        if (mode === "pair") return "作成";
        return "定期作成";
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
    function resolvePairingIds(): number[] {
        return matchPairingIds(
            pairings,
            mentorId ? [mentorId] : [],
            targets.map((t) => ({ kind: t.kind, targetId: t.targetId })),
        );
    }
    function toggleAllPairRows(checked: boolean) {
        selectedPairingIds = checked ? [...displayedPairingIds] : [];
    }

    async function loadPairingsIfNeeded() {
        if (pairingsLoaded || pairingsLoading) return;
        pairingsLoading = true;
        try {
            pairings = await api<MeetingPairingItem[]>(
                "GET",
                "/api/v1/pairings",
            );
            pairingsLoaded = true;
        } finally {
            pairingsLoading = false;
        }
    }

    async function onModeChange(next: CreateMode) {
        mode = next;
        if (next === "pair") {
            selectedPairingIds = [];
            await loadPairingsIfNeeded();
        } else if (next === "monthly") {
            await loadPairingsIfNeeded();
        }
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
    });

    async function submit() {
        error = "";
        if (!form.templateId) {
            error = "テンプレートを選択してください。";
            return;
        }
        submitting = true;
        try {
            const sharerPayload = sharerUserIds.length
                ? sharerUserIds
                : undefined;

            if (mode === "planned") {
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
                console.log(
                    {
                        mentorId,
                        menteeIds,
                        templateId: form.templateId,
                        implementedAt: form.implementedAt,
                        sharerUserIds: sharerPayload,
                    },
                    form,
                );
                await api("POST", "/api/v1/meetings", {
                    mentorId,
                    menteeIds,
                    templateId: form.templateId,
                    implementedAt: form.implementedAt,
                    sharerUserIds: sharerPayload,
                });
            } else if (mode === "pair") {
                await loadPairingsIfNeeded();
                if (!pairTableRows.length) {
                    error = "ペアがありません。";
                    return;
                }
                if (!selectedPairingIds.length) {
                    error = "ペアを選択してください。";
                    return;
                }
                if (!form.implementedAt) {
                    error = "実施日を入力してください。";
                    return;
                }
                const pairingIds = [...new Set(selectedPairingIds)];
                await api("POST", "/api/v1/meetings/all", {
                    templateId: form.templateId,
                    implementedAt: form.implementedAt,
                    sharerUserIds: sharerPayload,
                    pairingIds,
                });
            } else {
                if (!mentorId) {
                    error = "メンターを選択してください。";
                    return;
                }
                if (!userTargets.length) {
                    error = "メンティを選択してください。";
                    return;
                }
                await loadPairingsIfNeeded();
                const pairingIds = resolvePairingIds();
                if (!pairingIds.length) {
                    error = "該当するペアがありません。";
                    return;
                }
                if (!form.baseMonth) {
                    error = "基準月を入力してください。";
                    return;
                }
                if (form.week < 1 || form.week > 5) {
                    error = "第N週は 1〜5 で指定してください。";
                    return;
                }
                await api("POST", "/api/v1/meetings/monthly", {
                    templateId: form.templateId,
                    baseDate: `${form.baseMonth}-01`,
                    week: form.week,
                    weekday: form.weekday,
                    sharerUserIds: sharerPayload,
                    pairingIds,
                });
            }
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

<main class="mx-auto w-full max-w-6xl min-w-0 p-4 md:p-8">
    <h1 class="font-display text-2xl font-semibold">1on1ミーティング作成</h1>
    <p class="mt-1 text-sm text-base-content/60">
        作成方法を選び、メンター・メンティ・テンプレートを指定します
    </p>
    <div class="mt-4 flex flex-wrap gap-2">
        <button
            type="button"
            class="btn btn-sm"
            class:btn-primary={mode === "planned"}
            class:btn-outline={mode !== "planned"}
            onclick={() => onModeChange("planned")}>作成</button>
        <button
            type="button"
            class="btn btn-sm"
            class:btn-primary={mode === "pair"}
            class:btn-outline={mode !== "pair"}
            onclick={() => onModeChange("pair")}>ペア指定</button>
        <button
            type="button"
            class="btn btn-sm"
            class:btn-primary={mode === "monthly"}
            class:btn-outline={mode !== "monthly"}
            onclick={() => onModeChange("monthly")}>定期面談作成</button>
    </div>
    <div
        class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
        {#if error}
            <div class="alert alert-error text-sm">{error}</div>
        {/if}
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
                    }} />
                {#if form.templateId}
                    <button
                        type="button"
                        class="btn btn-ghost btn-xs btn-circle absolute right-1 top-1/2 -translate-y-1/2"
                        aria-label="クリア"
                        onclick={() => {
                            form.templateId = 0;
                            templatePickQuery = "";
                            templatePickOpen = false;
                        }}>×</button>
                {/if}
                {#if templatePickOpen}
                    <ul
                        class="absolute left-0 top-full z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg">
                        {#each templatePickFiltered as t (t.id)}
                            <li>
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-sm w-full justify-start font-normal"
                                    onclick={() => {
                                        form.templateId = t.id;
                                        templatePickQuery = "";
                                        templatePickOpen = false;
                                    }}>{t.title}</button>
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
        {#if mode === "pair"}
            <div>
                <span class="label-text font-medium">ペアを選択</span>
                <input
                    type="search"
                    class="input input-bordered mt-2 w-full"
                    placeholder="メンター・メンティ・部署で検索"
                    bind:value={pairSearch} />
                <div
                    class="mt-2 overflow-hidden rounded-2xl border border-base-300">
                    {#if pairingsLoading}
                        <div class="p-6 space-y-3">
                            {#each Array(4) as _, i (i)}
                                <div class="skeleton h-10 w-full rounded-xl">
                                </div>
                            {/each}
                        </div>
                    {:else if pairTableRows.length === 0}
                        <div
                            class="p-8 text-center text-sm text-base-content/50">
                            ペアがありません
                        </div>
                    {:else if filteredPairTableRows.length === 0}
                        <div
                            class="p-8 text-center text-sm text-base-content/50">
                            該当するペアがありません
                        </div>
                    {:else}
                        <div class="max-h-80 overflow-y-auto">
                            <table class="table table-pin-rows">
                                <thead>
                                    <tr class="text-base-content/70">
                                        <th class="w-10">
                                            <input
                                                type="checkbox"
                                                class="checkbox checkbox-sm"
                                                checked={allPairRowsSelected}
                                                onchange={(e) =>
                                                    toggleAllPairRows(
                                                        (
                                                            e.currentTarget as HTMLInputElement
                                                        ).checked,
                                                    )} />
                                        </th>
                                        <th>メンター</th>
                                        <th>メンティ</th>
                                        <th>部署</th>
                                        <th>チーム</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {#each filteredPairTableRows as row (row.rowKey)}
                                        <tr class="hover:bg-base-200/50">
                                            <td>
                                                <input
                                                    type="checkbox"
                                                    class="checkbox checkbox-sm"
                                                    value={row.pairing.id}
                                                    bind:group={
                                                        selectedPairingIds
                                                    } />
                                            </td>
                                            <td>{row.mentorName}</td>
                                            <td>{row.menteeName}</td>
                                            <td
                                                >{departmentNameForUser(
                                                    row.menteeUserId,
                                                    deptMaps,
                                                ) || "—"}</td>
                                            <td
                                                >{teamNameForUser(
                                                    row.menteeUserId,
                                                    teamMaps,
                                                )}</td>
                                        </tr>
                                    {/each}
                                </tbody>
                            </table>
                        </div>
                    {/if}
                </div>
            </div>
        {:else}
            <div class="grid items-end gap-4 md:grid-cols-3">
                <div class="form-control">
                    <span class="label-text font-medium">メンター</span>
                    <SearchableSelect
                        pickRootId="meeting-new-mentor"
                        items={users}
                        bind:value={mentorId}
                        labelOf={userLabelOf}
                        selectedLabel={userLabel}
                        placeholder="名前で検索" />
                </div>
                <div class="md:col-span-2">
                    <span class="label-text font-medium">メンティ</span>
                    <div class="form-control" data-pick-root="user-target">
                        <UserChipSearchField
                            chips={menteeChips}
                            className="border border-base-content border-opacity-25 h-12"
                            placeholder={menteeChips.length
                                ? "追加で検索"
                                : "名前で検索して追加"}
                            bind:query={userTargetPickQuery}
                            bind:open={userTargetPickOpen}
                            filteredUsers={userTargetPickFiltered}
                            onPick={pickUserTarget}
                            onRemove={(id) => removeTarget("user", id)} />
                    </div>
                    {#if managedTeams.length > 0}
                        <div
                            class="form-control mt-3"
                            data-pick-root="managed-team">
                            <span class="label-text font-medium"
                                >管理チームからまとめて追加</span>
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
                                    }} />
                                {#if managedTeamPickOpen}
                                    <ul
                                        class="absolute left-0 top-full z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg">
                                        {#each managedTeamPickFiltered as team (team.id)}
                                            <li>
                                                <button
                                                    type="button"
                                                    class="btn btn-ghost btn-sm w-full justify-start font-normal"
                                                    onclick={() =>
                                                        pickManagedTeam(team)}
                                                    >{team.name}（{team.members
                                                        ?.length ??
                                                        0}人）</button>
                                            </li>
                                        {/each}
                                    </ul>
                                {/if}
                            </div>
                        </div>
                    {/if}
                </div>
            </div>
            {#if mode === "monthly" && mentorId && userTargets.length}
                <p class="text-sm text-base-content/60">
                    一致するペア: <span class="font-medium text-base-content"
                        >{matchedPairingCount}</span> 件（各ペアに面談を1件作成）
                </p>
            {/if}
        {/if}
        <div class="form-control" data-pick-root="sharer">
            <span class="label-text font-medium">共有者（任意・複数）</span>
            <UserChipSearchField
                chips={sharerChips}
                className="border border-base-content border-opacity-25 h-12"
                placeholder={sharerChips.length
                    ? "追加で検索"
                    : "名前で検索して追加"}
                bind:query={sharerPickQuery}
                bind:open={sharerPickOpen}
                filteredUsers={sharerPickFiltered}
                onPick={addSharer}
                onRemove={removeSharer} />
        </div>
        {#if mode !== "monthly"}
            <label class="form-control">
                <span class="label-text font-medium">実施日</span>
                <input
                    class="input input-bordered"
                    type="date"
                    bind:value={form.implementedAt} />
            </label>
        {:else}
            <div class="flex justify-between gap-2 space-y-4 items-baseline">
                <label class="form-control w-2/6">
                    <span class="label-text font-medium">年月（から）</span>
                    <input
                        class="input input-bordered"
                        type="month"
                        bind:value={form.baseMonth} />
                </label>
                <label class="form-control w-2/6">
                    <span class="label-text font-medium">年月（まで）</span>
                    <input
                        class="input input-bordered"
                        type="month"
                        bind:value={form.baseMonth} />
                </label>
                <label class="form-control w-1/6">
                    <span class="label-text font-medium">実施週</span>
                    <input
                        class="input input-bordered"
                        type="number"
                        min="1"
                        max="5"
                        bind:value={form.week} />
                </label>
                <label class="form-control w-1/6">
                    <span class="label-text font-medium">曜日</span>
                    <select
                        class="select select-bordered"
                        bind:value={form.weekday}>
                        <option value={0}>日曜</option>
                        <option value={1}>月曜</option>
                        <option value={2}>火曜</option>
                        <option value={3}>水曜</option>
                        <option value={4}>木曜</option>
                        <option value={5}>金曜</option>
                        <option value={6}>土曜</option>
                    </select>
                </label>
            </div>
        {/if}
        <div class="flex flex-wrap gap-2 pt-2">
            <button
                type="button"
                class="btn btn-primary"
                disabled={submitting}
                onclick={submit}>{submitLabel()}</button>
            <a href="/planning/meetings" class="btn btn-ghost">戻る</a>
        </div>
    </div>
</main>

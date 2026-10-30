<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import UserChipSearchField from "$lib/components/UserChipSearchField.svelte";
    import SearchableSelect from "$lib/components/SearchableSelect.svelte";
    import { filterItems } from "$lib/searchable-select";

    type SelectItem = { id: number };
    type Target = { kind: string; targetId: number; label: string };
    type DirectoryUser = {
        id: number;
        name: string;
        email: string;
        roleName?: string;
        departmentNames?: string;
        teamNames?: string;
        departmentId?: number;
    };
    type DirectoryTeam = { id: number; name: string; departmentId: number };
    type PairingInitial = {
        uuid?: string;
        mentor?: { id?: number };
        mentee?: { id?: number; name?: string };
        targets?: { kind: string; targetId: number; label?: string }[];
    };

    let {
        mode,
        initialPairing = null,
    }: {
        mode: "create" | "edit";
        initialPairing?: PairingInitial | null;
    } = $props();

    let mentorId = $state(0);
    let targets = $state<Target[]>([]);
    let pairingUUID = $state("");
    let users = $state<DirectoryUser[]>([]);
    let departments = $state<{ id: number; name: string }[]>([]);
    let teams = $state<DirectoryTeam[]>([]);

    let nameSearch = $state("");
    let filterDeptId = $state(0);
    let filterTeamId = $state(0);
    let sortKey = $state<"role" | "name" | "department" | "team">("name");
    let sortDir = $state<"asc" | "desc">("asc");
    let selectedMenteeIds = $state<Set<number>>(new Set());

    let userTargetPickOpen = $state(false);
    let userTargetPickQuery = $state("");

    let userTargets = $derived(targets.filter((t) => t.kind === "user"));
    let menteeChips = $derived(
        userTargets.map((t) => ({
            id: t.targetId,
            name: displayNameForTarget(t),
        })),
    );

    function displayNameForTarget(t: Target) {
        const u = users.find((x) => x.id === t.targetId);
        if (u) return u.name;
        const idx = t.label.indexOf(" (");
        if (idx > 0) return t.label.slice(0, idx);
        return t.label || `#${t.targetId}`;
    }
    function userLabel(id: number) {
        const u = users.find((x) => x.id === id);
        return u ? u.name : `#${id}`;
    }
    function userLabelOf(item: SelectItem) {
        const u = users.find((x) => x.id === item.id);
        return u ? `${u.name} ${u.email}` : `#${item.id}`;
    }
    function targetKey(kind: string, targetId: number) {
        return `${kind}:${targetId}`;
    }
    function addUserTarget(u: DirectoryUser) {
        const key = targetKey("user", u.id);
        if (targets.some((t) => targetKey(t.kind, t.targetId) === key)) return;
        targets = [
            ...targets,
            { kind: "user", targetId: u.id, label: `${u.name} (${u.email})` },
        ];
    }

    let teamOptions = $derived.by(() => {
        if (!filterDeptId) return teams;
        return teams.filter((t) => t.departmentId === filterDeptId);
    });

    $effect(() => {
        if (!filterTeamId) return;
        const ok = teamOptions.some((t) => t.id === filterTeamId);
        if (!ok) filterTeamId = 0;
    });

    let userTargetPickFiltered = $derived.by(() => {
        const picked = new Set(userTargets.map((t) => t.targetId));
        const list = users.filter((u) => !picked.has(u.id));
        return filterItems(
            list,
            userTargetPickQuery,
            (u) => `${u.name} (${u.email})`,
        );
    });

    let tableRows = $derived.by(() => {
        let list = users.filter((u) => u.id !== mentorId);
        const q = nameSearch.trim().toLowerCase();
        if (q) list = list.filter((u) => u.name.toLowerCase().includes(q));
        if (filterDeptId) {
            list = list.filter((u) => u.departmentId === filterDeptId);
        }
        if (filterTeamId) {
            const team = teams.find((t) => t.id === filterTeamId);
            const teamName = team?.name || "";
            list = list.filter((u) => (u.teamNames || "").includes(teamName));
        }
        const dir = sortDir === "asc" ? 1 : -1;
        list.sort((a, b) => {
            let av = "";
            let bv = "";
            switch (sortKey) {
                case "role":
                    av = a.roleName || "";
                    bv = b.roleName || "";
                    break;
                case "department":
                    av = a.departmentNames || "";
                    bv = b.departmentNames || "";
                    break;
                case "team":
                    av = a.teamNames || "";
                    bv = b.teamNames || "";
                    break;
                default:
                    av = a.name;
                    bv = b.name;
            }
            return dir * av.localeCompare(bv, "ja");
        });
        return list;
    });

    function toggleMentee(id: number) {
        const next = new Set(selectedMenteeIds);
        if (next.has(id)) next.delete(id);
        else next.add(id);
        selectedMenteeIds = next;
    }

    function switchMentee(id: number, flag: boolean) {
        const next = new Set(selectedMenteeIds);
        if (!flag) next.delete(id);
        else next.add(id);
        selectedMenteeIds = next;
    }

    function toggleSort(key: "role" | "name" | "department" | "team") {
        if (sortKey === key) sortDir = sortDir === "asc" ? "desc" : "asc";
        else {
            sortKey = key;
            sortDir = "asc";
        }
    }

    function toggleAllUsers(ev: Event) {
        let target = ev.currentTarget;
        if (!(target instanceof HTMLInputElement)) return;
        tableRows.forEach((t) => {
            switchMentee(t.id, target.checked);
        });
    }

    function sortMark(key: string) {
        if (sortKey !== key) return "";
        return sortDir === "asc" ? " ↑" : " ↓";
    }

    $effect(() => {
        if (!initialPairing) return;
        pairingUUID = initialPairing.uuid || "";
        mentorId = initialPairing.mentor?.id || 0;
        const loaded = (initialPairing.targets || [])
            .filter((t) => t.kind === "user")
            .map((t) => ({
                kind: "user",
                targetId: t.targetId,
                label: t.label || "",
            }));
        if (
            initialPairing.mentee?.id &&
            !loaded.some((t) => t.targetId === initialPairing.mentee?.id)
        ) {
            loaded.push({
                kind: "user",
                targetId: initialPairing.mentee.id,
                label: initialPairing.mentee.name || "",
            });
        }
        targets = loaded;
    });

    onMount(async () => {
        const data = await api<{
            users?: DirectoryUser[];
            departments?: { id: number; name: string }[];
            teams?: DirectoryTeam[];
        }>("GET", "/api/v1/planning/directory");
        users = data.users || [];
        departments = data.departments || [];
        teams = (data.teams || []).map((t) => ({
            id: t.id,
            name: t.name,
            departmentId: t.departmentId,
        }));

        const onMouseDown = (ev: MouseEvent) => {
            const t = ev.target;
            if (!(t instanceof Element)) return;
            if (!t.closest('[data-pick-root="user-target"]'))
                userTargetPickOpen = false;
        };
        window.addEventListener("mousedown", onMouseDown);
        // return () => window.removeEventListener("mousedown", onMouseDown);
    });

    function pickUserTarget(u: DirectoryUser) {
        addUserTarget(u);
        userTargetPickQuery = "";
        userTargetPickOpen = false;
    }
    function removeUserTarget(targetId: number) {
        targets = targets.filter(
            (t) => !(t.kind === "user" && t.targetId === targetId),
        );
    }

    async function requestSave() {
        let userOnly: { kind: string; targetId: number }[];
        if (mode === "create") {
            userOnly = [...selectedMenteeIds].map((id) => ({
                kind: "user",
                targetId: id,
            }));
        } else {
            userOnly = targets
                .filter((t) => t.kind === "user")
                .map((t) => ({ kind: "user", targetId: t.targetId }));
        }
        if (!mentorId || !userOnly.length) {
            showToast("メンターとメンティを選択してください", "error");
            return;
        }
        const body = { mentorId, targets: userOnly };
        try {
            if (mode === "create") {
                await api("POST", "/api/v1/pairings", body);
            } else {
                await api("PUT", `/api/v1/pairings/${pairingUUID}`, body);
            }
            showToast("保存しました");
            await goto("/planning/pairs/users");
        } catch {
            showToast("保存に失敗しました", "error");
        }
    }
</script>

<svelte:head>
    <title>{mode === "create" ? "ペア作成" : "ペア編集"} - Growlink</title>
</svelte:head>

<main class="mx-auto w-full max-w-5xl min-w-0 p-4 md:p-8">
    <h1 class="font-display text-2xl font-semibold">
        {mode === "create" ? "ペア新規作成" : "ペア編集"}
    </h1>
    <p class="mt-1 text-sm text-base-content/60">
        メンターとメンティを{mode === "create" ? "登録" : "変更"}します
    </p>
    <div
        class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
        <div class="form-control max-w-md">
            <span class="label-text font-medium">メンター</span>
            <SearchableSelect
                pickRootId="pairing-form-mentor"
                items={users}
                bind:value={mentorId}
                labelOf={userLabelOf}
                selectedLabel={userLabel}
                placeholder="名前で検索" />
        </div>

        {#if mode === "create"}
            <div class="flex flex-wrap gap-3">
                <label
                    class="input input-bordered flex items-center gap-2 max-w-xs">
                    <input
                        type="search"
                        class="grow"
                        placeholder="氏名で検索"
                        bind:value={nameSearch} />
                </label>
                <select
                    class="select select-bordered select-sm"
                    bind:value={filterDeptId}>
                    <option value={0}>部署: すべて</option>
                    {#each departments as d (d.id)}
                        <option value={d.id}>{d.name}</option>
                    {/each}
                </select>
                <select
                    class="select select-bordered select-sm"
                    bind:value={filterTeamId}>
                    <option value={0}>チーム: すべて</option>
                    {#each teamOptions as t (t.id)}
                        <option value={t.id}>{t.name}</option>
                    {/each}
                </select>
            </div>
            <div class="overflow-x-auto rounded-2xl border border-base-300">
                <table class="table table-sm">
                    <thead>
                        <tr>
                            <th>
                                <input
                                    class="checkbox"
                                    type="checkbox"
                                    onclick={(e: Event) => toggleAllUsers(e)} />
                            </th>
                            <th>
                                <button
                                    type="button"
                                    class="font-semibold"
                                    onclick={() => toggleSort("role")}
                                    >役職{sortMark("role")}</button>
                            </th>
                            <th>
                                <button
                                    type="button"
                                    class="font-semibold"
                                    onclick={() => toggleSort("name")}
                                    >氏名{sortMark("name")}</button>
                            </th>
                            <th>
                                <button
                                    type="button"
                                    class="font-semibold"
                                    onclick={() => toggleSort("department")}
                                    >部署{sortMark("department")}</button>
                            </th>
                            <th>
                                <button
                                    type="button"
                                    class="font-semibold"
                                    onclick={() => toggleSort("team")}
                                    >チーム{sortMark("team")}</button>
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each tableRows as u (u.id)}
                            <tr class="hover:bg-base-200/40">
                                <td>
                                    <input
                                        type="checkbox"
                                        class="checkbox checkbox-sm"
                                        checked={selectedMenteeIds.has(u.id)}
                                        onchange={() => toggleMentee(u.id)} />
                                </td>
                                <td>{u.roleName || "—"}</td>
                                <td>{u.name}</td>
                                <td>{u.departmentNames || "—"}</td>
                                <td>{u.teamNames || "—"}</td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
            <p class="text-xs text-base-content/60">
                選択中: {selectedMenteeIds.size} 人
            </p>
        {:else}
            <div class="form-control" data-pick-root="user-target">
                <span class="label-text font-medium">メンティ</span>
                <UserChipSearchField
                    chips={menteeChips}
                    placeholder={menteeChips.length
                        ? "追加で検索"
                        : "名前で検索して追加"}
                    bind:query={userTargetPickQuery}
                    bind:open={userTargetPickOpen}
                    filteredUsers={userTargetPickFiltered}
                    onPick={pickUserTarget}
                    onRemove={removeUserTarget} />
            </div>
        {/if}
        <div class="flex flex-wrap gap-2 pt-2">
            <button type="button" class="btn btn-primary" onclick={requestSave}
                >{mode === "create" ? "作成" : "保存"}</button>
            <a href="/planning/pairs/users" class="btn btn-ghost">戻る</a>
        </div>
    </div>
</main>

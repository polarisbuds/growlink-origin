<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import { getAvatar } from "$lib/session";
    import MultiCheckSelect from "$lib/components/MultiCheckSelect.svelte";
    import Icon from "$lib/components/Icon.svelte";
    import {
        type PairingListItem,
        expandUserPairings,
    } from "$lib/pairing-helpers";
    import {
        buildUserDepartmentMaps,
        departmentNameForUser,
        teamNameForUser,
        buildUserTeamMaps,
        type PlanningDirectoryUser,
    } from "$lib/planning-directory";
    import {
        compareJa,
        sortIndicator,
        toggleSortKey,
        type SortDir,
    } from "$lib/list-sort";

    type SortKey =
        | "mentor"
        | "mentorDept"
        | "mentee"
        | "menteeDept"
        | "menteeTeam";
    type ViewMode = "individual" | "department";
    type DeptSortState = { sortKey: SortKey; sortDir: SortDir };

    const DEFAULT_DEPT_SORT: DeptSortState = {
        sortKey: "mentor",
        sortDir: "asc",
    };

    let pairings = $state<PairingListItem[]>([]);
    let directoryUsers = $state<PlanningDirectoryUser[]>([]);

    let filterMentorUuids = $state<string[]>([]);
    let filterMenteeUuids = $state<string[]>([]);
    let viewMode = $state<ViewMode>("individual");
    let sortKey = $state<SortKey>("mentor");
    let sortDir = $state<SortDir>("asc");
    let deptSorts = $state<Record<string, DeptSortState>>({});
    let selected = $state<PairingListItem | null>(null);
    let deleteOpen = $state(false);
    let isSystemAdmin = $state(false);

    let deptMaps = $derived(buildUserDepartmentMaps(directoryUsers));
    let teamMaps = $derived(buildUserTeamMaps(directoryUsers));

    let displayRows = $derived(expandUserPairings(pairings));

    function avatarURL(user: PlanningDirectoryUser | undefined) {
        if (!user?.imageUri) return "";
        return `/uploads/${user.imageUri}`;
    }

    function userOptionLabel(u: PlanningDirectoryUser) {
        return `${u.name} (${u.email})`;
    }

    function sortRows(
        rows: ReturnType<typeof expandUserPairings>,
        key: SortKey,
        dir: SortDir,
    ) {
        const list = [...rows];
        const mult = dir === "asc" ? 1 : -1;
        list.sort((a, b) => {
            let av = "";
            let bv = "";
            switch (key) {
                case "mentor":
                    av = a.mentorName;
                    bv = b.mentorName;
                    break;
                case "mentorDept":
                    av = departmentNameForUser(a.mentorUserUuid, deptMaps);
                    bv = departmentNameForUser(b.mentorUserUuid, deptMaps);
                    break;
                case "mentee":
                    av = a.menteeName;
                    bv = b.menteeName;
                    break;
                case "menteeDept":
                    av = departmentNameForUser(a.menteeUserUuid, deptMaps);
                    bv = departmentNameForUser(b.menteeUserUuid, deptMaps);
                    break;
                case "menteeTeam":
                    av = teamNameForUser(a.menteeUserUuid, teamMaps);
                    bv = teamNameForUser(b.menteeUserUuid, teamMaps);
                    break;
            }
            return mult * compareJa(av, bv);
        });
        return list;
    }

    let filteredRowsUnsorted = $derived.by(() => {
        let rows = displayRows;
        if (isSystemAdmin) {
            if (filterMentorUuids.length) {
                const set = new Set(filterMentorUuids);
                rows = rows.filter((r) => set.has(r.mentorUserUuid));
            }
            if (filterMenteeUuids.length) {
                const set = new Set(filterMenteeUuids);
                rows = rows.filter((r) => set.has(r.menteeUserUuid));
            }
        }
        return rows;
    });

    let filteredRows = $derived(
        sortRows(filteredRowsUnsorted, sortKey, sortDir),
    );

    let departmentGroups = $derived.by(() => {
        const groups = new Map<string, ReturnType<typeof expandUserPairings>>();
        for (const row of filteredRowsUnsorted) {
            const label =
                departmentNameForUser(row.menteeUserUuid, deptMaps) || "";
            if (!groups.has(label)) groups.set(label, []);
            groups.get(label)!.push(row);
        }
        const keys = [...groups.keys()].sort((a, b) => {
            if (!a) return 1;
            if (!b) return -1;
            return compareJa(a, b);
        });
        return keys.map((label) => {
            const sort = deptSorts[label] ?? DEFAULT_DEPT_SORT;
            return {
                label,
                sortKey: sort.sortKey,
                sortDir: sort.sortDir,
                rows: sortRows(
                    groups.get(label) || [],
                    sort.sortKey,
                    sort.sortDir,
                ),
            };
        });
    });

    onMount(async () => {
        isSystemAdmin = getAvatar()?.isSystemAdmin ?? false;
        const [pairingList, dir] = await Promise.all([
            api<PairingListItem[]>("GET", "/api/v1/pairings"),
            api<{ users?: PlanningDirectoryUser[] }>(
                "GET",
                "/api/v1/planning/directory",
            ),
        ]);
        pairings = pairingList;
        directoryUsers = dir.users || [];
    });

    function toggleSort(key: SortKey) {
        const next = toggleSortKey(sortKey, sortDir, key);
        sortKey = next.sortKey;
        sortDir = next.sortDir;
    }

    function toggleDeptSort(deptLabel: string, key: SortKey) {
        const current = deptSorts[deptLabel] ?? DEFAULT_DEPT_SORT;
        const next = toggleSortKey(current.sortKey, current.sortDir, key);
        deptSorts = { ...deptSorts, [deptLabel]: next };
    }

    function askDelete(p: PairingListItem) {
        selected = p;
        deleteOpen = true;
    }
    async function confirmDelete() {
        if (!selected) return;
        try {
            await api("DELETE", `/api/v1/pairings/${selected.uuid}`);
            deleteOpen = false;
            showToast("ペアを削除しました");
            pairings = await api<PairingListItem[]>("GET", "/api/v1/pairings");
        } catch {
            showToast("削除に失敗しました", "error");
        }
    }
</script>

<svelte:window
    onkeydown={(e: KeyboardEvent) =>
        e.key === "Escape" && (deleteOpen = false)} />

<div class="mt-4 flex flex-wrap items-center gap-2 justify-between">
    <span class="text-xs text-base-content/60">表示</span>
    <div class="join">
        <button
            type="button"
            class="btn btn-sm join-item"
            class:btn-primary={viewMode === "individual"}
            class:btn-outline={viewMode !== "individual"}
            onclick={() => (viewMode = "individual")}>個人</button>
        <button
            type="button"
            class="btn btn-sm join-item"
            class:btn-primary={viewMode === "department"}
            class:btn-outline={viewMode !== "department"}
            onclick={() => (viewMode = "department")}>部署</button>
    </div>
</div>

{#if isSystemAdmin}
    <div class="mt-4 grid gap-3 sm:grid-cols-2">
        <label class="form-control">
            <span class="label-text text-xs">メンター</span>
            <MultiCheckSelect
                pickRootId="pair-list-mentor"
                items={directoryUsers}
                bind:selectedUuids={filterMentorUuids}
                labelOf={userOptionLabel}
                placeholder="すべて" />
        </label>
        <label class="form-control">
            <span class="label-text text-xs">メンティ</span>
            <MultiCheckSelect
                pickRootId="pair-list-mentee"
                items={directoryUsers}
                bind:selectedUuids={filterMenteeUuids}
                labelOf={userOptionLabel}
                placeholder="すべて" />
        </label>
    </div>
{/if}

{#snippet tableHead(
    activeKey: SortKey,
    activeDir: SortDir,
    onToggle: (key: SortKey) => void,
)}
    <thead>
        <tr class="text-base-content/70">
            <th>
                <button
                    type="button"
                    class="font-semibold hover:text-primary"
                    onclick={() => onToggle("mentorDept")}
                    >部署{sortIndicator(
                        "mentorDept",
                        activeKey,
                        activeDir,
                    )}</button>
            </th>
            <th>
                <button
                    type="button"
                    class="font-semibold hover:text-primary"
                    onclick={() => onToggle("mentor")}
                    >メンター{sortIndicator(
                        "mentor",
                        activeKey,
                        activeDir,
                    )}</button>
            </th>
            <th>
                <button
                    type="button"
                    class="font-semibold hover:text-primary"
                    onclick={() => onToggle("menteeDept")}
                    >部署{sortIndicator(
                        "menteeDept",
                        activeKey,
                        activeDir,
                    )}</button>
            </th>
            <th>
                <button
                    type="button"
                    class="font-semibold hover:text-primary"
                    onclick={() => onToggle("mentee")}
                    >メンティ{sortIndicator(
                        "mentee",
                        activeKey,
                        activeDir,
                    )}</button>
            </th>
            <th>
                <button
                    type="button"
                    class="font-semibold hover:text-primary"
                    onclick={() => onToggle("menteeTeam")}
                    >チーム{sortIndicator(
                        "menteeTeam",
                        activeKey,
                        activeDir,
                    )}</button>
            </th>
            <th></th>
        </tr>
    </thead>
{/snippet}

{#snippet tableBody(rows: ReturnType<typeof expandUserPairings>)}
    <tbody>
        {#each rows as row (row.rowKey)}
            {@const mentorUser = directoryUsers.find(
                (u) => u.uuid === row.mentorUserUuid,
            )}
            {@const menteeUser = directoryUsers.find(
                (u) => u.uuid === row.menteeUserUuid,
            )}
            <tr class="hover:bg-base-200/40">
                <td class="text-sm"
                    >{departmentNameForUser(row.mentorUserUuid, deptMaps)}</td>
                <td>
                    <span class="flex gap-2 items-center">
                        <Icon
                            url={avatarURL(mentorUser)}
                            text={row.mentorName}
                            width={28}
                            height={28} />
                        {row.mentorName}
                    </span>
                </td>
                <td class="text-sm"
                    >{departmentNameForUser(row.menteeUserUuid, deptMaps)}</td>
                <td>
                    <span class="flex gap-2 items-center">
                        <Icon
                            url={avatarURL(menteeUser)}
                            text={row.menteeName}
                            width={28}
                            height={28} />
                        {row.menteeName}
                    </span>
                </td>
                <td class="text-sm">
                    {menteeUser?.teamNames || ""}
                </td>
                <td class="space-x-1 whitespace-nowrap text-right">
                    <a
                        href="/planning/pairs/{row.pairing.uuid}/edit"
                        class="btn btn-xs btn-outline">編集</a>
                    <button
                        type="button"
                        class="btn btn-xs btn-error btn-outline"
                        onclick={() => askDelete(row.pairing)}>削除</button>
                </td>
            </tr>
        {/each}
    </tbody>
{/snippet}

{#if filteredRowsUnsorted.length === 0}
    <div
        class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 p-12 text-center text-sm text-base-content/50 shadow-sm">
        該当するペアがありません
    </div>
{:else if viewMode === "individual"}
    <div
        class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
        <table class="table">
            {@render tableHead(sortKey, sortDir, toggleSort)}
            {@render tableBody(filteredRows)}
        </table>
    </div>
{:else}
    <div class="mt-6 space-y-6">
        {#each departmentGroups as group (group.label || "__empty__")}
            <section
                class="overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
                {#if group.label}
                    <div
                        class="border-b border-base-300 bg-base-200/40 px-4 py-2 text-sm font-semibold">
                        {group.label}
                    </div>
                {:else}
                    <div
                        class="border-b border-base-300 bg-base-200/40 px-4 py-2 text-sm font-semibold text-base-content/60">
                        部署未設定
                    </div>
                {/if}
                <table class="table">
                    {@render tableHead(group.sortKey, group.sortDir, (key) =>
                        toggleDeptSort(group.label, key),
                    )}
                    {@render tableBody(group.rows)}
                </table>
            </section>
        {/each}
    </div>
{/if}

<dialog class="modal" class:modal-open={deleteOpen}>
    <div class="modal-box">
        <h3 class="font-bold text-lg">ペアを削除しますか？</h3>
        <div class="modal-action">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (deleteOpen = false)}>キャンセル</button>
            <button type="button" class="btn btn-error" onclick={confirmDelete}
                >削除</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (deleteOpen = false)}
            >close</button>
    </form>
</dialog>

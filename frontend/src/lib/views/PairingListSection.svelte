<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import { getAvatar } from "$lib/session";
    import SearchableSelect from "$lib/components/SearchableSelect.svelte";
    import {
        type PairingListItem,
        expandUserPairings,
    } from "$lib/pairing-helpers";
    import {
        buildUserDepartmentMaps,
        departmentNameForUser,
        type PlanningDirectoryUser,
    } from "$lib/planning-directory";
    import {
        compareJa,
        sortIndicator,
        toggleSortKey,
        type SortDir,
    } from "$lib/list-sort";

    type SelectItem = { id: number };
    type SortKey = "mentor" | "mentorDept" | "mentee" | "menteeDept";
    type ViewMode = "individual" | "department";

    let pairings = $state<PairingListItem[]>([]);
    let directoryUsers = $state<PlanningDirectoryUser[]>([]);

    let filterMentorId = $state(0);
    let filterMenteeId = $state(0);
    let viewMode = $state<ViewMode>("individual");
    let sortKey = $state<SortKey>("mentor");
    let sortDir = $state<SortDir>("asc");
    let selected = $state<PairingListItem | null>(null);
    let deleteOpen = $state(false);
    let isSystemAdmin = $state(false);

    let deptMaps = $derived(buildUserDepartmentMaps(directoryUsers));

    let displayRows = $derived(expandUserPairings(pairings));

    function userLabel(id: number) {
        const u = directoryUsers.find((x) => x.id === id);
        return u ? u.name : `#${id}`;
    }

    function sortRows(rows: ReturnType<typeof expandUserPairings>) {
        const list = [...rows];
        const dir = sortDir === "asc" ? 1 : -1;
        list.sort((a, b) => {
            let av = "";
            let bv = "";
            switch (sortKey) {
                case "mentor":
                    av = a.mentorName;
                    bv = b.mentorName;
                    break;
                case "mentorDept":
                    av = departmentNameForUser(a.mentorUserId, deptMaps);
                    bv = departmentNameForUser(b.mentorUserId, deptMaps);
                    break;
                case "mentee":
                    av = a.menteeName;
                    bv = b.menteeName;
                    break;
                case "menteeDept":
                    av = departmentNameForUser(a.menteeUserId, deptMaps);
                    bv = departmentNameForUser(b.menteeUserId, deptMaps);
                    break;
            }
            return dir * compareJa(av, bv);
        });
        return list;
    }

    let filteredRows = $derived.by(() => {
        let rows = displayRows;
        if (isSystemAdmin) {
            if (filterMentorId) {
                rows = rows.filter((r) => r.mentorUserId === filterMentorId);
            }
            if (filterMenteeId) {
                rows = rows.filter((r) => r.menteeUserId === filterMenteeId);
            }
        }
        return sortRows(rows);
    });

    let departmentGroups = $derived.by(() => {
        const groups = new Map<string, ReturnType<typeof expandUserPairings>>();
        for (const row of filteredRows) {
            const label =
                departmentNameForUser(row.menteeUserId, deptMaps) || "";
            if (!groups.has(label)) groups.set(label, []);
            groups.get(label)!.push(row);
        }
        const keys = [...groups.keys()].sort((a, b) => {
            if (!a) return 1;
            if (!b) return -1;
            return compareJa(a, b);
        });
        return keys.map((label) => ({
            label,
            rows: groups.get(label) || [],
        }));
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

    function userLabelOf(item: SelectItem) {
        const u = directoryUsers.find((x) => x.id === item.id);
        if (!u) return `#${item.id}`;
        return `${u.name} ${u.email}`;
    }

    function askDelete(p: PairingListItem) {
        selected = p;
        deleteOpen = true;
    }
    async function confirmDelete() {
        if (!selected) return;
        try {
            await api(
                "DELETE",
                `/api/v1/pairings/${selected.uuid || selected.id}`,
            );
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

<div class="mt-4 flex flex-wrap items-center gap-2">
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
            <SearchableSelect
                pickRootId="pair-list-mentor"
                items={directoryUsers}
                bind:value={filterMentorId}
                labelOf={userLabelOf}
                selectedLabel={userLabel}
                placeholder="すべて"
                allowClear={true}
                inputClass="input input-bordered input-sm w-full placeholder-current" />
        </label>
        <label class="form-control">
            <span class="label-text text-xs">メンティ</span>
            <SearchableSelect
                pickRootId="pair-list-mentee"
                items={directoryUsers}
                bind:value={filterMenteeId}
                labelOf={userLabelOf}
                selectedLabel={userLabel}
                placeholder="すべて"
                allowClear={true}
                inputClass="input input-bordered input-sm w-full placeholder-current" />
        </label>
    </div>
{/if}

{#snippet tableHead()}
    <thead>
        <tr class="text-base-content/70">
            <th>
                <button
                    type="button"
                    class="font-semibold hover:text-primary"
                    onclick={() => toggleSort("mentor")}
                    >メンター{sortIndicator(
                        "mentor",
                        sortKey,
                        sortDir,
                    )}</button>
            </th>
            <th>
                <button
                    type="button"
                    class="font-semibold hover:text-primary"
                    onclick={() => toggleSort("mentorDept")}
                    >部署{sortIndicator(
                        "mentorDept",
                        sortKey,
                        sortDir,
                    )}</button>
            </th>
            <th>
                <button
                    type="button"
                    class="font-semibold hover:text-primary"
                    onclick={() => toggleSort("mentee")}
                    >メンティ{sortIndicator(
                        "mentee",
                        sortKey,
                        sortDir,
                    )}</button>
            </th>
            <th>
                <button
                    type="button"
                    class="font-semibold hover:text-primary"
                    onclick={() => toggleSort("menteeDept")}
                    >部署{sortIndicator(
                        "menteeDept",
                        sortKey,
                        sortDir,
                    )}</button>
            </th>
            <th></th>
        </tr>
    </thead>
{/snippet}

{#snippet tableBody(rows: ReturnType<typeof expandUserPairings>)}
    <tbody>
        {#each rows as row (row.rowKey)}
            <tr class="hover:bg-base-200/40">
                <td>{row.mentorName}</td>
                <td class="text-sm"
                    >{departmentNameForUser(row.mentorUserId, deptMaps)}</td>
                <td class="text-sm">{row.menteeName}</td>
                <td class="text-sm"
                    >{departmentNameForUser(row.menteeUserId, deptMaps)}</td>
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

{#if filteredRows.length === 0}
    <div
        class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 p-12 text-center text-sm text-base-content/50 shadow-sm">
        該当するペアがありません
    </div>
{:else if viewMode === "individual"}
    <div
        class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
        <table class="table">
            {@render tableHead()}
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
                    {@render tableHead()}
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

<script lang="ts">
    import { fade } from "svelte/transition";
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import { getAvatar } from "$lib/session";
    import MultiCheckSelect from "$lib/components/MultiCheckSelect.svelte";
    import Icon from "$lib/components/Icon.svelte";
    import PairSharerModal, {
        type PairShareTarget,
    } from "$lib/components/PairSharerModal.svelte";
    import { ChevronDown } from "lucide-svelte";
    import type { PairingDisplayRow } from "$lib/pairing-helpers";
    import {
        type PairingListItem,
        expandUserPairings,
    } from "$lib/pairing-helpers";
    import {
        buildUserDepartmentMaps,
        departmentNameForUser,
        type PlanningDirectoryUser,
    } from "$lib/planning-directory";
    import { compareJa } from "$lib/list-sort";
    import {
        buildPairingsQuery,
        lastPageForTotal,
        visiblePageNumbers,
        type PairingsPageResponse,
        type PairingListQuery,
    } from "$lib/planning-pairings-pager";

    type ViewMode = "individual" | "department";

    type Dept = { uuid: string; name: string };
    type Team = { uuid: string; name: string; departmentUuid: string };

    type SelectedPairing = {
        uuid: string;
        mentorUserUuid: string;
        menteeUserUuid: string;
        mentorName: string;
        menteeName: string;
    };

    const PAGE_SIZE = 20;

    let directoryUsers = $state<PlanningDirectoryUser[]>([]);

    let filterMentorUuids = $state<string[]>([]);
    let filterMenteeUuids = $state<string[]>([]);
    let viewMode = $state<ViewMode>("individual");
    let selectedPairings = $state<SelectedPairing[]>([]);
    let bulkDeleteOpen = $state(false);
    let shareOpen = $state(false);
    let shareTargets = $state<PairShareTarget[]>([]);
    let isSystemAdmin = $state(false);
    let departments = $state<Dept[]>([]);
    let selectDepartmentUuids = $state<string[]>([]);
    let teams = $state<Team[]>([]);
    let selectTeamUuids = $state<string[]>([]);

    let currentPage = $state(1);
    let total = $state(0);
    let pageCache = $state<Map<number, ReturnType<typeof expandUserPairings>>>(
        new Map(),
    );
    let listLoading = $state(false);
    let loadGeneration = 0;

    let departmentRows = $state<ReturnType<typeof expandUserPairings>>([]);
    let departmentLoading = $state(false);

    let deptMaps = $derived(buildUserDepartmentMaps(directoryUsers));

    let displayedRows = $derived(pageCache.get(currentPage) ?? []);
    let lastPage = $derived(lastPageForTotal(total, PAGE_SIZE));
    let pagerPages = $derived(visiblePageNumbers(currentPage, lastPage));

    let canBulkPairAction = $derived(selectedPairings.length > 0);

    function rowToSelected(row: PairingDisplayRow): SelectedPairing {
        return {
            uuid: row.pairing.uuid,
            mentorUserUuid: row.mentorUserUuid,
            menteeUserUuid: row.menteeUserUuid,
            mentorName: row.mentorName,
            menteeName: row.menteeName,
        };
    }

    function isPairSelected(uuid: string) {
        return selectedPairings.some((p) => p.uuid === uuid);
    }

    function togglePair(row: PairingDisplayRow, checked: boolean) {
        const uuid = row.pairing.uuid;
        if (checked) {
            if (isPairSelected(uuid)) return;
            selectedPairings = [...selectedPairings, rowToSelected(row)];
        } else {
            selectedPairings = selectedPairings.filter((p) => p.uuid !== uuid);
        }
    }

    function allRowsSelected(rows: PairingDisplayRow[]) {
        return (
            rows.length > 0 &&
            rows.every((r) => isPairSelected(r.pairing.uuid))
        );
    }

    function toggleAllPairsInTable(rows: PairingDisplayRow[], checked: boolean) {
        const rowUuids = rows.map((r) => r.pairing.uuid);
        if (checked) {
            const set = new Set(selectedPairings.map((p) => p.uuid));
            const next = [...selectedPairings];
            for (let i = 0; i < rows.length; i++) {
                const uuid = rowUuids[i];
                if (set.has(uuid)) continue;
                set.add(uuid);
                next.push(rowToSelected(rows[i]));
            }
            selectedPairings = next;
        } else {
            const remove = new Set(rowUuids);
            selectedPairings = selectedPairings.filter(
                (p) => !remove.has(p.uuid),
            );
        }
    }

    function filterQueryBase(): Omit<PairingListQuery, "page"> {
        return {
            mentorUuids: isSystemAdmin ? filterMentorUuids : [],
            menteeUuids: isSystemAdmin ? filterMenteeUuids : [],
            departmentUuids: isSystemAdmin ? selectDepartmentUuids : [],
            teamUuids: isSystemAdmin ? selectTeamUuids : [],
        };
    }

    function listQuery(page: number): PairingListQuery {
        return { ...filterQueryBase(), page };
    }

    function trimCache(around: number) {
        const next = new Map<number, ReturnType<typeof expandUserPairings>>();
        for (const p of [around - 1, around, around + 1]) {
            const rows = pageCache.get(p);
            if (p >= 1 && rows) next.set(p, rows);
        }
        pageCache = next;
    }

    async function fetchPaged(
        page: number,
        gen: number,
    ): Promise<PairingsPageResponse | null> {
        const qs = buildPairingsQuery(listQuery(page));
        try {
            const res = await api<PairingsPageResponse>(
                "GET",
                `/api/v1/pairings?${qs}`,
            );
            if (gen !== loadGeneration) return null;
            return res;
        } catch {
            if (gen !== loadGeneration) return null;
            return null;
        }
    }

    async function storePage(page: number, gen: number) {
        const res = await fetchPaged(page, gen);
        if (!res || gen !== loadGeneration) return;
        total = res.total;
        const next = new Map(pageCache);
        next.set(page, expandUserPairings(res.items as PairingListItem[]));
        pageCache = next;
    }

    async function prefetchNeighbors(page: number, gen: number) {
        const last = lastPageForTotal(total, PAGE_SIZE);
        const tasks: Promise<void>[] = [];
        if (page > 1) tasks.push(storePage(page - 1, gen));
        if (page < last) tasks.push(storePage(page + 1, gen));
        await Promise.all(tasks);
        if (gen !== loadGeneration) return;
        trimCache(page);
    }

    async function goToPage(page: number) {
        if (page < 1 || page > lastPageForTotal(total, PAGE_SIZE)) return;
        const gen = ++loadGeneration;
        const cached = pageCache.get(page);
        if (cached) {
            currentPage = page;
            void prefetchNeighbors(page, gen);
            return;
        }
        listLoading = true;
        currentPage = page;
        await storePage(page, gen);
        listLoading = false;
        if (gen !== loadGeneration) return;
        await prefetchNeighbors(page, gen);
    }

    async function reloadFromPageOne() {
        const gen = ++loadGeneration;
        pageCache = new Map();
        currentPage = 1;
        listLoading = true;
        await storePage(1, gen);
        listLoading = false;
        if (gen !== loadGeneration) return;
        await prefetchNeighbors(1, gen);
    }

    async function reloadDepartmentView() {
        const gen = ++loadGeneration;
        departmentLoading = true;
        const qs = buildPairingsQuery(filterQueryBase());
        try {
            const res = await api<PairingListItem[]>(
                "GET",
                qs ? `/api/v1/pairings?${qs}` : "/api/v1/pairings",
            );
            if (gen !== loadGeneration) return;
            departmentRows = expandUserPairings(res);
        } catch {
            if (gen === loadGeneration) departmentRows = [];
        } finally {
            if (gen === loadGeneration) departmentLoading = false;
        }
    }

    let departmentGroups = $derived.by(() => {
        const groups = new Map<string, ReturnType<typeof expandUserPairings>>();
        for (const row of departmentRows) {
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
        return keys.map((label) => ({
            label,
            rows: groups.get(label) || [],
        }));
    });

    let filterSignature = $derived(
        JSON.stringify({
            filterMentorUuids,
            filterMenteeUuids,
            selectDepartmentUuids,
            selectTeamUuids,
            viewMode,
        }),
    );
    let filterSignatureMounted = "";

    $effect(() => {
        const sig = filterSignature;
        if (filterSignatureMounted === "") return;
        if (sig === filterSignatureMounted) return;
        filterSignatureMounted = sig;
        if (viewMode === "individual") {
            void reloadFromPageOne();
        } else {
            void reloadDepartmentView();
        }
    });

    function avatarURL(user: PlanningDirectoryUser | undefined) {
        if (!user?.imageUri) return "";
        return `/uploads/${user.imageUri}`;
    }

    function userOptionLabel(u: PlanningDirectoryUser) {
        return `${u.name} (${u.email})`;
    }

    onMount(async () => {
        isSystemAdmin = getAvatar()?.isSystemAdmin ?? false;
        const dir = await api<{
            users?: PlanningDirectoryUser[];
            departments?: Dept[];
            teams?: Team[];
        }>("GET", "/api/v1/planning/directory");
        directoryUsers = dir.users || [];
        departments = dir.departments || [];
        teams = dir.teams || [];
        filterSignatureMounted = filterSignature;
        await reloadFromPageOne();
    });

    function rowToShareTarget(row: PairingDisplayRow): PairShareTarget {
        return {
            pairingUuid: row.pairing.uuid,
            mentorUserUuid: row.mentorUserUuid,
            menteeUserUuid: row.menteeUserUuid,
            mentorName: row.mentorName,
            menteeName: row.menteeName,
        };
    }

    function selectedToShareTargets(): PairShareTarget[] {
        const out: PairShareTarget[] = [];
        for (let i = 0; i < selectedPairings.length; i++) {
            const p = selectedPairings[i];
            out[i] = {
                pairingUuid: p.uuid,
                mentorUserUuid: p.mentorUserUuid,
                menteeUserUuid: p.menteeUserUuid,
                mentorName: p.mentorName,
                menteeName: p.menteeName,
            };
        }
        return out;
    }

    function openShare(row: PairingDisplayRow) {
        shareTargets = [rowToShareTarget(row)];
        shareOpen = true;
    }

    function openBulkShare() {
        if (selectedPairings.length === 0) return;
        shareTargets = selectedToShareTargets();
        shareOpen = true;
    }

    function closeShare() {
        shareOpen = false;
        shareTargets = [];
    }

    function askBulkDelete() {
        if (selectedPairings.length === 0) return;
        bulkDeleteOpen = true;
    }

    async function confirmBulkDelete() {
        if (selectedPairings.length === 0) return;
        const uuids = selectedPairings.map((p) => p.uuid);
        bulkDeleteOpen = false;
        try {
            for (let i = 0; i < uuids.length; i++) {
                await api("DELETE", `/api/v1/pairings/${uuids[i]}`);
            }
            selectedPairings = [];
            showToast(
                uuids.length === 1
                    ? "ペアを削除しました"
                    : `${uuids.length} 件のペアを削除しました`,
            );
            if (viewMode === "individual") {
                await reloadFromPageOne();
            } else {
                await reloadDepartmentView();
            }
        } catch {
            showToast("削除に失敗しました", "error");
            if (viewMode === "individual") {
                await reloadFromPageOne();
            } else {
                await reloadDepartmentView();
            }
        }
    }

    function setViewMode(mode: ViewMode) {
        viewMode = mode;
    }
</script>

<svelte:head><title>ペア一覧 - Growlink</title></svelte:head>

<svelte:window
    onkeydown={(e: KeyboardEvent) => {
        if (e.key !== "Escape") return;
        if (document.querySelector("[data-pick-panel]")) return;
        bulkDeleteOpen = false;
        const openModals = document.querySelectorAll(".modal.modal-open");
        if (shareOpen && openModals.length <= 1) closeShare();
    }} />

{#snippet pairBulkActions()}
    {#if selectedPairings.length > 0}
        <span class="text-sm text-base-content/60" transition:fade
            >{selectedPairings.length} 件選択中</span>
    {/if}
    <div class="dropdown dropdown-end">
        <button
            type="button"
            tabindex="0"
            class="btn btn-outline btn-sm gap-1 shrink-0">
            操作
            <ChevronDown class="h-4 w-4 opacity-60" />
        </button>
        <ul
            role="menu"
            tabindex="0"
            class="menu dropdown-content z-50 mt-2 w-44 rounded-box border border-base-300 bg-base-100 p-2 shadow-lg">
            <li>
                <button
                    type="button"
                    disabled={!canBulkPairAction}
                    onclick={openBulkShare}>共有者設定</button>
            </li>
            <li>
                <button
                    type="button"
                    class="text-error"
                    disabled={!canBulkPairAction}
                    onclick={askBulkDelete}>削除</button>
            </li>
        </ul>
    </div>
{/snippet}

<div class="mt-2 flex flex-wrap items-center gap-2 justify-end">
    <div class="join">
        <button
            type="button"
            class="btn btn-sm join-item"
            class:btn-primary={viewMode === "individual"}
            class:btn-outline={viewMode !== "individual"}
            onclick={() => setViewMode("individual")}>個人</button>
        <button
            type="button"
            class="btn btn-sm join-item"
            class:btn-primary={viewMode === "department"}
            class:btn-outline={viewMode !== "department"}
            onclick={() => setViewMode("department")}>部署</button>
    </div>
    {#if !isSystemAdmin}
        <div class="flex items-center gap-3">
            {@render pairBulkActions()}
        </div>
    {/if}
</div>

{#if isSystemAdmin}
    <div class="mt-2 flex flex-wrap items-end gap-3">
        <div class="grid flex-1 min-w-0 gap-3 sm:grid-cols-4">
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
            <label class="form-control">
                <span class="label-text text-xs">部署</span>
                <MultiCheckSelect
                    pickRootId="pair-list-dept"
                    items={departments}
                    bind:selectedUuids={selectDepartmentUuids}
                    labelOf={(item: Dept) => `${item.name}`}
                    placeholder="すべて" />
            </label>
            <label class="form-control">
                <span class="label-text text-xs">チーム</span>
                <MultiCheckSelect
                    pickRootId="pair-list-team"
                    items={teams}
                    bind:selectedUuids={selectTeamUuids}
                    labelOf={(item: Team) => item.name}
                    placeholder="すべて" />
            </label>
        </div>
        <div class="flex items-center gap-3 shrink-0 pb-0.5">
            {@render pairBulkActions()}
        </div>
    </div>
{/if}

{#snippet tableHead(rows: PairingDisplayRow[])}
    <thead>
        <tr class="text-base-content/70">
            <th class="w-10">
                <input
                    type="checkbox"
                    class="checkbox checkbox-sm"
                    checked={allRowsSelected(rows)}
                    onchange={(e) =>
                        toggleAllPairsInTable(
                            rows,
                            (e.currentTarget as HTMLInputElement).checked,
                        )} />
            </th>
            <th class="font-semibold">部署</th>
            <th class="font-semibold">メンター</th>
            <th class="font-semibold">部署</th>
            <th class="font-semibold">メンティ</th>
            <th class="font-semibold">チーム</th>
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
            <tr
                class="cursor-pointer hover:bg-base-200/40"
                onclick={() => openShare(row)}>
                <!-- svelte-ignore a11y_no_static_element_interactions -->
                <td onclick={(e) => e.stopPropagation()}>
                    <input
                        type="checkbox"
                        class="checkbox checkbox-sm"
                        checked={isPairSelected(row.pairing.uuid)}
                        onchange={(e) =>
                            togglePair(
                                row,
                                (e.currentTarget as HTMLInputElement).checked,
                            )} />
                </td>
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
            </tr>
        {/each}
    </tbody>
{/snippet}

{#if viewMode === "individual"}
    {#if listLoading && displayedRows.length === 0}
        <div
            class="mt-4 overflow-hidden rounded-3xl border border-base-300 bg-base-100 p-12 text-center text-sm text-base-content/50 shadow-sm">
            読み込み中…
        </div>
    {:else if total === 0}
        <div
            class="mt-4 overflow-hidden rounded-3xl border border-base-300 bg-base-100 p-12 text-center text-sm text-base-content/50 shadow-sm">
            該当するペアがありません
        </div>
    {:else}
        <div
            class="mt-4 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
            <table class="table">
                {@render tableHead(displayedRows)}
                {@render tableBody(displayedRows)}
            </table>
        </div>
        <div
            class="mt-4 flex flex-wrap items-center justify-center gap-2 text-sm">
            <button
                type="button"
                class="btn btn-sm btn-ghost"
                disabled={currentPage <= 1}
                onclick={() => goToPage(currentPage - 1)}>前へ</button>
            {#each pagerPages.entries() as [idx, p] (typeof p === "number" ? p : `e-${p}-${idx}`)}
                {#if p === "ellipsis"}
                    <span class="px-1 w-8 text-base-content/50">…</span>
                {:else}
                    <button
                        type="button"
                        class="btn btn-sm {p === currentPage
                            ? 'btn-primary'
                            : 'btn-ghost'}"
                        onclick={() => goToPage(p)}>{p}</button>
                {/if}
            {/each}
            <button
                type="button"
                class="btn btn-sm btn-ghost"
                disabled={currentPage >= lastPage}
                onclick={() => goToPage(currentPage + 1)}>次へ</button>
            <span class="ml-2 text-base-content/60">全 {total} 件</span>
        </div>
    {/if}
{:else if departmentLoading}
    <div
        class="mt-4 overflow-hidden rounded-3xl border border-base-300 bg-base-100 p-12 text-center text-sm text-base-content/50 shadow-sm">
        読み込み中…
    </div>
{:else if departmentRows.length === 0}
    <div
        class="mt-4 overflow-hidden rounded-3xl border border-base-300 bg-base-100 p-12 text-center text-sm text-base-content/50 shadow-sm">
        該当するペアがありません
    </div>
{:else}
    <div class="mt-4 space-y-6">
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
                    {@render tableHead(group.rows)}
                    {@render tableBody(group.rows)}
                </table>
            </section>
        {/each}
    </div>
{/if}

<PairSharerModal
    open={shareOpen}
    targets={shareTargets}
    users={directoryUsers}
    {departments}
    {teams}
    onClose={closeShare} />

<dialog class="modal" class:modal-open={bulkDeleteOpen}>
    <div class="modal-box">
        <h3 class="font-bold text-lg">
            選択した {selectedPairings.length} 件のペアを削除しますか？
        </h3>
        <div class="modal-action">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (bulkDeleteOpen = false)}>キャンセル</button>
            <button
                type="button"
                class="btn btn-error"
                onclick={confirmBulkDelete}>削除</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (bulkDeleteOpen = false)}
            >close</button>
    </form>
</dialog>

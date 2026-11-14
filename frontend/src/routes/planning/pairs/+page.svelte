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
    type Team = { uuid: string; name: string };

    const PAGE_SIZE = 20;

    let directoryUsers = $state<PlanningDirectoryUser[]>([]);

    let filterMentorUuids = $state<string[]>([]);
    let filterMenteeUuids = $state<string[]>([]);
    let viewMode = $state<ViewMode>("individual");
    let selected = $state<PairingListItem | null>(null);
    let deleteOpen = $state(false);
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
            if (viewMode === "individual") {
                await reloadFromPageOne();
            } else {
                await reloadDepartmentView();
            }
        } catch {
            showToast("削除に失敗しました", "error");
        }
    }

    function setViewMode(mode: ViewMode) {
        viewMode = mode;
    }
</script>

<svelte:head><title>ペア一覧 - Growlink</title></svelte:head>

<svelte:window
    onkeydown={(e: KeyboardEvent) =>
        e.key === "Escape" && (deleteOpen = false)} />

<div class="mt-4 flex flex-wrap items-center gap-2 justify-end">
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
</div>

{#if isSystemAdmin}
    <div class="mt-4 grid gap-3 sm:grid-cols-4">
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
                pickRootId="pair-list-tea,"
                items={teams}
                bind:selectedUuids={selectTeamUuids}
                labelOf={(item: Team) => item.name}
                placeholder="すべて" />
        </label>
    </div>
{/if}

{#snippet tableHead()}
    <thead>
        <tr class="text-base-content/70">
            <th class="font-semibold">部署</th>
            <th class="font-semibold">メンター</th>
            <th class="font-semibold">部署</th>
            <th class="font-semibold">メンティ</th>
            <th class="font-semibold">チーム</th>
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
                    <button
                        type="button"
                        class="btn btn-xs btn-error btn-outline"
                        onclick={() => askDelete(row.pairing)}>削除</button>
                </td>
            </tr>
        {/each}
    </tbody>
{/snippet}

{#if viewMode === "individual"}
    {#if listLoading && displayedRows.length === 0}
        <div
            class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 p-12 text-center text-sm text-base-content/50 shadow-sm">
            読み込み中…
        </div>
    {:else if total === 0}
        <div
            class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 p-12 text-center text-sm text-base-content/50 shadow-sm">
            該当するペアがありません
        </div>
    {:else}
        <div
            class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
            <table class="table">
                {@render tableHead()}
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
            {#each pagerPages as p (typeof p === "number" ? p : `e-${p}`)}
                {#if p === "ellipsis"}
                    <span class="px-1 text-base-content/50">…</span>
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
        class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 p-12 text-center text-sm text-base-content/50 shadow-sm">
        読み込み中…
    </div>
{:else if departmentRows.length === 0}
    <div
        class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 p-12 text-center text-sm text-base-content/50 shadow-sm">
        該当するペアがありません
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
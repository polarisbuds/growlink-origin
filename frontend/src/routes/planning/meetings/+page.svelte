<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import PlanningTabs from "$lib/components/PlanningTabs.svelte";
    import MultiCheckSelect from "$lib/components/MultiCheckSelect.svelte";
    import ChevronRight from "lucide-svelte/icons/chevron-right";
    import { formatDate } from "$lib/meeting-helpers";
    import { statusBadge, statusLabel } from "$lib/status";
    import {
        buildUserDepartmentMaps,
        departmentNameForUser,
        type PlanningDirectoryUser,
    } from "$lib/planning-directory";
    import {
        buildMeetingsAllQuery,
        lastPageForTotal,
        visiblePageNumbers,
        type MeetingsPageResponse,
        type MeetingListQuery,
    } from "$lib/planning-meetings-pager";
    import { sortIndicator, type SortDir } from "$lib/list-sort";
    import { goto } from "$app/navigation";
    import Icon from "$lib/components/Icon.svelte";

    type Meeting = {
        uuid: string;
        status: number;
        implementedAt?: string;
        template?: { title?: string };
        pairing?: {
            mentor?: { uuid?: string; name?: string };
            mentee?: { uuid?: string; name?: string };
        };
    };

    const PAGE_SIZE = 20;

    let directoryUsers = $state<PlanningDirectoryUser[]>([]);
    let departments = $state<{ uuid: string; name: string }[]>([]);

    let filterMentorUuids = $state<string[]>([]);
    let filterMenteeUuids = $state<string[]>([]);
    let filterDeptUuids = $state<string[]>([]);
    let filterStatus = $state("");
    let filterDateFrom = $state("");
    let filterDateTo = $state("");
    let dateDir = $state<SortDir>("desc");

    let currentPage = $state(1);
    let total = $state(0);
    let pageCache = $state<Map<number, Meeting[]>>(new Map());
    let listLoading = $state(false);
    let loadGeneration = 0;

    let deptMaps = $derived(buildUserDepartmentMaps(directoryUsers));

    let displayedMeetings = $derived(pageCache.get(currentPage) ?? []);
    let lastPage = $derived(lastPageForTotal(total, PAGE_SIZE));
    let pagerPages = $derived(visiblePageNumbers(currentPage, lastPage));

    function listQuery(page: number): MeetingListQuery {
        return {
            page,
            dir: dateDir,
            mentorUuids: filterMentorUuids,
            menteeUuids: filterMenteeUuids,
            departmentUuids: filterDeptUuids,
            status: filterStatus,
            from: filterDateFrom,
            to: filterDateTo,
        };
    }

    function trimCache(around: number) {
        const next = new Map<number, Meeting[]>();
        for (const p of [around - 1, around, around + 1]) {
            const rows = pageCache.get(p);
            if (p >= 1 && rows) next.set(p, rows);
        }
        pageCache = next;
    }

    async function fetchPage(
        page: number,
        gen: number,
    ): Promise<MeetingsPageResponse | null> {
        const qs = buildMeetingsAllQuery(listQuery(page));
        try {
            const res = await api<MeetingsPageResponse>(
                "GET",
                `/api/v1/meetings/all?${qs}`,
            );
            if (gen !== loadGeneration) return null;
            return res;
        } catch {
            if (gen !== loadGeneration) return null;
            return null;
        }
    }

    async function storePage(page: number, gen: number) {
        const res = await fetchPage(page, gen);
        if (!res || gen !== loadGeneration) return;
        total = res.total;
        const next = new Map(pageCache);
        next.set(page, res.items as Meeting[]);
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

    function toggleDateDir() {
        dateDir = dateDir === "asc" ? "desc" : "asc";
        void reloadFromPageOne();
    }

    let filterSignature = $derived(
        JSON.stringify({
            filterMentorUuids,
            filterMenteeUuids,
            filterDeptUuids,
            filterStatus,
            filterDateFrom,
            filterDateTo,
        }),
    );
    let filterSignatureMounted = "";

    $effect(() => {
        const sig = filterSignature;
        if (filterSignatureMounted === "") return;
        if (sig === filterSignatureMounted) return;
        filterSignatureMounted = sig;
        void reloadFromPageOne();
    });

    function avatarURL(user: PlanningDirectoryUser | undefined) {
        if (!user?.imageUri) return "";
        return `/uploads/${user.imageUri}`;
    }

    function userOptionLabel(u: PlanningDirectoryUser) {
        return `${u.name} (${u.email})`;
    }

    function departmentOptionLabel(d: { uuid: string; name: string }) {
        return d.name;
    }

    onMount(async () => {
        const dir = await api<{
            users?: PlanningDirectoryUser[];
            departments?: { uuid: string; name: string }[];
        }>("GET", "/api/v1/planning/directory");
        directoryUsers = dir.users || [];
        departments = dir.departments || [];
        filterSignatureMounted = filterSignature;
        await reloadFromPageOne();
    });
</script>

<svelte:head><title>全ての面談 - Growlink</title></svelte:head>

<div class="flex flex-wrap items-center justify-between gap-2">
    <div>
        <h1 class="font-display text-2xl font-semibold">面談一覧</h1>
        <p class="mt-1 text-sm text-base-content/60">組織内の面談を一覧表示</p>
    </div>
    <div class="flex justify-end items-end">
        <a href="/planning/meetings/new" class="btn btn-sm btn-primary"
            >1on1を追加</a>
    </div>
</div>

<PlanningTabs active="meetings" />

<div class="mt-4 grid gap-3 sm:grid-cols-3 lg:grid-cols-6">
    <label class="form-control">
        <span class="label-text text-xs">メンター</span>
        <MultiCheckSelect
            pickRootId="meeting-list-mentor"
            items={directoryUsers}
            bind:selectedUuids={filterMentorUuids}
            labelOf={userOptionLabel}
            placeholder="すべて" />
    </label>
    <label class="form-control">
        <span class="label-text text-xs">メンティ</span>
        <MultiCheckSelect
            pickRootId="meeting-list-mentee"
            items={directoryUsers}
            bind:selectedUuids={filterMenteeUuids}
            labelOf={userOptionLabel}
            placeholder="すべて" />
    </label>
    <label class="form-control">
        <span class="label-text text-xs">部署</span>
        <MultiCheckSelect
            pickRootId="meeting-list-dept"
            items={departments}
            bind:selectedUuids={filterDeptUuids}
            labelOf={departmentOptionLabel}
            placeholder="すべて" />
    </label>
    <label class="form-control">
        <span class="label-text text-xs">状態</span>
        <select
            class="select select-bordered select-sm"
            bind:value={filterStatus}>
            <option value="">すべて</option>
            <option value="1">未着手</option>
            <option value="2">下書き</option>
            <option value="3">完了</option>
        </select>
    </label>
    <div class="grid grid-cols-2 gap-2 sm:col-span-2">
        <label class="form-control">
            <span class="label-text text-xs">実施日（から）</span>
            <input
                type="date"
                class="input input-bordered input-sm"
                bind:value={filterDateFrom} />
        </label>
        <label class="form-control">
            <span class="label-text text-xs">実施日（まで）</span>
            <input
                type="date"
                class="input input-bordered input-sm"
                bind:value={filterDateTo} />
        </label>
    </div>
</div>

<div
    class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
    {#if listLoading && displayedMeetings.length === 0}
        <div class="p-12 text-center text-sm text-base-content/50">読み込み中…</div>
    {:else if total === 0}
        <div class="p-12 text-center text-sm text-base-content/50">
            該当する面談がありません
        </div>
    {:else}
        <table class="table table-pin-rows max-h-full">
            <thead>
                <tr class="text-base-content/70">
                    <th class="font-semibold">部署</th>
                    <th class="font-semibold">メンター</th>
                    <th class="font-semibold">部署</th>
                    <th class="font-semibold">メンティ</th>
                    <th class="font-semibold">テンプレート</th>
                    <th class="font-semibold">状態</th>
                    <th>
                        <button
                            type="button"
                            class="font-semibold hover:text-primary"
                            onclick={toggleDateDir}
                            >実施日{sortIndicator(
                                "date",
                                "date",
                                dateDir,
                            )}</button>
                    </th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                {#each displayedMeetings as m (m.uuid)}
                    {@const mentorUser = directoryUsers.find(
                        (u) => u.uuid === m.pairing?.mentor?.uuid,
                    )}
                    {@const menteeUser = directoryUsers.find(
                        (u) => u.uuid === m.pairing?.mentee?.uuid,
                    )}
                    <tr
                        class="hover:bg-base-200/40 cursor-pointer"
                        onclick={() => goto("/meeting/" + m.uuid)}>
                        <td class="text-sm"
                            >{departmentNameForUser(
                                m.pairing?.mentor?.uuid,
                                deptMaps,
                            )}
                        </td>
                        <td>
                            <span class="flex gap-2 items-center">
                                <Icon
                                    url={avatarURL(mentorUser)}
                                    text={m.pairing?.mentor?.name}
                                    width={28}
                                    height={28} />
                                {m.pairing?.mentor?.name}
                            </span>
                        </td>
                        <td class="text-sm"
                            >{departmentNameForUser(
                                m.pairing?.mentee?.uuid,
                                deptMaps,
                            )}
                        </td>
                        <td>
                            <span class="flex gap-2 items-center">
                                <Icon
                                    url={avatarURL(menteeUser)}
                                    text={m.pairing?.mentee?.name}
                                    width={28}
                                    height={28} />
                                {m.pairing?.mentee?.name}
                            </span>
                        </td>
                        <td>{m.template?.title || "—"}</td>
                        <td>
                            <span class="badge badge-sm {statusBadge(m.status)}"
                                >{statusLabel(m.status)}</span>
                        </td>
                        <td>{formatDate(m.implementedAt)}</td>
                        <td>
                            <ChevronRight></ChevronRight>
                        </td></tr>
                {/each}
            </tbody>
        </table>
    {/if}
</div>

{#if total > 0}
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
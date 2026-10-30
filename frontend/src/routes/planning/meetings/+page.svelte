<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import PlanningTabs from "$lib/components/PlanningTabs.svelte";
    import SearchableSelect from "$lib/components/SearchableSelect.svelte";
    import { formatDate } from "$lib/meeting-helpers";
    import { statusBadge, statusLabel } from "$lib/status";
    import {
        buildUserDepartmentMaps,
        departmentIdForUser,
        departmentNameForUser,
        type PlanningDirectoryUser,
    } from "$lib/planning-directory";
    import {
        compareJa,
        sortIndicator,
        toggleSortKey,
        type SortDir,
    } from "$lib/list-sort";

    type Meeting = {
        uuid: string;
        status: number;
        implementedAt?: string;
        template?: { title?: string };
        pairing?: {
            mentor?: { id?: number; name?: string };
            mentee?: { id?: number; name?: string };
        };
    };

    type SortKey =
        | "mentor"
        | "mentorDept"
        | "mentee"
        | "menteeDept"
        | "template"
        | "status"
        | "date";

    let meetings = $state<Meeting[]>([]);
    let directoryUsers = $state<PlanningDirectoryUser[]>([]);
    let departments = $state<{ id: number; name: string }[]>([]);

    let filterMentorId = $state(0);
    let filterMenteeId = $state(0);
    let filterDeptId = $state(0);
    let filterStatus = $state("");
    let filterDateFrom = $state("");
    let filterDateTo = $state("");
    let sortKey = $state<SortKey>("date");
    let sortDir = $state<SortDir>("desc");

    let deptMaps = $derived(buildUserDepartmentMaps(directoryUsers));

    type SelectItem = { id: number };

    function userLabel(id: number) {
        const u = directoryUsers.find((x) => x.id === id);
        return u ? u.name : `#${id}`;
    }

    function userLabelOf(item: SelectItem) {
        const u = directoryUsers.find((x) => x.id === item.id);
        return u ? `${u.name} ${u.email}` : `#${item.id}`;
    }

    onMount(async () => {
        const [meetingList, dir] = await Promise.all([
            api<Meeting[]>("GET", "/api/v1/meetings/all"),
            api<{
                users?: PlanningDirectoryUser[];
                departments?: { id: number; name: string }[];
            }>("GET", "/api/v1/planning/directory"),
        ]);
        meetings = meetingList;
        directoryUsers = dir.users || [];
        departments = dir.departments || [];
    });

    let displayedMeetings = $derived.by(() => {
        let list = [...meetings];
        if (filterMentorId) {
            list = list.filter((m) => m.pairing?.mentor?.id === filterMentorId);
        }
        if (filterMenteeId) {
            list = list.filter((m) => m.pairing?.mentee?.id === filterMenteeId);
        }
        if (filterDeptId) {
            list = list.filter(
                (m) =>
                    departmentIdForUser(m.pairing?.mentor?.id, deptMaps) ===
                        filterDeptId ||
                    departmentIdForUser(m.pairing?.mentee?.id, deptMaps) ===
                        filterDeptId,
            );
        }
        if (filterStatus !== "") {
            const st = Number(filterStatus);
            list = list.filter((m) => m.status === st);
        }
        if (filterDateFrom) {
            const from = new Date(filterDateFrom);
            list = list.filter((m) => {
                const d = new Date(m.implementedAt || "");
                return !Number.isNaN(d.getTime()) && d >= from;
            });
        }
        if (filterDateTo) {
            const to = new Date(filterDateTo);
            to.setHours(23, 59, 59, 999);
            list = list.filter((m) => {
                const d = new Date(m.implementedAt || "");
                return !Number.isNaN(d.getTime()) && d <= to;
            });
        }
        const dir = sortDir === "asc" ? 1 : -1;
        list.sort((a, b) => {
            switch (sortKey) {
                case "mentor":
                    return (
                        dir *
                        compareJa(
                            a.pairing?.mentor?.name || "",
                            b.pairing?.mentor?.name || "",
                        )
                    );
                case "mentorDept":
                    return (
                        dir *
                        compareJa(
                            departmentNameForUser(
                                a.pairing?.mentor?.id,
                                deptMaps,
                            ),
                            departmentNameForUser(
                                b.pairing?.mentor?.id,
                                deptMaps,
                            ),
                        )
                    );
                case "mentee":
                    return (
                        dir *
                        compareJa(
                            a.pairing?.mentee?.name || "",
                            b.pairing?.mentee?.name || "",
                        )
                    );
                case "menteeDept":
                    return (
                        dir *
                        compareJa(
                            departmentNameForUser(
                                a.pairing?.mentee?.id,
                                deptMaps,
                            ),
                            departmentNameForUser(
                                b.pairing?.mentee?.id,
                                deptMaps,
                            ),
                        )
                    );
                case "template":
                    return (
                        dir *
                        compareJa(
                            a.template?.title || "",
                            b.template?.title || "",
                        )
                    );
                case "status":
                    return dir * ((a.status ?? 0) - (b.status ?? 0));
                case "date": {
                    const da = new Date(a.implementedAt || "").getTime() || 0;
                    const db = new Date(b.implementedAt || "").getTime() || 0;
                    return dir * (da - db);
                }
                default:
                    return 0;
            }
        });
        return list;
    });

    function toggleSort(key: SortKey) {
        const next = toggleSortKey(sortKey, sortDir, key);
        sortKey = next.sortKey;
        sortDir = next.sortDir;
    }
</script>

<svelte:head><title>全ての面談 - Growlink</title></svelte:head>

<div class="flex flex-wrap items-start justify-between gap-3">
    <div>
        <h1 class="font-display text-2xl font-semibold">面談一覧</h1>
        <p class="mt-1 text-sm text-base-content/60">組織内の面談を一覧表示</p>
    </div>
    <a href="/planning/meetings/new" class="btn btn-sm btn-primary"
        >1on1を追加</a>
</div>

<PlanningTabs active="meetings" />

<div class="mt-4 grid gap-3 sm:grid-cols-3 lg:grid-cols-6">
    <label class="form-control">
        <span class="label-text text-xs">メンター</span>
        <SearchableSelect
            pickRootId="meeting-list-mentor"
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
            pickRootId="meeting-list-mentee"
            items={directoryUsers}
            bind:value={filterMenteeId}
            labelOf={userLabelOf}
            selectedLabel={userLabel}
            placeholder="すべて"
            allowClear={true}
            inputClass="input input-bordered input-sm w-full placeholder-current" />
    </label>
    <label class="form-control">
        <span class="label-text text-xs">部署</span>
        <select
            class="select select-bordered select-sm"
            bind:value={filterDeptId}>
            <option value={0}>すべて</option>
            {#each departments as d (d.id)}
                <option value={d.id}>{d.name}</option>
            {/each}
        </select>
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
    {#if displayedMeetings.length === 0}
        <div class="p-12 text-center text-sm text-base-content/50">
            該当する面談がありません
        </div>
    {:else}
        <table class="table table-pin-rows max-h-full">
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
                    <th>
                        <button
                            type="button"
                            class="font-semibold hover:text-primary"
                            onclick={() => toggleSort("template")}
                            >テンプレート{sortIndicator(
                                "template",
                                sortKey,
                                sortDir,
                            )}</button>
                    </th>
                    <th>
                        <button
                            type="button"
                            class="font-semibold hover:text-primary"
                            onclick={() => toggleSort("status")}
                            >状態{sortIndicator(
                                "status",
                                sortKey,
                                sortDir,
                            )}</button>
                    </th>
                    <th>
                        <button
                            type="button"
                            class="font-semibold hover:text-primary"
                            onclick={() => toggleSort("date")}
                            >実施日{sortIndicator(
                                "date",
                                sortKey,
                                sortDir,
                            )}</button>
                    </th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                {#each displayedMeetings as m (m.uuid)}
                    <tr class="hover:bg-base-200/40">
                        <td>{m.pairing?.mentor?.name}</td>
                        <td class="text-sm"
                            >{departmentNameForUser(
                                m.pairing?.mentor?.id,
                                deptMaps,
                            )}</td>
                        <td>{m.pairing?.mentee?.name}</td>
                        <td class="text-sm"
                            >{departmentNameForUser(
                                m.pairing?.mentee?.id,
                                deptMaps,
                            )}</td>
                        <td>{m.template?.title || "—"}</td>
                        <td>
                            <span class="badge badge-sm {statusBadge(m.status)}"
                                >{statusLabel(m.status)}</span>
                        </td>
                        <td>{formatDate(m.implementedAt)}</td>
                        <td
                            ><a
                                href="/meeting/{m.uuid}"
                                class="btn btn-xs btn-ghost gap-1">閲覧</a
                            ></td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {/if}
</div>

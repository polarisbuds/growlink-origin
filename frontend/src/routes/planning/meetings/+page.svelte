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
        departmentUuidForUser,
        departmentNameForUser,
        type PlanningDirectoryUser,
    } from "$lib/planning-directory";
    import {
        compareJa,
        sortIndicator,
        toggleSortKey,
        type SortDir,
    } from "$lib/list-sort";
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
    let departments = $state<{ uuid: string; name: string }[]>([]);

    let filterMentorUuids = $state<string[]>([]);
    let filterMenteeUuids = $state<string[]>([]);
    let filterDeptUuids = $state<string[]>([]);
    let filterStatus = $state("");
    let filterDateFrom = $state("");
    let filterDateTo = $state("");
    let sortKey = $state<SortKey>("date");
    let sortDir = $state<SortDir>("desc");

    let deptMaps = $derived(buildUserDepartmentMaps(directoryUsers));

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
        const [meetingList, dir] = await Promise.all([
            api<Meeting[]>("GET", "/api/v1/meetings/all"),
            api<{
                users?: PlanningDirectoryUser[];
                departments?: { uuid: string; name: string }[];
            }>("GET", "/api/v1/planning/directory"),
        ]);
        meetings = meetingList;
        directoryUsers = dir.users || [];
        departments = dir.departments || [];
    });

    let displayedMeetings = $derived.by(() => {
        let list = [...meetings];
        if (filterMentorUuids.length) {
            const set = new Set(filterMentorUuids);
            list = list.filter((m) => {
                const uuid = m.pairing?.mentor?.uuid;
                return uuid != null && set.has(uuid);
            });
        }
        if (filterMenteeUuids.length) {
            const set = new Set(filterMenteeUuids);
            list = list.filter((m) => {
                const uuid = m.pairing?.mentee?.uuid;
                return uuid != null && set.has(uuid);
            });
        }
        if (filterDeptUuids.length) {
            const set = new Set(filterDeptUuids);
            list = list.filter((m) => {
                const mentorDept = departmentUuidForUser(
                    m.pairing?.mentor?.uuid,
                    deptMaps,
                );
                const menteeDept = departmentUuidForUser(
                    m.pairing?.mentee?.uuid,
                    deptMaps,
                );
                return (
                    (mentorDept && set.has(mentorDept)) ||
                    (menteeDept && set.has(menteeDept))
                );
            });
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
                                a.pairing?.mentor?.uuid,
                                deptMaps,
                            ),
                            departmentNameForUser(
                                b.pairing?.mentor?.uuid,
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
                                a.pairing?.mentee?.uuid,
                                deptMaps,
                            ),
                            departmentNameForUser(
                                b.pairing?.mentee?.uuid,
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

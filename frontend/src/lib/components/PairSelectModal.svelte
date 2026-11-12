<script lang="ts">
    import MultiCheckSelect from "$lib/components/MultiCheckSelect.svelte";
    import Icon from "$lib/components/Icon.svelte";
    import { type PairingDisplayRow } from "$lib/pairing-helpers";
    import {
        buildUserDepartmentMaps,
        buildUserTeamMaps,
        departmentNameForUser,
        teamNameForUser,
        type PlanningDirectoryUser,
    } from "$lib/planning-directory";

    const PAIR_PAGE_SIZE = 10;

    type Dept = { uuid: string; name: string };
    type Team = { uuid: string; name: string; departmentUuid: string };

    type PickerUser = PlanningDirectoryUser & {
        imageUri?: string;
    };

    let {
        open = false,
        title = "ペアを選択",
        rows,
        users,
        departments = [],
        teams = [],
        selectedUuids = [],
        pickRootId = "pair-select-modal-name",
        onApply,
        onClose,
    }: {
        open?: boolean;
        title?: string;
        rows: PairingDisplayRow[];
        users: PickerUser[];
        departments?: Dept[];
        teams?: Team[];
        selectedUuids?: string[];
        pickRootId?: string;
        onApply: (uuids: string[]) => void;
        onClose: () => void;
    } = $props();

    let dialogEl = $state<HTMLDialogElement | null>(null);

    let draftUuids = $state<Set<string>>(new Set());
    let filterNameUuids = $state<string[]>([]);
    let filterDeptUuids = $state<string[]>([]);
    let filterTeamUuids = $state<string[]>([]);
    let page = $state(1);

    let deptMaps = $derived(buildUserDepartmentMaps(users));
    let teamMaps = $derived(buildUserTeamMaps(users));

    let nameFilterUsers = $derived.by(() => {
        const seen = new Set<string>();
        const list: PickerUser[] = [];
        for (const row of rows) {
            for (const uuid of [row.mentorUserUuid, row.menteeUserUuid]) {
                if (!uuid || seen.has(uuid)) continue;
                const u = users.find((x) => x.uuid === uuid);
                if (u) {
                    seen.add(uuid);
                    list.push(u);
                }
            }
        }
        list.sort((a, b) => a.name.localeCompare(b.name, "ja"));
        return list;
    });

    let teamOptions = $derived.by(() => {
        if (filterDeptUuids.length === 0) return teams;
        const deptSet = new Set(filterDeptUuids);
        return teams.filter((t) => deptSet.has(t.departmentUuid));
    });

    $effect(() => {
        const ok = new Set(teamOptions.map((t) => t.uuid));
        const pruned = filterTeamUuids.filter((id) => ok.has(id));
        if (pruned.length !== filterTeamUuids.length) {
            filterTeamUuids = pruned;
        }
    });

    $effect(() => {
        if (open) {
            draftUuids = new Set(selectedUuids);
            filterNameUuids = [];
            filterDeptUuids = [];
            filterTeamUuids = [];
            page = 1;
        }
    });

    $effect(() => {
        filterNameUuids;
        filterDeptUuids;
        filterTeamUuids;
        page = 1;
    });

    function rowMatchesDept(row: PairingDisplayRow, departmentName: string) {
        if (!departmentName) return true;
        const mentorDept = departmentNameForUser(row.mentorUserUuid, deptMaps);
        const menteeDept = departmentNameForUser(row.menteeUserUuid, deptMaps);
        return (
            mentorDept.includes(departmentName) ||
            menteeDept.includes(departmentName)
        );
    }

    function rowMatchesTeam(row: PairingDisplayRow, teamName: string) {
        if (!teamName) return true;
        const u = users.find((x) => x.uuid === row.menteeUserUuid);
        return (u?.teamNames || "").includes(teamName);
    }

    let filteredRows = $derived.by(() => {
        let list = [...rows];

        if (filterNameUuids.length > 0) {
            const nameSet = new Set(filterNameUuids);
            list = list.filter(
                (r) =>
                    (r.mentorUserUuid != null && nameSet.has(r.mentorUserUuid)) ||
                    (r.menteeUserUuid != null && nameSet.has(r.menteeUserUuid)),
            );
        }

        if (filterDeptUuids.length > 0) {
            const departmentNames = departments
                .filter((d) => filterDeptUuids.includes(d.uuid))
                .map((d) => d.name)
                .filter(Boolean);
            list = list.filter((r) =>
                departmentNames.some((name) => rowMatchesDept(r, name)),
            );
        }

        if (filterTeamUuids.length > 0) {
            const teamNames = teams
                .filter((t) => filterTeamUuids.includes(t.uuid))
                .map((t) => t.name)
                .filter(Boolean);
            list = list.filter((r) =>
                teamNames.some((name) => rowMatchesTeam(r, name)),
            );
        }

        list.sort((a, b) => {
            const da = departmentNameForUser(a.menteeUserUuid, deptMaps) || "";
            const db = departmentNameForUser(b.menteeUserUuid, deptMaps) || "";
            const c = da.localeCompare(db, "ja");
            if (c !== 0) return c;
            return a.menteeName.localeCompare(b.menteeName, "ja");
        });
        return list;
    });

    let pageCount = $derived(
        Math.max(1, Math.ceil(filteredRows.length / PAIR_PAGE_SIZE)),
    );

    $effect(() => {
        if (page > pageCount) page = pageCount;
    });

    let pageRows = $derived(
        filteredRows.slice(
            (page - 1) * PAIR_PAGE_SIZE,
            page * PAIR_PAGE_SIZE,
        ),
    );

    let pagePairingUuids = $derived([
        ...new Set(pageRows.map((r) => r.pairing.uuid)),
    ]);

    let allPageSelected = $derived(
        pagePairingUuids.length > 0 &&
            pagePairingUuids.every((id) => draftUuids.has(id)),
    );

    function userLabelOf(u: PickerUser) {
        return `${u.name} (${u.email})`;
    }

    function avatarURL(uuid: string) {
        const u = users.find((x) => x.uuid === uuid);
        if (!u?.imageUri) return "";
        return `/uploads/${u.imageUri}`;
    }

    function pairingUuid(row: PairingDisplayRow) {
        return row.pairing.uuid;
    }

    function toggleDraft(uuid: string) {
        const next = new Set(draftUuids);
        if (next.has(uuid)) next.delete(uuid);
        else next.add(uuid);
        draftUuids = next;
    }

    function setDraft(uuid: string, checked: boolean) {
        const next = new Set(draftUuids);
        if (checked) next.add(uuid);
        else next.delete(uuid);
        draftUuids = next;
    }

    function toggleAllOnPage(checked: boolean) {
        const next = new Set(draftUuids);
        for (const id of pagePairingUuids) {
            if (checked) next.add(id);
            else next.delete(id);
        }
        draftUuids = next;
    }

    function apply() {
        onApply([...draftUuids]);
    }

    function discardClose() {
        onClose();
    }

    function onWindowKeydown(e: KeyboardEvent) {
        if (!open || e.key !== "Escape") return;
        if (!dialogEl) return;
        if (dialogEl.querySelector("[data-pick-root] ul")) return;
        discardClose();
    }
</script>

<svelte:window onkeydown={onWindowKeydown} />

<dialog bind:this={dialogEl} class="modal" class:modal-open={open}>
    <div class="modal-box flex max-h-[85vh] max-w-6xl flex-col p-0">
        <div
            class="flex shrink-0 items-center justify-between gap-2 border-b border-base-200 px-6 py-4">
            <h3 class="text-lg font-bold">{title}</h3>
            <button
                type="button"
                class="btn btn-sm btn-circle btn-ghost"
                aria-label="閉じる"
                onclick={discardClose}>×</button>
        </div>

        <div class="min-h-0 flex-1 overflow-y-auto px-6 py-4 space-y-4">
            <div class="grid gap-3 sm:grid-cols-3">
                <label class="form-control">
                    <span class="label-text text-xs">氏名</span>
                    <MultiCheckSelect
                        {pickRootId}
                        items={nameFilterUsers}
                        bind:selectedUuids={filterNameUuids}
                        labelOf={userLabelOf}
                        placeholder="すべて" />
                </label>
                <label class="form-control">
                    <span class="label-text text-xs">部署</span>
                    <MultiCheckSelect
                        pickRootId={`${pickRootId}-dept`}
                        items={departments}
                        bind:selectedUuids={filterDeptUuids}
                        labelOf={(d: Dept) => d.name}
                        placeholder="すべて" />
                </label>
                <label class="form-control">
                    <span class="label-text text-xs">チーム</span>
                    <MultiCheckSelect
                        pickRootId={`${pickRootId}-team`}
                        items={teamOptions}
                        bind:selectedUuids={filterTeamUuids}
                        labelOf={(t: Team) => t.name}
                        placeholder="すべて" />
                </label>
            </div>

            <div class="overflow-x-auto rounded-2xl border border-base-300">
                {#if filteredRows.length === 0}
                    <div class="p-8 text-center text-sm text-base-content/50">
                        該当するペアがありません
                    </div>
                {:else}
                    <table class="table table-sm">
                        <thead>
                            <tr>
                                <th class="w-10">
                                    <input
                                        type="checkbox"
                                        class="checkbox checkbox-sm"
                                        checked={allPageSelected}
                                        onchange={(e) =>
                                            toggleAllOnPage(
                                                (
                                                    e.currentTarget as HTMLInputElement
                                                ).checked,
                                            )} />
                                </th>
                                <th>メンター</th>
                                <th>部署</th>
                                <th>メンティ</th>
                                <th>部署</th>
                                <th>チーム</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each pageRows as row (row.rowKey)}
                                {@const pid = pairingUuid(row)}
                                <tr
                                    class="cursor-pointer hover:bg-base-200/40"
                                    onclick={() => toggleDraft(pid)}>
                                    <td onclick={(e) => e.stopPropagation()}>
                                        <input
                                            type="checkbox"
                                            class="checkbox checkbox-sm"
                                            checked={draftUuids.has(pid)}
                                            onchange={(e) =>
                                                setDraft(
                                                    pid,
                                                    (
                                                        e.currentTarget as HTMLInputElement
                                                    ).checked,
                                                )} />
                                    </td>
                                    <td>
                                        <span class="flex items-center gap-2">
                                            <Icon
                                                url={avatarURL(
                                                    row.mentorUserUuid,
                                                )}
                                                text={row.mentorName}
                                                width={28}
                                                height={28} />
                                            {row.mentorName}
                                        </span>
                                    </td>
                                    <td
                                        >{departmentNameForUser(
                                            row.mentorUserUuid,
                                            deptMaps,
                                        ) || "—"}</td>
                                    <td>
                                        <span class="flex items-center gap-2">
                                            <Icon
                                                url={avatarURL(
                                                    row.menteeUserUuid,
                                                )}
                                                text={row.menteeName}
                                                width={28}
                                                height={28} />
                                            {row.menteeName}
                                        </span>
                                    </td>
                                    <td
                                        >{departmentNameForUser(
                                            row.menteeUserUuid,
                                            deptMaps,
                                        ) || "—"}</td>
                                    <td
                                        >{teamNameForUser(
                                            row.menteeUserUuid,
                                            teamMaps,
                                        ) || "—"}</td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                {/if}
            </div>

            <div
                class="flex flex-wrap items-center justify-between gap-2 text-sm">
                <span class="text-base-content/60">
                    選択中: {draftUuids.size} ペア（表示 {filteredRows.length}
                    件）
                </span>
                {#if filteredRows.length > 0}
                    <div class="join">
                        <button
                            type="button"
                            class="btn btn-sm join-item"
                            disabled={page <= 1}
                            onclick={() => (page = Math.max(1, page - 1))}
                            >前へ</button>
                        <span
                            class="btn btn-sm join-item btn-disabled pointer-events-none"
                            >{page} / {pageCount}</span>
                        <button
                            type="button"
                            class="btn btn-sm join-item"
                            disabled={page >= pageCount}
                            onclick={() =>
                                (page = Math.min(pageCount, page + 1))}
                            >次へ</button>
                    </div>
                {/if}
            </div>
        </div>

        <div
            class="modal-action shrink-0 border-t border-base-200 p-4 flex justify-end">
            <button type="button" class="btn btn-primary" onclick={apply}
                >適用</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={discardClose}>close</button>
    </form>
</dialog>

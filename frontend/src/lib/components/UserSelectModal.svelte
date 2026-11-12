<script lang="ts">
    import MultiCheckSelect from "$lib/components/MultiCheckSelect.svelte";
    import Icon from "$lib/components/Icon.svelte";

    const USER_PAGE_SIZE = 10;

    export type UserPickerUser = {
        uuid: string;
        name: string;
        email: string;
        departmentNames?: string;
        teamNames?: string;
        imageUri?: string;
    };

    type Dept = { uuid: string; name: string };
    type Team = { uuid: string; name: string; departmentUuid: string };

    let {
        open = false,
        title = "ユーザを選択",
        users,
        departments = [],
        teams = [],
        selectedUuids = [],
        pickRootId = "user-select-modal-name",
        multiple = true,
        onApply,
        onClose,
    }: {
        open?: boolean;
        title?: string;
        users: UserPickerUser[];
        departments?: Dept[];
        teams?: Team[];
        selectedUuids?: string[];
        pickRootId?: string;
        multiple?: boolean;
        onApply: (uuids: string[]) => void;
        onClose: () => void;
    } = $props();

    let dialogEl = $state<HTMLDialogElement | null>(null);

    let draftUuids = $state<Set<string>>(new Set());
    let filterNameUuids = $state<string[]>([]);
    let filterDeptUuids = $state<string[]>([]);
    let filterTeamUuids = $state<string[]>([]);
    let page = $state(1);

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

    let filteredRows = $derived.by(() => {
        let list = [...users];

        if (filterNameUuids.length > 0) {
            const nameSet = new Set(filterNameUuids);
            list = list.filter((u) => nameSet.has(u.uuid));
        }

        if (filterDeptUuids.length > 0) {
            const nameSet = new Set(
                departments
                    .filter((d) => filterDeptUuids.includes(d.uuid))
                    .map((d) => d.name),
            );
            list = list.filter((u) =>
                [...nameSet].some((n) => (u.departmentNames || "").includes(n)),
            );
        }
        if (filterTeamUuids.length > 0) {
            const nameSet = new Set(
                teams
                    .filter((t) => filterTeamUuids.includes(t.uuid))
                    .map((t) => t.name),
            );
            list = list.filter((u) =>
                [...nameSet].some((n) => (u.teamNames || "").includes(n)),
            );
        }

        list.sort((a, b) => {
            const da = a.departmentNames || "";
            const db = b.departmentNames || "";
            const c = da.localeCompare(db, "ja");
            if (c !== 0) return c;
            return a.name.localeCompare(b.name, "ja");
        });
        return list;
    });

    let pageCount = $derived(
        Math.max(1, Math.ceil(filteredRows.length / USER_PAGE_SIZE)),
    );

    $effect(() => {
        if (page > pageCount) page = pageCount;
    });

    let pageRows = $derived(
        filteredRows.slice(
            (page - 1) * USER_PAGE_SIZE,
            page * USER_PAGE_SIZE,
        ),
    );

    let allPageSelected = $derived(
        pageRows.length > 0 && pageRows.every((u) => draftUuids.has(u.uuid)),
    );

    function userLabelOf(u: UserPickerUser) {
        return `${u.name} (${u.email})`;
    }

    function avatarURL(u: UserPickerUser | undefined) {
        if (!u?.imageUri) return "";
        return `/uploads/${u.imageUri}`;
    }

    function toggleDraft(uuid: string) {
        const next = new Set(draftUuids);
        if (!multiple) {
            if (next.has(uuid)) next.clear();
            else {
                next.clear();
                next.add(uuid);
            }
        } else if (next.has(uuid)) next.delete(uuid);
        else next.add(uuid);
        draftUuids = next;
    }

    function setDraft(uuid: string, checked: boolean) {
        if (!multiple) {
            if (checked) draftUuids = new Set([uuid]);
            else {
                const next = new Set(draftUuids);
                next.delete(uuid);
                draftUuids = next;
            }
            return;
        }
        const next = new Set(draftUuids);
        if (checked) next.add(uuid);
        else next.delete(uuid);
        draftUuids = next;
    }

    function toggleAllOnPage(checked: boolean) {
        const next = new Set(draftUuids);
        for (const u of pageRows) {
            if (checked) next.add(u.uuid);
            else next.delete(u.uuid);
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
    <div class="modal-box flex max-h-[85vh] max-w-4xl flex-col p-0">
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
                    <span class="label-text text-xs">名前</span>
                    <MultiCheckSelect
                        pickRootId={pickRootId}
                        items={users}
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
                        該当する人員がありません
                    </div>
                {:else}
                    <table class="table table-sm">
                        <thead>
                            <tr>
                                <th class="w-10">
                                    {#if multiple}
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
                                    {/if}
                                </th>
                                <th>部署</th>
                                <th>氏名</th>
                                <th>チーム</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each pageRows as u (u.uuid)}
                                <tr
                                    class="cursor-pointer hover:bg-base-200/40"
                                    onclick={() => toggleDraft(u.uuid)}>
                                    <td onclick={(e) => e.stopPropagation()}>
                                        <input
                                            type="checkbox"
                                            class="checkbox checkbox-sm"
                                            checked={draftUuids.has(u.uuid)}
                                            onchange={(e) =>
                                                setDraft(
                                                    u.uuid,
                                                    (
                                                        e.currentTarget as HTMLInputElement
                                                    ).checked,
                                                )} />
                                    </td>
                                    <td>{u.departmentNames || "所属無し"}</td>
                                    <td>
                                        <span class="flex items-center gap-2">
                                            <Icon
                                                url={avatarURL(u)}
                                                text={u.name}
                                                width={28}
                                                height={28} />
                                            {u.name}
                                        </span>
                                    </td>
                                    <td>{u.teamNames || ""}</td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                {/if}
            </div>

            <div
                class="flex flex-wrap items-center justify-between gap-2 text-sm">
                <span class="text-base-content/60">
                    選択中: {draftUuids.size} 人（表示 {filteredRows.length} 人）
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

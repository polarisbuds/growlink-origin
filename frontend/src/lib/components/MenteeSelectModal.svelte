<script lang="ts">
    import SearchableSelect from "$lib/components/SearchableSelect.svelte";
    import Icon from "$lib/components/Icon.svelte";

    const MENTEE_PAGE_SIZE = 10;

    export type MenteePickerUser = {
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
        users,
        departments = [],
        teams = [],
        selectedUuids = [],
        pickRootId = "mentee-select-modal-name",
        onApply,
        onClose,
    }: {
        open?: boolean;
        users: MenteePickerUser[];
        departments?: Dept[];
        teams?: Team[];
        selectedUuids?: string[];
        pickRootId?: string;
        onApply: (uuids: string[]) => void;
        onClose: () => void;
    } = $props();

    let draftUuids = $state<Set<string>>(new Set());
    let filterNameUuid = $state("");
    let filterDeptUuid = $state("");
    let filterTeamUuid = $state("");
    let page = $state(1);

    let teamOptions = $derived.by(() => {
        if (!filterDeptUuid) return teams;
        return teams.filter((t) => t.departmentUuid === filterDeptUuid);
    });

    $effect(() => {
        if (!filterTeamUuid) return;
        const ok = teamOptions.some((t) => t.uuid === filterTeamUuid);
        if (!ok) filterTeamUuid = "";
    });

    $effect(() => {
        if (open) {
            draftUuids = new Set(selectedUuids);
            filterNameUuid = "";
            filterDeptUuid = "";
            filterTeamUuid = "";
            page = 1;
        }
    });

    $effect(() => {
        filterNameUuid;
        filterDeptUuid;
        filterTeamUuid;
        page = 1;
    });

    let filteredRows = $derived.by(() => {
        let list = [...users];

        if (filterNameUuid) {
            list = list.filter((u) => u.uuid === filterNameUuid);
        }

        if (filterDeptUuid) {
            const department = departments.find((d) => d.uuid === filterDeptUuid);
            const departmentName = department?.name || "";
            list = list.filter((u) => u.departmentNames === departmentName);
        }
        if (filterTeamUuid) {
            const team = teams.find((t) => t.uuid === filterTeamUuid);
            const teamName = team?.name || "";
            list = list.filter((u) => (u.teamNames || "").includes(teamName));
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
        Math.max(1, Math.ceil(filteredRows.length / MENTEE_PAGE_SIZE)),
    );

    $effect(() => {
        if (page > pageCount) page = pageCount;
    });

    let pageRows = $derived(
        filteredRows.slice(
            (page - 1) * MENTEE_PAGE_SIZE,
            page * MENTEE_PAGE_SIZE,
        ),
    );

    let allPageSelected = $derived(
        pageRows.length > 0 &&
            pageRows.every((u) => draftUuids.has(u.uuid)),
    );

    function userLabelOf(u: MenteePickerUser) {
        return `${u.name} (${u.email})`;
    }
    function userLabel(uuid: string) {
        const u = users.find((x) => x.uuid === uuid);
        return u ? u.name : uuid;
    }

    function avatarURL(u: MenteePickerUser | undefined) {
        if (!u?.imageUri) return "";
        return `/uploads/${u.imageUri}`;
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
        const modal = document.querySelector("[data-mentee-select-modal]");
        if (!modal) return;
        if (modal.querySelector("[data-pick-root] ul")) return;
        discardClose();
    }
</script>

<svelte:window onkeydown={onWindowKeydown} />

<dialog class="modal" class:modal-open={open} data-mentee-select-modal>
    <div class="modal-box flex max-h-[85vh] max-w-4xl flex-col p-0">
        <div
            class="flex shrink-0 items-center justify-between gap-2 border-b border-base-200 px-6 py-4">
            <h3 class="text-lg font-bold">メンティを選択</h3>
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
                    <SearchableSelect
                        {pickRootId}
                        items={users}
                        bind:value={filterNameUuid}
                        labelOf={userLabelOf}
                        selectedLabel={userLabel}
                        placeholder="名前で検索"
                        allowClear={true}
                        inputSize="input-sm" />
                </label>
                <label class="form-control">
                    <span class="label-text text-xs">部署</span>
                    <select
                        class="select select-bordered select-sm w-full"
                        bind:value={filterDeptUuid}>
                        <option value="">すべて</option>
                        {#each departments as d (d.uuid)}
                            <option value={d.uuid}>{d.name}</option>
                        {/each}
                    </select>
                </label>
                <label class="form-control">
                    <span class="label-text text-xs">チーム</span>
                    <select
                        class="select select-bordered select-sm w-full"
                        bind:value={filterTeamUuid}>
                        <option value="">すべて</option>
                        {#each teamOptions as t (t.uuid)}
                            <option value={t.uuid}>{t.name}</option>
                        {/each}
                    </select>
                </label>
            </div>

            <div class="overflow-x-auto rounded-2xl border border-base-300">
                {#if filteredRows.length === 0}
                    <div
                        class="p-8 text-center text-sm text-base-content/50">
                        該当する人員がありません
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
            <button
                type="button"
                class="btn btn-primary"
                onclick={apply}>適用</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={discardClose}>close</button>
    </form>
</dialog>

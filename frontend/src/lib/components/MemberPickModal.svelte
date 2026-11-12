<script lang="ts">
    import Icon from "$lib/components/Icon.svelte";

    const PAGE_SIZE = 10;

    export type MemberPickerUser = {
        uuid: string;
        name: string;
        email: string;
        imageUri?: string;
    };

    let {
        open = false,
        title = "メンバーを選択",
        users,
        selectedUuids = [],
        excludeUuids = [],
        multiple = true,
        onApply,
        onClose,
    }: {
        open?: boolean;
        title?: string;
        users: MemberPickerUser[];
        selectedUuids?: string[];
        excludeUuids?: string[];
        multiple?: boolean;
        onApply: (uuids: string[]) => void;
        onClose: () => void;
    } = $props();

    let dialogEl = $state<HTMLDialogElement | null>(null);
    let draftUuids = $state<Set<string>>(new Set());
    let searchQuery = $state("");
    let page = $state(1);

    $effect(() => {
        if (open) {
            draftUuids = new Set(selectedUuids);
            searchQuery = "";
            page = 1;
        }
    });

    $effect(() => {
        searchQuery;
        page = 1;
    });

    let filteredRows = $derived.by(() => {
        const q = searchQuery.trim().toLowerCase();
        let list = [...users]
            .filter(
                (u) => !(excludeUuids || []).some((uuid) => uuid === u.uuid),
            )
            .filter(Boolean);
        if (q) {
            list = list.filter(
                (u) =>
                    (u.name || "").toLowerCase().includes(q) ||
                    (u.email || "").toLowerCase().includes(q),
            );
        }

        list.sort((a, b) => a.name.localeCompare(b.name, "ja"));
        return list;
    });

    let pageCount = $derived(
        Math.max(1, Math.ceil(filteredRows.length / PAGE_SIZE)),
    );

    $effect(() => {
        if (page > pageCount) page = pageCount;
    });

    let pageRows = $derived(
        filteredRows.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE),
    );

    let allPageSelected = $derived(
        pageRows.length > 0 && pageRows.every((u) => draftUuids.has(u.uuid)),
    );

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

    function avatarURL(u: MemberPickerUser | undefined) {
        if (!u?.imageUri) return "";
        return `/uploads/${u.imageUri}`;
    }

    function onWindowKeydown(e: KeyboardEvent) {
        if (!open || e.key !== "Escape") return;
        if (!dialogEl) return;
        discardClose();
    }
</script>

<svelte:window onkeydown={onWindowKeydown} />

<dialog bind:this={dialogEl} class="modal z-[1100]" class:modal-open={open}>
    <div class="modal-box flex max-h-[85vh] max-w-3xl flex-col p-0">
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
            <label class="form-control">
                <span class="label-text text-xs">氏名・メールアドレス</span>
                <input
                    type="search"
                    class="input input-bordered input-sm w-full"
                    placeholder="検索"
                    bind:value={searchQuery}
                    autocomplete="off" />
            </label>

            <div class="overflow-x-auto rounded-2xl border border-base-300">
                {#if filteredRows.length === 0}
                    <div class="p-8 text-center text-sm text-base-content/50">
                        該当するユーザがありません
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
                                <th>氏名</th>
                                <th>メールアドレス</th>
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
                                    <td class="text-base-content/80"
                                        >{u.email}</td>
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
            class="modal-action shrink-0 border-t border-base-200 p-4 flex justify-end gap-2">
            <button type="button" class="btn btn-ghost" onclick={discardClose}
                >キャンセル</button>
            <button type="button" class="btn btn-primary" onclick={apply}
                >適用</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={discardClose}>close</button>
    </form>
</dialog>

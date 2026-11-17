<script lang="ts">
    const PAGE_SIZE = 8;

    export type TemplatePickerTopic = {
        uuid?: string;
        displayNo?: number;
        question?: string;
        description?: string;
        isRequire?: boolean;
    };

    export type TemplatePickerTemplate = {
        uuid: string;
        title?: string;
        description?: string;
        isDefault?: boolean;
        topics?: TemplatePickerTopic[];
    };

    let {
        open = false,
        title = "テンプレートを選択",
        templates,
        selectedUuid = "",
        class: boxClass = "min-h-[24rem]",
        onApply,
        onClose,
    }: {
        open?: boolean;
        title?: string;
        templates: TemplatePickerTemplate[];
        selectedUuid?: string;
        class?: string;
        onApply: (uuid: string) => void;
        onClose: () => void;
    } = $props();

    let dialogEl = $state<HTMLDialogElement | null>(null);
    let draftUuid = $state("");
    let query = $state("");
    let page = $state(1);

    $effect(() => {
        if (open) {
            draftUuid = selectedUuid;
            query = "";
            page = 1;
        }
    });

    let filtered = $derived.by(() => {
        const q = query.trim().toLowerCase();
        const list = q
            ? templates.filter((t) => (t.title || "").toLowerCase().includes(q))
            : [...templates];
        list.sort((a, b) => (a.title || "").localeCompare(b.title || "", "ja"));
        return list;
    });

    let pageCount = $derived(
        Math.max(1, Math.ceil(filtered.length / PAGE_SIZE)),
    );
    $effect(() => {
        query;
        page = 1;
    });
    $effect(() => {
        if (page > pageCount) page = pageCount;
    });

    let pageRows = $derived(
        filtered.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE),
    );

    let preview = $derived(templates.find((t) => t.uuid === draftUuid));
    let previewTopics = $derived(
        [...(preview?.topics || [])].sort(
            (a, b) => (a.displayNo || 0) - (b.displayNo || 0),
        ),
    );

    function onWindowKeydown(e: KeyboardEvent) {
        if (!open || e.key !== "Escape" || !dialogEl) return;
        onClose();
    }
</script>

<svelte:window onkeydown={onWindowKeydown} />

<dialog bind:this={dialogEl} class="modal" class:modal-open={open}>
    <div
        class="modal-box flex max-h-[85vh] max-w-5xl flex-col !overflow-hidden p-0 {boxClass}">
        <div
            class="flex shrink-0 items-center justify-between gap-2 border-b border-base-200 px-6 py-4">
            <h3 class="text-lg font-bold">{title}</h3>
            <button
                type="button"
                class="btn btn-sm btn-circle btn-ghost"
                aria-label="閉じる"
                onclick={onClose}>×</button>
        </div>

        <div
            class="grid min-h-0 flex-1 gap-4 overflow-hidden px-6 py-4 lg:grid-cols-2">
            <div class="flex min-h-0 flex-col gap-3">
                <label class="form-control">
                    <span class="label-text text-xs">テンプレート名</span>
                    <input
                        class="input input-bordered input-sm"
                        placeholder="名前で検索"
                        bind:value={query} />
                </label>
                <div
                    class="min-h-0 flex-1 overflow-y-auto rounded-2xl border border-base-300">
                    {#if filtered.length === 0}
                        <div
                            class="p-8 text-center text-sm text-base-content/50">
                            該当するテンプレートがありません
                        </div>
                    {:else}
                        <ul>
                            {#each pageRows as t (t.uuid)}
                                <li>
                                    <button
                                        type="button"
                                        class="flex w-full items-start gap-2 border-b border-base-200 px-3 py-2 text-left hover:bg-base-200/40"
                                        class:bg-primary={draftUuid === t.uuid}
                                        class:text-primary-content={draftUuid ===
                                            t.uuid}
                                        onclick={() => (draftUuid = t.uuid)}>
                                        <input
                                            type="radio"
                                            class="radio radio-sm mt-0.5"
                                            checked={draftUuid === t.uuid}
                                            tabindex="-1" />
                                        <span class="min-w-0">
                                            <span
                                                class="block text-sm font-medium">
                                                {t.title || "無題"}
                                            </span>
                                            <span
                                                class="block text-xs opacity-70">
                                                トピック {(t.topics || [])
                                                    .length} 件
                                            </span>
                                        </span>
                                    </button>
                                </li>
                            {/each}
                        </ul>
                    {/if}
                </div>
                {#if filtered.length > 0}
                    <div class="join self-end">
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

            <div
                class="min-h-0 overflow-y-auto rounded-2xl border border-base-300 p-4">
                {#if !preview}
                    <p class="text-sm text-base-content/50">
                        テンプレートを選ぶと質問と説明が表示されます
                    </p>
                {:else}
                    <h4 class="font-medium">{preview.title || "無題"}</h4>
                    {#if preview.description}
                        <p
                            class="mt-1 whitespace-pre-wrap text-sm text-base-content/60">
                            {preview.description}
                        </p>
                    {/if}
                    <div class="mt-3 space-y-2">
                        {#each previewTopics as topic (topic.uuid || topic.displayNo)}
                            <div
                                class="rounded-xl border border-base-300 bg-base-200/30 p-3">
                                <p class="text-sm">
                                    <span
                                        class="badge badge-primary badge-outline badge-sm mr-1"
                                        >Q{topic.displayNo}</span>
                                    {topic.question}
                                    {#if topic.isRequire}
                                        <span
                                            class="badge badge-error badge-outline badge-xs ml-1"
                                            >必須</span>
                                    {/if}
                                </p>
                                {#if topic.description}
                                    <p
                                        class="mt-1 whitespace-pre-wrap text-xs text-base-content/60">
                                        {topic.description}
                                    </p>
                                {/if}
                            </div>
                        {:else}
                            <p class="text-sm text-base-content/50">
                                質問はありません
                            </p>
                        {/each}
                    </div>
                {/if}
            </div>
        </div>

        <div class="modal-action shrink-0 border-t border-base-200 p-4">
            <button
                type="button"
                class="btn btn-primary"
                disabled={!draftUuid}
                onclick={() => onApply(draftUuid)}>適用</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={onClose}>close</button>
    </form>
</dialog>

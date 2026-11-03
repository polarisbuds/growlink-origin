<script lang="ts" generics="T extends { id: number }">
    import { onMount } from "svelte";
    import { filterItems } from "$lib/searchable-select";

    let {
        items,
        selectedIds = $bindable([] as number[]),
        labelOf,
        placeholder = "すべて",
        pickRootId = "multi-check-select",
        inputClass = "input input-bordered input-sm w-full text-left font-normal placeholder-current",
    }: {
        items: T[];
        selectedIds?: number[];
        labelOf: (item: T) => string;
        placeholder?: string;
        pickRootId?: string;
        inputClass?: string;
    } = $props();

    let open = $state(false);
    let query = $state("");

    let filtered = $derived(filterItems(items, query, labelOf));
    let selectedSet = $derived(new Set(selectedIds));

    let summaryLabel = $derived.by(() => {
        if (selectedIds.length === 0) return placeholder;
        if (selectedIds.length === 1) {
            const item = items.find((i) => i.id === selectedIds[0]);
            return item ? labelOf(item) : placeholder;
        }
        return `${selectedIds.length}件選択`;
    });

    let allFilteredSelected = $derived(
        filtered.length > 0 &&
            filtered.every((item) => selectedSet.has(item.id)),
    );
    let someFilteredSelected = $derived(
        filtered.some((item) => selectedSet.has(item.id)) &&
            !allFilteredSelected,
    );

    let selectAllEl = $state<HTMLInputElement | null>(null);
    $effect(() => {
        if (selectAllEl) {
            selectAllEl.indeterminate = someFilteredSelected;
        }
    });

    function toggleItem(id: number, checked: boolean) {
        if (checked) {
            if (selectedSet.has(id)) return;
            selectedIds = [...selectedIds, id];
        } else {
            selectedIds = selectedIds.filter((x) => x !== id);
        }
    }

    function toggleAllFiltered(checked: boolean) {
        const filteredIds = filtered.map((i) => i.id);
        if (checked) {
            const next = new Set(selectedIds);
            for (const id of filteredIds) next.add(id);
            selectedIds = [...next];
        } else {
            const remove = new Set(filteredIds);
            selectedIds = selectedIds.filter((id) => !remove.has(id));
        }
    }

    onMount(() => {
        const onMouseDown = (ev: MouseEvent) => {
            const t = ev.target;
            if (!(t instanceof Element)) return;
            if (!t.closest(`[data-pick-root="${pickRootId}"]`)) open = false;
        };
        const onKeyDown = (ev: KeyboardEvent) => {
            if (ev.key === "Escape") open = false;
        };
        window.addEventListener("mousedown", onMouseDown);
        window.addEventListener("keydown", onKeyDown);
        return () => {
            window.removeEventListener("mousedown", onMouseDown);
            window.removeEventListener("keydown", onKeyDown);
        };
    });
</script>

<div class="relative" data-pick-root={pickRootId}>
    <button
        type="button"
        class="{inputClass} flex min-h-9 items-center justify-between gap-2"
        aria-haspopup="listbox"
        aria-expanded={open}
        onclick={() => (open = !open)}>
        <span class="truncate {selectedIds.length === 0 ? 'text-base-content/50' : ''}"
            >{summaryLabel}</span>
        <span class="text-base-content/40" aria-hidden="true">▾</span>
    </button>
    {#if open}
        <div
            class="absolute left-0 top-full z-20 mt-1 w-full overflow-hidden rounded-2xl border border-base-300 bg-base-100 shadow-lg">
            <div class="border-b border-base-300 p-2">
                <input
                    type="search"
                    class="input input-bordered input-sm w-full placeholder-current"
                    placeholder="検索"
                    autocomplete="off"
                    bind:value={query} />
            </div>
            {#if filtered.length > 0}
                <label
                    class="flex cursor-pointer items-center gap-2 border-b border-base-200 px-3 py-2 text-sm hover:bg-base-200/50">
                    <input
                        type="checkbox"
                        class="checkbox checkbox-sm"
                        bind:this={selectAllEl}
                        checked={allFilteredSelected}
                        onchange={(e) =>
                            toggleAllFiltered(
                                (e.currentTarget as HTMLInputElement).checked,
                            )} />
                    <span class="font-medium">全選択</span>
                </label>
            {/if}
            <ul
                class="max-h-48 overflow-y-auto py-1"
                role="listbox"
                aria-multiselectable="true">
                {#each filtered as item (item.id)}
                    <li>
                        <label
                            class="flex cursor-pointer items-center gap-2 px-3 py-2 text-sm hover:bg-base-200/50">
                            <input
                                type="checkbox"
                                class="checkbox checkbox-sm"
                                checked={selectedSet.has(item.id)}
                                onchange={(e) =>
                                    toggleItem(
                                        item.id,
                                        (e.currentTarget as HTMLInputElement)
                                            .checked,
                                    )} />
                            <span class="truncate">{labelOf(item)}</span>
                        </label>
                    </li>
                {/each}
                {#if filtered.length === 0}
                    <li class="px-4 py-3 text-sm text-base-content/50">
                        該当なし
                    </li>
                {/if}
            </ul>
        </div>
    {/if}
</div>

<script lang="ts" generics="T extends { uuid: string }">
    import { onMount } from "svelte";
    import { filterItems } from "$lib/searchable-select";

    let {
        items,
        selectedUuids = $bindable([] as string[]),
        labelOf,
        placeholder = "すべて",
        pickRootId = "multi-check-select",
        inputClass = "input input-bordered input-sm w-full text-left font-normal placeholder-current",
    }: {
        items: T[];
        selectedUuids?: string[];
        labelOf: (item: T) => string;
        placeholder?: string;
        pickRootId?: string;
        inputClass?: string;
    } = $props();

    let open = $state(false);
    let query = $state("");

    let filtered = $derived(filterItems(items, query, labelOf));
    let selectedSet = $derived(new Set(selectedUuids));

    let summaryLabel = $derived.by(() => {
        if (selectedUuids.length === 0) return placeholder;
        if (selectedUuids.length === 1) {
            const item = items.find((i) => i.uuid === selectedUuids[0]);
            return item ? labelOf(item) : placeholder;
        }
        return `${selectedUuids.length}件選択`;
    });

    let allFilteredSelected = $derived(
        filtered.length > 0 &&
            filtered.every((item) => selectedSet.has(item.uuid)),
    );
    let someFilteredSelected = $derived(
        filtered.some((item) => selectedSet.has(item.uuid)) &&
            !allFilteredSelected,
    );

    let selectAllEl = $state<HTMLInputElement | null>(null);
    $effect(() => {
        if (selectAllEl) {
            selectAllEl.indeterminate = someFilteredSelected;
        }
    });

    function toggleItem(uuid: string, checked: boolean) {
        if (checked) {
            if (selectedSet.has(uuid)) return;
            selectedUuids = [...selectedUuids, uuid];
        } else {
            selectedUuids = selectedUuids.filter((x) => x !== uuid);
        }
    }

    function toggleAllFiltered(checked: boolean) {
        const filteredUuids = filtered.map((i) => i.uuid);
        if (checked) {
            const next = new Set(selectedUuids);
            for (const id of filteredUuids) next.add(id);
            selectedUuids = [...next];
        } else {
            const remove = new Set(filteredUuids);
            selectedUuids = selectedUuids.filter((id) => !remove.has(id));
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
        <span class="truncate {selectedUuids.length === 0 ? 'text-base-content/50' : ''}"
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
                {#each filtered as item (item.uuid)}
                    <li>
                        <label
                            class="flex cursor-pointer items-center gap-2 px-3 py-2 text-sm hover:bg-base-200/50">
                            <input
                                type="checkbox"
                                class="checkbox checkbox-sm"
                                checked={selectedSet.has(item.uuid)}
                                onchange={(e) =>
                                    toggleItem(
                                        item.uuid,
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

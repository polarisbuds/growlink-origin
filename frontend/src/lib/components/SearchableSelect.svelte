<script lang="ts">
    import { onMount } from "svelte";
    import { filterItems } from "$lib/searchable-select";

    type SelectItem = { id: number };

    let {
        items,
        value = $bindable(0),
        labelOf,
        selectedLabel: selectedLabelFn,
        placeholder = "名前で検索",
        allowClear = false,
        inputClass = "input input-bordered w-full placeholder-current",
        pickRootId = "searchable-select",
        inputSize = "",
    }: {
        items: SelectItem[];
        value?: number;
        labelOf: (item: SelectItem) => string;
        selectedLabel?: (id: number) => string;
        placeholder?: string;
        allowClear?: boolean;
        inputClass?: string;
        pickRootId?: string;
        inputSize?: string;
    } = $props();

    let open = $state(false);
    let query = $state("");

    let filtered = $derived(filterItems(items, query, labelOf));

    function displaySelected(): string {
        if (!value) return "";
        if (selectedLabelFn) return selectedLabelFn(value);
        const item = items.find((i) => i.id === value);
        return item ? labelOf(item) : "";
    }

    function clear() {
        value = 0;
        query = "";
        open = false;
    }

    onMount(() => {
        const onMouseDown = (ev: MouseEvent) => {
            const t = ev.target;
            if (!(t instanceof Element)) return;
            if (!t.closest(`[data-pick-root="${pickRootId}"]`)) open = false;
        };
        window.addEventListener("mousedown", onMouseDown);
        return () => window.removeEventListener("mousedown", onMouseDown);
    });
</script>

<div class="relative" data-pick-root={pickRootId}>
    <input
        type="text"
        class="{inputClass} {inputSize} {allowClear && value ? 'pr-10' : ''}"
        autocomplete="off"
        placeholder={value ? "" : placeholder}
        value={open ? query : displaySelected()}
        onfocus={() => {
            open = true;
            if (value) query = "";
        }}
        oninput={(e) => {
            query = e.currentTarget.value;
            open = true;
        }} />
    {#if allowClear && value}
        <button
            type="button"
            class="btn btn-ghost btn-xs btn-circle absolute right-1 top-1/2 -translate-y-1/2"
            aria-label="クリア"
            onclick={clear}>×</button>
    {/if}
    {#if open}
        <ul
            class="absolute left-0 top-full z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg">
            {#each filtered as item (item.id)}
                <li>
                    <button
                        type="button"
                        class="btn btn-ghost btn-sm w-full justify-start font-normal"
                        onclick={() => {
                            value = item.id;
                            query = "";
                            open = false;
                        }}>{labelOf(item)}</button>
                </li>
            {/each}
            {#if filtered.length === 0}
                <li class="px-4 py-3 text-sm text-base-content/50">該当なし</li>
            {/if}
        </ul>
    {/if}
</div>

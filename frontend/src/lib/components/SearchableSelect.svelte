<script lang="ts" generics="T extends { uuid: string }">
    import { onMount } from "svelte";
    import { filterItems } from "$lib/searchable-select";
    import PickOptionList from "$lib/components/PickOptionList.svelte";

    let {
        items,
        value = $bindable(""),
        labelOf,
        selectedLabel: selectedLabelFn,
        placeholder = "名前で検索",
        allowClear = false,
        clearOnSelect = false,
        onSelect,
        inputClass = "input input-bordered w-full placeholder-current",
        pickRootId = "searchable-select",
        inputSize = "",
    }: {
        items: T[];
        value?: string;
        labelOf: (item: T) => string;
        selectedLabel?: (uuid: string) => string;
        placeholder?: string;
        allowClear?: boolean;
        clearOnSelect?: boolean;
        onSelect?: (item: T) => void;
        inputClass?: string;
        pickRootId?: string;
        inputSize?: string;
    } = $props();

    let open = $state(false);
    let query = $state("");

    let filtered = $derived(filterItems(items, query, labelOf));

    function displaySelected(): string {
        if (clearOnSelect || !value) return "";
        if (selectedLabelFn) return selectedLabelFn(value);
        const item = items.find((i) => i.uuid === value);
        return item ? labelOf(item) : "";
    }

    function clear() {
        value = "";
        query = "";
        open = false;
    }

    function pick(item: T) {
        onSelect?.(item);
        if (clearOnSelect) {
            query = "";
            open = false;
        } else {
            value = item.uuid;
            query = "";
            open = false;
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
    <input
        type="text"
        class="{inputClass} {inputSize} {allowClear && value && !clearOnSelect ? 'pr-10' : ''}"
        autocomplete="off"
        placeholder={displaySelected() ? "" : placeholder}
        value={open ? query : displaySelected()}
        onfocus={() => {
            open = true;
            if (value && !clearOnSelect) query = "";
        }}
        oninput={(e) => {
            query = e.currentTarget.value;
            open = true;
        }} />
    {#if allowClear && value && !clearOnSelect}
        <button
            type="button"
            class="btn btn-ghost btn-xs btn-circle absolute right-1 top-1/2 -translate-y-1/2"
            aria-label="クリア"
            onclick={clear}>×</button>
    {/if}
    {#if open}
        <PickOptionList items={filtered} {labelOf} onSelect={pick} />
    {/if}
</div>

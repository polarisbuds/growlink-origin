<script lang="ts" generics="T extends { uuid: string }">
    import { onMount } from "svelte";
    import { filterItems } from "$lib/searchable-select";
    import { floating, portal } from "$lib/floating-panel";

    let {
        items,
        value = $bindable(""),
        labelOf,
        selectedLabel: selectedLabelFn,
        placeholder = "名前で検索",
        allowClear = false,
        clearOnSelect = false,
        onSelect,
        inputClass =
            "input input-bordered input-sm w-full text-left font-normal placeholder-current",
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
    let triggerEl = $state<HTMLButtonElement | null>(null);
    let searchEl = $state<HTMLInputElement | null>(null);

    let filtered = $derived(filterItems(items, query, labelOf));

    let summaryLabel = $derived.by(() => {
        if (clearOnSelect || !value) return placeholder;
        if (selectedLabelFn) return selectedLabelFn(value);
        const item = items.find((i) => i.uuid === value);
        return item ? labelOf(item) : placeholder;
    });

    let hasSelection = $derived(!clearOnSelect && !!value);

    function closePanel() {
        open = false;
        query = "";
    }

    function toggleOpen() {
        if (open) {
            closePanel();
        } else {
            query = "";
            open = true;
        }
    }

    function clear() {
        value = "";
        closePanel();
    }

    function pick(item: T) {
        onSelect?.(item);
        if (clearOnSelect) {
            closePanel();
        } else {
            value = item.uuid;
            closePanel();
        }
    }

    $effect(() => {
        if (open && searchEl) {
            searchEl.focus();
        }
    });

    onMount(() => {
        const onMouseDown = (ev: MouseEvent) => {
            const t = ev.target;
            if (!(t instanceof Element)) return;
            if (!t.closest(`[data-pick-root="${pickRootId}"]`)) {
                if (open) closePanel();
            }
        };
        const onKeyDown = (ev: KeyboardEvent) => {
            if (ev.key !== "Escape" || !open) return;
            closePanel();
            ev.stopPropagation();
        };
        window.addEventListener("mousedown", onMouseDown);
        window.addEventListener("keydown", onKeyDown);
        return () => {
            window.removeEventListener("mousedown", onMouseDown);
            window.removeEventListener("keydown", onKeyDown);
        };
    });
</script>

<div class="relative h-full" data-pick-root={pickRootId}>
    <button
        bind:this={triggerEl}
        type="button"
        class="{inputClass} {inputSize} flex min-h-9 items-center justify-between gap-2 h-full {allowClear &&
        hasSelection
            ? 'pr-8'
            : ''}"
        aria-haspopup="listbox"
        aria-expanded={open}
        onclick={toggleOpen}>
        <span
            class="truncate {hasSelection
                ? ''
                : 'text-base-content/50'}">{summaryLabel}</span>
        <span class="text-base-content/40 shrink-0" aria-hidden="true">▾</span>
    </button>
    {#if allowClear && hasSelection}
        <button
            type="button"
            class="btn btn-ghost btn-xs btn-circle absolute right-7 top-1/2 z-10 -translate-y-1/2"
            aria-label="クリア"
            onclick={(e) => {
                e.stopPropagation();
                clear();
            }}>×</button>
    {/if}
    {#if open}
        <div
            use:portal
            use:floating={{ getTrigger: () => triggerEl }}
            data-pick-root={pickRootId}
            data-pick-panel=""
            class="fixed z-[1300] overflow-hidden rounded-2xl border border-base-300 bg-base-100 shadow-lg">
            <div class="border-b border-base-300 p-2">
                <input
                    bind:this={searchEl}
                    type="search"
                    class="input input-bordered input-sm w-full placeholder-current"
                    placeholder="検索"
                    autocomplete="off"
                    bind:value={query} />
            </div>
            <ul
                class="max-h-48 overflow-y-auto py-1"
                role="listbox">
                {#each filtered as item (item.uuid)}
                    <li>
                        <button
                            type="button"
                            role="option"
                            aria-selected={!clearOnSelect && value === item.uuid}
                            class="flex w-full cursor-pointer items-center px-3 py-2 text-left text-sm hover:bg-base-200/50 {value ===
                                item.uuid && !clearOnSelect
                                ? 'bg-base-200/70 font-medium'
                                : ''}"
                            onclick={() => pick(item)}>
                            <span class="truncate">{labelOf(item)}</span>
                        </button>
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
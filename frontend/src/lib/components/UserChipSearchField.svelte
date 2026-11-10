<script lang="ts" generics="T extends { uuid: string }">
    import { onMount } from "svelte";
    import { filterItems } from "$lib/searchable-select";
    import PickOptionList from "$lib/components/PickOptionList.svelte";

    type Chip = { uuid: string; name: string };

    let {
        chips,
        items,
        excludeUuids = [],
        placeholder = "名前で検索して追加",
        query = $bindable(""),
        className,
        labelOf,
        pickRootId = "user-chip-search",
        onPick,
        onRemove,
    }: {
        chips: Chip[];
        items: T[];
        excludeUuids?: string[];
        placeholder?: string;
        query?: string;
        className?: string;
        labelOf?: (item: T) => string;
        pickRootId?: string;
        onPick: (item: T) => void;
        onRemove: (uuid: string) => void;
    } = $props();

    let open = $state(false);

    function itemLabel(item: T): string {
        if (labelOf) return labelOf(item);
        const u = item as { name?: string; email?: string };
        return `${u.name ?? ""} (${u.email ?? ""})`;
    }

    let excluded = $derived(
        new Set([...excludeUuids, ...chips.map((c) => c.uuid)]),
    );
    let availableItems = $derived(items.filter((i) => !excluded.has(i.uuid)));
    let filtered = $derived(filterItems(availableItems, query, itemLabel));

    function pick(item: T) {
        onPick(item);
        query = "";
        open = false;
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
    <div
        class="flex min-h-11 w-full flex-wrap items-center gap-1.5 rounded-lg border border-base-300 bg-base-100 px-3 py-2 focus-within:border-primary {className}">
        {#each chips as chip (chip.uuid)}
            <span
                class="inline-flex max-w-full shrink-0 items-center gap-0.5 rounded-full bg-base-200 pl-2.5 pr-0.5 py-0.5 text-sm text-base-content">
                <span class="max-w-[11rem] truncate" title={chip.name}
                    >{chip.name}</span>
                <button
                    type="button"
                    class="btn btn-ghost btn-xs btn-circle h-5 min-h-5 w-5 shrink-0 opacity-70 hover:opacity-100"
                    aria-label="{chip.name}を削除"
                    onclick={() => onRemove(chip.uuid)}>×</button>
            </span>
        {/each}
        <input
            type="text"
            class="min-w-[7rem] flex-1 border-0 bg-transparent py-1 text-sm outline-none focus:outline-none placeholder-current"
            {placeholder}
            autocomplete="off"
            bind:value={query}
            onfocus={() => (open = true)}
            oninput={() => (open = true)} />
    </div>
    {#if open}
        <PickOptionList items={filtered} labelOf={itemLabel} onSelect={pick} />
    {/if}
</div>

<script lang="ts" generics="T extends { uuid: string }">
    import { floating, portal } from "$lib/floating-panel";

    let {
        items,
        labelOf,
        onSelect,
        pickRootId = "",
        getTrigger,
    }: {
        items: T[];
        labelOf: (item: T) => string;
        onSelect: (item: T) => void;
        pickRootId?: string;
        getTrigger: () => HTMLElement | null;
    } = $props();
</script>

<div
    use:portal
    use:floating={{ getTrigger }}
    data-pick-root={pickRootId}
    data-pick-panel=""
    class="fixed z-[1300] overflow-hidden rounded-2xl border border-base-300 bg-base-100 shadow-lg">
    <ul
        class="max-h-48 overflow-y-auto py-1"
        role="listbox">
        {#each items as item (item.uuid)}
            <li>
                <button
                    type="button"
                    class="btn btn-ghost btn-sm w-full justify-start font-normal"
                    onclick={() => onSelect(item)}>{labelOf(item)}</button>
            </li>
        {/each}
        {#if items.length === 0}
            <li class="px-4 py-3 text-sm text-base-content/50">該当なし</li>
        {/if}
    </ul>
</div>
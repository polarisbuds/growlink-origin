<script lang="ts" generics="T extends { uuid: string }">
    let {
        items,
        labelOf,
        onSelect,
    }: {
        items: T[];
        labelOf: (item: T) => string;
        onSelect: (item: T) => void;
    } = $props();
</script>

<ul
    class="absolute left-0 top-full z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg">
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

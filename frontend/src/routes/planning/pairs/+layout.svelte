<script lang="ts">
    import PlanningTabs from "$lib/components/PlanningTabs.svelte";
    import { page } from "$app/state";

    let { children } = $props();

    let showListHeader = $derived(
        /^\/planning\/pairs(\/users)?\/?$/.test(page.url.pathname),
    );
</script>

{#if showListHeader}
    <div class="flex flex-wrap items-start justify-between gap-3">
        <div>
            <h1 class="font-display text-2xl font-semibold">ペア一覧</h1>
            <p class="mt-1 text-sm text-base-content/60">
                メンターとメンティの 1:1 ペア
            </p>
        </div>
        <a href="/planning/pairs/new" class="btn btn-primary btn-sm">+ 作成</a>
    </div>
{/if}
<PlanningTabs active="pairs" />
{@render children()}

<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import PlanningTabs from "$lib/components/PlanningTabs.svelte";
    import { ensureSession, getAvatar } from "$lib/session";
    import { page } from "$app/state";

    let { children } = $props();

    let allowed = $state(false);

    let showListHeader = $derived(
        allowed && /^\/planning\/pairs?\/?$/.test(page.url.pathname),
    );

    onMount(async () => {
        const me = getAvatar() ?? (await ensureSession());
        if (!me?.canAccessPairs) {
            await goto("/planning/meetings");
            return;
        }
        allowed = true;
    });
</script>

{#if allowed}
    {#if showListHeader}
        <div class="flex flex-wrap items-start justify-between gap-3">
            <div>
                <h1 class="font-display text-2xl font-semibold">ペア一覧</h1>
                <p class="mt-1 text-sm text-base-content/60">
                    メンターとメンティの 1:1 ペア
                </p>
            </div>
        </div>
    {/if}
    <PlanningTabs active="pairs" />
    {@render children()}
{/if}

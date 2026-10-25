<script lang="ts">
	import PlanningTabs from '$lib/components/PlanningTabs.svelte';
	import { page } from '$app/state';

	let { children } = $props();

	let showListHeader = $derived(page.url.pathname.includes('/planning/pairs/users'));
</script>

<PlanningTabs active="pairs" />
{#if showListHeader}
	<div class="mt-4 flex flex-wrap items-start justify-between gap-3">
		<div>
			<h1 class="font-display text-2xl font-semibold">ペア一覧</h1>
			<p class="mt-1 text-sm text-base-content/60">メンターとメンティの 1:1 ペア</p>
		</div>
		<a href="/planning/pairs/new" class="btn btn-primary btn-sm">+ 作成</a>
	</div>
{/if}
{@render children()}

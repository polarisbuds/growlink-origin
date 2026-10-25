<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import ProgressChart from '$lib/components/ProgressChart.svelte';
	import type { ProgressBreakdownRow } from '$lib/progress-types';

	type ProgressSummary = {
		pending: number;
		draft: number;
		done: number;
		total: number;
	};

	type BreakdownGroup = 'department' | 'team';

	let summary = $state<ProgressSummary>({ pending: 0, draft: 0, done: 0, total: 0 });
	let breakdownGroup = $state<BreakdownGroup>('department');
	let breakdownRows = $state<ProgressBreakdownRow[]>([]);
	let breakdownLoading = $state(false);

	let groupLabel = $derived(breakdownGroup === 'department' ? '部署' : 'チーム');

	async function loadBreakdown() {
		breakdownLoading = true;
		try {
			breakdownRows = await api<ProgressBreakdownRow[]>(
				'GET',
				`/api/v1/progress/breakdown?group=${breakdownGroup}`
			);
		} catch {
			breakdownRows = [];
		} finally {
			breakdownLoading = false;
		}
	}

	onMount(async () => {
		try {
			summary = await api<ProgressSummary>('GET', '/api/v1/progress/summary');
		} catch {
			/* keep zeros */
		}
		await loadBreakdown();
	});

	async function setBreakdownGroup(group: BreakdownGroup) {
		if (breakdownGroup === group) return;
		breakdownGroup = group;
		await loadBreakdown();
	}
</script>

<svelte:head><title>進捗 - Growlink</title></svelte:head>

<h1 class="font-display text-2xl font-semibold">進捗</h1>
<p class="mt-1 text-sm text-base-content/60">面談ステータス別の件数</p>
<div class="mt-8 grid gap-4 md:grid-cols-2 xl:grid-cols-4">
	<div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
		<p class="text-sm text-base-content/60">未着手</p>
		<p class="mt-2 text-3xl font-semibold text-primary">{summary.pending}</p>
	</div>
	<div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
		<p class="text-sm text-base-content/60">下書き</p>
		<p class="mt-2 text-3xl font-semibold text-warning">{summary.draft}</p>
	</div>
	<div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
		<p class="text-sm text-base-content/60">完了</p>
		<p class="mt-2 text-3xl font-semibold text-success">{summary.done}</p>
	</div>
	<div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
		<p class="text-sm text-base-content/60">合計</p>
		<p class="mt-2 text-3xl font-semibold">{summary.total}</p>
	</div>
</div>

<section class="mt-10 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div>
			<h2 class="text-lg font-semibold">ステータス内訳</h2>
			<p class="mt-1 text-sm text-base-content/60">{groupLabel}ごとの面談件数</p>
		</div>
		<div class="tabs tabs-boxed">
			<button
				type="button"
				class="tab"
				class:tab-active={breakdownGroup === 'department'}
				onclick={() => setBreakdownGroup('department')}>部署</button
			>
			<button
				type="button"
				class="tab"
				class:tab-active={breakdownGroup === 'team'}
				onclick={() => setBreakdownGroup('team')}>チーム</button
			>
		</div>
	</div>
	<div class="mt-6">
		{#if breakdownLoading}
			<p class="py-12 text-center text-sm text-base-content/50">読み込み中…</p>
		{:else}
			<ProgressChart rows={breakdownRows} {groupLabel} />
		{/if}
	</div>
</section>

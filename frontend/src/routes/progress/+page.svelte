<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

	type ProgressSummary = {
		pending: number;
		draft: number;
		done: number;
		total: number;
	};

	let summary = $state<ProgressSummary>({ pending: 0, draft: 0, done: 0, total: 0 });

	onMount(async () => {
		try {
			summary = await api<ProgressSummary>('GET', '/api/v1/progress/summary');
		} catch {
			/* keep zeros */
		}
	});
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

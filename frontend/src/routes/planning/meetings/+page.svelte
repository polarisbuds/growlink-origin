<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { api } from '$lib/api';
	import PlanningTabs from '$lib/components/PlanningTabs.svelte';
	import { formatDate } from '$lib/meeting-helpers';
	import { statusBadge, statusLabel } from '$lib/status';

	type Meeting = {
		uuid: string;
		status: number;
		implementedAt?: string;
		template?: { title?: string };
		pairing?: { mentor?: { name?: string }; mentee?: { name?: string } };
	};

	let meetings = $state<Meeting[]>([]);
	let search = $state('');
	let simple = $derived(page.url.searchParams.get('view') === 'simple');

	let filteredMeetings = $derived.by(() => {
		const q = search.trim().toLowerCase();
		if (!q) return meetings;
		return meetings.filter((m) => {
			const hay = [m.template?.title, m.pairing?.mentor?.name, m.pairing?.mentee?.name]
				.filter(Boolean)
				.join(' ')
				.toLowerCase();
			return hay.includes(q);
		});
	});

	onMount(async () => {
		meetings = await api<Meeting[]>('GET', '/api/v1/meetings/all');
	});
</script>

<svelte:head><title>全ての面談 - Growlink</title></svelte:head>

<PlanningTabs active="meetings" />
<div class="mt-4 flex flex-wrap items-start justify-between gap-3">
	<div>
		<h1 class="font-display text-2xl font-semibold">面談一覧</h1>
		<p class="mt-1 text-sm text-base-content/60">組織内の面談を一覧表示</p>
	</div>
	<div class="flex flex-wrap items-center gap-2">
		<a href="/planning/meetings/new" class="btn btn-sm btn-primary">1on1を追加</a>
		<a href={simple ? '/planning/meetings' : '/planning/meetings?view=simple'} class="btn btn-sm btn-outline"
			>{simple ? '詳細表示' : '簡易表示'}</a
		>
	</div>
</div>
<label class="input input-bordered mt-4 flex max-w-md items-center gap-2">
	<input type="search" class="grow" placeholder="検索" bind:value={search} />
</label>
<div class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
	{#if filteredMeetings.length === 0}
		<div class="p-12 text-center text-sm text-base-content/50">該当する面談がありません</div>
	{:else}
		<table class="table">
			<thead>
				<tr class="text-base-content/70">
					<th>メンター</th>
					<th>メンティ</th>
					<th>状態</th>
					<th>実施日</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each filteredMeetings as m (m.uuid)}
					<tr class="hover:bg-base-200/40">
						<td>{m.pairing?.mentor?.name}</td>
						<td>{m.pairing?.mentee?.name}</td>
						<td>
							<span class="badge badge-sm {statusBadge(m.status)}">{statusLabel(m.status)}</span>
						</td>
						<td>{formatDate(m.implementedAt)}</td>
						<td><a href="/meeting/{m.uuid}" class="btn btn-xs btn-ghost gap-1">閲覧</a></td>
					</tr>
					{#if simple}
						<tr class="bg-base-200/40">
							<td colspan="5" class="text-sm text-base-content/70">{m.template?.title}</td>
						</tr>
					{/if}
				{/each}
			</tbody>
		</table>
	{/if}
</div>

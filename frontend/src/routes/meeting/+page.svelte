<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { statusBadge, statusLabel } from '$lib/status';
	import { formatDate } from '$lib/meeting-helpers';

	let meetings = $state<Record<string, unknown>[]>([]);
	let search = $state('');

	let filtered = $derived.by(() => {
		const q = search.trim().toLowerCase();
		if (!q) return meetings;
		return meetings.filter((m) => {
			const template = m.template as { title?: string };
			const pairing = m.pairing as {
				mentor?: { name?: string };
				mentee?: { name?: string };
			};
			const hay = [template?.title, pairing?.mentor?.name, pairing?.mentee?.name, statusLabel(m.status as number)]
				.filter(Boolean)
				.join(' ')
				.toLowerCase();
			return hay.includes(q);
		});
	});

	onMount(async () => {
		meetings = await api('GET', '/api/v1/meetings/involved');
	});
</script>

<svelte:head><title>面談一覧 - Growlink</title></svelte:head>

<h1 class="font-display text-2xl font-semibold">面談一覧</h1>
<input
	type="search"
	class="input input-bordered mt-4 w-full max-w-md"
	placeholder="検索"
	bind:value={search}
/>
<div class="mt-6 overflow-x-auto rounded-2xl border border-base-300 bg-base-100">
	<table class="table">
		<thead>
			<tr>
				<th>テンプレート</th>
				<th>ペア</th>
				<th>実施日</th>
				<th>状態</th>
				<th></th>
			</tr>
		</thead>
		<tbody>
			{#each filtered as m (m.uuid)}
				<tr>
					<td>{(m.template as { title?: string })?.title || '—'}</td>
					<td>
						{(m.pairing as { mentor?: { name?: string } })?.mentor?.name || '—'} /
						{(m.pairing as { mentee?: { name?: string } })?.mentee?.name || '—'}
					</td>
					<td>{formatDate(m.implementedAt as string)}</td>
					<td>
						<span class="badge {statusBadge(m.status as number)}">{statusLabel(m.status as number)}</span>
					</td>
					<td class="text-right">
						<a href="/meeting/{m.uuid}" class="btn btn-ghost btn-xs">開く</a>
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>

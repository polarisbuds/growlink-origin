<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { showToast } from '$lib/toast';
	import {
		type PairingListItem,
		type PairingSection,
		expandPairingsForSection,
		rowSearchHaystack,
		sectionSecondColumnLabel
	} from '$lib/pairing-helpers';

	let { section }: { section: PairingSection } = $props();

	let pairings = $state<PairingListItem[]>([]);
	let search = $state('');
	let selected = $state<PairingListItem | null>(null);
	let deleteOpen = $state(false);

	let displayRows = $derived(expandPairingsForSection(pairings, section));

	let filteredRows = $derived.by(() => {
		const q = search.trim().toLowerCase();
		if (!q) return displayRows;
		return displayRows.filter((r) => rowSearchHaystack(r).includes(q));
	});

	const secondCol = $derived(sectionSecondColumnLabel(section));

	onMount(async () => {
		pairings = await api<PairingListItem[]>('GET', '/api/v1/pairings');
	});

	function askDelete(p: PairingListItem) {
		selected = p;
		deleteOpen = true;
	}
	async function confirmDelete() {
		if (!selected) return;
		try {
			await api('DELETE', `/api/v1/pairings/${selected.uuid || selected.id}`);
			deleteOpen = false;
			showToast('ペアを削除しました');
			pairings = await api<PairingListItem[]>('GET', '/api/v1/pairings');
		} catch {
			showToast('削除に失敗しました', 'error');
		}
	}
</script>

<svelte:window onkeydown={(e) => e.key === 'Escape' && (deleteOpen = false)} />

<label class="input input-bordered mt-4 flex max-w-md items-center gap-2">
	<input
		type="search"
		class="grow"
		placeholder="メンター・{secondCol}で検索"
		bind:value={search}
	/>
</label>
<div class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
	{#if filteredRows.length === 0}
		<div class="p-12 text-center text-sm text-base-content/50">該当するペアがありません</div>
	{:else}
		<table class="table">
			<thead>
				<tr class="text-base-content/70">
					<th>メンター</th>
					<th>{secondCol}</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each filteredRows as row (row.rowKey)}
					<tr class="hover:bg-base-200/40">
						<td>{row.mentorName}</td>
						<td class="text-sm">{row.counterpartyName}</td>
						<td class="space-x-1 whitespace-nowrap text-right">
							<a
								href="/planning/pairs/{row.pairing.uuid}/edit"
								class="btn btn-xs btn-outline">編集</a
							>
							<button
								type="button"
								class="btn btn-xs btn-error btn-outline"
								onclick={() => askDelete(row.pairing)}>削除</button
							>
							<a
								href="/planning/pairs/{row.pairing.uuid}/edit"
								class="btn btn-ghost btn-xs btn-square inline-flex align-middle">›</a
							>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>

<dialog class="modal" class:modal-open={deleteOpen}>
	<div class="modal-box max-h-[85vh]">
		<h3 class="text-lg font-bold">削除の確認</h3>
		<p class="py-4 text-sm text-base-content/80">このペアを削除しますか？</p>
		<div class="modal-action">
			<button type="button" class="btn btn-ghost" onclick={() => (deleteOpen = false)}>キャンセル</button>
			<button type="button" class="btn btn-error" onclick={confirmDelete}>削除確定</button>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button type="button" onclick={() => (deleteOpen = false)}>close</button>
	</form>
</dialog>

<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { showToast } from '$lib/toast';
	import PlanningTabs from '$lib/components/PlanningTabs.svelte';

	type Pairing = {
		id: number;
		uuid: string;
		mentor?: { name?: string };
		mentee?: { name?: string };
		targets?: { label?: string }[];
	};

	let pairings = $state<Pairing[]>([]);
	let search = $state('');
	let selected = $state<Pairing | null>(null);
	let deleteOpen = $state(false);

	function searchHaystack(p: Pairing) {
		const targetLabels = (p.targets || []).map((t) => t.label).join(' ');
		return [p.mentor?.name, p.mentee?.name, targetLabels].filter(Boolean).join(' ').toLowerCase();
	}
	function targetsLabel(p: Pairing) {
		const labels = (p.targets || []).map((t) => t.label).filter(Boolean);
		if (labels.length) return labels.join('、');
		return p.mentee?.name || '—';
	}

	let filteredPairings = $derived.by(() => {
		const q = search.trim().toLowerCase();
		if (!q) return pairings;
		return pairings.filter((p) => searchHaystack(p).includes(q));
	});

	onMount(async () => {
		pairings = await api<Pairing[]>('GET', '/api/v1/pairings');
	});

	function askDelete(p: Pairing) {
		selected = p;
		deleteOpen = true;
	}
	async function confirmDelete() {
		if (!selected) return;
		try {
			await api('DELETE', `/api/v1/pairings/${selected.uuid || selected.id}`);
			deleteOpen = false;
			showToast('ペアを削除しました');
			pairings = await api<Pairing[]>('GET', '/api/v1/pairings');
		} catch {
			showToast('削除に失敗しました', 'error');
		}
	}
</script>

<svelte:head><title>ペア一覧 - Growlink</title></svelte:head>

<svelte:window onkeydown={(e) => e.key === 'Escape' && (deleteOpen = false)} />

<PlanningTabs active="pairs" />
<div class="mt-4 flex flex-wrap items-start justify-between gap-3">
	<div>
		<h1 class="font-display text-2xl font-semibold">ペア一覧</h1>
		<p class="mt-1 text-sm text-base-content/60">メンターとメンティのペア</p>
	</div>
	<a href="/planning/pairs/new" class="btn btn-primary btn-sm">+ 作成</a>
</div>
<label class="input input-bordered mt-4 flex max-w-md items-center gap-2">
	<input type="search" class="grow" placeholder="メンター・対象で検索" bind:value={search} />
</label>
<div class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
	{#if filteredPairings.length === 0}
		<div class="p-12 text-center text-sm text-base-content/50">該当するペアがありません</div>
	{:else}
		<table class="table">
			<thead>
				<tr class="text-base-content/70">
					<th>メンター</th>
					<th>対象</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each filteredPairings as p (p.id)}
					<tr class="hover:bg-base-200/40">
						<td>{p.mentor?.name}</td>
						<td class="text-sm">{targetsLabel(p)}</td>
						<td class="space-x-1 whitespace-nowrap text-right">
							<a href="/planning/pairs/{p.uuid}/edit" class="btn btn-xs btn-outline">編集</a>
							<button type="button" class="btn btn-xs btn-error btn-outline" onclick={() => askDelete(p)}>削除</button>
							<a href="/planning/pairs/{p.uuid}/edit" class="btn btn-ghost btn-xs btn-square inline-flex align-middle">›</a>
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

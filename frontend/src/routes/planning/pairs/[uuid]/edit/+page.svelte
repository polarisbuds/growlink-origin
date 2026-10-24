<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import PairingFormView from '$lib/views/PairingFormView.svelte';
	import { api } from '$lib/api';

	let pairing = $state<Record<string, unknown> | null>(null);
	let uuid = $derived(page.params.uuid);

	onMount(async () => {
		pairing = await api('GET', `/api/v1/pairings/${uuid}`);
	});
</script>

{#if pairing}
	<PairingFormView mode="edit" initialPairing={pairing as never} />
{:else}
	<p class="text-sm text-base-content/60">読み込み中…</p>
{/if}

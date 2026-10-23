<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { beforeNavigate } from '$app/navigation';
	import { page } from '$app/stores';
	import { ensureSession } from '$lib/session';
	import type { AvatarMe } from '$lib/types';
	import AppShell from '$lib/components/AppShell.svelte';

	let me = $state<AvatarMe | null>(null);
	let loading = $state(true);

	onMount(async () => {
		me = await ensureSession();
		loading = false;
	});

	beforeNavigate(async (nav) => {
		if (nav.to?.url.pathname === '/logout') return;
		const session = await ensureSession();
		if (!session) nav.cancel();
	});
</script>

<svelte:head>
	<title>Growlink</title>
</svelte:head>

{#if loading}
	<div class="flex min-h-screen items-center justify-center bg-base-200 text-sm text-base-content/60">
		読み込み中…
	</div>
{:else if me}
	<AppShell {me} pathname={$page.url.pathname} breadcrumb={$page.data?.breadcrumb}>
		<slot />
	</AppShell>
{/if}

<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api } from '$lib/api';
	import { showToast } from '$lib/toast';
	import type { MeetingPayload } from '$lib/meeting-helpers';

	let date = $state('');
	let uuid = $derived($page.params.uuid);

	onMount(async () => {
		const payload = await api<MeetingPayload>('GET', `/api/v1/meetings/${uuid}`);
		const implementedAt = payload.meeting?.implementedAt as string | undefined;
		const d = implementedAt ? new Date(implementedAt) : new Date();
		date = d.toISOString().slice(0, 10);
	});

	function save() {
		showToast('実施日の API は今後接続されます（表示のみ更新可）', 'info');
	}
</script>

<svelte:head><title>ミーティング設定 - Growlink</title></svelte:head>

<h1 class="font-display text-2xl font-semibold">実施予定日</h1>
<p class="mt-1 text-sm text-base-content/60">この面談の実施予定日を変更します</p>
<div class="mt-6 max-w-lg rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
	<label class="form-control">
		<span class="label-text font-medium">実施予定日</span>
		<input type="date" class="input input-bordered" bind:value={date} />
	</label>
	<div class="mt-6 flex flex-wrap gap-2">
		<button type="button" class="btn btn-primary" onclick={save}>保存</button>
		<a href="/meeting" class="btn btn-ghost">戻る</a>
	</div>
</div>

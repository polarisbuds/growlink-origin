<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { api } from '$lib/api';
	import { showToast } from '$lib/toast';
	import type { MeetingPayload } from '$lib/meeting-helpers';

	type Answer = { id: number; content?: string; topic?: { displayNo?: number; question?: string; description?: string; isRequire?: boolean } };

	let answers = $state<Answer[]>([]);
	let drafts = $state<Record<number, string>>({});
	let uuid = $derived($page.params.uuid);

	onMount(async () => {
		const payload = await api<MeetingPayload>('GET', `/api/v1/meetings/${uuid}`);
		answers = (payload.answers || []) as Answer[];
		const next: Record<number, string> = {};
		for (const a of answers) next[a.id] = a.content || '';
		drafts = next;
	});

	function answerPayload() {
		return answers.map((a) => ({ ...a, content: drafts[a.id] || '' }));
	}

	async function saveDraft() {
		await api('POST', `/api/v1/meetings/${uuid}/draft`, { answers: answerPayload() });
		showToast('下書きを保存しました');
	}
	async function submit() {
		await api('POST', `/api/v1/meetings/${uuid}/submit`, { answers: answerPayload() });
		showToast('提出しました');
		await goto('/meeting');
	}
	async function cancelSubmit() {
		await api('POST', `/api/v1/meetings/${uuid}/draft/cancel`, null);
		showToast('提出を取り消しました');
	}
</script>

<svelte:head><title>面談記入 - Growlink</title></svelte:head>

<div class="flex max-w-2xl items-center justify-between gap-4">
	<div>
		<h1 class="font-display text-2xl font-semibold">面談記入</h1>
		<p class="mt-1 text-sm text-base-content/60">各質問への回答を入力し、下書き保存または提出できます</p>
	</div>
	<a href="/meeting" class="btn btn-ghost btn-sm shrink-0">戻る</a>
</div>
<div class="mt-6 max-w-2xl space-y-4">
	{#each answers as a, i (a.id)}
		<div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
			<div class="flex items-start gap-2">
				<span class="badge badge-primary badge-outline shrink-0">Q{a.topic?.displayNo || i + 1}</span>
				<div class="min-w-0 flex-1">
					<p class="font-medium leading-relaxed">{a.topic?.question}</p>
					{#if a.topic?.description}
						<p class="mt-1 whitespace-pre-wrap text-sm text-base-content/60">{a.topic.description}</p>
					{/if}
				</div>
				{#if a.topic?.isRequire}
					<span class="badge badge-error badge-outline badge-sm shrink-0">必須</span>
				{/if}
			</div>
			<textarea class="textarea textarea-bordered mt-3 w-full" rows="3" bind:value={drafts[a.id]}></textarea>
		</div>
	{/each}
</div>
<div class="mt-6 flex max-w-2xl flex-wrap gap-2">
	<button type="button" class="btn btn-outline" onclick={saveDraft}>下書き保存</button>
	<button type="button" class="btn btn-primary" onclick={submit}>提出</button>
	<button type="button" class="btn btn-ghost" onclick={cancelSubmit}>提出キャンセル</button>
</div>

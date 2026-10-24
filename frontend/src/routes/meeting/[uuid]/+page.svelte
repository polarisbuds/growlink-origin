<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { api } from '$lib/api';
	import { showToast } from '$lib/toast';
	import { formatDate, displayTopics, type MeetingPayload } from '$lib/meeting-helpers';
	import { statusBadge, statusLabel } from '$lib/status';

	let payload = $state<MeetingPayload | null>(null);
	let comment = $state('');

	let rows = $derived(payload ? displayTopics(payload) : []);
	let meeting = $derived((payload?.meeting || {}) as Record<string, unknown>);
	let template = $derived(meeting.template as { title?: string } | undefined);
	let canEdit = $derived(!!payload?.canEdit);
	let uuid = $derived(page.params.uuid);

	function rowKey(row: { topic?: unknown; displayNo: number }) {
		const topic = row.topic as { id?: number } | undefined;
		return topic?.id ?? row.displayNo;
	}

	onMount(async () => {
		payload = await api<MeetingPayload>('GET', `/api/v1/meetings/${uuid}`);
	});

	async function postComment(isLimited: boolean) {
		if (!comment.trim() || !payload) return;
		await api('POST', `/api/v1/meetings/${uuid}/meeting-comment`, {
			content: comment,
			isLimited
		});
		comment = '';
		showToast('投稿しました');
	}
</script>

<svelte:head><title>面談閲覧 - Growlink</title></svelte:head>

{#if payload}
	<div class="flex max-w-2xl min-w-0 w-full flex-wrap items-start justify-between gap-4">
		<div class="min-w-0">
			<h1 class="truncate font-display text-2xl font-semibold">{template?.title || '面談'}</h1>
			<div class="mt-2 flex flex-wrap items-center gap-2 text-sm text-base-content/70">
				<span>実施日: {formatDate(meeting.implementedAt as string)}</span>
				<span class="badge badge-sm {statusBadge(meeting.status as number)}">{statusLabel(meeting.status as number)}</span>
			</div>
			<p class="mt-1 text-sm text-base-content/60">提出済みの回答とコメントを確認できます</p>
		</div>
		<div class="flex shrink-0 flex-wrap gap-2">
			{#if canEdit}
				<a href="/meeting/{uuid}/setting" class="btn btn-outline btn-sm">設定</a>
			{/if}
			<a href="/meeting/{uuid}/answer" class="btn btn-outline btn-sm">更新</a>
			<a href="/meeting" class="btn btn-ghost btn-sm">戻る</a>
		</div>
	</div>
	<div class="mt-6 max-w-2xl space-y-4">
		{#each rows as row (rowKey(row))}
			<div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
				<div class="flex items-start gap-2">
					<span class="badge badge-primary badge-outline shrink-0">Q{row.displayNo}</span>
					<div class="min-w-0 flex-1">
						<p class="font-medium leading-relaxed">{(row.topic as { question?: string })?.question || '質問'}</p>
						{#if (row.topic as { description?: string })?.description}
							<p class="mt-1 whitespace-pre-wrap text-sm text-base-content/60">{(row.topic as { description?: string }).description}</p>
						{/if}
					</div>
				</div>
				{#if (row.answer as { content?: string })?.content}
					<p class="mt-3 whitespace-pre-wrap text-sm text-base-content/80">{(row.answer as { content?: string }).content}</p>
				{:else}
					<p class="mt-3 text-sm italic text-base-content/40">（未回答）</p>
				{/if}
			</div>
		{/each}
	</div>
	<div class="mt-6 max-w-2xl rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
		<h2 class="font-semibold">ミーティングコメント</h2>
		<p class="mt-1 text-xs text-base-content/60">ペア以外への共有や限定公開が選べます</p>
		<textarea class="textarea textarea-bordered mt-3 w-full" rows="2" bind:value={comment} placeholder="コメント"></textarea>
		<div class="mt-2 flex flex-wrap gap-2">
			<button type="button" class="btn btn-sm btn-primary" onclick={() => postComment(false)}>コメント投稿</button>
			<button type="button" class="btn btn-sm btn-outline" onclick={() => postComment(true)}>限定公開コメント</button>
		</div>
		<p class="mt-2 text-xs text-base-content/50">限定公開はペア（メンター・メンティ）のみが閲覧できます</p>
	</div>
{:else}
	<p class="text-sm text-base-content/60">読み込み中…</p>
{/if}

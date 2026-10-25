<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { statusBadge, statusLabel } from '$lib/status';
	import { formatDate, pairLabel, splitByToday, type MeetingListItem } from '$lib/meeting-helpers';

	let tab = $state<'own' | 'shared'>('own');
	let untilToday = $state<MeetingListItem[]>([]);
	let fromTomorrow = $state<MeetingListItem[]>([]);
	let sharedMeetings = $state<MeetingListItem[]>([]);

	onMount(async () => {
		const [own, shared] = await Promise.all([
			api<MeetingListItem[]>('GET', '/api/v1/meetings/involved'),
			api<MeetingListItem[]>('GET', '/api/v1/meetings/shared')
		]);
		const split = splitByToday(own);
		untilToday = split.untilToday;
		fromTomorrow = split.fromTomorrow;
		sharedMeetings = shared;
	});
</script>

<svelte:head><title>ワークスペース - Growlink</title></svelte:head>

<h1 class="font-display text-2xl font-semibold">ワークスペース</h1>
<p class="mt-1 text-sm text-base-content/60">あなたの1on1と共有された1on1</p>
<div class="tabs tabs-boxed mt-6 w-fit">
	<button type="button" class="tab" class:tab-active={tab === 'own'} onclick={() => (tab = 'own')}>
		あなたの1on1
	</button>
	<button type="button" class="tab" class:tab-active={tab === 'shared'} onclick={() => (tab = 'shared')}>
		共有された1on1
	</button>
</div>

{#if tab === 'own'}
	<section class="mt-8">
		<h2 class="text-sm font-semibold text-base-content/70">今日まで</h2>
		<div class="mt-4 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
			{#each untilToday as item (item.uuid)}
				<div
					class="group rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm hover:border-primary/40 transition-colors"
				>
					<div class="flex items-start justify-between gap-3">
						<div class="min-w-0">
							<p class="text-sm text-base-content/60 truncate">{pairLabel(item)}</p>
							<h3 class="mt-1 font-semibold">
								{item.template?.title || '1on1'}
							</h3>
						</div>
						<span class="badge shrink-0 {statusBadge(item.status ?? 0)}">
							{statusLabel(item.status ?? 0)}
						</span>
					</div>
					<p class="mt-3 text-sm text-base-content/70">
						実施予定日: {formatDate(item.implementedAt)}
					</p>
					<div class="mt-4 flex flex-wrap gap-2 items-center">
						<a href="/meeting/{item.uuid}/answer" class="btn btn-sm btn-primary">記入</a>
						<a href="/meeting/{item.uuid}" class="btn btn-sm btn-ghost">閲覧</a>
						<a href="/meeting/{item.uuid}/setting" class="btn btn-sm btn-ghost">設定</a>
					</div>
				</div>
			{/each}
			{#if untilToday.length === 0}
				<div
					class="rounded-3xl border border-dashed border-base-300 bg-base-100 p-10 text-center text-sm text-base-content/50 md:col-span-2 xl:col-span-3"
				>
					今日までの面談はありません
				</div>
			{/if}
		</div>
		<h2 class="mt-10 text-sm font-semibold text-base-content/70">明日以降</h2>
		<div class="mt-4 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
			{#each fromTomorrow as item (item.uuid)}
				<div class="group rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
					<p class="text-sm text-base-content/60">{pairLabel(item)}</p>
					<h3 class="mt-1 font-semibold">{item.template?.title || '1on1'}</h3>
					<p class="mt-3 text-sm text-base-content/70">
						実施予定日: {formatDate(item.implementedAt)}
					</p>
					<div class="mt-4 flex items-center gap-2">
						<a href="/meeting/{item.uuid}/answer" class="btn btn-sm btn-primary">記入</a>
						<a href="/meeting/{item.uuid}" class="btn btn-sm btn-ghost">閲覧</a>
					</div>
				</div>
			{/each}
			{#if fromTomorrow.length === 0}
				<div
					class="rounded-3xl border border-dashed border-base-300 bg-base-100 p-10 text-center text-sm text-base-content/50 md:col-span-2 xl:col-span-3"
				>
					明日以降の面談はありません
				</div>
			{/if}
		</div>
	</section>
{:else}
	<section class="mt-8">
		<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
			{#each sharedMeetings as item (item.uuid)}
				<div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
					<p class="text-sm text-base-content/60">{pairLabel(item)}</p>
					<h3 class="mt-1 font-semibold">{item.template?.title || '1on1'}</h3>
					<p class="mt-3 text-sm text-base-content/70">
						実施予定日: {formatDate(item.implementedAt)}
					</p>
					<div class="mt-4 flex gap-2">
						<span class="badge badge-sm {statusBadge(item.status ?? 0)}">
							{statusLabel(item.status ?? 0)}
						</span>
						<a href="/meeting/{item.uuid}" class="btn btn-sm btn-ghost ml-auto">閲覧</a>
					</div>
				</div>
			{/each}
			{#if sharedMeetings.length === 0}
				<div
					class="rounded-3xl border border-dashed border-base-300 bg-base-100 p-10 text-center text-sm text-base-content/50 md:col-span-2 xl:col-span-3"
				>
					共有された1on1はありません
				</div>
			{/if}
		</div>
	</section>
{/if}

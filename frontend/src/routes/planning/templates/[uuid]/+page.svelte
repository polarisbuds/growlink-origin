<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api } from '$lib/api';

	type Topic = { displayNo?: number; question?: string; description?: string; isRequire?: boolean };
	type Template = { title?: string; topics?: Topic[] };

	let tpl = $state<Template | null>(null);
	let uuid = $derived($page.params.uuid);

	onMount(async () => {
		tpl = await api<Template>('GET', `/api/v1/templates/${uuid}`);
	});
</script>

<svelte:head><title>テンプレート詳細 - Growlink</title></svelte:head>

{#if tpl}
	<h1 class="max-w-2xl font-display text-2xl font-semibold">{tpl.title}</h1>
	<p class="mt-1 max-w-2xl text-sm text-base-content/60">質問一覧</p>
	<ul class="mt-6 max-w-2xl space-y-3">
		{#each tpl.topics || [] as topic (topic.displayNo)}
			<li class="rounded-2xl border border-base-300 bg-base-100 p-4 shadow-sm">
				<div class="flex items-start gap-2">
					<span class="badge badge-primary badge-outline shrink-0">Q{topic.displayNo}</span>
					<div class="min-w-0 flex-1">
						<p class="font-medium leading-relaxed">{topic.question}</p>
						{#if topic.description}
							<p class="mt-1 whitespace-pre-wrap text-sm text-base-content/60">{topic.description}</p>
						{/if}
					</div>
					{#if topic.isRequire}
						<span class="badge badge-error badge-outline badge-sm shrink-0">必須</span>
					{/if}
				</div>
			</li>
		{/each}
	</ul>
	<a href="/planning/templates" class="btn btn-ghost mt-6">一覧へ</a>
{:else}
	<p class="text-sm text-base-content/60">読み込み中…</p>
{/if}

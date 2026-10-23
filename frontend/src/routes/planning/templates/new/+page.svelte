<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { showToast } from '$lib/toast';

	type Topic = { question: string; description: string; isRequire: boolean; displayNo: number };

	let title = $state('');
	let topics = $state<Topic[]>([{ question: '', description: '', isRequire: true, displayNo: 1 }]);

	function addTopic() {
		topics = [...topics, { question: '', description: '', isRequire: true, displayNo: topics.length + 1 }];
	}
	function removeTopic(i: number) {
		if (topics.length <= 1) return;
		topics = topics.filter((_, idx) => idx !== i).map((t, idx) => ({ ...t, displayNo: idx + 1 }));
	}
	async function save() {
		try {
			await api('POST', '/api/v1/templates', { title, topics });
			showToast('テンプレートを作成しました');
			await goto('/planning/templates');
		} catch {
			showToast('保存に失敗しました', 'error');
		}
	}
</script>

<svelte:head><title>テンプレート作成 - Growlink</title></svelte:head>

<main class="mx-auto w-full max-w-3xl p-4 md:p-8">
	<h1 class="font-display text-2xl font-semibold">テンプレート新規作成</h1>
	<p class="mt-1 text-sm text-base-content/60">質問を追加して面談のひな形を作ります</p>
	<div class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
		<label class="form-control">
			<span class="label-text font-medium">タイトル</span>
			<input class="input input-bordered" bind:value={title} />
		</label>
		{#each topics as topic, i (i)}
			<div class="flex items-center gap-2 rounded-2xl border border-base-300 p-2">
				<span class="badge badge-ghost shrink-0">Q{i + 1}</span>
				<div class="min-w-0 flex-1 space-y-1">
					<input class="input input-bordered input-sm w-full" bind:value={topic.question} placeholder="質問文" />
					<textarea
						class="textarea textarea-bordered textarea-sm min-h-[2.5rem] w-full"
						rows="2"
						bind:value={topic.description}
						placeholder="説明（改行可能）"
					></textarea>
				</div>
				<label class="label shrink-0 cursor-pointer gap-1 py-0">
					<input type="checkbox" class="checkbox checkbox-sm" bind:checked={topic.isRequire} />
					<span class="label-text text-xs">必須</span>
				</label>
				<button
					type="button"
					class="btn btn-ghost btn-xs btn-square text-error"
					title="削除"
					disabled={topics.length <= 1}
					onclick={() => removeTopic(i)}>×</button
				>
			</div>
		{/each}
		<button type="button" class="btn btn-outline btn-sm w-full" onclick={addTopic}>質問を追加</button>
		<button type="button" class="btn btn-primary" onclick={save}>保存</button>
	</div>
</main>

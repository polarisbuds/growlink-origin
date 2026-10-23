<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { showToast } from '$lib/toast';

	type Profile = {
		name?: string;
		email?: string;
		lineWorksId?: string;
		imageUri?: string | null;
		remindEnabled?: boolean;
	};

	function apiErrorMessage(e: unknown, fallback: string): string {
		if (e instanceof Error) {
			try {
				const j = JSON.parse(e.message) as { error?: string };
				if (j?.error) return j.error;
			} catch {
				/* not json */
			}
			if (e.message && e.message !== 'unauthorized') return e.message;
		}
		return fallback;
	}

	let form = $state({
		name: '',
		email: '',
		lineWorksId: '',
		remindEnabled: false
	});
	let imageUri = $state<string | null>(null);
	let fieldError = $state('');

	let avatarUrl = $derived(imageUri ? `/uploads/${imageUri}` : null);
	let initial = $derived((form.name || '?').slice(0, 1));

	onMount(async () => {
		const profile = await api<Profile>('GET', '/api/v1/setting/profile');
		form = {
			name: profile.name || '',
			email: profile.email || '',
			lineWorksId: profile.lineWorksId || '',
			remindEnabled: !!profile.remindEnabled
		};
		imageUri = profile.imageUri ?? null;
	});

	async function saveProfile() {
		fieldError = '';
		try {
			await api('PUT', '/api/v1/setting', form);
			showToast('保存しました');
		} catch (e) {
			fieldError = apiErrorMessage(e, '保存に失敗しました');
			showToast(fieldError, 'error');
		}
	}

	async function uploadIcon(e: Event) {
		const input = e.currentTarget as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		const fd = new FormData();
		fd.append('icon', file);
		try {
			await api('POST', '/api/v1/setting/icon', fd);
			const data = await api<Profile>('GET', '/api/v1/setting/profile');
			imageUri = data.imageUri ?? null;
			showToast('アイコンをアップロードしました');
		} catch {
			showToast('アップロードに失敗しました', 'error');
		}
		input.value = '';
	}

	async function deleteIcon() {
		try {
			await api('DELETE', '/api/v1/setting/icon');
			imageUri = null;
			showToast('アイコンを削除しました');
		} catch {
			showToast('削除に失敗しました', 'error');
		}
	}
</script>

<svelte:head><title>設定 - Growlink</title></svelte:head>

<div class="max-w-3xl mx-auto w-full">
	<h1 class="font-display text-2xl font-semibold">個人設定</h1>
	<p class="text-sm text-base-content/60 mt-1">プロフィールと通知の設定</p>

	<div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-5">
		<h2 class="font-semibold">基本情報</h2>
		<div class="flex flex-col sm:flex-row items-center sm:items-start gap-5">
			{#if avatarUrl}
				<img src={avatarUrl} alt="" class="w-24 h-24 rounded-full object-cover border-2 border-base-300 shadow-sm" />
			{:else}
				<div
					class="w-24 h-24 rounded-full bg-primary text-primary-content flex items-center justify-center text-3xl font-semibold shadow-sm"
				>
					{initial}
				</div>
			{/if}
			<div class="flex flex-col gap-2 w-full sm:flex-1">
				<input
					type="file"
					class="file-input file-input-bordered file-input-sm w-full"
					accept="image/*"
					onchange={uploadIcon}
				/>
				<button type="button" class="btn btn-ghost btn-sm w-fit" onclick={deleteIcon}>アイコンを削除</button>
			</div>
		</div>
		<label class="form-control">
			<span class="label-text font-medium">氏名</span>
			<input class="input input-bordered" bind:value={form.name} />
		</label>
		<label class="form-control">
			<span class="label-text font-medium">メール</span>
			<input class="input input-bordered" bind:value={form.email} type="email" />
		</label>
		<label class="form-control">
			<span class="label-text font-medium">LINE WORKS ID</span>
			<div class="join w-full">
				<input
					class="input input-bordered join-item w-full"
					class:input-error={!!fieldError}
					bind:value={form.lineWorksId}
					placeholder="ローカル部のみ"
				/>
				<span class="btn join-item no-animation pointer-events-none bg-base-200">@cap</span>
			</div>
			<span class="label-text-alt text-base-content/50 mt-1">ドメイン @cap は固定です。変更できません。</span>
			{#if fieldError}
				<span class="label-text-alt text-error mt-1">{fieldError}</span>
			{/if}
		</label>
		<button type="button" class="btn btn-primary w-full sm:w-auto" onclick={saveProfile}>保存</button>
	</div>

	<div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
		<h2 class="font-semibold">通知</h2>
		<label class="label cursor-pointer justify-start gap-3 mt-3 rounded-2xl border border-base-300 px-4 py-3">
			<input type="checkbox" class="toggle toggle-primary" bind:checked={form.remindEnabled} />
			<span class="label-text">トピック未提出のリマインドを受け取る</span>
		</label>
		<p class="text-xs text-base-content/50 mt-3">送信処理はバックエンドのジョブで実行されます（設定の保存のみ）</p>
		<button type="button" class="btn btn-outline btn-sm mt-4" onclick={saveProfile}>通知設定を保存</button>
	</div>
</div>

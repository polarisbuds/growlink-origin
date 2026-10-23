<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { filterItems } from '$lib/searchable-select';

	type Pairing = {
		id: number;
		mentor?: { name?: string };
		mentee?: { name?: string };
		targets?: { label?: string }[];
	};
	type Template = { id: number; title?: string };
	type User = { id: number; name: string; email: string };

	function todayISO() {
		const d = new Date();
		return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
	}
	function currentMonthISO() {
		const d = new Date();
		return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
	}
	function pairingSummary(p: Pairing) {
		const mentor = p.mentor?.name || '—';
		const targetLabels = (p.targets || []).map((t) => t.label).filter(Boolean);
		const targets = targetLabels.length ? targetLabels.join('、') : p.mentee?.name || '—';
		return `${mentor} / ${targets}`;
	}

	let pairings = $state<Pairing[]>([]);
	let templates = $state<Template[]>([]);
	let users = $state<User[]>([]);
	let mode = $state<'single' | 'all' | 'monthly'>('single');
	let sharerUserIds = $state<number[]>([]);
	let submitting = $state(false);
	let error = $state('');

	let form = $state({
		pairingId: 0,
		templateId: 0,
		implementedAt: todayISO(),
		baseMonth: currentMonthISO(),
		week: 2,
		weekday: 3
	});

	let pairPickOpen = $state(false);
	let pairPickQuery = $state('');
	let templatePickOpen = $state(false);
	let templatePickQuery = $state('');
	let sharerPickOpen = $state(false);
	let sharerPickQuery = $state('');

	let pairPickFiltered = $derived(filterItems(pairings, pairPickQuery, pairingSummary));
	let templatePickFiltered = $derived(
		filterItems(templates, templatePickQuery, (t) => t.title || '')
	);
	let sharerPickFiltered = $derived.by(() => {
		const picked = new Set(sharerUserIds);
		const list = users.filter((u) => !picked.has(u.id));
		return filterItems(list, sharerPickQuery, (u) => `${u.name} (${u.email})`);
	});

	function userLabel(id: number) {
		const u = users.find((x) => x.id === id);
		return u ? u.name : `#${id}`;
	}
	function closePicks() {
		pairPickOpen = false;
		templatePickOpen = false;
		sharerPickOpen = false;
	}
	function submitLabel() {
		if (submitting) return '作成中…';
		if (mode === 'single') return '作成';
		if (mode === 'all') return '全ペアに作成';
		return '全ペアに定期作成';
	}
	function addSharer(u: User) {
		if (!u?.id || sharerUserIds.includes(u.id)) return;
		sharerUserIds = [...sharerUserIds, u.id];
		sharerPickQuery = '';
		sharerPickOpen = false;
	}
	function removeSharer(id: number) {
		sharerUserIds = sharerUserIds.filter((x) => x !== id);
	}

	onMount(async () => {
		const [p, t, dir] = await Promise.all([
			api<Pairing[]>('GET', '/api/v1/pairings'),
			api<Template[]>('GET', '/api/v1/templates'),
			api<{ users?: User[] }>('GET', '/api/v1/planning/directory')
		]);
		pairings = p;
		templates = t;
		users = dir.users || [];
		form.templateId = templates[0]?.id || 0;

		const onMouseDown = (ev: MouseEvent) => {
			const t = ev.target;
			if (!(t instanceof Element)) return;
			if (!t.closest('[data-pick-root="pair"]')) pairPickOpen = false;
			if (!t.closest('[data-pick-root="template"]')) templatePickOpen = false;
			if (!t.closest('[data-pick-root="sharer"]')) sharerPickOpen = false;
		};
		window.addEventListener('mousedown', onMouseDown);
		return () => window.removeEventListener('mousedown', onMouseDown);
	});

	async function submit() {
		error = '';
		if (!form.templateId) {
			error = 'テンプレートを選択してください。';
			return;
		}
		submitting = true;
		try {
			let url: string;
			let body: Record<string, unknown>;
			const sharerPayload = sharerUserIds.length ? sharerUserIds : undefined;
			if (mode === 'single') {
				if (!form.pairingId) {
					error = 'ペアを選択してください。';
					return;
				}
				if (!form.implementedAt) {
					error = '実施日を入力してください。';
					return;
				}
				url = '/api/v1/meetings';
				body = {
					pairingId: form.pairingId,
					templateId: form.templateId,
					implementedAt: form.implementedAt,
					sharerUserIds: sharerPayload
				};
			} else if (mode === 'all') {
				if (!form.implementedAt) {
					error = '実施日を入力してください。';
					return;
				}
				url = '/api/v1/meetings/all';
				body = {
					templateId: form.templateId,
					implementedAt: form.implementedAt,
					sharerUserIds: sharerPayload
				};
			} else {
				if (!form.baseMonth) {
					error = '基準月を入力してください。';
					return;
				}
				if (form.week < 1 || form.week > 5) {
					error = '第N週は 1〜5 で指定してください。';
					return;
				}
				url = '/api/v1/meetings/monthly';
				body = {
					templateId: form.templateId,
					baseDate: `${form.baseMonth}-01`,
					week: form.week,
					weekday: form.weekday,
					sharerUserIds: sharerPayload
				};
			}
			await api('POST', url, body);
			await goto('/planning/meetings');
		} catch (e) {
			error = e instanceof Error ? e.message : '作成に失敗しました。';
		} finally {
			submitting = false;
		}
	}
</script>

<svelte:window onkeydown={(e) => e.key === 'Escape' && closePicks()} />

<svelte:head><title>1on1作成 - Growlink</title></svelte:head>

<main class="mx-auto w-full max-w-3xl min-w-0 p-4 md:p-8">
	<h1 class="font-display text-2xl font-semibold">1on1ミーティング作成</h1>
	<p class="mt-1 text-sm text-base-content/60">ペアとテンプレートを選び、実施日または定期スケジュールを指定します</p>
	<div class="mt-4 flex flex-wrap gap-2">
		<button type="button" class="btn btn-sm" class:btn-primary={mode === 'single'} class:btn-outline={mode !== 'single'} onclick={() => (mode = 'single')}>1件作成</button>
		<button type="button" class="btn btn-sm" class:btn-primary={mode === 'all'} class:btn-outline={mode !== 'all'} onclick={() => (mode = 'all')}>全ペア一括</button>
		<button type="button" class="btn btn-sm" class:btn-primary={mode === 'monthly'} class:btn-outline={mode !== 'monthly'} onclick={() => (mode = 'monthly')}>定期作成</button>
	</div>
	<div class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
		{#if error}
			<div class="alert alert-error text-sm">{error}</div>
		{/if}
		{#if mode === 'single'}
			<div class="form-control relative" data-pick-root="pair">
				<span class="label-text font-medium">ペア</span>
				<div class="relative mt-1">
					<input
						type="text"
						class="input input-bordered w-full pr-10"
						autocomplete="off"
						placeholder={form.pairingId ? '' : '名前で検索して選択'}
						value={pairPickOpen
							? pairPickQuery
							: form.pairingId
								? pairingSummary(pairings.find((p) => p.id === form.pairingId)!)
								: ''}
						onfocus={() => {
							pairPickOpen = true;
							if (form.pairingId) pairPickQuery = '';
						}}
						oninput={(e) => {
							pairPickQuery = e.currentTarget.value;
							pairPickOpen = true;
						}}
					/>
					{#if form.pairingId}
						<button
							type="button"
							class="btn btn-ghost btn-xs btn-circle absolute right-1 top-1/2 -translate-y-1/2"
							aria-label="クリア"
							onclick={() => {
								form.pairingId = 0;
								pairPickQuery = '';
								pairPickOpen = false;
							}}>×</button
						>
					{/if}
				</div>
				{#if pairPickOpen}
					<ul class="absolute z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg">
						{#each pairPickFiltered as p (p.id)}
							<li>
								<button
									type="button"
									class="btn btn-ghost btn-sm w-full justify-start font-normal"
									onclick={() => {
										form.pairingId = p.id;
										pairPickQuery = '';
										pairPickOpen = false;
									}}>{pairingSummary(p)}</button
								>
							</li>
						{/each}
						{#if pairPickFiltered.length === 0}
							<li class="px-4 py-3 text-sm text-base-content/50">該当なし</li>
						{/if}
					</ul>
				{/if}
			</div>
		{/if}
		<div class="form-control relative" data-pick-root="template">
			<span class="label-text font-medium">テンプレート</span>
			<div class="relative mt-1">
				<input
					type="text"
					class="input input-bordered w-full pr-10"
					autocomplete="off"
					placeholder={form.templateId ? '' : '名前で検索して選択'}
					value={templatePickOpen
						? templatePickQuery
						: form.templateId
							? templates.find((t) => t.id === form.templateId)?.title || ''
							: ''}
					onfocus={() => {
						templatePickOpen = true;
						if (form.templateId) templatePickQuery = '';
					}}
					oninput={(e) => {
						templatePickQuery = e.currentTarget.value;
						templatePickOpen = true;
					}}
				/>
				{#if form.templateId}
					<button
						type="button"
						class="btn btn-ghost btn-xs btn-circle absolute right-1 top-1/2 -translate-y-1/2"
						aria-label="クリア"
						onclick={() => {
							form.templateId = 0;
							templatePickQuery = '';
							templatePickOpen = false;
						}}>×</button
					>
				{/if}
			</div>
			{#if templatePickOpen}
				<ul class="absolute z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg">
					{#each templatePickFiltered as t (t.id)}
						<li>
							<button
								type="button"
								class="btn btn-ghost btn-sm w-full justify-start font-normal"
								onclick={() => {
									form.templateId = t.id;
									templatePickQuery = '';
									templatePickOpen = false;
								}}>{t.title}</button
							>
						</li>
					{/each}
					{#if templatePickFiltered.length === 0}
						<li class="px-4 py-3 text-sm text-base-content/50">該当なし</li>
					{/if}
				</ul>
			{/if}
		</div>
		<div class="form-control relative" data-pick-root="sharer">
			<span class="label-text font-medium">共有者（任意・複数）</span>
			<div class="mt-2 flex flex-wrap gap-2">
				{#each sharerUserIds as id (id)}
					<span class="badge badge-lg gap-1">
						<span>{userLabel(id)}</span>
						<button type="button" class="btn btn-ghost btn-xs btn-circle" onclick={() => removeSharer(id)}>×</button>
					</span>
				{/each}
			</div>
			<div class="relative mt-2">
				<input
					type="text"
					class="input input-bordered w-full"
					placeholder="名前で検索して追加"
					autocomplete="off"
					value={sharerPickOpen ? sharerPickQuery : ''}
					onfocus={() => (sharerPickOpen = true)}
					oninput={(e) => {
						sharerPickQuery = e.currentTarget.value;
						sharerPickOpen = true;
					}}
				/>
			</div>
			{#if sharerPickOpen}
				<ul class="absolute z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg">
					{#each sharerPickFiltered as u (u.id)}
						<li>
							<button type="button" class="btn btn-ghost btn-sm w-full justify-start font-normal" onclick={() => addSharer(u)}
								>{u.name} ({u.email})</button
							>
						</li>
					{/each}
				</ul>
			{/if}
		</div>
		{#if mode !== 'monthly'}
			<label class="form-control">
				<span class="label-text font-medium">実施日</span>
				<input class="input input-bordered" type="date" bind:value={form.implementedAt} />
			</label>
		{:else}
			<div class="space-y-4">
				<label class="form-control">
					<span class="label-text font-medium">基準月（その月の第N週を計算）</span>
					<input class="input input-bordered" type="month" bind:value={form.baseMonth} />
				</label>
				<label class="form-control">
					<span class="label-text font-medium">第N週（1〜5）</span>
					<input class="input input-bordered" type="number" min="1" max="5" bind:value={form.week} />
				</label>
				<label class="form-control">
					<span class="label-text font-medium">曜日</span>
					<select class="select select-bordered" bind:value={form.weekday}>
						<option value={0}>日曜</option>
						<option value={1}>月曜</option>
						<option value={2}>火曜</option>
						<option value={3}>水曜</option>
						<option value={4}>木曜</option>
						<option value={5}>金曜</option>
						<option value={6}>土曜</option>
					</select>
				</label>
			</div>
		{/if}
		<div class="flex flex-wrap gap-2 pt-2">
			<button type="button" class="btn btn-primary" disabled={submitting} onclick={submit}>{submitLabel()}</button>
			<a href="/planning/meetings" class="btn btn-ghost">戻る</a>
		</div>
	</div>
</main>

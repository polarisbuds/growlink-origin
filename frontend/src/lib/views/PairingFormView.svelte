<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { showToast } from '$lib/toast';
	import { filterItems } from '$lib/searchable-select';
	import UserChipSearchField from '$lib/components/UserChipSearchField.svelte';

	type Target = { kind: string; targetId: number; label: string };
	type User = { id: number; name: string; email: string };
	type PairingInitial = {
		uuid?: string;
		mentor?: { id?: number };
		mentee?: { id?: number; name?: string };
		targets?: { kind: string; targetId: number; label?: string }[];
	};

	let {
		mode,
		initialPairing = null
	}: {
		mode: 'create' | 'edit';
		initialPairing?: PairingInitial | null;
	} = $props();

	let mentorId = $state(0);
	let targets = $state<Target[]>([]);
	let pairingUUID = $state('');
	let users = $state<User[]>([]);

	let mentorPickOpen = $state(false);
	let mentorPickQuery = $state('');
	let userTargetPickOpen = $state(false);
	let userTargetPickQuery = $state('');

	let userTargets = $derived(targets.filter((t) => t.kind === 'user'));
	let menteeChips = $derived(
		userTargets.map((t) => ({ id: t.targetId, name: displayNameForTarget(t) }))
	);

	function displayNameForTarget(t: Target) {
		const u = users.find((x) => x.id === t.targetId);
		if (u) return u.name;
		const idx = t.label.indexOf(' (');
		if (idx > 0) return t.label.slice(0, idx);
		return t.label || `#${t.targetId}`;
	}
	function userLabel(id: number) {
		const u = users.find((x) => x.id === id);
		return u ? u.name : `#${id}`;
	}
	function targetKey(kind: string, targetId: number) {
		return `${kind}:${targetId}`;
	}
	function addUserTarget(u: User) {
		const key = targetKey('user', u.id);
		if (targets.some((t) => targetKey(t.kind, t.targetId) === key)) return;
		targets = [...targets, { kind: 'user', targetId: u.id, label: `${u.name} (${u.email})` }];
	}

	let mentorPickFiltered = $derived(filterItems(users, mentorPickQuery, (u) => `${u.name} (${u.email})`));
	let userTargetPickFiltered = $derived.by(() => {
		const picked = new Set(userTargets.map((t) => t.targetId));
		const list = users.filter((u) => !picked.has(u.id));
		return filterItems(list, userTargetPickQuery, (u) => `${u.name} (${u.email})`);
	});

	$effect(() => {
		if (!initialPairing) return;
		pairingUUID = initialPairing.uuid || '';
		mentorId = initialPairing.mentor?.id || 0;
		const loaded = (initialPairing.targets || [])
			.filter((t) => t.kind === 'user')
			.map((t) => ({
				kind: 'user',
				targetId: t.targetId,
				label: t.label || ''
			}));
		if (
			initialPairing.mentee?.id &&
			!loaded.some((t) => t.targetId === initialPairing.mentee?.id)
		) {
			loaded.push({
				kind: 'user',
				targetId: initialPairing.mentee.id,
				label: initialPairing.mentee.name || ''
			});
		}
		targets = loaded;
	});

	onMount(async () => {
		const data = await api<{ users?: User[] }>('GET', '/api/v1/planning/directory');
		users = data.users || [];
		targets = targets.map((t) => {
			if (t.label) return t;
			if (t.kind === 'user') return { ...t, label: userLabel(t.targetId) };
			return t;
		});

		const onMouseDown = (ev: MouseEvent) => {
			const t = ev.target;
			if (!(t instanceof Element)) return;
			if (!t.closest('[data-pick-root="mentor"]')) mentorPickOpen = false;
			if (!t.closest('[data-pick-root="user-target"]')) userTargetPickOpen = false;
		};
		window.addEventListener('mousedown', onMouseDown);
		return () => window.removeEventListener('mousedown', onMouseDown);
	});

	function pickUserTarget(u: User) {
		addUserTarget(u);
		userTargetPickQuery = '';
		userTargetPickOpen = false;
	}
	function removeUserTarget(targetId: number) {
		targets = targets.filter((t) => !(t.kind === 'user' && t.targetId === targetId));
	}

	async function requestSave() {
		const userOnly = targets.filter((t) => t.kind === 'user');
		if (!mentorId || !userOnly.length) {
			showToast('メンターとメンティを選択してください', 'error');
			return;
		}
		const body = {
			mentorId,
			targets: userOnly.map((t) => ({ kind: 'user', targetId: t.targetId }))
		};
		try {
			if (mode === 'create') {
				await api('POST', '/api/v1/pairings', body);
			} else {
				await api('PUT', `/api/v1/pairings/${pairingUUID}`, body);
			}
			showToast('保存しました');
			await goto('/planning/pairs/users');
		} catch {
			showToast('保存に失敗しました', 'error');
		}
	}
</script>

<svelte:head>
	<title>{mode === 'create' ? 'ペア作成' : 'ペア編集'} - Growlink</title>
</svelte:head>

<main class="mx-auto w-full max-w-3xl min-w-0 p-4 md:p-8">
	<h1 class="font-display text-2xl font-semibold">{mode === 'create' ? 'ペア新規作成' : 'ペア編集'}</h1>
	<p class="mt-1 text-sm text-base-content/60">メンターとメンティを{mode === 'create' ? '登録' : '変更'}します</p>
	<div class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
		<div class="grid items-start gap-4 md:grid-cols-2">
			<div class="form-control" data-pick-root="mentor">
				<span class="label-text font-medium">メンター</span>
				<div class="relative mt-1">
					<input
						type="text"
						class="input input-bordered w-full"
						autocomplete="off"
						placeholder={mentorId ? '' : '名前で検索'}
						value={mentorPickOpen ? mentorPickQuery : mentorId ? userLabel(mentorId) : ''}
						onfocus={() => {
							mentorPickOpen = true;
							if (mentorId) mentorPickQuery = '';
						}}
						oninput={(e) => {
							mentorPickQuery = e.currentTarget.value;
							mentorPickOpen = true;
						}}
					/>
					{#if mentorPickOpen}
						<ul class="absolute left-0 top-full z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg">
							{#each mentorPickFiltered as u (u.id)}
								<li>
									<button
										type="button"
										class="btn btn-ghost btn-sm w-full justify-start font-normal"
										onclick={() => {
											mentorId = u.id;
											mentorPickQuery = '';
											mentorPickOpen = false;
										}}>{u.name} ({u.email})</button
									>
								</li>
							{/each}
						</ul>
					{/if}
				</div>
			</div>
			<div class="form-control" data-pick-root="user-target">
				<span class="label-text font-medium">メンティ</span>
				<UserChipSearchField
					chips={menteeChips}
					placeholder={menteeChips.length ? '追加で検索' : '名前で検索して追加'}
					bind:query={userTargetPickQuery}
					bind:open={userTargetPickOpen}
					filteredUsers={userTargetPickFiltered}
					onPick={pickUserTarget}
					onRemove={removeUserTarget}
				/>
			</div>
		</div>
		<div class="flex flex-wrap gap-2 pt-2">
			<button type="button" class="btn btn-primary" onclick={requestSave}>{mode === 'create' ? '作成' : '保存'}</button>
			<a href="/planning/pairs/users" class="btn btn-ghost">戻る</a>
		</div>
	</div>
</main>

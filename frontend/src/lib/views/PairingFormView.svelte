<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { showToast } from '$lib/toast';
	import { filterItems } from '$lib/searchable-select';

	type Target = { kind: string; targetId: number; label: string };
	type User = { id: number; name: string; email: string };
	type Team = { id: number; name: string };
	type Department = { id: number; name: string };
	type PairingInitial = {
		uuid?: string;
		mentor?: { id?: number };
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
	let teams = $state<Team[]>([]);
	let departments = $state<Department[]>([]);
	let teamPickId = $state(0);
	let departmentPickId = $state(0);

	let mentorPickOpen = $state(false);
	let mentorPickQuery = $state('');
	let userTargetPickOpen = $state(false);
	let userTargetPickQuery = $state('');

	function userLabel(id: number) {
		const u = users.find((x) => x.id === id);
		return u ? u.name : `#${id}`;
	}
	function teamLabel(id: number) {
		const t = teams.find((x) => x.id === id);
		return t ? t.name : `#${id}`;
	}
	function departmentLabel(id: number) {
		const d = departments.find((x) => x.id === id);
		return d ? d.name : `#${id}`;
	}
	function targetKey(kind: string, targetId: number) {
		return `${kind}:${targetId}`;
	}
	function targetKindLabel(kind: string) {
		if (kind === 'user') return 'ユーザー';
		if (kind === 'team') return 'チーム';
		if (kind === 'department') return '部署';
		return kind;
	}
	function addTarget(kind: string, targetId: number, label: string) {
		const key = targetKey(kind, targetId);
		if (targets.some((t) => targetKey(t.kind, t.targetId) === key)) return;
		targets = [...targets, { kind, targetId, label }];
	}

	let mentorPickFiltered = $derived(filterItems(users, mentorPickQuery, (u) => `${u.name} (${u.email})`));
	let userTargetPickFiltered = $derived.by(() => {
		const picked = new Set(targets.filter((t) => t.kind === 'user').map((t) => t.targetId));
		const list = users.filter((u) => !picked.has(u.id));
		return filterItems(list, userTargetPickQuery, (u) => `${u.name} (${u.email})`);
	});

	$effect(() => {
		if (!initialPairing) return;
		pairingUUID = initialPairing.uuid || '';
		mentorId = initialPairing.mentor?.id || 0;
		targets = (initialPairing.targets || []).map((t) => ({
			kind: t.kind,
			targetId: t.targetId,
			label: t.label || ''
		}));
	});

	onMount(async () => {
		const data = await api<{ users?: User[]; teams?: Team[]; departments?: Department[] }>(
			'GET',
			'/api/v1/planning/directory'
		);
		users = data.users || [];
		teams = data.teams || [];
		departments = data.departments || [];
		targets = targets.map((t) => {
			if (t.label) return t;
			if (t.kind === 'user') return { ...t, label: userLabel(t.targetId) };
			if (t.kind === 'team') return { ...t, label: teamLabel(t.targetId) };
			if (t.kind === 'department') return { ...t, label: departmentLabel(t.targetId) };
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
		addTarget('user', u.id, `${u.name} (${u.email})`);
		userTargetPickQuery = '';
		userTargetPickOpen = false;
	}
	function removeTarget(index: number) {
		targets = targets.filter((_, i) => i !== index);
	}
	function addTeamTarget() {
		const id = Number(teamPickId);
		if (!id) return;
		addTarget('team', id, teamLabel(id));
		teamPickId = 0;
	}
	function addDepartmentTarget() {
		const id = Number(departmentPickId);
		if (!id) return;
		addTarget('department', id, departmentLabel(id));
		departmentPickId = 0;
	}

	async function requestSave() {
		if (!mentorId || !targets.length) {
			showToast('メンターと対象を選択してください', 'error');
			return;
		}
		const body = {
			mentorId,
			targets: targets.map((t) => ({ kind: t.kind, targetId: t.targetId }))
		};
		try {
			if (mode === 'create') {
				await api('POST', '/api/v1/pairings', body);
			} else {
				await api('PUT', `/api/v1/pairings/${pairingUUID}`, body);
			}
			showToast('保存しました');
			await goto('/planning/pairs');
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
	<p class="mt-1 text-sm text-base-content/60">メンターと対象（ユーザー・チーム・部署）を{mode === 'create' ? '登録' : '変更'}します</p>
	<div class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
		<div class="form-control relative" data-pick-root="mentor">
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
			</div>
			{#if mentorPickOpen}
				<ul class="absolute z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg">
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
		<div>
			<span class="label-text font-medium">対象</span>
			<div class="mt-2 flex flex-wrap gap-2">
				{#each targets as t, i (targetKindLabel(t.kind) + '-' + t.targetId)}
					<span class="badge badge-lg gap-1">
						<span class="text-xs opacity-60">{targetKindLabel(t.kind)}</span>
						<span>{t.label}</span>
						<button type="button" class="btn btn-ghost btn-xs btn-circle" onclick={() => removeTarget(i)}>×</button>
					</span>
				{/each}
			</div>
		</div>
		<div class="form-control relative" data-pick-root="user-target">
			<span class="label-text font-medium">ユーザーを追加</span>
			<div class="relative mt-1">
				<input
					type="text"
					class="input input-bordered w-full"
					placeholder="名前で検索"
					autocomplete="off"
					value={userTargetPickOpen ? userTargetPickQuery : ''}
					onfocus={() => (userTargetPickOpen = true)}
					oninput={(e) => {
						userTargetPickQuery = e.currentTarget.value;
						userTargetPickOpen = true;
					}}
				/>
			</div>
			{#if userTargetPickOpen}
				<ul class="absolute z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg">
					{#each userTargetPickFiltered as u (u.id)}
						<li>
							<button type="button" class="btn btn-ghost btn-sm w-full justify-start font-normal" onclick={() => pickUserTarget(u)}
								>{u.name} ({u.email})</button
							>
						</li>
					{/each}
				</ul>
			{/if}
		</div>
		<div class="flex flex-wrap items-end gap-2">
			<label class="form-control min-w-[12rem] flex-1">
				<span class="label-text font-medium">チームを追加</span>
				<select class="select select-bordered" bind:value={teamPickId}>
					<option value={0}>選択</option>
					{#each teams as t (t.id)}
						<option value={t.id}>{t.name}</option>
					{/each}
				</select>
			</label>
			<button type="button" class="btn btn-outline" onclick={addTeamTarget}>追加</button>
		</div>
		<div class="flex flex-wrap items-end gap-2">
			<label class="form-control min-w-[12rem] flex-1">
				<span class="label-text font-medium">部署を追加</span>
				<select class="select select-bordered" bind:value={departmentPickId}>
					<option value={0}>選択</option>
					{#each departments as d (d.id)}
						<option value={d.id}>{d.name}</option>
					{/each}
				</select>
			</label>
			<button type="button" class="btn btn-outline" onclick={addDepartmentTarget}>追加</button>
		</div>
		<div class="flex flex-wrap gap-2 pt-2">
			<button type="button" class="btn btn-primary" onclick={requestSave}>{mode === 'create' ? '作成' : '保存'}</button>
			<a href="/planning/pairs" class="btn btn-ghost">戻る</a>
		</div>
	</div>
</main>

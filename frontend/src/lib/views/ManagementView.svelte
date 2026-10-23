<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { showToast } from '$lib/toast';
	import Users from 'lucide-svelte/icons/users';
	import Building2 from 'lucide-svelte/icons/building-2';
	import Bell from 'lucide-svelte/icons/bell';
	import Search from 'lucide-svelte/icons/search';
	import ChevronRight from 'lucide-svelte/icons/chevron-right';
	import Info from 'lucide-svelte/icons/info';

	const AUTHORITY_OPTIONS = [
		{ value: 'general', label: 'メンバー' },
		{ value: 'boss', label: 'ボス' },
		{ value: 'owner', label: 'オーナー' }
	];
	const SYSTEM_ROLE_OPTIONS = [
		{ value: 'none', label: 'なし' },
		{ value: 'system_admin', label: 'システム管理者' },
		{ value: 'analyst', label: '分析者' }
	];
	const OPERATION_OPTIONS = [
		{ value: 'none', label: 'なし' },
		{ value: 'department', label: '部署' },
		{ value: 'team', label: 'チーム' },
		{ value: 'subordinate', label: '部下' }
	];
	const WEEKDAYS = [
		{ value: 1, label: '月曜日' },
		{ value: 2, label: '火曜日' },
		{ value: 3, label: '水曜日' },
		{ value: 4, label: '木曜日' },
		{ value: 5, label: '金曜日' },
		{ value: 6, label: '土曜日' },
		{ value: 7, label: '日曜日' }
	];

	type ManagementUser = {
		id: number;
		name: string;
		email: string;
		authority: string;
		systemRole: string;
		operationTarget: string;
		targetIds?: number[];
		departmentId: number;
		roleId: number;
		suspended: boolean;
		imageUri?: string;
	};

	type OrgTreeNode = {
		kind: 'team' | 'department';
		id: number;
		name: string;
		depth: number;
		memberCount?: number;
		leaderName?: string;
	};

	type OrgMember = { userId: number; userName: string; title: string };

	type OrgDetail = {
		id: number;
		name: string;
		kind: 'team' | 'department';
		isRoot?: boolean;
		members?: OrgMember[];
		childIds?: number[];
		teamIds?: number[];
		userIds?: number[];
		memberIds?: number[];
		leaderId?: number;
	};

	type UserForm = {
		id: number;
		name: string;
		authority: string;
		systemRole: string;
		operationTarget: string;
		targetIds: number[];
		departmentId: number;
		roleId: number;
		suspended: boolean;
	};

	type Options = {
		departments: { id: number; name: string }[];
		teams: { id: number; name: string }[];
		roles: { id: number; name: string }[];
	};

	function toggleInList(list: number[], id: number | string): number[] {
		const n = Number(id);
		const i = list.indexOf(n);
		if (i >= 0) list.splice(i, 1);
		else list.push(n);
		return [...list];
	}

	function filterByQuery<T extends Record<string, unknown>>(
		items: T[],
		query: string,
		labelKey = 'name'
	): T[] {
		const q = (query || '').trim().toLowerCase();
		if (!q) return items;
		return items.filter((item) => String(item[labelKey] || '').toLowerCase().includes(q));
	}

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

	let tab = $state<'users' | 'departments' | 'remind'>('users');
	let users = $state<ManagementUser[]>([]);
	let orgTree = $state<OrgTreeNode[]>([]);
	let options = $state<Options>({ departments: [], teams: [], roles: [] });
	let selectedUserIds = $state<number[]>([]);
	let userSearch = $state('');
	let usersLoading = $state(true);
	let orgLoading = $state(true);
	let remind = $state({
		remindAutoSend: false,
		remindMode: 'prev_weekday',
		remindDaysBefore: 1,
		remindWeekday: 1,
		remindTime: '09:00'
	});
	let targetSearch = $state('');
	let orgChildSearch = $state('');
	let orgTeamSearch = $state('');
	let orgUserSearch = $state('');
	let orgMemberSearch = $state('');
	let userModalOpen = $state(false);
	let userForm = $state<UserForm | null>(null);
	let authoritySelect = $state('general');
	let systemRoleSelect = $state('none');
	let createUserOpen = $state(false);
	let createUserForm = $state({ name: '', email: '', departmentId: 0, roleId: 0, password: '' });
	let createDeptOpen = $state(false);
	let createDeptForm = $state({ name: '', parentDepartmentId: 0 });
	let orgModalOpen = $state(false);
	let orgDetail = $state<OrgDetail | null>(null);
	let orgEditMode = $state(false);
	let orgEdit = $state({
		name: '',
		childIds: [] as number[],
		teamIds: [] as number[],
		userIds: [] as number[],
		memberIds: [] as number[],
		leaderId: 0
	});
	let confirmOpen = $state(false);
	let confirmTitle = $state('');
	let confirmMessage = $state('');
	let confirmAction = $state<(() => void | Promise<void>) | null>(null);
	let confirmOnCancel = $state<(() => void) | null>(null);
	let confirmDanger = $state(false);

	let filteredUsers = $derived.by(() => {
		const q = userSearch.trim().toLowerCase();
		if (!q) return users;
		return users.filter(
			(u) =>
				(u.name || '').toLowerCase().includes(q) || (u.email || '').toLowerCase().includes(q)
		);
	});

	let allUsersSelected = $derived(
		filteredUsers.length > 0 && filteredUsers.every((u) => selectedUserIds.includes(u.id))
	);

	function authorityLabel(v: string) {
		return AUTHORITY_OPTIONS.find((o) => o.value === v)?.label || v;
	}

	function systemRoleLabel(v: string) {
		return SYSTEM_ROLE_OPTIONS.find((o) => o.value === v)?.label || v;
	}

	function userName(id: number | string) {
		const u = users.find((x) => x.id === Number(id));
		return u ? u.name : '';
	}

	function userAvatar(u: ManagementUser) {
		return u.imageUri ? `/uploads/${u.imageUri}` : null;
	}

	function userInitial(u: ManagementUser) {
		return (u.name || '?').slice(0, 1);
	}

	function filteredDepartments() {
		return filterByQuery(options.departments || [], targetSearch);
	}

	function filteredTeams() {
		return filterByQuery(options.teams || [], targetSearch);
	}

	function filteredUsersForTarget() {
		return filterByQuery(users, targetSearch);
	}

	function filteredOrgDepartments() {
		return filterByQuery(options.departments || [], orgChildSearch);
	}

	function filteredOrgTeams() {
		return filterByQuery(options.teams || [], orgTeamSearch);
	}

	function filteredOrgUsers() {
		return filterByQuery(users, orgUserSearch);
	}

	function filteredOrgMembers() {
		return filterByQuery(users, orgMemberSearch);
	}

	function isInList(list: number[] | undefined, id: number | string) {
		return (list || []).map(Number).includes(Number(id));
	}

	async function loadUsers() {
		usersLoading = true;
		try {
			users = await api<ManagementUser[]>('GET', '/api/v1/management/users');
		} catch {
			users = [];
		}
		usersLoading = false;
	}

	async function loadOptions() {
		try {
			options = await api<Options>('GET', '/api/v1/management/options');
		} catch {
			/* keep defaults */
		}
	}

	async function loadOrgTree() {
		orgLoading = true;
		try {
			orgTree = await api<OrgTreeNode[]>('GET', '/api/v1/management/org/tree');
		} catch {
			orgTree = [];
		}
		orgLoading = false;
	}

	async function loadRemind() {
		try {
			const data = await api<{
				remindAutoSend?: boolean;
				remindMode?: string;
				remindDaysBefore?: number;
				remindWeekday?: number;
				remindTime?: string;
			}>('GET', '/api/v1/management/settings/remind');
			remind = {
				remindAutoSend: !!data.remindAutoSend,
				remindMode: data.remindMode || 'prev_weekday',
				remindDaysBefore: data.remindDaysBefore && data.remindDaysBefore > 0 ? data.remindDaysBefore : 1,
				remindWeekday: data.remindWeekday || 1,
				remindTime: (data.remindTime || '09:00').slice(0, 5)
			};
		} catch {
			/* keep defaults */
		}
	}

	function toggleAllUsers(checked: boolean) {
		selectedUserIds = checked ? filteredUsers.map((u) => u.id) : [];
	}

	function toggleUserTarget(id: number) {
		if (!userForm) return;
		userForm.targetIds = toggleInList(userForm.targetIds || [], id);
	}

	function toggleOrgChild(id: number) {
		orgEdit.childIds = toggleInList(orgEdit.childIds || [], id);
	}

	function toggleOrgTeam(id: number) {
		orgEdit.teamIds = toggleInList(orgEdit.teamIds || [], id);
	}

	function toggleOrgUser(id: number) {
		orgEdit.userIds = toggleInList(orgEdit.userIds || [], id);
	}

	function toggleOrgMember(id: number) {
		orgEdit.memberIds = toggleInList(orgEdit.memberIds || [], id);
		if (!orgEdit.memberIds.includes(Number(orgEdit.leaderId))) {
			orgEdit.leaderId = orgEdit.memberIds[0] || 0;
		}
	}

	function setOperationTarget(value: string) {
		if (!userForm) return;
		userForm.operationTarget = value;
		userForm.targetIds = [];
	}

	function openUserModal(u: ManagementUser) {
		targetSearch = '';
		userForm = {
			id: u.id,
			name: u.name,
			authority: u.authority,
			systemRole: u.systemRole,
			operationTarget: u.operationTarget,
			targetIds: (u.targetIds || []).map(Number),
			departmentId: u.departmentId,
			roleId: u.roleId,
			suspended: u.suspended
		};
		authoritySelect = u.authority;
		systemRoleSelect = u.systemRole;
		userModalOpen = true;
	}

	function onAuthoritySelectChange() {
		if (!userForm) return;
		const newVal = authoritySelect;
		const oldVal = userForm.authority;
		if (newVal === oldVal) return;
		confirmDanger = false;
		confirmTitle = '権限の変更';
		confirmMessage = `権限を「${authorityLabel(newVal)}」に変更しますか？`;
		confirmAction = () => {
			if (userForm) userForm.authority = newVal;
		};
		confirmOnCancel = () => {
			authoritySelect = userForm?.authority ?? authoritySelect;
		};
		confirmOpen = true;
	}

	function onSystemRoleSelectChange() {
		if (!userForm) return;
		const newVal = systemRoleSelect;
		const oldVal = userForm.systemRole;
		if (newVal === oldVal) return;
		confirmDanger = false;
		confirmTitle = 'ロールの変更';
		confirmMessage = `ロールを「${systemRoleLabel(newVal)}」に変更しますか？`;
		confirmAction = () => {
			if (userForm) userForm.systemRole = newVal;
		};
		confirmOnCancel = () => {
			systemRoleSelect = userForm?.systemRole ?? systemRoleSelect;
		};
		confirmOpen = true;
	}

	function openCreateUserModal() {
		const deptId = options.departments?.[0]?.id || 0;
		const roleId = options.roles?.[0]?.id || 0;
		createUserForm = { name: '', email: '', departmentId: deptId, roleId, password: '' };
		createUserOpen = true;
	}

	async function submitCreateUser() {
		try {
			await api('POST', '/api/v1/management/users', createUserForm);
		} catch (e) {
			showToast(apiErrorMessage(e, '作成に失敗しました'), 'error');
			return;
		}
		createUserOpen = false;
		await loadUsers();
		showToast('ユーザを追加しました');
	}

	function openCreateDeptModal() {
		createDeptForm = { name: '', parentDepartmentId: 0 };
		createDeptOpen = true;
	}

	async function submitCreateDept() {
		try {
			await api('POST', '/api/v1/management/departments', createDeptForm);
		} catch (e) {
			showToast(apiErrorMessage(e, '作成に失敗しました'), 'error');
			return;
		}
		createDeptOpen = false;
		await Promise.all([loadOrgTree(), loadOptions()]);
		showToast('部署を追加しました');
	}

	async function saveUser() {
		if (!userForm) return;
		const body = {
			authority: userForm.authority,
			systemRole: userForm.systemRole,
			operationTarget: userForm.operationTarget,
			targetIds: (userForm.targetIds || []).map(Number),
			departmentId: userForm.departmentId,
			roleId: userForm.roleId,
			suspended: userForm.suspended
		};
		try {
			await api('PUT', `/api/v1/management/users/${userForm.id}`, body);
		} catch (e) {
			showToast(apiErrorMessage(e, '保存に失敗しました'), 'error');
			return;
		}
		userModalOpen = false;
		await loadUsers();
		showToast('ユーザ情報を保存しました');
	}

	async function deleteUserIcon() {
		if (!userForm) return;
		try {
			await api('DELETE', `/api/v1/management/users/${userForm.id}/icon`);
		} catch {
			showToast('アイコンの削除に失敗しました', 'error');
			return;
		}
		await loadUsers();
		showToast('アイコンを削除しました');
	}

	function toggleSuspend(unsuspend: boolean) {
		if (!userForm) return;
		userForm.suspended = !unsuspend;
	}

	function askSuspend() {
		confirmDanger = true;
		confirmOnCancel = null;
		confirmTitle = 'アカウント停止';
		confirmMessage = `選択した ${selectedUserIds.length} 件のアカウントを停止しますか？`;
		confirmAction = async () => {
			try {
				await api('POST', '/api/v1/management/users/suspend', { ids: selectedUserIds });
			} catch {
				showToast('停止に失敗しました', 'error');
				return;
			}
			selectedUserIds = [];
			await loadUsers();
			showToast('アカウントを停止しました');
		};
		confirmOpen = true;
	}

	async function openOrgModal(n: OrgTreeNode) {
		orgEditMode = false;
		orgChildSearch = '';
		orgTeamSearch = '';
		orgUserSearch = '';
		orgMemberSearch = '';
		const url =
			n.kind === 'team'
				? `/api/v1/management/teams/${n.id}`
				: `/api/v1/management/departments/${n.id}`;
		let detail: OrgDetail;
		try {
			detail = await api<OrgDetail>('GET', url);
		} catch {
			showToast('詳細の取得に失敗しました', 'error');
			return;
		}
		orgDetail = detail;
		orgEdit = {
			name: detail.name,
			childIds: (detail.childIds || []).map(Number),
			teamIds: (detail.teamIds || []).map(Number),
			userIds: (detail.userIds || []).map(Number),
			memberIds: (detail.memberIds || []).map(Number),
			leaderId: detail.leaderId || 0
		};
		orgModalOpen = true;
	}

	async function saveOrg() {
		if (!orgDetail) return;
		try {
			if (orgDetail.kind === 'team') {
				await api('PUT', `/api/v1/management/teams/${orgDetail.id}`, {
					name: orgEdit.name,
					memberIds: orgEdit.memberIds.map(Number),
					leaderId: orgEdit.leaderId
				});
			} else {
				await api('PUT', `/api/v1/management/departments/${orgDetail.id}`, {
					name: orgEdit.name,
					childDepartmentIds: orgEdit.childIds.map(Number),
					teamIds: orgEdit.teamIds.map(Number),
					userIds: orgEdit.userIds.map(Number)
				});
			}
		} catch (e) {
			showToast(apiErrorMessage(e, '保存に失敗しました'), 'error');
			return;
		}
		orgModalOpen = false;
		await loadOrgTree();
		showToast('保存しました');
	}

	function askDeleteOrg() {
		if (!orgDetail) return;
		confirmDanger = true;
		confirmOnCancel = null;
		confirmTitle = '削除の確認';
		confirmMessage = `${orgDetail.name} を削除しますか？`;
		confirmAction = async () => {
			if (!orgDetail) return;
			const url =
				orgDetail.kind === 'team'
					? `/api/v1/management/teams/${orgDetail.id}`
					: `/api/v1/management/departments/${orgDetail.id}`;
			try {
				await api('DELETE', url);
			} catch (e) {
				showToast(apiErrorMessage(e, '削除に失敗しました'), 'error');
				return;
			}
			orgModalOpen = false;
			await loadOrgTree();
			showToast('削除しました');
		};
		confirmOpen = true;
	}

	async function saveRemind() {
		try {
			await api('PUT', '/api/v1/management/settings/remind', remind);
		} catch {
			showToast('保存に失敗しました', 'error');
			return;
		}
		showToast('通知設定を保存しました');
	}

	function switchTab(name: 'users' | 'departments' | 'remind') {
		tab = name;
	}

	async function runConfirm() {
		confirmOpen = false;
		const action = confirmAction;
		confirmAction = null;
		confirmOnCancel = null;
		if (action) await action();
	}

	function cancelConfirm() {
		confirmOpen = false;
		if (confirmOnCancel) confirmOnCancel();
		confirmAction = null;
		confirmOnCancel = null;
	}

	function closeModals() {
		if (confirmOpen) {
			cancelConfirm();
			return;
		}
		userModalOpen = false;
		orgModalOpen = false;
	}

	function onWindowKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') closeModals();
	}

	onMount(() => {
		void Promise.all([loadUsers(), loadOptions(), loadOrgTree(), loadRemind()]);
	});
</script>

<svelte:window onkeydown={onWindowKeydown} />

<svelte:head><title>管理 - Growlink</title></svelte:head>

<h1 class="font-display text-2xl font-semibold">管理</h1>
<p class="mt-1 text-sm text-base-content/60">ユーザ・組織・通知の設定を行います</p>

<div role="tablist" class="tabs tabs-bordered mt-6 max-w-3xl">
	<button
		type="button"
		role="tab"
		class="tab gap-2"
		class:tab-active={tab === 'users'}
		onclick={() => switchTab('users')}
	>
		<Users class="h-4 w-4" />ユーザ一覧
	</button>
	<button
		type="button"
		role="tab"
		class="tab gap-2"
		class:tab-active={tab === 'departments'}
		onclick={() => switchTab('departments')}
	>
		<Building2 class="h-4 w-4" />部署一覧
	</button>
	<button
		type="button"
		role="tab"
		class="tab gap-2"
		class:tab-active={tab === 'remind'}
		onclick={() => switchTab('remind')}
	>
		<Bell class="h-4 w-4" />1on1通知
	</button>
</div>

{#if tab === 'users'}
	<div class="mt-6">
		<div class="flex flex-col sm:flex-row sm:items-center gap-3 mb-4">
			<label class="input input-bordered flex items-center gap-2 flex-1 max-w-md">
				<Search class="h-4 w-4 opacity-50" />
				<input type="search" class="grow" placeholder="氏名・メールで検索" bind:value={userSearch} />
			</label>
			<button type="button" class="btn btn-primary btn-sm shrink-0" onclick={openCreateUserModal}
				>ユーザを追加</button
			>
			{#if selectedUserIds.length > 0}
				<div class="flex items-center gap-2">
					<span class="text-sm text-base-content/60">{selectedUserIds.length} 件選択中</span>
					<button type="button" class="btn btn-error btn-sm" onclick={askSuspend}>アカウント停止</button>
				</div>
			{/if}
		</div>
		<div class="overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
			{#if usersLoading}
				<div class="p-6 space-y-3">
					{#each Array(5) as _, i (i)}
						<div class="skeleton h-10 w-full rounded-xl"></div>
					{/each}
				</div>
			{:else if filteredUsers.length === 0}
				<div class="p-12 text-center text-base-content/50 text-sm">該当するユーザがいません</div>
			{:else}
				<table class="table">
					<thead>
						<tr class="text-base-content/70">
							<th class="w-10">
								<input
									type="checkbox"
									class="checkbox checkbox-sm"
									checked={allUsersSelected}
									onchange={(e) =>
										toggleAllUsers((e.currentTarget as HTMLInputElement).checked)}
								/>
							</th>
							<th>ユーザ</th><th>権限</th><th>ロール</th><th>状態</th><th class="w-8"></th>
						</tr>
					</thead>
					<tbody>
						{#each filteredUsers as u (u.id)}
							<tr class="hover:bg-base-200/50 cursor-pointer" onclick={() => openUserModal(u)}>
								<!-- svelte-ignore a11y_no_static_element_interactions -->
								<td onclick={(e) => e.stopPropagation()}>
									<input
										type="checkbox"
										class="checkbox checkbox-sm"
										value={u.id}
										bind:group={selectedUserIds}
									/>
								</td>
								<td>
									<div class="flex items-center gap-3">
										{#if userAvatar(u)}
											<img
												src={userAvatar(u)!}
												alt=""
												class="h-9 w-9 rounded-full object-cover border border-base-300"
											/>
										{:else}
											<span
												class="grid h-9 w-9 place-items-center rounded-full bg-primary text-sm font-semibold text-primary-content"
												>{userInitial(u)}</span
											>
										{/if}
										<div>
											<div class="font-medium">{u.name}</div>
											<div class="text-xs text-base-content/50">{u.email}</div>
										</div>
									</div>
								</td>
								<td
									><span class="badge badge-outline badge-sm">{authorityLabel(u.authority)}</span></td
								>
								<td
									><span class="badge badge-ghost badge-sm">{systemRoleLabel(u.systemRole)}</span></td
								>
								<td>
									<span
										class="badge badge-sm"
										class:badge-error={u.suspended}
										class:badge-success={!u.suspended}>{u.suspended ? '停止' : '有効'}</span
									>
								</td>
								<td><ChevronRight class="h-4 w-4 opacity-40" /></td>
							</tr>
						{/each}
					</tbody>
				</table>
			{/if}
		</div>
	</div>
{/if}

{#if tab === 'departments'}
	<div class="mt-6">
		<div class="flex justify-end mb-4">
			<button type="button" class="btn btn-primary btn-sm" onclick={openCreateDeptModal}
				>部署を追加</button
			>
		</div>
		<div class="rounded-3xl border border-base-300 bg-base-100 shadow-sm overflow-hidden">
			{#if orgLoading}
				<div class="p-6 space-y-3">
					{#each Array(6) as _, i (i)}
						<div class="skeleton h-10 w-full rounded-xl"></div>
					{/each}
				</div>
			{:else if orgTree.length === 0}
				<div class="p-12 text-center text-base-content/50 text-sm">部署・チームが登録されていません</div>
			{:else}
				<div class="divide-y divide-base-200">
					{#each orgTree as n (n.kind + '-' + n.id)}
						<button
							type="button"
							class="w-full text-left px-4 py-3 hover:bg-base-200/60 flex items-center gap-2 transition-colors"
							style="padding-left: {1 + n.depth * 1.25}rem"
							onclick={() => openOrgModal(n)}
						>
							{#if n.kind === 'team'}
								<Users class="h-4 w-4 shrink-0 text-primary/80" />
							{:else}
								<Building2 class="h-4 w-4 shrink-0 text-primary/80" />
							{/if}
							<span class="font-medium">{n.name}</span>
							{#if n.kind === 'department'}
								<span class="badge badge-ghost badge-sm">{n.memberCount}人</span>
							{/if}
							{#if n.kind === 'team' && n.leaderName}
								<span class="badge badge-secondary badge-sm">SM</span>
								<span class="text-sm text-base-content/60">{n.leaderName}</span>
							{/if}
							<ChevronRight class="h-4 w-4 ml-auto opacity-40" />
						</button>
					{/each}
				</div>
			{/if}
		</div>
	</div>
{/if}

{#if tab === 'remind'}
	<div class="mt-6 max-w-xl">
		<div class="rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-5">
			<div role="alert" class="alert alert-info text-sm">
				<Info class="h-5 w-5 shrink-0" />
				<span>設定の保存のみ行います。実際の通知送信はバックエンドジョブで別途実装されます。</span>
			</div>
			<label class="label cursor-pointer justify-start gap-3 rounded-2xl border border-base-300 px-4 py-3">
				<input type="checkbox" class="toggle toggle-primary" bind:checked={remind.remindAutoSend} />
				<span class="label-text font-medium">1on1の通知を自動送信する</span>
			</label>
			<div>
				<span class="label-text font-medium">通知タイミング</span>
				<div class="flex flex-col gap-2 mt-2">
					<label class="label cursor-pointer justify-start gap-3 rounded-2xl border border-base-300 px-4 py-3">
						<input
							type="radio"
							name="remindMode"
							class="radio radio-primary"
							value="days_before"
							bind:group={remind.remindMode}
						/>
						<span class="label-text">実施予定日の N 日前</span>
					</label>
					<label class="label cursor-pointer justify-start gap-3 rounded-2xl border border-base-300 px-4 py-3">
						<input
							type="radio"
							name="remindMode"
							class="radio radio-primary"
							value="prev_weekday"
							bind:group={remind.remindMode}
						/>
						<span class="label-text">前の週の指定曜日</span>
					</label>
				</div>
			</div>
			{#if remind.remindMode === 'days_before'}
				<label class="form-control">
					<span class="label-text font-medium">何日前（1〜30）</span>
					<input
						type="number"
						min="1"
						max="30"
						class="input input-bordered"
						bind:value={remind.remindDaysBefore}
					/>
				</label>
			{/if}
			{#if remind.remindMode === 'prev_weekday'}
				<label class="form-control">
					<span class="label-text font-medium">曜日</span>
					<select class="select select-bordered" bind:value={remind.remindWeekday}>
						{#each WEEKDAYS as w (w.value)}
							<option value={w.value}>{w.label}</option>
						{/each}
					</select>
				</label>
			{/if}
			<label class="form-control">
				<span class="label-text font-medium">送信時刻</span>
				<input type="time" class="input input-bordered" bind:value={remind.remindTime} />
			</label>
			<button type="button" class="btn btn-primary w-full sm:w-auto" onclick={saveRemind}>保存</button>
		</div>
	</div>
{/if}

<dialog class="modal" class:modal-open={userModalOpen}>
	<div class="modal-box max-w-2xl max-h-[85vh] flex flex-col p-0">
		<div class="flex justify-between items-start gap-2 p-6 pb-2 border-b border-base-200 shrink-0">
			<div>
				<h3 class="font-bold text-lg">{userForm?.name}</h3>
				<p class="text-sm text-base-content/60 mt-0.5">権限・操作範囲・所属</p>
			</div>
			<button
				type="button"
				class="btn btn-sm btn-circle btn-ghost"
				aria-label="閉じる"
				onclick={() => (userModalOpen = false)}>×</button
			>
		</div>
		{#if userForm}
			<div class="overflow-y-auto px-6 py-4 space-y-5 flex-1">
				<div class="grid sm:grid-cols-2 gap-4">
					<label class="form-control">
						<span class="label-text font-medium">権限</span>
						<p class="text-xs text-base-content/50 mt-0.5">
							現在: <span>{authorityLabel(userForm.authority)}</span>
						</p>
						<select
							class="select select-bordered select-sm mt-1"
							bind:value={authoritySelect}
							onchange={onAuthoritySelectChange}
						>
							{#each AUTHORITY_OPTIONS as opt (opt.value)}
								<option value={opt.value}>{opt.label}</option>
							{/each}
						</select>
					</label>
					<label class="form-control">
						<span class="label-text font-medium">ロール</span>
						<p class="text-xs text-base-content/50 mt-0.5">
							現在: <span>{systemRoleLabel(userForm.systemRole)}</span>
						</p>
						<select
							class="select select-bordered select-sm mt-1"
							bind:value={systemRoleSelect}
							onchange={onSystemRoleSelectChange}
						>
							{#each SYSTEM_ROLE_OPTIONS as opt (opt.value)}
								<option value={opt.value}>{opt.label}</option>
							{/each}
						</select>
					</label>
				</div>
				<div>
					<span class="label-text font-medium">操作対象</span>
					<div class="flex flex-wrap gap-2 mt-2">
						{#each OPERATION_OPTIONS as opt (opt.value)}
							<button
								type="button"
								class="btn btn-sm"
								class:btn-primary={userForm.operationTarget === opt.value}
								class:btn-outline={userForm.operationTarget !== opt.value}
								onclick={() => setOperationTarget(opt.value)}>{opt.label}</button
							>
						{/each}
					</div>
				</div>
				{#if userForm.operationTarget === 'department'}
					<div class="space-y-2">
						<input
							class="input input-bordered input-sm w-full"
							placeholder="部署を検索"
							bind:value={targetSearch}
						/>
						<div class="border border-base-300 rounded-2xl max-h-44 overflow-y-auto p-2 space-y-0.5">
							{#each filteredDepartments() as d (d.id)}
								<label class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer">
									<input
										type="checkbox"
										class="checkbox checkbox-sm"
										checked={isInList(userForm.targetIds, d.id)}
										onchange={() => toggleUserTarget(d.id)}
									/>
									<span>{d.name}</span>
								</label>
							{/each}
						</div>
					</div>
				{/if}
				{#if userForm.operationTarget === 'team'}
					<div class="space-y-2">
						<input
							class="input input-bordered input-sm w-full"
							placeholder="チームを検索"
							bind:value={targetSearch}
						/>
						<div class="border border-base-300 rounded-2xl max-h-44 overflow-y-auto p-2 space-y-0.5">
							{#each filteredTeams() as t (t.id)}
								<label class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer">
									<input
										type="checkbox"
										class="checkbox checkbox-sm"
										checked={isInList(userForm.targetIds, t.id)}
										onchange={() => toggleUserTarget(t.id)}
									/>
									<span>{t.name}</span>
								</label>
							{/each}
						</div>
					</div>
				{/if}
				{#if userForm.operationTarget === 'subordinate'}
					<div class="space-y-2">
						<input
							class="input input-bordered input-sm w-full"
							placeholder="部下を検索"
							bind:value={targetSearch}
						/>
						<div class="border border-base-300 rounded-2xl max-h-44 overflow-y-auto p-2 space-y-0.5">
							{#each filteredUsersForTarget() as u (`sub-${u.id}`)}
								{#if u.id !== userForm.id}
									<label
										class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer"
									>
										<input
											type="checkbox"
											class="checkbox checkbox-sm"
											checked={isInList(userForm.targetIds, u.id)}
											onchange={() => toggleUserTarget(u.id)}
										/>
										<span>{u.name}</span>
									</label>
								{/if}
							{/each}
						</div>
					</div>
				{/if}
				<div class="grid sm:grid-cols-2 gap-4">
					<label class="form-control"
						><span class="label-text font-medium">所属部署</span>
						<select class="select select-bordered" bind:value={userForm.departmentId}>
							{#each options.departments as d (`dept-${d.id}`)}
								<option value={d.id}>{d.name}</option>
							{/each}
						</select>
					</label>
					<label class="form-control"
						><span class="label-text font-medium">役職（マスタ）</span>
						<select class="select select-bordered" bind:value={userForm.roleId}>
							{#each options.roles as r (r.id)}
								<option value={r.id}>{r.name}</option>
							{/each}
						</select>
					</label>
				</div>
				<div class="border-t border-base-300 pt-4 flex flex-wrap gap-2">
					<button type="button" class="btn btn-outline btn-sm" onclick={deleteUserIcon}>アイコン削除</button>
					{#if !userForm.suspended}
						<button type="button" class="btn btn-warning btn-sm" onclick={() => toggleSuspend(false)}
							>アカウント停止</button
						>
					{/if}
					{#if userForm.suspended}
						<button type="button" class="btn btn-success btn-sm" onclick={() => toggleSuspend(true)}
							>停止解除</button
						>
					{/if}
				</div>
			</div>
		{/if}
		<div class="modal-action border-t border-base-200 p-4 shrink-0">
			<button type="button" class="btn btn-ghost" onclick={() => (userModalOpen = false)}>閉じる</button>
			<button type="button" class="btn btn-primary" onclick={saveUser}>保存</button>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop"
		><button type="button" onclick={() => (userModalOpen = false)}>close</button></form
	>
</dialog>

<dialog class="modal" class:modal-open={orgModalOpen}>
	<div class="modal-box max-w-2xl max-h-[85vh] flex flex-col p-0">
		<div class="flex justify-between items-start gap-2 p-6 pb-2 border-b border-base-200 shrink-0">
			<div>
				<h3 class="font-bold text-lg">{orgDetail?.name}</h3>
				<p class="text-sm text-base-content/60">{orgDetail?.kind === 'team' ? 'チーム' : '部署'}</p>
			</div>
			<div class="flex gap-1">
				{#if !orgEditMode}
					<button type="button" class="btn btn-sm btn-outline" onclick={() => (orgEditMode = true)}
						>編集</button
					>
				{/if}
				<button
					type="button"
					class="btn btn-sm btn-circle btn-ghost"
					aria-label="閉じる"
					onclick={() => (orgModalOpen = false)}>×</button
				>
			</div>
		</div>
		{#if !orgEditMode}
			<div class="overflow-y-auto px-6 py-4 flex-1">
				{#if (orgDetail?.members || []).length}
					<ul class="space-y-2">
						{#each orgDetail?.members || [] as m (m.userId)}
							<li class="flex justify-between items-center rounded-xl border border-base-300 px-3 py-2">
								<span class="font-medium">{m.userName}</span>
								<span class="badge badge-ghost badge-sm">{m.title}</span>
							</li>
						{/each}
					</ul>
				{:else}
					<p class="text-sm text-base-content/50 text-center py-8">所属メンバーはいません</p>
				{/if}
			</div>
		{/if}
		{#if orgEditMode}
			<div class="overflow-y-auto px-6 py-4 space-y-4 flex-1">
				<label class="form-control"
					><span class="label-text font-medium">名前</span><input
						class="input input-bordered"
						bind:value={orgEdit.name}
					/></label
				>
				{#if orgDetail?.kind === 'department' && orgDetail?.isRoot}
					<div class="space-y-2">
						<span class="label-text font-medium">子部署</span>
						<input
							class="input input-bordered input-sm w-full"
							placeholder="検索"
							bind:value={orgChildSearch}
						/>
						<div class="border border-base-300 rounded-2xl max-h-44 overflow-y-auto p-2">
							{#each filteredOrgDepartments() as d (`c-${d.id}`)}
								{#if d.id !== orgDetail?.id}
									<label
										class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer"
									>
										<input
											type="checkbox"
											class="checkbox checkbox-sm"
											checked={isInList(orgEdit.childIds, d.id)}
											onchange={() => toggleOrgChild(d.id)}
										/>
										<span>{d.name}</span>
									</label>
								{/if}
							{/each}
						</div>
					</div>
				{/if}
				{#if orgDetail?.kind === 'department' && !orgDetail?.isRoot}
					<div class="space-y-4">
						<div class="space-y-2">
							<span class="label-text font-medium">チーム</span>
							<input
								class="input input-bordered input-sm w-full"
								placeholder="検索"
								bind:value={orgTeamSearch}
							/>
							<div class="border border-base-300 rounded-2xl max-h-36 overflow-y-auto p-2">
								{#each filteredOrgTeams() as t (`t-${t.id}`)}
									<label
										class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer"
									>
										<input
											type="checkbox"
											class="checkbox checkbox-sm"
											checked={isInList(orgEdit.teamIds, t.id)}
											onchange={() => toggleOrgTeam(t.id)}
										/>
										<span>{t.name}</span>
									</label>
								{/each}
							</div>
						</div>
						<div class="space-y-2">
							<span class="label-text font-medium">所属ユーザ</span>
							<input
								class="input input-bordered input-sm w-full"
								placeholder="検索"
								bind:value={orgUserSearch}
							/>
							<div class="border border-base-300 rounded-2xl max-h-36 overflow-y-auto p-2">
								{#each filteredOrgUsers() as u (`ou-${u.id}`)}
									<label
										class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer"
									>
										<input
											type="checkbox"
											class="checkbox checkbox-sm"
											checked={isInList(orgEdit.userIds, u.id)}
											onchange={() => toggleOrgUser(u.id)}
										/>
										<span>{u.name}</span>
									</label>
								{/each}
							</div>
						</div>
					</div>
				{/if}
				{#if orgDetail?.kind === 'team'}
					<div class="space-y-4">
						<div class="space-y-2">
							<span class="label-text font-medium">メンバー</span>
							<input
								class="input input-bordered input-sm w-full"
								placeholder="検索"
								bind:value={orgMemberSearch}
							/>
							<div class="border border-base-300 rounded-2xl max-h-36 overflow-y-auto p-2">
								{#each filteredOrgMembers() as u (`tm-${u.id}`)}
									<label
										class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer"
									>
										<input
											type="checkbox"
											class="checkbox checkbox-sm"
											checked={isInList(orgEdit.memberIds, u.id)}
											onchange={() => toggleOrgMember(u.id)}
										/>
										<span>{u.name}</span>
									</label>
								{/each}
							</div>
						</div>
						{#if orgEdit.memberIds.length}
							<div>
								<span class="label-text font-medium">リーダー（SM）</span>
								<div class="flex flex-wrap gap-2 mt-2">
									{#each orgEdit.memberIds as uid (`ld-${uid}`)}
										<button
											type="button"
											class="btn btn-sm"
											class:btn-secondary={Number(orgEdit.leaderId) === Number(uid)}
											class:btn-outline={Number(orgEdit.leaderId) !== Number(uid)}
											onclick={() => (orgEdit.leaderId = Number(uid))}>{userName(uid)}</button
										>
									{/each}
								</div>
							</div>
						{/if}
					</div>
				{/if}
				<div class="border-t border-base-300 pt-4">
					<button type="button" class="btn btn-error btn-outline btn-sm" onclick={askDeleteOrg}
						>削除</button
					>
				</div>
			</div>
		{/if}
		{#if !orgEditMode}
			<div class="modal-action border-t border-base-200 p-4 shrink-0">
				<button type="button" class="btn btn-error btn-outline btn-sm mr-auto" onclick={askDeleteOrg}
					>削除</button
				>
				<button type="button" class="btn btn-ghost" onclick={() => (orgModalOpen = false)}>閉じる</button>
			</div>
		{/if}
		{#if orgEditMode}
			<div class="modal-action border-t border-base-200 p-4 shrink-0">
				<button type="button" class="btn btn-ghost" onclick={() => (orgModalOpen = false)}>閉じる</button>
				<button type="button" class="btn btn-primary" onclick={saveOrg}>保存</button>
			</div>
		{/if}
	</div>
	<form method="dialog" class="modal-backdrop"
		><button type="button" onclick={() => (orgModalOpen = false)}>close</button></form
	>
</dialog>

<dialog class="modal" class:modal-open={createUserOpen}>
	<div class="modal-box max-w-lg max-h-[85vh] flex flex-col p-0">
		<div class="p-6 pb-2 border-b border-base-200 shrink-0">
			<h3 class="font-bold text-lg">ユーザを追加</h3>
			<p class="text-sm text-base-content/60 mt-0.5">権限はメンバー、ロールはなしで作成されます</p>
		</div>
		<div class="overflow-y-auto px-6 py-4 space-y-4 flex-1">
			<label class="form-control"
				><span class="label-text font-medium">氏名</span><input
					class="input input-bordered"
					bind:value={createUserForm.name}
				/></label
			>
			<label class="form-control"
				><span class="label-text font-medium">メール</span><input
					type="email"
					class="input input-bordered"
					bind:value={createUserForm.email}
				/></label
			>
			<label class="form-control"
				><span class="label-text font-medium">所属部署</span>
				<select class="select select-bordered" bind:value={createUserForm.departmentId}>
					<option value={0}>選択してください</option>
					{#each options.departments as d (`new-dept-${d.id}`)}
						<option value={d.id}>{d.name}</option>
					{/each}
				</select>
			</label>
			<label class="form-control"
				><span class="label-text font-medium">役職（マスタ）</span>
				<select class="select select-bordered" bind:value={createUserForm.roleId}>
					<option value={0}>選択してください</option>
					{#each options.roles as r (`new-role-${r.id}`)}
						<option value={r.id}>{r.name}</option>
					{/each}
				</select>
			</label>
			<label class="form-control"
				><span class="label-text font-medium">パスワード（任意・Basic認証用）</span><input
					type="password"
					class="input input-bordered"
					bind:value={createUserForm.password}
					autocomplete="new-password"
				/></label
			>
		</div>
		<div class="modal-action border-t border-base-200 p-4 shrink-0">
			<button type="button" class="btn btn-ghost" onclick={() => (createUserOpen = false)}>キャンセル</button>
			<button type="button" class="btn btn-primary" onclick={submitCreateUser}>作成</button>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop"
		><button type="button" onclick={() => (createUserOpen = false)}>close</button></form
	>
</dialog>

<dialog class="modal" class:modal-open={createDeptOpen}>
	<div class="modal-box max-w-md max-h-[85vh] flex flex-col p-0">
		<div class="p-6 pb-2 border-b border-base-200 shrink-0">
			<h3 class="font-bold text-lg">部署を追加</h3>
		</div>
		<div class="overflow-y-auto px-6 py-4 space-y-4 flex-1">
			<label class="form-control"
				><span class="label-text font-medium">部署名</span><input
					class="input input-bordered"
					bind:value={createDeptForm.name}
				/></label
			>
			<label class="form-control"
				><span class="label-text font-medium">親部署（任意）</span>
				<select class="select select-bordered" bind:value={createDeptForm.parentDepartmentId}>
					<option value={0}>なし（ルート）</option>
					{#each options.departments as d (`parent-${d.id}`)}
						<option value={d.id}>{d.name}</option>
					{/each}
				</select>
			</label>
		</div>
		<div class="modal-action border-t border-base-200 p-4 shrink-0">
			<button type="button" class="btn btn-ghost" onclick={() => (createDeptOpen = false)}>キャンセル</button>
			<button type="button" class="btn btn-primary" onclick={submitCreateDept}>作成</button>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop"
		><button type="button" onclick={() => (createDeptOpen = false)}>close</button></form
	>
</dialog>

<dialog class="modal" class:modal-open={confirmOpen}>
	<div class="modal-box max-h-[85vh]">
		<h3 class="font-bold">{confirmTitle}</h3>
		<p class="py-4 text-sm text-base-content/80">{confirmMessage}</p>
		<div class="modal-action">
			<button type="button" class="btn btn-ghost" onclick={cancelConfirm}>キャンセル</button>
			<button
				type="button"
				class="btn"
				class:btn-error={confirmDanger}
				class:btn-primary={!confirmDanger}
				onclick={runConfirm}>{confirmDanger ? '確定' : '変更する'}</button
			>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop"
		><button type="button" onclick={cancelConfirm}>close</button></form
	>
</dialog>

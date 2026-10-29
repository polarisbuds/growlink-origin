<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import Users from "lucide-svelte/icons/users";
    import Building2 from "lucide-svelte/icons/building-2";
    import Bell from "lucide-svelte/icons/bell";
    import Search from "lucide-svelte/icons/search";
    import ChevronRight from "lucide-svelte/icons/chevron-right";
    import Info from "lucide-svelte/icons/info";

    const AUTHORITY_OPTIONS = [
        { value: "general", label: "メンバー" },
        { value: "boss", label: "ボス" },
        { value: "owner", label: "オーナー" },
    ];
    const SYSTEM_ROLE_OPTIONS = [
        { value: "none", label: "なし" },
        { value: "system_admin", label: "システム管理者" },
        { value: "analyst", label: "分析者" },
    ];
    const OPERATION_OPTIONS = [
        { value: "none", label: "なし" },
        { value: "department", label: "部署" },
        { value: "team", label: "チーム" },
        { value: "subordinate", label: "部下" },
    ];
    const WEEKDAYS = [
        { value: 1, label: "月曜日" },
        { value: 2, label: "火曜日" },
        { value: 3, label: "水曜日" },
        { value: 4, label: "木曜日" },
        { value: 5, label: "金曜日" },
        { value: 6, label: "土曜日" },
        { value: 7, label: "日曜日" },
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
        kind: "team" | "department";
        id: number;
        name: string;
        depth: number;
        departmentId?: number;
        memberCount?: number;
        teamMemberCount?: number;
        leaderName?: string;
        headUserName?: string;
        teamCount?: number;
        departmentCount?: number;
    };

    type OrgMember = { userId: number; userName: string; title: string };

    type OrgDetail = {
        id: number;
        name: string;
        kind: "team" | "department";
        isRoot?: boolean;
        members?: OrgMember[];
        childIds?: number[];
        teamIds?: number[];
        userIds?: number[];
        memberIds?: number[];
        leaderId?: number;
        headUserId?: number;
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
        roles: { id: number; name: string; rank?: number }[];
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
        labelKey = "name",
    ): T[] {
        const q = (query || "").trim().toLowerCase();
        if (!q) return items;
        return items.filter((item) =>
            String(item[labelKey] || "")
                .toLowerCase()
                .includes(q),
        );
    }

    function apiErrorMessage(e: unknown, fallback: string): string {
        if (e instanceof Error) {
            try {
                const j = JSON.parse(e.message) as { error?: string };
                if (j?.error) return j.error;
            } catch {
                /* not json */
            }
            if (e.message && e.message !== "unauthorized") return e.message;
        }
        return fallback;
    }

    let tab = $state<"users" | "departments" | "remind">("users");
    let users = $state<ManagementUser[]>([]);
    let orgTree = $state<OrgTreeNode[]>([]);
    let options = $state<Options>({ departments: [], teams: [], roles: [] });
    let selectedUserIds = $state<number[]>([]);
    let userSearch = $state("");
    let usersLoading = $state(true);
    let orgLoading = $state(true);
    let remind = $state({
        remindAutoSend: false,
        remindMode: "prev_weekday",
        remindDaysBefore: 1,
        remindWeekday: 1,
        remindTime: "09:00",
    });
    let targetSearch = $state("");
    let orgChildSearch = $state("");
    let orgTeamSearch = $state("");
    let orgUserSearch = $state("");
    let orgMemberSearch = $state("");
    let userModalOpen = $state(false);
    let userForm = $state<UserForm | null>(null);
    let authoritySelect = $state("general");
    let systemRoleSelect = $state("none");
    let createUserOpen = $state(false);
    let createUserForm = $state({
        name: "",
        email: "",
        departmentId: 0,
        roleId: 0,
        password: "",
    });
    let createDeptOpen = $state(false);
    let createDeptForm = $state({ name: "", parentDepartmentId: 0 });
    let createTeamOpen = $state(false);
    let createTeamForm = $state({ name: "", departmentId: 0 });
    let createRoleOpen = $state(false);
    let createRoleForm = $state({ name: "", rank: 50 });
    let expandedDeptIds = $state<Set<number>>(new Set());
    let orgModalOpen = $state(false);
    let orgDetail = $state<OrgDetail | null>(null);
    let orgEditMode = $state(false);
    let orgEdit = $state({
        name: "",
        childIds: [] as number[],
        teamIds: [] as number[],
        userIds: [] as number[],
        memberIds: [] as number[],
        leaderId: 0,
        headUserId: 0,
    });

    let headRankThreshold = $derived(
        (options.roles || []).find((r) => r.name === "部長")?.rank ?? 100,
    );

    let headUserCandidates = $derived(
        users.filter((u) => {
            const role = (options.roles || []).find((r) => r.id === u.roleId);
            return (role?.rank ?? 0) >= headRankThreshold;
        }),
    );

    let visibleOrgNodes = $derived.by(() => {
        const out: OrgTreeNode[] = [];
        let hiddenBelow: number | null = null;

        for (const n of orgTree) {
            if (hiddenBelow != null && n.depth > hiddenBelow) continue;
            hiddenBelow = null;

            out.push(n);
            if (n.kind === "department" && !expandedDeptIds.has(n.id)) {
                hiddenBelow = n.depth;
            }
        }
        return out;
    });

    function toggleDeptExpand(id: number) {
        const next = new Set(expandedDeptIds);
        if (next.has(id)) next.delete(id);
        else next.add(id);
        expandedDeptIds = next;
    }
    let confirmOpen = $state(false);
    let confirmTitle = $state("");
    let confirmMessage = $state("");
    let confirmAction = $state<(() => void | Promise<void>) | null>(null);
    let confirmOnCancel = $state<(() => void) | null>(null);
    let confirmDanger = $state(false);

    let filteredUsers = $derived.by(() => {
        const q = userSearch.trim().toLowerCase();
        if (!q) return users;
        return users.filter(
            (u) =>
                (u.name || "").toLowerCase().includes(q) ||
                (u.email || "").toLowerCase().includes(q),
        );
    });

    let allUsersSelected = $derived(
        filteredUsers.length > 0 &&
            filteredUsers.every((u) => selectedUserIds.includes(u.id)),
    );

    function authorityLabel(v: string) {
        return AUTHORITY_OPTIONS.find((o) => o.value === v)?.label || v;
    }

    function systemRoleLabel(v: string) {
        return SYSTEM_ROLE_OPTIONS.find((o) => o.value === v)?.label || v;
    }

    function userName(id: number | string) {
        const u = users.find((x) => x.id === Number(id));
        return u ? u.name : "";
    }

    function userAvatar(u: ManagementUser) {
        return u.imageUri ? `/uploads/${u.imageUri}` : null;
    }

    function userInitial(u: ManagementUser) {
        return (u.name || "?").slice(0, 1);
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
            users = await api<ManagementUser[]>(
                "GET",
                "/api/v1/management/users",
            );
        } catch {
            users = [];
        }
        usersLoading = false;
    }

    async function loadOptions() {
        try {
            options = await api<Options>("GET", "/api/v1/management/options");
        } catch {
            /* keep defaults */
        }
    }

    async function loadOrgTree() {
        orgLoading = true;
        try {
            orgTree = await api<OrgTreeNode[]>(
                "GET",
                "/api/v1/management/org/tree",
            );
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
            }>("GET", "/api/v1/management/settings/remind");
            remind = {
                remindAutoSend: !!data.remindAutoSend,
                remindMode: data.remindMode || "prev_weekday",
                remindDaysBefore:
                    data.remindDaysBefore && data.remindDaysBefore > 0
                        ? data.remindDaysBefore
                        : 1,
                remindWeekday: data.remindWeekday || 1,
                remindTime: (data.remindTime || "09:00").slice(0, 5),
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
        targetSearch = "";
        userForm = {
            id: u.id,
            name: u.name,
            authority: u.authority,
            systemRole: u.systemRole,
            operationTarget: u.operationTarget,
            targetIds: (u.targetIds || []).map(Number),
            departmentId: u.departmentId,
            roleId: u.roleId,
            suspended: u.suspended,
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
        confirmTitle = "権限の変更";
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
        confirmTitle = "ロールの変更";
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
        createUserForm = {
            name: "",
            email: "",
            departmentId: deptId,
            roleId,
            password: "",
        };
        createUserOpen = true;
    }

    async function submitCreateUser() {
        try {
            await api("POST", "/api/v1/management/users", createUserForm);
        } catch (e) {
            showToast(apiErrorMessage(e, "作成に失敗しました"), "error");
            return;
        }
        createUserOpen = false;
        await loadUsers();
        showToast("ユーザを追加しました");
    }

    function openCreateDeptModal() {
        createDeptForm = { name: "", parentDepartmentId: 0 };
        createDeptOpen = true;
    }

    async function submitCreateDept() {
        try {
            await api("POST", "/api/v1/management/departments", createDeptForm);
        } catch (e) {
            showToast(apiErrorMessage(e, "作成に失敗しました"), "error");
            return;
        }
        createDeptOpen = false;
        await Promise.all([loadOrgTree(), loadOptions()]);
        showToast("部署を追加しました");
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
            suspended: userForm.suspended,
        };
        try {
            await api("PUT", `/api/v1/management/users/${userForm.id}`, body);
        } catch (e) {
            showToast(apiErrorMessage(e, "保存に失敗しました"), "error");
            return;
        }
        userModalOpen = false;
        await loadUsers();
        showToast("ユーザ情報を保存しました");
    }

    async function deleteUserIcon() {
        if (!userForm) return;
        try {
            await api("DELETE", `/api/v1/management/users/${userForm.id}/icon`);
        } catch {
            showToast("アイコンの削除に失敗しました", "error");
            return;
        }
        await loadUsers();
        showToast("アイコンを削除しました");
    }

    async function setSuspendedForUsers(ids: number[], suspended: boolean) {
        try {
            await api("POST", "/api/v1/management/users/suspend", {
                ids,
                suspended,
            });
        } catch {
            showToast(
                suspended ? "停止に失敗しました" : "停止解除に失敗しました",
                "error",
            );
            return false;
        }
        return true;
    }

    async function toggleSuspendInModal(unsuspend: boolean) {
        if (!userForm) return;
        const ok = await setSuspendedForUsers([userForm.id], !unsuspend);
        if (!ok) return;
        userForm.suspended = !unsuspend;
        await loadUsers();
        showToast(
            unsuspend
                ? "アカウントの停止を解除しました"
                : "アカウントを停止しました",
        );
    }

    function askSuspend() {
        confirmDanger = true;
        confirmOnCancel = null;
        confirmTitle = "アカウント停止";
        confirmMessage = `選択した ${selectedUserIds.length} 件のアカウントを停止しますか？`;
        confirmAction = async () => {
            const ok = await setSuspendedForUsers(selectedUserIds, true);
            if (!ok) return;
            selectedUserIds = [];
            await loadUsers();
            showToast("アカウントを停止しました");
        };
        confirmOpen = true;
    }

    function askUnsuspend() {
        confirmDanger = false;
        confirmOnCancel = null;
        confirmTitle = "停止解除";
        confirmMessage = `選択した ${selectedUserIds.length} 件のアカウントの停止を解除しますか？`;
        confirmAction = async () => {
            const ok = await setSuspendedForUsers(selectedUserIds, false);
            if (!ok) return;
            selectedUserIds = [];
            await loadUsers();
            showToast("アカウントの停止を解除しました");
        };
        confirmOpen = true;
    }

    async function openOrgModal(n: OrgTreeNode) {
        orgEditMode = false;
        orgChildSearch = "";
        orgTeamSearch = "";
        orgUserSearch = "";
        orgMemberSearch = "";
        const url =
            n.kind === "team"
                ? `/api/v1/management/teams/${n.id}`
                : `/api/v1/management/departments/${n.id}`;
        let detail: OrgDetail;
        try {
            detail = await api<OrgDetail>("GET", url);
        } catch {
            showToast("詳細の取得に失敗しました", "error");
            return;
        }
        orgDetail = detail;
        orgEdit = {
            name: detail.name,
            childIds: (detail.childIds || []).map(Number),
            teamIds: (detail.teamIds || []).map(Number),
            userIds: (detail.userIds || []).map(Number),
            memberIds: (detail.memberIds || []).map(Number),
            leaderId: detail.leaderId || 0,
            headUserId: detail.headUserId || 0,
        };
        orgModalOpen = true;
    }

    async function saveCreateTeam() {
        try {
            await api("POST", "/api/v1/management/teams", createTeamForm);
        } catch (e) {
            showToast(
                apiErrorMessage(e, "チームの作成に失敗しました"),
                "error",
            );
            return;
        }
        createTeamOpen = false;
        await Promise.all([loadOrgTree(), loadOptions()]);
        showToast("チームを追加しました");
    }

    async function saveCreateRole() {
        try {
            await api("POST", "/api/v1/management/roles", createRoleForm);
        } catch (e) {
            showToast(apiErrorMessage(e, "役職の登録に失敗しました"), "error");
            return;
        }
        createRoleOpen = false;
        await loadOptions();
        showToast("役職を登録しました");
    }

    function openCreateTeamModal(departmentId: number) {
        createTeamForm = { name: "", departmentId };
        createTeamOpen = true;
    }

    async function saveOrg() {
        if (!orgDetail) return;
        try {
            if (orgDetail.kind === "team") {
                await api("PUT", `/api/v1/management/teams/${orgDetail.id}`, {
                    name: orgEdit.name,
                    memberIds: orgEdit.memberIds.map(Number),
                    leaderId: orgEdit.leaderId,
                });
            } else {
                await api(
                    "PUT",
                    `/api/v1/management/departments/${orgDetail.id}`,
                    {
                        name: orgEdit.name,
                        childDepartmentIds: orgEdit.childIds.map(Number),
                        teamIds: orgEdit.teamIds.map(Number),
                        userIds: orgEdit.userIds.map(Number),
                        headUserId: orgEdit.headUserId,
                    },
                );
            }
        } catch (e) {
            showToast(apiErrorMessage(e, "保存に失敗しました"), "error");
            return;
        }
        orgModalOpen = false;
        await loadOrgTree();
        showToast("保存しました");
    }

    function askDeleteOrg() {
        if (!orgDetail) return;
        confirmDanger = true;
        confirmOnCancel = null;
        confirmTitle = "削除の確認";
        confirmMessage = `${orgDetail.name} を削除しますか？`;
        confirmAction = async () => {
            if (!orgDetail) return;
            const url =
                orgDetail.kind === "team"
                    ? `/api/v1/management/teams/${orgDetail.id}`
                    : `/api/v1/management/departments/${orgDetail.id}`;
            try {
                await api("DELETE", url);
            } catch (e) {
                showToast(apiErrorMessage(e, "削除に失敗しました"), "error");
                return;
            }
            orgModalOpen = false;
            await loadOrgTree();
            showToast("削除しました");
        };
        confirmOpen = true;
    }

    async function saveRemind() {
        try {
            await api("PUT", "/api/v1/management/settings/remind", remind);
        } catch {
            showToast("保存に失敗しました", "error");
            return;
        }
        showToast("通知設定を保存しました");
    }

    function switchTab(name: "users" | "departments" | "remind") {
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
        if (e.key === "Escape") closeModals();
    }

    onMount(() => {
        void Promise.all([
            loadUsers(),
            loadOptions(),
            loadOrgTree(),
            loadRemind(),
        ]);
    });
</script>

<svelte:window onkeydown={onWindowKeydown} />

<svelte:head><title>管理 - Growlink</title></svelte:head>

<h1 class="font-display text-2xl font-semibold">管理</h1>
<p class="mt-1 text-sm text-base-content/60">
    ユーザ・組織・通知の設定を行います
</p>

<div role="tablist" class="tabs tabs-bordered mt-6 max-w-3xl">
    <button
        type="button"
        role="tab"
        class="tab gap-2"
        class:tab-active={tab === "users"}
        onclick={() => switchTab("users")}>
        <Users class="h-4 w-4" />ユーザ一覧
    </button>
    <button
        type="button"
        role="tab"
        class="tab gap-2"
        class:tab-active={tab === "departments"}
        onclick={() => switchTab("departments")}>
        <Building2 class="h-4 w-4" />組織一覧
    </button>
    <button
        type="button"
        role="tab"
        class="tab gap-2"
        class:tab-active={tab === "remind"}
        onclick={() => switchTab("remind")}>
        <Bell class="h-4 w-4" />1on1通知
    </button>
</div>

{#if tab === "users"}
    <div class="mt-6">
        <div class="flex flex-col sm:flex-row sm:items-center gap-3 mb-4">
            <label
                class="input input-bordered flex items-center gap-2 flex-1 max-w-md">
                <Search class="h-4 w-4 opacity-50" />
                <input
                    type="search"
                    class="grow placeholder-current"
                    placeholder="氏名・メールで検索"
                    bind:value={userSearch} />
            </label>
            <button
                type="button"
                class="btn btn-primary btn-sm shrink-0"
                onclick={openCreateUserModal}>ユーザを追加</button>
            <button
                type="button"
                class="btn btn-outline btn-sm shrink-0"
                onclick={() => (createRoleOpen = true)}>役職を登録</button>
            {#if selectedUserIds.length > 0}
                <div class="flex items-center gap-2">
                    <span class="text-sm text-base-content/60"
                        >{selectedUserIds.length} 件選択中</span>
                    <button
                        type="button"
                        class="btn btn-error btn-sm"
                        onclick={askSuspend}>アカウント停止</button>
                    <button
                        type="button"
                        class="btn btn-success btn-sm"
                        onclick={askUnsuspend}>停止解除</button>
                </div>
            {/if}
        </div>
        <div
            class="overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
            {#if usersLoading}
                <div class="p-6 space-y-3">
                    {#each Array(5) as _, i (i)}
                        <div class="skeleton h-10 w-full rounded-xl"></div>
                    {/each}
                </div>
            {:else if filteredUsers.length === 0}
                <div class="p-12 text-center text-base-content/50 text-sm">
                    該当するユーザがいません
                </div>
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
                                        toggleAllUsers(
                                            (
                                                e.currentTarget as HTMLInputElement
                                            ).checked,
                                        )} />
                            </th>
                            <th>ユーザ</th><th>権限</th><th>ロール</th><th
                                >状態</th
                            ><th class="w-8"></th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each filteredUsers as u (u.id)}
                            <tr
                                class="hover:bg-base-200/50 cursor-pointer"
                                onclick={() => openUserModal(u)}>
                                <!-- svelte-ignore a11y_no_static_element_interactions -->
                                <td onclick={(e) => e.stopPropagation()}>
                                    <input
                                        type="checkbox"
                                        class="checkbox checkbox-sm"
                                        value={u.id}
                                        bind:group={selectedUserIds} />
                                </td>
                                <td>
                                    <div class="flex items-center gap-3">
                                        {#if userAvatar(u)}
                                            <img
                                                src={userAvatar(u)!}
                                                alt=""
                                                class="h-9 w-9 rounded-full object-cover border border-base-300" />
                                        {:else}
                                            <span
                                                class="grid h-9 w-9 place-items-center rounded-full bg-primary text-sm font-semibold text-primary-content"
                                                >{userInitial(u)}</span>
                                        {/if}
                                        <div>
                                            <div class="font-medium">
                                                {u.name}
                                            </div>
                                            <div
                                                class="text-xs text-base-content/50">
                                                {u.email}
                                            </div>
                                        </div>
                                    </div>
                                </td>
                                <td
                                    ><span class="badge badge-outline badge-sm"
                                        >{authorityLabel(u.authority)}</span
                                    ></td>
                                <td
                                    ><span class="badge badge-ghost badge-sm"
                                        >{systemRoleLabel(u.systemRole)}</span
                                    ></td>
                                <td>
                                    <span
                                        class="badge badge-sm"
                                        class:badge-error={u.suspended}
                                        class:badge-success={!u.suspended}
                                        >{u.suspended ? "停止" : "有効"}</span>
                                </td>
                                <td
                                    ><ChevronRight
                                        class="h-4 w-4 opacity-40" /></td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            {/if}
        </div>
    </div>
{/if}

{#if tab === "departments"}
    <div class="mt-6">
        <div class="flex justify-end mb-4">
            <button
                type="button"
                class="btn btn-primary btn-sm"
                onclick={openCreateDeptModal}>部署を追加</button>
        </div>
        <div
            class="rounded-3xl border border-base-300 bg-base-100 shadow-sm overflow-hidden">
            {#if orgLoading}
                <div class="p-6 space-y-3">
                    {#each Array(6) as _, i (i)}
                        <div class="skeleton h-10 w-full rounded-xl"></div>
                    {/each}
                </div>
            {:else if orgTree.length === 0}
                <div class="p-12 text-center text-base-content/50 text-sm">
                    部署・チームが登録されていません
                </div>
            {:else}
                <div class="divide-y divide-base-200">
                    {#each visibleOrgNodes as n (n.kind + "-" + n.id)}
                        <div
                            class="flex items-center gap-2 px-4 py-3 hover:bg-base-200/40 transition-colors"
                            style="padding-left: {1 + n.depth * 1.25}rem">
                            {#if n.kind === "department"}
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-xs btn-square shrink-0"
                                    aria-label="チームを展開"
                                    onclick={() => toggleDeptExpand(n.id)}>
                                    <ChevronRight
                                        class="h-4 w-4 transition-transform {expandedDeptIds.has(
                                            n.id,
                                        )
                                            ? 'rotate-90'
                                            : ''}" />
                                </button>
                                <Building2
                                    class="h-4 w-4 shrink-0 text-primary/80" />
                                <button
                                    type="button"
                                    class="min-w-0 flex-1 text-left"
                                    onclick={() => openOrgModal(n)}>
                                    <span class="font-medium">{n.name}</span>
                                    {#if n.departmentCount}
                                        <span
                                            class="ml-2 text-sm text-base-content/60"
                                            >{n.departmentCount} 部署</span>
                                    {/if}
                                    {#if n.teamCount}
                                        <span
                                            class="ml-2 text-sm text-base-content/60"
                                            >{n.teamCount} チーム</span>
                                    {/if}
                                </button>
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-xs"
                                    onclick={() => openCreateTeamModal(n.id)}
                                    >チームを追加</button>
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-xs"
                                    onclick={() => openOrgModal(n)}
                                    >編集</button>
                            {:else}
                                <span class="w-8 shrink-0"></span>
                                <Users
                                    class="h-4 w-4 shrink-0 text-primary/80" />
                                <button
                                    type="button"
                                    class="min-w-0 flex-1 text-left"
                                    onclick={() => openOrgModal(n)}>
                                    <span class="font-medium">{n.name}</span>
                                    <span
                                        class="ml-2 text-sm text-base-content/60"
                                        >{n.teamMemberCount ?? 0}人</span>
                                    {#if n.leaderName}
                                        <span
                                            class="ml-2 badge badge-secondary badge-sm"
                                            >SM {n.leaderName}</span>
                                    {/if}
                                </button>
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-xs"
                                    onclick={() => openOrgModal(n)}
                                    >編集</button>
                            {/if}
                        </div>
                    {/each}
                </div>
            {/if}
        </div>
    </div>
{/if}

{#if tab === "remind"}
    <div class="mt-6 max-w-xl">
        <div
            class="rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-5">
            <div role="alert" class="alert alert-info text-sm">
                <Info class="h-5 w-5 shrink-0" />
                <span
                    >設定の保存のみ行います。実際の通知送信はバックエンドジョブで別途実装されます。</span>
            </div>
            <label
                class="label cursor-pointer justify-start gap-3 rounded-2xl border border-base-300 px-4 py-3">
                <input
                    type="checkbox"
                    class="toggle toggle-primary"
                    bind:checked={remind.remindAutoSend} />
                <span class="label-text font-medium"
                    >1on1の通知を自動送信する</span>
            </label>
            <div>
                <span class="label-text font-medium">通知タイミング</span>
                <div class="flex flex-col gap-2 mt-2">
                    <label
                        class="label cursor-pointer justify-start gap-3 rounded-2xl border border-base-300 px-4 py-3">
                        <input
                            type="radio"
                            name="remindMode"
                            class="radio radio-primary"
                            value="days_before"
                            bind:group={remind.remindMode} />
                        <span class="label-text">実施予定日の N 日前</span>
                    </label>
                    <label
                        class="label cursor-pointer justify-start gap-3 rounded-2xl border border-base-300 px-4 py-3">
                        <input
                            type="radio"
                            name="remindMode"
                            class="radio radio-primary"
                            value="prev_weekday"
                            bind:group={remind.remindMode} />
                        <span class="label-text">前の週の指定曜日</span>
                    </label>
                </div>
            </div>
            {#if remind.remindMode === "days_before"}
                <label class="form-control">
                    <span class="label-text font-medium">何日前（1〜30）</span>
                    <input
                        type="number"
                        min="1"
                        max="30"
                        class="input input-bordered"
                        bind:value={remind.remindDaysBefore} />
                </label>
            {/if}
            {#if remind.remindMode === "prev_weekday"}
                <label class="form-control">
                    <span class="label-text font-medium">曜日</span>
                    <select
                        class="select select-bordered"
                        bind:value={remind.remindWeekday}>
                        {#each WEEKDAYS as w (w.value)}
                            <option value={w.value}>{w.label}</option>
                        {/each}
                    </select>
                </label>
            {/if}
            <label class="form-control">
                <span class="label-text font-medium">送信時刻</span>
                <input
                    type="time"
                    class="input input-bordered"
                    bind:value={remind.remindTime} />
            </label>
            <button
                type="button"
                class="btn btn-primary w-full sm:w-auto"
                onclick={saveRemind}>保存</button>
        </div>
    </div>
{/if}

<dialog class="modal" class:modal-open={userModalOpen}>
    <div class="modal-box max-w-2xl max-h-[85vh] flex flex-col p-0">
        <div
            class="flex justify-between items-start gap-2 p-6 pb-2 border-b border-base-200 shrink-0">
            <div>
                <h3 class="font-bold text-lg">{userForm?.name}</h3>
                <p class="text-sm text-base-content/60 mt-0.5">
                    権限・操作範囲・所属
                </p>
            </div>
            <button
                type="button"
                class="btn btn-sm btn-circle btn-ghost"
                aria-label="閉じる"
                onclick={() => (userModalOpen = false)}>×</button>
        </div>
        {#if userForm}
            <div class="overflow-y-auto px-6 py-4 space-y-5 flex-1">
                <div class="grid sm:grid-cols-2 gap-4">
                    <label class="form-control">
                        <span class="label-text font-medium">権限</span>
                        <p class="text-xs text-base-content/50 mt-0.5">
                            現在: <span
                                >{authorityLabel(userForm.authority)}</span>
                        </p>
                        <select
                            class="select select-bordered select-sm mt-1"
                            bind:value={authoritySelect}
                            onchange={onAuthoritySelectChange}>
                            {#each AUTHORITY_OPTIONS as opt (opt.value)}
                                <option value={opt.value}>{opt.label}</option>
                            {/each}
                        </select>
                    </label>
                    <label class="form-control">
                        <span class="label-text font-medium">ロール</span>
                        <p class="text-xs text-base-content/50 mt-0.5">
                            現在: <span
                                >{systemRoleLabel(userForm.systemRole)}</span>
                        </p>
                        <select
                            class="select select-bordered select-sm mt-1"
                            bind:value={systemRoleSelect}
                            onchange={onSystemRoleSelectChange}>
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
                                class:btn-primary={userForm.operationTarget ===
                                    opt.value}
                                class:btn-outline={userForm.operationTarget !==
                                    opt.value}
                                onclick={() => setOperationTarget(opt.value)}
                                >{opt.label}</button>
                        {/each}
                    </div>
                </div>
                {#if userForm.operationTarget === "department"}
                    <div class="space-y-2">
                        <input
                            class="input input-bordered input-sm w-full placeholder-current"
                            placeholder="部署を検索"
                            bind:value={targetSearch} />
                        <div
                            class="border border-base-300 rounded-2xl max-h-44 overflow-y-auto p-2 space-y-0.5">
                            {#each filteredDepartments() as d (d.id)}
                                <label
                                    class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer">
                                    <input
                                        type="checkbox"
                                        class="checkbox checkbox-sm"
                                        checked={isInList(
                                            userForm.targetIds,
                                            d.id,
                                        )}
                                        onchange={() =>
                                            toggleUserTarget(d.id)} />
                                    <span>{d.name}</span>
                                </label>
                            {/each}
                        </div>
                    </div>
                {/if}
                {#if userForm.operationTarget === "team"}
                    <div class="space-y-2">
                        <input
                            class="input input-bordered input-sm w-full placeholder-current"
                            placeholder="チームを検索"
                            bind:value={targetSearch} />
                        <div
                            class="border border-base-300 rounded-2xl max-h-44 overflow-y-auto p-2 space-y-0.5">
                            {#each filteredTeams() as t (t.id)}
                                <label
                                    class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer">
                                    <input
                                        type="checkbox"
                                        class="checkbox checkbox-sm"
                                        checked={isInList(
                                            userForm.targetIds,
                                            t.id,
                                        )}
                                        onchange={() =>
                                            toggleUserTarget(t.id)} />
                                    <span>{t.name}</span>
                                </label>
                            {/each}
                        </div>
                    </div>
                {/if}
                {#if userForm.operationTarget === "subordinate"}
                    <div class="space-y-2">
                        <input
                            class="input input-bordered input-sm w-full placeholder-current"
                            placeholder="部下を検索"
                            bind:value={targetSearch} />
                        <div
                            class="border border-base-300 rounded-2xl max-h-44 overflow-y-auto p-2 space-y-0.5">
                            {#each filteredUsersForTarget() as u (`sub-${u.id}`)}
                                {#if u.id !== userForm.id}
                                    <label
                                        class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer">
                                        <input
                                            type="checkbox"
                                            class="checkbox checkbox-sm"
                                            checked={isInList(
                                                userForm.targetIds,
                                                u.id,
                                            )}
                                            onchange={() =>
                                                toggleUserTarget(u.id)} />
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
                        <select
                            class="select select-bordered"
                            bind:value={userForm.departmentId}>
                            {#each options.departments as d (`dept-${d.id}`)}
                                <option value={d.id}>{d.name}</option>
                            {/each}
                        </select>
                    </label>
                    <label class="form-control"
                        ><span class="label-text font-medium"
                            >役職（マスタ）</span>
                        <select
                            class="select select-bordered"
                            bind:value={userForm.roleId}>
                            {#each options.roles as r (r.id)}
                                <option value={r.id}>{r.name}</option>
                            {/each}
                        </select>
                    </label>
                </div>
                <div class="border-t border-base-300 pt-4 flex flex-wrap gap-2">
                    <button
                        type="button"
                        class="btn btn-outline btn-sm"
                        onclick={deleteUserIcon}>アイコン削除</button>
                    {#if !userForm.suspended}
                        <button
                            type="button"
                            class="btn btn-warning btn-sm"
                            onclick={() => toggleSuspendInModal(false)}
                            >アカウント停止</button>
                    {/if}
                    {#if userForm.suspended}
                        <button
                            type="button"
                            class="btn btn-success btn-sm"
                            onclick={() => toggleSuspendInModal(true)}
                            >停止解除</button>
                    {/if}
                </div>
            </div>
        {/if}
        <div class="modal-action border-t border-base-200 p-4 shrink-0">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (userModalOpen = false)}>閉じる</button>
            <button type="button" class="btn btn-primary" onclick={saveUser}
                >保存</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (userModalOpen = false)}
            >close</button>
    </form>
</dialog>

<dialog class="modal" class:modal-open={orgModalOpen}>
    <div class="modal-box max-w-2xl max-h-[85vh] flex flex-col p-0">
        <div
            class="flex justify-between items-start gap-2 p-6 pb-2 border-b border-base-200 shrink-0">
            <div>
                <h3 class="font-bold text-lg">{orgDetail?.name}</h3>
                <p class="text-sm text-base-content/60">
                    {orgDetail?.kind === "team" ? "チーム" : "部署"}
                </p>
            </div>
            <div class="flex gap-1">
                {#if !orgEditMode}
                    <button
                        type="button"
                        class="btn btn-sm btn-outline"
                        onclick={() => (orgEditMode = true)}>編集</button>
                {/if}
                <button
                    type="button"
                    class="btn btn-sm btn-circle btn-ghost"
                    aria-label="閉じる"
                    onclick={() => (orgModalOpen = false)}>×</button>
            </div>
        </div>
        {#if !orgEditMode}
            <div class="overflow-y-auto px-6 py-4 flex-1 space-y-4">
                {#if orgDetail?.kind === "department" && orgDetail.headUserId}
                    <div
                        class="rounded-xl border border-base-300 bg-base-200/30 px-3 py-2 text-sm">
                        <span class="text-base-content/60">部長</span>
                        <span class="ml-2 font-medium"
                            >{userName(orgDetail.headUserId)}</span>
                    </div>
                {/if}
                {#if orgDetail?.kind === "team" && orgDetail.leaderId}
                    <div
                        class="rounded-xl border border-base-300 bg-base-200/30 px-3 py-2 text-sm">
                        <span class="text-base-content/60">SM</span>
                        <span class="ml-2 font-medium"
                            >{userName(orgDetail.leaderId)}</span>
                    </div>
                {/if}
                {#if (orgDetail?.members || []).length}
                    <p class="text-xs font-medium text-base-content/60">
                        メンバー（{(orgDetail?.members || []).length}人）
                    </p>
                    <ul class="space-y-2">
                        {#each orgDetail?.members || [] as m (m.userId)}
                            <li
                                class="flex justify-between items-center rounded-xl border border-base-300 px-3 py-2">
                                <span class="font-medium">{m.userName}</span>
                                <span class="badge badge-ghost badge-sm"
                                    >{m.title}</span>
                            </li>
                        {/each}
                    </ul>
                {:else}
                    <p class="text-sm text-base-content/50 text-center py-8">
                        所属メンバーはいません
                    </p>
                {/if}
            </div>
        {/if}
        {#if orgEditMode}
            <div class="overflow-y-auto px-6 py-4 space-y-4 flex-1">
                <label class="form-control"
                    ><span class="label-text font-medium">名前</span><input
                        class="input input-bordered"
                        bind:value={orgEdit.name} /></label>
                {#if orgDetail?.kind === "department" && orgDetail?.isRoot}
                    <div class="space-y-2">
                        <span class="label-text font-medium">子部署</span>
                        <input
                            class="input input-bordered input-sm w-full placeholder-current"
                            placeholder="検索"
                            bind:value={orgChildSearch} />
                        <div
                            class="border border-base-300 rounded-2xl max-h-44 overflow-y-auto p-2">
                            {#each filteredOrgDepartments() as d (`c-${d.id}`)}
                                {#if d.id !== orgDetail?.id}
                                    <label
                                        class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer">
                                        <input
                                            type="checkbox"
                                            class="checkbox checkbox-sm"
                                            checked={isInList(
                                                orgEdit.childIds,
                                                d.id,
                                            )}
                                            onchange={() =>
                                                toggleOrgChild(d.id)} />
                                        <span>{d.name}</span>
                                    </label>
                                {/if}
                            {/each}
                        </div>
                    </div>
                {/if}
                {#if orgDetail?.kind === "department" && !orgDetail?.isRoot}
                    <div class="space-y-4">
                        <label class="form-control">
                            <span class="label-text font-medium">部長</span>
                            <select
                                class="select select-bordered"
                                bind:value={orgEdit.headUserId}>
                                <option value={0}>未設定</option>
                                {#each headUserCandidates as u (`head-${u.id}`)}
                                    <option value={u.id}>{u.name}</option>
                                {/each}
                            </select>
                        </label>
                        <div class="space-y-2">
                            <span class="label-text font-medium">チーム</span>
                            <input
                                class="input input-bordered input-sm w-full placeholder-current"
                                placeholder="検索"
                                bind:value={orgTeamSearch} />
                            <div
                                class="border border-base-300 rounded-2xl max-h-36 overflow-y-auto p-2">
                                {#each filteredOrgTeams() as t (`t-${t.id}`)}
                                    <label
                                        class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer">
                                        <input
                                            type="checkbox"
                                            class="checkbox checkbox-sm"
                                            checked={isInList(
                                                orgEdit.teamIds,
                                                t.id,
                                            )}
                                            onchange={() =>
                                                toggleOrgTeam(t.id)} />
                                        <span>{t.name}</span>
                                    </label>
                                {/each}
                            </div>
                        </div>
                        <div class="space-y-2">
                            <span class="label-text font-medium"
                                >所属ユーザ</span>
                            <input
                                class="input input-bordered input-sm w-full placeholder-current"
                                placeholder="検索"
                                bind:value={orgUserSearch} />
                            <div
                                class="border border-base-300 rounded-2xl max-h-36 overflow-y-auto p-2">
                                {#each filteredOrgUsers() as u (`ou-${u.id}`)}
                                    <label
                                        class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer">
                                        <input
                                            type="checkbox"
                                            class="checkbox checkbox-sm"
                                            checked={isInList(
                                                orgEdit.userIds,
                                                u.id,
                                            )}
                                            onchange={() =>
                                                toggleOrgUser(u.id)} />
                                        <span>{u.name}</span>
                                    </label>
                                {/each}
                            </div>
                        </div>
                    </div>
                {/if}
                {#if orgDetail?.kind === "team"}
                    <div class="space-y-4">
                        <div class="space-y-2">
                            <span class="label-text font-medium">メンバー</span>
                            <input
                                class="input input-bordered input-sm w-full placeholder-current"
                                placeholder="検索"
                                bind:value={orgMemberSearch} />
                            <div
                                class="border border-base-300 rounded-2xl max-h-36 overflow-y-auto p-2">
                                {#each filteredOrgMembers() as u (`tm-${u.id}`)}
                                    <label
                                        class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer">
                                        <input
                                            type="checkbox"
                                            class="checkbox checkbox-sm"
                                            checked={isInList(
                                                orgEdit.memberIds,
                                                u.id,
                                            )}
                                            onchange={() =>
                                                toggleOrgMember(u.id)} />
                                        <span>{u.name}</span>
                                    </label>
                                {/each}
                            </div>
                        </div>
                        {#if orgEdit.memberIds.length}
                            <div>
                                <span class="label-text font-medium"
                                    >リーダー（SM）</span>
                                <div class="flex flex-wrap gap-2 mt-2">
                                    {#each orgEdit.memberIds as uid (`ld-${uid}`)}
                                        <button
                                            type="button"
                                            class="btn btn-sm"
                                            class:btn-secondary={Number(
                                                orgEdit.leaderId,
                                            ) === Number(uid)}
                                            class:btn-outline={Number(
                                                orgEdit.leaderId,
                                            ) !== Number(uid)}
                                            onclick={() =>
                                                (orgEdit.leaderId =
                                                    Number(uid))}
                                            >{userName(uid)}</button>
                                    {/each}
                                </div>
                            </div>
                        {/if}
                    </div>
                {/if}
                <div class="border-t border-base-300 pt-4">
                    <button
                        type="button"
                        class="btn btn-error btn-outline btn-sm"
                        onclick={askDeleteOrg}>削除</button>
                </div>
            </div>
        {/if}
        {#if !orgEditMode}
            <div class="modal-action border-t border-base-200 p-4 shrink-0">
                <button
                    type="button"
                    class="btn btn-error btn-outline btn-sm mr-auto"
                    onclick={askDeleteOrg}>削除</button>
                <button
                    type="button"
                    class="btn btn-ghost"
                    onclick={() => (orgModalOpen = false)}>閉じる</button>
            </div>
        {/if}
        {#if orgEditMode}
            <div class="modal-action border-t border-base-200 p-4 shrink-0">
                <button
                    type="button"
                    class="btn btn-ghost"
                    onclick={() => (orgModalOpen = false)}>閉じる</button>
                <button type="button" class="btn btn-primary" onclick={saveOrg}
                    >保存</button>
            </div>
        {/if}
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (orgModalOpen = false)}
            >close</button>
    </form>
</dialog>

<dialog class="modal" class:modal-open={createUserOpen}>
    <div class="modal-box max-w-lg max-h-[85vh] flex flex-col p-0">
        <div class="p-6 pb-2 border-b border-base-200 shrink-0">
            <h3 class="font-bold text-lg">ユーザを追加</h3>
            <p class="text-sm text-base-content/60 mt-0.5">
                権限はメンバー、ロールはなしで作成されます
            </p>
        </div>
        <div class="overflow-y-auto px-6 py-4 space-y-4 flex-1">
            <label class="form-control"
                ><span class="label-text font-medium">氏名</span><input
                    class="input input-bordered"
                    bind:value={createUserForm.name} /></label>
            <label class="form-control"
                ><span class="label-text font-medium">メール</span><input
                    type="email"
                    class="input input-bordered"
                    bind:value={createUserForm.email} /></label>
            <label class="form-control"
                ><span class="label-text font-medium">所属部署</span>
                <select
                    class="select select-bordered"
                    bind:value={createUserForm.departmentId}>
                    <option value={0}>選択してください</option>
                    {#each options.departments as d (`new-dept-${d.id}`)}
                        <option value={d.id}>{d.name}</option>
                    {/each}
                </select>
            </label>
            <label class="form-control"
                ><span class="label-text font-medium">役職（マスタ）</span>
                <select
                    class="select select-bordered"
                    bind:value={createUserForm.roleId}>
                    <option value={0}>選択してください</option>
                    {#each options.roles as r (`new-role-${r.id}`)}
                        <option value={r.id}>{r.name}</option>
                    {/each}
                </select>
            </label>
            <label class="form-control"
                ><span class="label-text font-medium"
                    >パスワード（任意・Basic認証用）</span
                ><input
                    type="password"
                    class="input input-bordered"
                    bind:value={createUserForm.password}
                    autocomplete="new-password" /></label>
        </div>
        <div class="modal-action border-t border-base-200 p-4 shrink-0">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (createUserOpen = false)}>キャンセル</button>
            <button
                type="button"
                class="btn btn-primary"
                onclick={submitCreateUser}>作成</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (createUserOpen = false)}
            >close</button>
    </form>
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
                    bind:value={createDeptForm.name} /></label>
            <label class="form-control"
                ><span class="label-text font-medium">親部署（任意）</span>
                <select
                    class="select select-bordered"
                    bind:value={createDeptForm.parentDepartmentId}>
                    <option value={0}>なし（ルート）</option>
                    {#each options.departments as d (`parent-${d.id}`)}
                        <option value={d.id}>{d.name}</option>
                    {/each}
                </select>
            </label>
        </div>
        <div class="modal-action border-t border-base-200 p-4 shrink-0">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (createDeptOpen = false)}>キャンセル</button>
            <button
                type="button"
                class="btn btn-primary"
                onclick={submitCreateDept}>作成</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (createDeptOpen = false)}
            >close</button>
    </form>
</dialog>

<dialog class="modal" class:modal-open={confirmOpen}>
    <div class="modal-box max-h-[85vh]">
        <h3 class="font-bold">{confirmTitle}</h3>
        <p class="py-4 text-sm text-base-content/80">{confirmMessage}</p>
        <div class="modal-action">
            <button type="button" class="btn btn-ghost" onclick={cancelConfirm}
                >キャンセル</button>
            <button
                type="button"
                class="btn"
                class:btn-error={confirmDanger}
                class:btn-primary={!confirmDanger}
                onclick={runConfirm}
                >{confirmDanger ? "確定" : "変更する"}</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={cancelConfirm}>close</button>
    </form>
</dialog>

<dialog class="modal" class:modal-open={createTeamOpen}>
    <div class="modal-box max-w-md">
        <h3 class="font-bold text-lg">チームを追加</h3>
        <label class="form-control mt-4">
            <span class="label-text font-medium">チーム名</span>
            <input
                class="input input-bordered"
                bind:value={createTeamForm.name} />
        </label>
        <div class="modal-action">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (createTeamOpen = false)}>キャンセル</button>
            <button
                type="button"
                class="btn btn-primary"
                onclick={saveCreateTeam}>作成</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (createTeamOpen = false)}
            >close</button>
    </form>
</dialog>

<dialog class="modal" class:modal-open={createRoleOpen}>
    <div class="modal-box max-w-md">
        <h3 class="font-bold text-lg">役職を登録</h3>
        <label class="form-control mt-4">
            <span class="label-text font-medium">役職名</span>
            <input
                class="input input-bordered"
                bind:value={createRoleForm.name} />
        </label>
        <label class="form-control mt-3">
            <span class="label-text font-medium">順位（大きいほど上位）</span>
            <input
                type="number"
                class="input input-bordered"
                bind:value={createRoleForm.rank} />
        </label>
        <div class="modal-action">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (createRoleOpen = false)}>キャンセル</button>
            <button
                type="button"
                class="btn btn-primary"
                onclick={saveCreateRole}>登録</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (createRoleOpen = false)}
            >close</button>
    </form>
</dialog>

<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import Users from "lucide-svelte/icons/users";
    import Building2 from "lucide-svelte/icons/building-2";
    import Bell from "lucide-svelte/icons/bell";
    import Briefcase from "lucide-svelte/icons/briefcase";
    import Search from "lucide-svelte/icons/search";
    import ChevronsUpDown from "lucide-svelte/icons/chevrons-up-down";
    import ChevronsDownUp from "lucide-svelte/icons/chevrons-down-up";
    import ChevronRight from "lucide-svelte/icons/chevron-right";
    import Info from "lucide-svelte/icons/info";
    import MemberPickModal from "$lib/components/MemberPickModal.svelte";
    import MultiCheckSelect from "$lib/components/MultiCheckSelect.svelte";
    import { fade } from "svelte/transition";
    import Icon from "$lib/components/Icon.svelte";
    import Trash2 from "lucide-svelte/icons/trash-2";
    import {
        lastPageForTotal,
        visiblePageNumbers,
    } from "$lib/planning-meetings-pager";

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
    const WEEKDAYS = [
        { value: 1, label: "月曜日" },
        { value: 2, label: "火曜日" },
        { value: 3, label: "水曜日" },
        { value: 4, label: "木曜日" },
        { value: 5, label: "金曜日" },
    ];

    type DeptAssignment = {
        departmentUuid: string;
        roleUuid: string;
        departmentName?: string;
    };
    type TeamAssignment = {
        teamUuid: string;
        isLeader?: boolean;
    };

    type ManagementUser = {
        uuid: string;
        name: string;
        email: string;
        authority: string;
        systemRole: string;
        departmentName: string;
        departmentUuid: string;
        roleUuid: string;
        departments?: DeptAssignment[];
        teams?: TeamAssignment[];
        suspended: boolean;
        imageUri?: string;
    };

    type OrgTreeNode = {
        kind: "team" | "department";
        uuid: string;
        name: string;
        depth: number;
        memberCount?: number;
        teamMemberCount?: number;
        leaderName?: string;
        leaderRole?: string;
        headUserName?: string;
        headUserRole?: string;
        teamCount?: number;
        departmentCount?: number;
    };

    type OrgMember = { userUuid: string; userName: string; title: string };

    type OrgDetail = {
        uuid: string;
        departmentUuid?: string;
        name: string;
        kind: "team" | "department";
        isRoot?: boolean;
        members?: OrgMember[];
        childDepartmentUuids?: string[];
        teamUuids?: string[];
        userUuids?: string[];
        memberUuids?: string[];
        leaderUuid?: string;
        headUserUuid?: string;
        managerUserUuids?: string[];
    };

    type UserForm = {
        uuid: string;
        name: string;
        authority: string;
        systemRole: string;
        imageUri: string | undefined;
        departments: DeptAssignment[];
        teams: TeamAssignment[];
        suspended: boolean;
    };

    type Options = {
        departments: { uuid: string; name: string; isRoot: boolean }[];
        teams: {
            uuid: string;
            name: string;
            departmentName: string;
            departmentUuid: string;
        }[];
        roles: { uuid: string; name: string }[];
    };

    type RoleItem = { uuid: string; name: string };

    function toggleInList(list: string[], id: string): string[] {
        const i = list.indexOf(id);
        if (i >= 0) return list.filter((_, idx) => idx !== i);
        return [...list, id];
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

    let tab = $state<"users" | "departments" | "roles" | "remind">("users");
    let users = $state<ManagementUser[]>([]);
    let orgTree = $state<OrgTreeNode[]>([]);
    let options = $state<Options>({ departments: [], teams: [], roles: [] });
    let selectedUserUuids = $state<string[]>([]);
    let userSearch = $state("");
    let userDeptFilter = $state<string[]>([]);
    let currentUserPage = $state(1);
    const USER_PAGE_SIZE = 20;
    let usersLoading = $state(true);
    let orgLoading = $state(true);
    let remind = $state({
        remindAutoSend: false,
        remindMode: "prev_weekday",
        remindDaysBefore: 1,
        remindWeekday: 1,
        remindTime: "09:00",
    });
    let orgChildSearch = $state("");
    let orgTeamSearch = $state("");
    let orgMemberSearch = $state("");
    let userModalOpen = $state(false);
    let userForm = $state<UserForm | null>(null);
    let authoritySelect = $state("general");
    let systemRoleSelect = $state("none");
    let createUserOpen = $state(false);
    let createUserForm = $state({
        name: "",
        email: "",
        password: "",
        departments: [] as DeptAssignment[],
        teams: [] as TeamAssignment[],
    });
    let createDeptOpen = $state(false);
    let createDeptForm = $state({
        name: "",
        parentDepartmentUuid: "",
        headUserUuid: "",
        managerUserUuids: [] as string[],
        memberUserUuids: [] as string[],
    });
    let createTeamOpen = $state(false);
    let createTeamForm = $state({
        name: "",
        departmentUuid: "",
        leaderUuid: "",
        memberUserUuids: [] as string[],
    });
    let roles = $state<RoleItem[]>([]);
    let rolesLoading = $state(true);
    let createRoleOpen = $state(false);
    let createRoleForm = $state({ name: "" });
    let editRoleOpen = $state(false);
    let editRoleForm = $state({ uuid: "", name: "" });
    let expandedAllDept = $state(false);
    let expandedDeptUuids = $state<Set<string>>(new Set());
    let orgModalOpen = $state(false);
    let orgDetail = $state<OrgDetail | null>(null);
    let orgEditMode = $state(false);
    let orgEdit = $state({
        name: "",
        childDepartmentUuids: [] as string[],
        teamUuids: [] as string[],
        leaderUuid: "",
        headUserUuid: "",
        managerUserUuids: [] as string[],
        memberUserUuids: [] as string[],
    });
    type OrgPickTarget = "member" | "manager" | "head" | "leader";
    type OrgPickOwner = "edit" | "createDept" | "createTeam";
    let orgMemberPickOpen = $state(false);
    let orgMemberPickTarget = $state<OrgPickTarget>("member");
    let orgPickOwner = $state<OrgPickOwner>("edit");
    let orgManagerSearch = $state("");

    let memberPickerUsers = $derived(
        users.map((u) => ({
            uuid: u.uuid,
            name: u.name,
            email: u.email,
            imageUri: u.imageUri,
        })),
    );

    function openOrgPick(target: OrgPickTarget, owner: OrgPickOwner = "edit") {
        orgMemberPickTarget = target;
        orgPickOwner = owner;
        orgMemberPickOpen = true;
    }

    function createTeamLeaderCandidates(departmentUuid: string) {
        if (!departmentUuid) return [];
        return teamLeaderCandidates.filter((u) => {
            if (u.departmentUuid === departmentUuid) return true;
            return (u.departments || []).some(
                (d) => d.departmentUuid === departmentUuid,
            );
        });
    }

    function toPickerUsers(candidates: ManagementUser[]) {
        return candidates.map((u) => ({
            uuid: u.uuid,
            name: u.name,
            email: u.email,
            imageUri: u.imageUri,
        }));
    }

    function uuidList(uuid: string): string[] {
        return uuid ? [uuid] : [];
    }

    let orgMemberPickUsers = $derived.by(() => {
        if (
            orgMemberPickTarget === "manager" ||
            orgMemberPickTarget === "head"
        ) {
            return managerPickerUsers;
        }
        if (orgMemberPickTarget === "leader") {
            const deptUuid =
                orgPickOwner === "createTeam"
                    ? createTeamForm.departmentUuid
                    : orgDetail?.departmentUuid || "";
            return toPickerUsers(createTeamLeaderCandidates(deptUuid));
        }
        return memberPickerUsers;
    });

    let orgMemberPickMultiple = $derived(
        orgMemberPickTarget !== "head" && orgMemberPickTarget !== "leader",
    );

    let orgMemberPickTitle = $derived(
        orgMemberPickTarget === "head"
            ? "責任者を選択"
            : orgMemberPickTarget === "leader"
              ? "チームリーダーを選択"
              : orgMemberPickTarget === "manager"
                ? "管理職を選択"
                : "メンバーを選択",
    );

    let orgMemberPickSelectedUuids = $derived.by(() => {
        if (orgPickOwner === "createDept") {
            if (orgMemberPickTarget === "head") {
                return uuidList(createDeptForm.headUserUuid);
            }
            if (orgMemberPickTarget === "manager") {
                return createDeptForm.managerUserUuids;
            }
            return createDeptForm.memberUserUuids;
        }
        if (orgPickOwner === "createTeam") {
            if (orgMemberPickTarget === "leader") {
                return uuidList(createTeamForm.leaderUuid);
            }
            return createTeamForm.memberUserUuids;
        }
        if (orgMemberPickTarget === "head") {
            return uuidList(orgEdit.headUserUuid);
        }
        if (orgMemberPickTarget === "leader") {
            return uuidList(orgEdit.leaderUuid);
        }
        if (orgMemberPickTarget === "manager") {
            return orgEdit.managerUserUuids;
        }
        return orgEdit.memberUserUuids;
    });

    let orgMemberPickExcludeUuids = $derived.by(() => {
        if (
            orgMemberPickTarget === "head" ||
            orgMemberPickTarget === "leader"
        ) {
            return [] as string[];
        }
        if (orgPickOwner === "createDept") {
            if (orgMemberPickTarget === "manager") {
                return createDeptForm.headUserUuid
                    ? [createDeptForm.headUserUuid]
                    : [];
            }
            return [
                createDeptForm.headUserUuid,
                ...createDeptForm.managerUserUuids,
            ].filter((id): id is string => !!id);
        }
        if (orgPickOwner === "createTeam") {
            return createTeamForm.leaderUuid ? [createTeamForm.leaderUuid] : [];
        }
        if (orgMemberPickTarget === "manager") {
            return orgEdit.headUserUuid ? [orgEdit.headUserUuid] : [];
        }
        return [
            orgEdit.headUserUuid,
            ...orgEdit.managerUserUuids,
            orgEdit.leaderUuid,
        ].filter((id): id is string => !!id);
    });

    function applyOrgMemberPick(uuids: string[]) {
        const single = uuids[0] || "";
        if (orgMemberPickTarget === "head") {
            if (orgPickOwner === "createDept") {
                createDeptForm.headUserUuid = single;
            } else {
                orgEdit.headUserUuid = single;
            }
        } else if (orgMemberPickTarget === "leader") {
            if (orgPickOwner === "createTeam") {
                createTeamForm.leaderUuid = single;
            } else {
                orgEdit.leaderUuid = single;
            }
        } else if (orgPickOwner === "createDept") {
            if (orgMemberPickTarget === "manager") {
                createDeptForm.managerUserUuids = uuids;
            } else {
                createDeptForm.memberUserUuids = uuids;
            }
        } else if (orgPickOwner === "createTeam") {
            createTeamForm.memberUserUuids = uuids;
        } else if (orgMemberPickTarget === "manager") {
            orgEdit.managerUserUuids = uuids;
        } else {
            orgEdit.memberUserUuids = uuids;
        }
        orgMemberPickOpen = false;
    }

    function uniqueUuids(groups: (string | undefined | null)[][]): string[] {
        return [...new Set(groups.flat().filter((id): id is string => !!id))];
    }

    function plainDepartmentMemberUuids(detail: OrgDetail): string[] {
        const head = detail.headUserUuid || "";
        const managers = new Set(detail.managerUserUuids || []);
        return (detail.userUuids || []).filter(
            (id) => id && id !== head && !managers.has(id),
        );
    }

    function plainTeamMemberUuids(detail: OrgDetail): string[] {
        const leader = detail.leaderUuid || "";
        return (detail.memberUuids || []).filter((id) => id && id !== leader);
    }

    function filteredListedUuids(uuids: string[], search: string): string[] {
        const q = search.trim().toLowerCase();
        if (!q) return uuids;
        return uuids.filter((id) => {
            const u = users.find((x) => x.uuid === id);
            if (!u) return false;
            return (
                (u.name || "").toLowerCase().includes(q) ||
                (u.email || "").toLowerCase().includes(q)
            );
        });
    }

    function userByUuid(uuid: string): ManagementUser | undefined {
        return users.find((x) => x.uuid === uuid);
    }

    function isDepartmentLeaderCandidate(u: ManagementUser): boolean {
        return u.authority !== "general";
    }

    let departmentLeaderCandidates = $derived(
        users.filter(isDepartmentLeaderCandidate),
    );
    let managerPickerUsers = $derived(
        departmentLeaderCandidates.map((u) => ({
            uuid: u.uuid,
            name: u.name,
            email: u.email,
            imageUri: u.imageUri,
        })),
    );
    let teamLeaderCandidates = $derived(users);

    let visibleOrgNodes = $derived.by(() => {
        const visibleNode: OrgTreeNode[] = [];
        let hiddenBelow: number | null = null;

        for (const n of orgTree) {
            if (hiddenBelow != null && n.depth > hiddenBelow) continue;
            hiddenBelow = null;

            visibleNode.push(n);
            if (n.kind === "department" && !expandedDeptUuids.has(n.uuid)) {
                hiddenBelow = n.depth;
            }
        }
        return visibleNode;
    });

    function toggleDeptExpand(uuid: string) {
        const next = new Set(expandedDeptUuids);
        if (next.has(uuid)) next.delete(uuid);
        else next.add(uuid);
        expandedDeptUuids = next;
    }
    let confirmOpen = $state(false);
    let confirmTitle = $state("");
    let confirmMessage = $state("");
    let confirmAction = $state<(() => void | Promise<void>) | null>(null);
    let confirmOnCancel = $state<(() => void) | null>(null);
    let confirmDanger = $state(false);

    let filteredUsers = $derived.by(() => {
        const q = userSearch.trim().toLowerCase();
        return users.filter((u) => {
            const matchText =
                !q ||
                (u.name || "").toLowerCase().includes(q) ||
                (u.email || "").toLowerCase().includes(q);
            const deptIds = new Set(userDeptFilter);
            const matchDept =
                deptIds.size === 0 ||
                deptIds.has(u.departmentUuid) ||
                (u.departments || []).some((d) =>
                    deptIds.has(d.departmentUuid),
                );
            return matchText && matchDept;
        });
    });

    function departmentLabels(u: ManagementUser): string[] {
        const fromRows = (u.departments || [])
            .map((d) => d.departmentName || "")
            .filter((name) => name !== "");
        const names =
            fromRows.length > 0
                ? fromRows
                : u.departmentName
                  ? [u.departmentName]
                  : [];
        return [...new Set(names)];
    }

    function teamLabels(u: ManagementUser): string[] {
        const byUuid = new Map(options.teams.map((t) => [t.uuid, t.name]));
        return [
            ...new Set(
                (u.teams || [])
                    .map((t) => byUuid.get(t.teamUuid) || "")
                    .filter((name) => name !== ""),
            ),
        ];
    }

    let pagedFilteredUsers = $derived.by(() => {
        const start = (currentUserPage - 1) * USER_PAGE_SIZE;
        return filteredUsers.slice(start, start + USER_PAGE_SIZE);
    });

    let userListTotal = $derived(filteredUsers.length);
    let userLastPage = $derived(lastPageForTotal(userListTotal, USER_PAGE_SIZE));
    let userPagerPages = $derived(
        visiblePageNumbers(currentUserPage, userLastPage),
    );

    let allUsersSelected = $derived(
        pagedFilteredUsers.length > 0 &&
            pagedFilteredUsers.every((u) => selectedUserUuids.includes(u.uuid)),
    );

    let userFilterSignature = $derived(
        JSON.stringify({
            q: userSearch.trim(),
            dept: userDeptFilter,
        }),
    );
    let userFilterSignatureMounted = "";

    $effect(() => {
        const sig = userFilterSignature;
        if (userFilterSignatureMounted === "") return;
        if (sig === userFilterSignatureMounted) return;
        userFilterSignatureMounted = sig;
        currentUserPage = 1;
    });

    function goToUserPage(page: number) {
        if (page < 1 || page > userLastPage) return;
        currentUserPage = page;
    }

    function authorityLabel(v: string) {
        return AUTHORITY_OPTIONS.find((o) => o.value === v)?.label || v;
    }

    function systemRoleLabel(v: string) {
        return SYSTEM_ROLE_OPTIONS.find((o) => o.value === v)?.label || v;
    }

    function userName(uuid: string) {
        const u = users.find((x) => x.uuid === uuid);
        return u ? u.name : "";
    }

    function managerNames(uuids: string[] | undefined): string {
        const ids = uuids || [];
        if (ids.length === 0) return "未設定";
        return ids
            .map((id) => userName(id))
            .filter(Boolean)
            .join("、");
    }

    function selectedDepartmentUuids(rows: DeptAssignment[]) {
        return rows.map((row) => row.departmentUuid).filter((id) => id !== "");
    }

    function departmentsHaveTeams(departments: DeptAssignment[]) {
        const deptIds = new Set(selectedDepartmentUuids(departments));
        return options.teams.some((team) => deptIds.has(team.departmentUuid));
    }

    function departmentChoices(rows: DeptAssignment[], index: number) {
        const taken = new Set(
            rows
                .map((row, i) => (i === index ? "" : row.departmentUuid))
                .filter((id) => id !== ""),
        );
        const current = rows[index]?.departmentUuid || "";
        return options.departments.filter(
            (dept) => !taken.has(dept.uuid) || dept.uuid === current,
        );
    }

    function teamChoices(
        departments: DeptAssignment[],
        teams: TeamAssignment[],
        index: number,
    ) {
        const deptIds = new Set(selectedDepartmentUuids(departments));
        const taken = new Set(
            teams
                .map((row, i) => (i === index ? "" : row.teamUuid))
                .filter((id) => id !== ""),
        );
        const current = teams[index]?.teamUuid || "";
        return options.teams.filter(
            (team) =>
                deptIds.has(team.departmentUuid) &&
                (!taken.has(team.uuid) || team.uuid === current),
        );
    }

    function pruneTeamRows(
        departments: DeptAssignment[],
        teams: TeamAssignment[],
    ) {
        const deptIds = new Set(selectedDepartmentUuids(departments));
        return teams.filter((row) => {
            const team = options.teams.find(
                (item) => item.uuid === row.teamUuid,
            );
            return !!team && deptIds.has(team.departmentUuid);
        });
    }

    function nextDepartmentRow(rows: DeptAssignment[]): DeptAssignment | null {
        const taken = new Set(rows.map((row) => row.departmentUuid));
        const dept = options.departments.find((item) => !taken.has(item.uuid));
        if (!dept) return null;
        return {
            departmentUuid: dept.uuid,
            roleUuid: "",
        };
    }

    function nextTeamRow(
        departments: DeptAssignment[],
        teams: TeamAssignment[],
    ): TeamAssignment | null {
        const choices = teamChoices(
            departments,
            [...teams, { teamUuid: "" }],
            teams.length,
        );
        const team = choices[0];
        if (!team) return null;
        return { teamUuid: team.uuid };
    }

    function onUserDepartmentChange() {
        if (!userForm) return;
        userForm.teams = pruneTeamRows(userForm.departments, userForm.teams);
    }

    function addUserDepartmentRow() {
        if (!userForm) return;
        const row = nextDepartmentRow(userForm.departments);
        if (!row) return;
        userForm.departments = [...userForm.departments, row];
    }

    function removeUserDepartmentRow(index: number) {
        if (!userForm) return;
        userForm.departments = userForm.departments.filter(
            (_, i) => i !== index,
        );
        onUserDepartmentChange();
    }

    function addUserTeamRow() {
        if (!userForm) return;
        const row = nextTeamRow(userForm.departments, userForm.teams);
        if (!row) return;
        userForm.teams = [...userForm.teams, row];
    }

    function removeUserTeamRow(index: number) {
        if (!userForm) return;
        userForm.teams = userForm.teams.filter((_, i) => i !== index);
    }

    function onCreateDepartmentChange() {
        createUserForm.teams = pruneTeamRows(
            createUserForm.departments,
            createUserForm.teams,
        );
    }

    function addCreateDepartmentRow() {
        const row = nextDepartmentRow(createUserForm.departments);
        if (!row) return;
        createUserForm.departments = [...createUserForm.departments, row];
    }

    function removeCreateDepartmentRow(index: number) {
        createUserForm.departments = createUserForm.departments.filter(
            (_, i) => i !== index,
        );
        onCreateDepartmentChange();
    }

    function addCreateTeamRow() {
        const row = nextTeamRow(
            createUserForm.departments,
            createUserForm.teams,
        );
        if (!row) return;
        createUserForm.teams = [...createUserForm.teams, row];
    }

    function removeCreateTeamRow(index: number) {
        createUserForm.teams = createUserForm.teams.filter(
            (_, i) => i !== index,
        );
    }

    function userAvatar(u: ManagementUser) {
        return u.imageUri ? `/uploads/${u.imageUri}` : null;
    }

    function userAvatarByUri(uri: string | undefined) {
        return uri ? `/uploads/${uri}` : undefined;
    }

    function userInitial(u: ManagementUser) {
        return (u.name || "?").slice(0, 1);
    }

    function filteredOrgDepartments() {
        return filterByQuery(options.departments || [], orgChildSearch);
    }

    function filteredOrgTeams() {
        return filterByQuery(options.teams || [], orgTeamSearch);
    }

    function isInList(list: string[] | undefined, id: string) {
        return (list || []).includes(id);
    }

    function initExpandDept() {
        expandedDeptUuids = new Set(
            orgTree
                .filter((n) => n.kind === "department" && n.depth < 1)
                .map((n) => n.uuid),
        );
    }

    function toggleExpandAllDept() {
        if (expandedAllDept) initExpandDept();
        else expandedDeptUuids = new Set(orgTree.map((n) => n.uuid));
        expandedAllDept = !expandedAllDept;
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

    async function loadRoles() {
        rolesLoading = true;
        try {
            roles = await api<RoleItem[]>("GET", "/api/v1/management/roles");
        } catch {
            roles = [];
        }
        rolesLoading = false;
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
        const pageUuids = pagedFilteredUsers.map((u) => u.uuid);
        if (checked) {
            const set = new Set(selectedUserUuids);
            for (const id of pageUuids) set.add(id);
            selectedUserUuids = [...set];
        } else {
            const pageSet = new Set(pageUuids);
            selectedUserUuids = selectedUserUuids.filter((id) => !pageSet.has(id));
        }
    }

    function toggleOrgChild(uuid: string) {
        orgEdit.childDepartmentUuids = toggleInList(
            orgEdit.childDepartmentUuids || [],
            uuid,
        );
    }

    function toggleOrgTeam(uuid: string) {
        orgEdit.teamUuids = toggleInList(orgEdit.teamUuids || [], uuid);
    }

    function openUserModal(u: ManagementUser) {
        const departments = (
            u.departments && u.departments.length
                ? u.departments
                : [
                      {
                          departmentUuid: u.departmentUuid || "",
                          roleUuid: u.roleUuid || "",
                      },
                  ]
        ).map((row) => ({
            departmentUuid: row.departmentUuid,
            roleUuid: row.roleUuid,
        }));
        const teams = (u.teams || []).map((row) => ({
            teamUuid: row.teamUuid,
            isLeader: !!row.isLeader,
        }));
        userForm = {
            uuid: u.uuid,
            name: u.name,
            authority: u.authority,
            systemRole: u.systemRole,
            imageUri: u.imageUri,
            departments,
            teams: pruneTeamRows(departments, teams),
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
        const row = nextDepartmentRow([]);
        createUserForm = {
            name: "",
            email: "",
            password: "",
            departments: row ? [row] : [],
            teams: [],
        };
        createUserOpen = true;
    }

    async function submitCreateUser() {
        const departments = createUserForm.departments.map((row) => ({
            departmentUuid: row.departmentUuid,
            roleUuid: row.roleUuid,
        }));
        if (
            departments.length === 0 ||
            departments.some((row) => !row.departmentUuid || !row.roleUuid)
        ) {
            showToast("所属部署と役職を選択してください", "error");
            return;
        }
        const teams = createUserForm.teams.map((row) => ({
            teamUuid: row.teamUuid,
        }));
        if (teams.some((row) => !row.teamUuid)) {
            showToast("チームを選択してください", "error");
            return;
        }
        try {
            await api("POST", "/api/v1/management/users", {
                name: createUserForm.name,
                email: createUserForm.email,
                password: createUserForm.password,
                departments,
                teams,
            });
        } catch (e) {
            showToast(apiErrorMessage(e, "作成に失敗しました"), "error");
            return;
        }
        createUserOpen = false;
        await Promise.all([loadUsers(), loadOrgTree()]);
        showToast("ユーザを追加しました");
    }

    function openCreateDeptModal(parentDepartmentUuid = "") {
        orgManagerSearch = "";
        orgMemberSearch = "";
        createDeptForm = {
            name: "",
            parentDepartmentUuid,
            headUserUuid: "",
            managerUserUuids: [],
            memberUserUuids: [],
        };
        createDeptOpen = true;
    }

    async function submitCreateDept() {
        const headUuid = createDeptForm.headUserUuid;
        const managerUserUuids = createDeptForm.managerUserUuids.filter(
            (id) => id && id !== headUuid,
        );
        const userUuids = uniqueUuids([
            createDeptForm.memberUserUuids,
            headUuid ? [headUuid] : [],
            managerUserUuids,
        ]);
        try {
            await api("POST", "/api/v1/management/departments", {
                name: createDeptForm.name,
                parentDepartmentUuid:
                    createDeptForm.parentDepartmentUuid || undefined,
                headUserUuid: headUuid || undefined,
                managerUserUuids,
                userUuids,
            });
        } catch (e) {
            showToast(apiErrorMessage(e, "作成に失敗しました"), "error");
            return;
        }
        createDeptOpen = false;
        await Promise.all([loadOrgTree(), loadOptions()]);
        toggleExpandAllDept();
        showToast("部署を追加しました");
    }

    async function saveUser() {
        if (!userForm) return;
        const departments = userForm.departments.map((row) => ({
            departmentUuid: row.departmentUuid,
            roleUuid: row.roleUuid,
        }));
        if (
            departments.length === 0 ||
            departments.some((row) => !row.departmentUuid || !row.roleUuid)
        ) {
            showToast("所属部署と役職を選択してください", "error");
            return;
        }
        const teams = userForm.teams.map((row) => ({
            teamUuid: row.teamUuid,
        }));
        if (teams.some((row) => !row.teamUuid)) {
            showToast("チームを選択してください", "error");
            return;
        }
        const body = {
            authority: userForm.authority,
            systemRole: userForm.systemRole,
            departments,
            teams,
            suspended: userForm.suspended,
        };
        try {
            await api("PUT", `/api/v1/management/users/${userForm.uuid}`, body);
        } catch (e) {
            showToast(apiErrorMessage(e, "保存に失敗しました"), "error");
            return;
        }
        userModalOpen = false;
        await Promise.all([loadUsers(), loadOrgTree()]);
        showToast("ユーザ情報を保存しました");
    }

    async function deleteUserIcon() {
        if (!userForm) return;
        try {
            await api(
                "DELETE",
                `/api/v1/management/users/${userForm.uuid}/icon`,
            );
        } catch {
            showToast("アイコンの削除に失敗しました", "error");
            return;
        }
        await loadUsers();
        showToast("アイコンを削除しました");
    }

    async function setSuspendedForUsers(uuids: string[], suspended: boolean) {
        try {
            await api("POST", "/api/v1/management/users/suspend", {
                uuids,
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
        const ok = await setSuspendedForUsers([userForm.uuid], !unsuspend);
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
        confirmMessage = `選択した ${selectedUserUuids.length} 件のアカウントを停止しますか？`;
        confirmAction = async () => {
            const ok = await setSuspendedForUsers(selectedUserUuids, true);
            if (!ok) return;
            selectedUserUuids = [];
            await loadUsers();
            showToast("アカウントを停止しました");
        };
        confirmOpen = true;
    }

    function askUnsuspend() {
        confirmDanger = false;
        confirmOnCancel = null;
        confirmTitle = "停止解除";
        confirmMessage = `選択した ${selectedUserUuids.length} 件のアカウントの停止を解除しますか？`;
        confirmAction = async () => {
            const ok = await setSuspendedForUsers(selectedUserUuids, false);
            if (!ok) return;
            selectedUserUuids = [];
            await loadUsers();
            showToast("アカウントの停止を解除しました");
        };
        confirmOpen = true;
    }

    async function openOrgModal(n: OrgTreeNode) {
        orgEditMode = false;
        orgChildSearch = "";
        orgTeamSearch = "";
        orgMemberSearch = "";
        orgManagerSearch = "";
        const url =
            n.kind === "team"
                ? `/api/v1/management/teams/${n.uuid}`
                : `/api/v1/management/departments/${n.uuid}`;
        let detail: OrgDetail;
        try {
            const result = await Promise.all([
                api<OrgDetail>("GET", url),
                loadUsers(),
            ]);
            detail = result[0];
        } catch {
            showToast("詳細の取得に失敗しました", "error");
            return;
        }
        orgDetail = detail;
        orgEdit = {
            name: detail.name,
            childDepartmentUuids: detail.childDepartmentUuids || [],
            teamUuids: detail.teamUuids || [],
            leaderUuid: detail.leaderUuid || "",
            headUserUuid: detail.headUserUuid || "",
            managerUserUuids: [...(detail.managerUserUuids || [])],
            memberUserUuids:
                detail.kind === "team"
                    ? plainTeamMemberUuids(detail)
                    : plainDepartmentMemberUuids(detail),
        };
        orgModalOpen = true;
    }

    async function saveCreateTeam() {
        if (!createTeamForm.departmentUuid) {
            showToast("部署を選択してください", "error");
            return;
        }
        const leaderUuid = createTeamForm.leaderUuid;
        const memberUuids = uniqueUuids([
            createTeamForm.memberUserUuids,
            leaderUuid ? [leaderUuid] : [],
        ]);
        try {
            await api("POST", "/api/v1/management/teams", {
                name: createTeamForm.name,
                departmentUuid: createTeamForm.departmentUuid,
                leaderUuid: leaderUuid || undefined,
                memberUuids,
            });
        } catch (e) {
            showToast(
                apiErrorMessage(e, "チームの作成に失敗しました"),
                "error",
            );
            return;
        }
        createTeamOpen = false;
        await Promise.all([loadOrgTree(), loadOptions()]);
        toggleExpandAllDept();
        showToast("チームを追加しました");
    }

    function openCreateRoleModal() {
        createRoleForm = { name: "" };
        createRoleOpen = true;
    }

    async function saveCreateRole() {
        try {
            await api("POST", "/api/v1/management/roles", createRoleForm);
        } catch (e) {
            showToast(apiErrorMessage(e, "役職の登録に失敗しました"), "error");
            return;
        }
        createRoleOpen = false;
        await Promise.all([loadRoles(), loadOptions()]);
        showToast("役職を登録しました");
    }

    function openEditRoleModal(role: RoleItem) {
        editRoleForm = { uuid: role.uuid, name: role.name };
        editRoleOpen = true;
    }

    async function saveEditRole() {
        try {
            await api("PUT", `/api/v1/management/roles/${editRoleForm.uuid}`, {
                name: editRoleForm.name,
            });
        } catch (e) {
            showToast(apiErrorMessage(e, "役職の更新に失敗しました"), "error");
            return;
        }
        editRoleOpen = false;
        await Promise.all([loadRoles(), loadOptions()]);
        showToast("役職を更新しました");
    }

    function askDeleteRole(role: RoleItem) {
        confirmDanger = true;
        confirmOnCancel = null;
        confirmTitle = "削除の確認";
        confirmMessage = `役職「${role.name}」を削除しますか？`;
        confirmAction = async () => {
            try {
                await api("DELETE", `/api/v1/management/roles/${role.uuid}`);
            } catch (e) {
                showToast(apiErrorMessage(e, "削除に失敗しました"), "error");
                return;
            }
            await Promise.all([loadRoles(), loadOptions()]);
            showToast("役職を削除しました");
        };
        confirmOpen = true;
    }

    function openCreateTeamModal(departmentUuid = "") {
        orgMemberSearch = "";
        createTeamForm = {
            name: "",
            departmentUuid,
            leaderUuid: "",
            memberUserUuids: [],
        };
        createTeamOpen = true;
    }

    function onCreateTeamDepartmentChange(event: Event) {
        const value = (event.currentTarget as HTMLSelectElement).value;
        createTeamForm.departmentUuid = value;
        createTeamForm.leaderUuid = "";
        createTeamForm.memberUserUuids = [];
    }

    async function saveOrg() {
        if (!orgDetail) return;
        try {
            if (orgDetail.kind === "team") {
                const leaderUuid = orgEdit.leaderUuid;
                const memberUuids = uniqueUuids([
                    orgEdit.memberUserUuids,
                    leaderUuid ? [leaderUuid] : [],
                ]);
                await api("PUT", `/api/v1/management/teams/${orgDetail.uuid}`, {
                    name: orgEdit.name,
                    leaderUuid,
                    memberUuids,
                });
            } else {
                const headUuid = orgEdit.headUserUuid;
                const managerUserUuids = orgEdit.managerUserUuids.filter(
                    (id) => id && id !== headUuid,
                );
                const userUuids = uniqueUuids([
                    orgEdit.memberUserUuids,
                    headUuid ? [headUuid] : [],
                    managerUserUuids,
                ]);
                await api(
                    "PUT",
                    `/api/v1/management/departments/${orgDetail.uuid}`,
                    {
                        name: orgEdit.name,
                        childDepartmentUuids: orgEdit.childDepartmentUuids,
                        teamUuids: orgEdit.teamUuids,
                        headUserUuid: headUuid,
                        managerUserUuids,
                        userUuids,
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
                    ? `/api/v1/management/teams/${orgDetail.uuid}`
                    : `/api/v1/management/departments/${orgDetail.uuid}`;
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

    function switchTab(name: "users" | "departments" | "roles" | "remind") {
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

    onMount(async () => {
        await Promise.all([
            loadUsers(),
            loadOptions(),
            loadOrgTree(),
            loadRoles(),
            loadRemind(),
        ]);
        initExpandDept();
        userFilterSignatureMounted = userFilterSignature;
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
        onclick={() => switchTab("users")}
    >
        <Users class="h-4 w-4" />ユーザ一覧
    </button>
    <button
        type="button"
        role="tab"
        class="tab gap-2"
        class:tab-active={tab === "departments"}
        onclick={() => switchTab("departments")}
    >
        <Building2 class="h-4 w-4" />組織一覧
    </button>
    <button
        type="button"
        role="tab"
        class="tab gap-2"
        class:tab-active={tab === "roles"}
        onclick={() => switchTab("roles")}
    >
        <Briefcase class="h-4 w-4" />役職
    </button>
    <button
        type="button"
        role="tab"
        class="tab gap-2"
        class:tab-active={tab === "remind"}
        onclick={() => switchTab("remind")}
    >
        <Bell class="h-4 w-4" />1on1通知
    </button>
</div>

{#if tab === "users"}
    <div class="mt-6">
        <div
            class="flex flex-col sm:flex-row sm:items-center gap-3 mb-4 justify-between"
        >
            <div class="flex flex-1 flex-col sm:flex-row gap-3 max-w-2xl">
                <label
                    class="input input-bordered flex items-center gap-2 flex-1"
                >
                    <Search class="h-4 w-4 opacity-50" />
                    <input
                        type="search"
                        class="grow placeholder-current"
                        placeholder="氏名・メールで検索"
                        bind:value={userSearch}
                    />
                </label>
                <div class="w-full sm:w-52 shrink-0">
                    <MultiCheckSelect
                        pickRootId="management-user-dept"
                        items={options.departments}
                        bind:selectedUuids={userDeptFilter}
                        labelOf={(d) => d.name}
                        placeholder="すべての部署"
                    />
                </div>
            </div>
            <div class="flex justify-end gap-3">
                {#if selectedUserUuids.length > 0}
                    <div class="flex items-center gap-2" transition:fade>
                        <span class="text-sm text-base-content/60"
                            >{selectedUserUuids.length} 件選択中</span
                        >
                        {#if users.some((u) => selectedUserUuids.indexOf(u.uuid) !== -1 && !u.suspended)}
                            <button
                                type="button"
                                class="btn btn-error btn-sm"
                                onclick={askSuspend}>アカウント停止</button
                            >
                        {/if}
                        {#if users.some((u) => selectedUserUuids.indexOf(u.uuid) !== -1 && u.suspended)}
                            <button
                                type="button"
                                class="btn btn-success btn-sm"
                                onclick={askUnsuspend}>停止解除</button
                            >
                        {/if}
                    </div>
                {/if}
                <button
                    type="button"
                    class="btn btn-primary btn-sm shrink-0"
                    onclick={openCreateUserModal}>ユーザを追加</button
                >
            </div>
        </div>
        <div
            class="overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm"
        >
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
                                        )}
                                />
                            </th>
                            <th>部署</th>
                            <th>チーム</th>
                            <th>ユーザ</th>
                            <th>権限</th>
                            <th>ロール</th>
                            <th>状態</th><th class="w-8"></th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each pagedFilteredUsers as u (u.uuid)}
                            {@const deptNames = departmentLabels(u)}
                            {@const teamNames = teamLabels(u)}
                            <tr
                                class="hover:bg-base-200/50 cursor-pointer"
                                onclick={() => openUserModal(u)}
                            >
                                <!-- svelte-ignore a11y_no_static_element_interactions -->
                                <td onclick={(e) => e.stopPropagation()}>
                                    <input
                                        type="checkbox"
                                        class="checkbox checkbox-sm"
                                        value={u.uuid}
                                        bind:group={selectedUserUuids}
                                    />
                                </td>
                                <td>
                                    <span title={deptNames.join("、")}>
                                        {deptNames[0] || ""}
                                        {#if deptNames.length > 1}
                                            <span class="text-base-content/50"
                                                >+{deptNames.length - 1}</span
                                            >
                                        {/if}
                                    </span>
                                </td>
                                <td>
                                    <span title={teamNames.join("、")}>
                                        {teamNames[0] || ""}
                                        {#if teamNames.length > 1}
                                            <span class="text-base-content/50"
                                                >+{teamNames.length - 1}</span
                                            >
                                        {/if}
                                    </span>
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
                                            <div class="font-medium">
                                                {u.name}
                                            </div>
                                            <div
                                                class="text-xs text-base-content/50"
                                            >
                                                {u.email}
                                            </div>
                                        </div>
                                    </div>
                                </td>
                                <td>
                                    <span class="badge badge-outline badge-sm"
                                        >{authorityLabel(u.authority)}</span
                                    >
                                </td>
                                <td
                                    ><span class="badge badge-ghost badge-sm"
                                        >{systemRoleLabel(u.systemRole)}</span
                                    ></td
                                >
                                <td>
                                    <span
                                        class="badge badge-sm"
                                        class:badge-error={u.suspended}
                                        class:badge-success={!u.suspended}
                                        >{u.suspended ? "停止" : "有効"}</span
                                    >
                                </td>
                                <td
                                    ><ChevronRight
                                        class="h-4 w-4 opacity-40"
                                    /></td
                                >
                            </tr>
                        {/each}
                    </tbody>
                </table>
            {/if}
        </div>
        {#if userListTotal > 0}
            <div
                class="mt-4 flex flex-wrap items-center justify-center gap-2 text-sm"
            >
                <button
                    type="button"
                    class="btn btn-sm btn-ghost"
                    disabled={currentUserPage <= 1}
                    onclick={() => goToUserPage(currentUserPage - 1)}
                    >前へ</button
                >
                {#each userPagerPages as p (typeof p === "number" ? p : `e-${p}`)}
                    {#if p === "ellipsis"}
                        <span class="px-1 text-base-content/50">…</span>
                    {:else}
                        <button
                            type="button"
                            class="btn btn-sm {p === currentUserPage
                                ? 'btn-primary'
                                : 'btn-ghost'}"
                            onclick={() => goToUserPage(p)}>{p}</button
                        >
                    {/if}
                {/each}
                <button
                    type="button"
                    class="btn btn-sm btn-ghost"
                    disabled={currentUserPage >= userLastPage}
                    onclick={() => goToUserPage(currentUserPage + 1)}
                    >次へ</button
                >
                <span class="ml-2 text-base-content/60"
                    >全 {userListTotal} 件</span
                >
            </div>
        {/if}
    </div>
{/if}

{#if tab === "departments"}
    <div class="mt-6">
        <div class="flex justify-between mb-4 items-center">
            <div class="flex justfy-start gap-2 items-center">
                {#if expandedAllDept}
                    <button
                        type="button"
                        class="btn btn-ghost btn-xs btn-square shrink-0 btn-outline"
                        aria-label="部署を展開"
                        onclick={() => toggleExpandAllDept()}
                    >
                        <ChevronsDownUp
                            class="h-4 w-4 transition-transform {expandedAllDept
                                ? 'rotate-180'
                                : ''}"
                        />
                    </button>
                    閉じる
                {:else}
                    <button
                        type="button"
                        class="btn btn-ghost btn-xs btn-square shrink-0 btn-outline"
                        aria-label="部署を展開"
                        onclick={() => toggleExpandAllDept()}
                    >
                        <ChevronsUpDown
                            class="h-4 w-4 transition-transform {expandedAllDept
                                ? 'rotate-180'
                                : ''}"
                        />
                    </button>
                    全て展開
                {/if}
            </div>
            <div class="flex gap-2">
                <button
                    type="button"
                    class="btn btn-primary btn-sm"
                    onclick={() => openCreateDeptModal()}>部署を追加</button
                >
                <button
                    type="button"
                    class="btn btn-outline btn-sm"
                    onclick={() => openCreateTeamModal()}>チームを追加</button
                >
            </div>
        </div>
        <div
            class="rounded-3xl border border-base-300 bg-base-100 shadow-sm overflow-hidden"
        >
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
                    {#each visibleOrgNodes.entries() as [idx, n] (n.kind + "-" + n.uuid + "-" + idx)}
                        <div
                            class="flex items-center gap-2 px-4 py-3 hover:bg-base-200/40 transition-colors"
                            style="padding-left: {1 + n.depth * 1.2}rem"
                        >
                            {#if n.kind === "department"}
                                {#if (n.departmentCount && n.departmentCount > 0) || (n.teamCount && n.teamCount > 0)}
                                    <button
                                        type="button"
                                        class="btn btn-ghost btn-xs btn-square shrink-0"
                                        aria-label="チームを展開"
                                        onclick={() => toggleDeptExpand(n.uuid)}
                                    >
                                        <ChevronRight
                                            class="h-4 w-4 transition-transform {expandedDeptUuids.has(
                                                n.uuid,
                                            )
                                                ? 'rotate-90'
                                                : ''}"
                                        />
                                    </button>
                                {:else}
                                    <span class="w-6 shrink-0"></span>
                                {/if}
                                <Building2
                                    class="h-4 w-4 shrink-0 text-primary/80"
                                />
                                <button
                                    type="button"
                                    class="min-w-0 flex-1 text-left"
                                    onclick={() => openOrgModal(n)}
                                >
                                    <span class="font-medium">{n.name}</span>
                                    <span
                                        class="ml-2 text-sm text-base-content/60"
                                        >{n.headUserName || ""}
                                        {n.headUserRole || ""}</span
                                    >
                                </button>
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-xs"
                                    onclick={() => openCreateDeptModal(n.uuid)}
                                    >部署を追加</button
                                >
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-xs"
                                    onclick={() => openCreateTeamModal(n.uuid)}
                                    >チームを追加</button
                                >
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-xs"
                                    onclick={() => openOrgModal(n)}>編集</button
                                >
                            {:else}
                                <span class="w-6 shrink-0"></span>
                                <Users
                                    class="h-4 w-4 shrink-0 text-primary/80"
                                />
                                <button
                                    type="button"
                                    class="min-w-0 flex-1 text-left"
                                    onclick={() => openOrgModal(n)}
                                >
                                    <span class="font-medium">{n.name}</span>
                                    <span
                                        class="ml-2 text-sm text-base-content/60"
                                        >{n.leaderName || ""}
                                        {n.leaderRole || ""}</span
                                    >
                                </button>
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-xs"
                                    onclick={() => openOrgModal(n)}>編集</button
                                >
                            {/if}
                        </div>
                    {/each}
                </div>
            {/if}
        </div>
    </div>
{/if}

{#if tab === "roles"}
    <div class="mt-6 max-w-7xl">
        <div class="flex justify-end mb-4">
            <button
                type="button"
                class="btn btn-primary btn-sm"
                onclick={openCreateRoleModal}>役職を追加</button
            >
        </div>
        <div
            class="overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm"
        >
            {#if rolesLoading}
                <div class="p-6 space-y-3">
                    {#each Array(3) as _, i (i)}
                        <div class="skeleton h-10 w-full rounded-xl"></div>
                    {/each}
                </div>
            {:else if roles.length === 0}
                <div class="p-12 text-center text-base-content/50 text-sm">
                    役職が登録されていません
                </div>
            {:else}
                <table class="table">
                    <thead>
                        <tr class="text-base-content/70">
                            <th>役職名</th>
                            <th class="w-40 text-right">操作</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each roles as role (`role-${role.uuid}`)}
                            <tr>
                                <td class="font-medium">{role.name}</td>
                                <td class="text-right">
                                    <div class="flex justify-end gap-2">
                                        <button
                                            type="button"
                                            class="btn btn-ghost btn-sm"
                                            onclick={() =>
                                                openEditRoleModal(role)}
                                            >名称変更</button
                                        >
                                        <button
                                            type="button"
                                            class="btn btn-ghost btn-sm text-error"
                                            onclick={() => askDeleteRole(role)}
                                            aria-label="削除"
                                        >
                                            <Trash2 class="h-4 w-4" />
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            {/if}
        </div>
    </div>
{/if}

{#if tab === "remind"}
    <div class="mt-6 max-w-xl">
        <div
            class="rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-5"
        >
            <div role="alert" class="alert alert-info text-sm">
                <Info class="h-5 w-5 shrink-0" />
                <span
                    >設定の保存のみ行います。実際の通知送信はバックエンドジョブで別途実装されます。</span
                >
            </div>
            <label
                class="label cursor-pointer justify-start gap-3 rounded-2xl border border-base-300 px-4 py-3"
            >
                <input
                    type="checkbox"
                    class="toggle toggle-primary"
                    bind:checked={remind.remindAutoSend}
                />
                <span class="label-text font-medium"
                    >1on1の通知を自動送信する</span
                >
            </label>
            <div>
                <span class="label-text font-medium">通知タイミング</span>
                <div class="flex flex-col gap-2 mt-2">
                    <label
                        class="label cursor-pointer justify-start gap-3 rounded-2xl border border-base-300 px-4 py-3"
                    >
                        <input
                            type="radio"
                            name="remindMode"
                            class="radio radio-primary"
                            value="days_before"
                            bind:group={remind.remindMode}
                        />
                        <span class="label-text">実施予定日の N 日前</span>
                    </label>
                    <label
                        class="label cursor-pointer justify-start gap-3 rounded-2xl border border-base-300 px-4 py-3"
                    >
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
            {#if remind.remindMode === "days_before"}
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
            {#if remind.remindMode === "prev_weekday"}
                <label class="form-control">
                    <span class="label-text font-medium">曜日</span>
                    <select
                        class="select select-bordered"
                        bind:value={remind.remindWeekday}
                    >
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
                    bind:value={remind.remindTime}
                />
            </label>
            <button
                type="button"
                class="btn btn-primary w-full sm:w-auto"
                onclick={saveRemind}>保存</button
            >
        </div>
    </div>
{/if}

<dialog class="modal" class:modal-open={userModalOpen}>
    <div class="modal-box max-w-3xl flex flex-col p-0">
        <div
            class="flex justify-between items-start gap-2 p-6 pb-2 border-b border-base-200 shrink-0"
        >
            <div>
                <h3 class="font-bold text-lg">{userForm?.name}</h3>
                <p class="text-sm text-base-content/60 mt-0.5">権限・所属</p>
            </div>
            <button
                type="button"
                class="btn btn-sm btn-circle btn-ghost"
                aria-label="閉じる"
                onclick={() => (userModalOpen = false)}>×</button
            >
        </div>
        {#if userForm}
            <div
                class="overflow-y-auto px-6 py-4 space-y-5 flex justify-around gap-2"
            >
                <div class="flex flex-col w-1/3 items-around gap-2 h-full">
                    <div class="flex justify-center w-full">
                        <Icon
                            url={userAvatarByUri(userForm.imageUri)}
                            text={userForm.name}
                            height="200px"
                            width="200px"
                            alt={userForm.name}
                        ></Icon>
                    </div>
                    <div class="pt-4 flex flex-wrap gap-2 justify-around">
                        <button
                            type="button"
                            class="btn btn-ghost btn-sm text-error hover:bg-error/10"
                            onclick={deleteUserIcon}
                        >
                            <Trash2 class="h-4 w-4" />
                            アイコンを削除
                        </button>
                        {#if !userForm.suspended}
                            <button
                                type="button"
                                class="btn btn-warning btn-sm"
                                onclick={() => toggleSuspendInModal(false)}
                                >アカウント停止</button
                            >
                        {/if}
                        {#if userForm.suspended}
                            <button
                                type="button"
                                class="btn btn-success btn-sm"
                                onclick={() => toggleSuspendInModal(true)}
                                >停止解除</button
                            >
                        {/if}
                    </div>
                </div>
                <div class="flex flex-col w-2/3">
                    <div class="grid sm:grid-cols-2 gap-4">
                        <label class="form-control">
                            <span class="label-text font-medium">権限</span>
                            <select
                                class="select select-bordered select-sm mt-1"
                                bind:value={authoritySelect}
                                onchange={onAuthoritySelectChange}
                            >
                                {#each AUTHORITY_OPTIONS as opt (opt.value)}
                                    <option value={opt.value}
                                        >{opt.label}</option
                                    >
                                {/each}
                            </select>
                        </label>
                        <label class="form-control">
                            <span class="label-text font-medium">ロール</span>
                            <select
                                class="select select-bordered select-sm mt-1"
                                bind:value={systemRoleSelect}
                                onchange={onSystemRoleSelectChange}
                            >
                                {#each SYSTEM_ROLE_OPTIONS as opt (opt.value)}
                                    <option value={opt.value}
                                        >{opt.label}</option
                                    >
                                {/each}
                            </select>
                        </label>
                    </div>
                    <div class="flex flex-col mt-4 gap-1">
                        <div class="flex items-center justify-between mt-1">
                            <span class="label-text font-medium">所属部署</span>
                            <button
                                type="button"
                                class="btn btn-outline btn-xs"
                                onclick={addUserDepartmentRow}>行を追加</button
                            >
                        </div>
                        {#each userForm.departments as row, i (`user-dept-${i}`)}
                            <div class="flex items-center gap-1">
                                <select
                                    class="select select-bordered select-sm flex-1"
                                    bind:value={row.departmentUuid}
                                    onchange={onUserDepartmentChange}
                                >
                                    <option value="">部署を選択</option>
                                    {#each departmentChoices(userForm.departments, i) as d (`user-dept-opt-${i}-${d.uuid}`)}
                                        <option value={d.uuid}>{d.name}</option>
                                    {/each}
                                </select>
                                <select
                                    class="select select-bordered select-sm flex-1"
                                    bind:value={row.roleUuid}
                                >
                                    <option value="">役職を選択</option>
                                    {#each options.roles as r (`user-role-opt-${i}-${r.uuid}`)}
                                        <option value={r.uuid}>{r.name}</option>
                                    {/each}
                                </select>
                                <button
                                    type="button"
                                    class="btn btn-ghost btn-sm"
                                    onclick={() => removeUserDepartmentRow(i)}
                                    >削除</button
                                >
                            </div>
                        {/each}
                    </div>
                    {#if departmentsHaveTeams(userForm.departments)}
                        <div class="flex flex-col gap-1 mt-4">
                            <div
                                class="flex items-center justify-between gap-2"
                            >
                                <span class="label-text font-medium"
                                    >所属チーム</span
                                >
                                <button
                                    type="button"
                                    class="btn btn-outline btn-xs"
                                    onclick={addUserTeamRow}>行を追加</button
                                >
                            </div>
                            {#each userForm.teams as row, i (`user-team-${i}`)}
                                <div class="flex items-center gap-2">
                                    <select
                                        class="select select-bordered select-sm flex-1"
                                        bind:value={row.teamUuid}
                                    >
                                        <option value={0} selected
                                            >チームを選択</option
                                        >
                                        {#each teamChoices(userForm.departments, userForm.teams, i) as t (`user-team-opt-${i}-${t.uuid}`)}
                                            <option value={t.uuid}
                                                >{`${t.departmentName} ${t.name}`}</option
                                            >
                                        {/each}
                                    </select>
                                    <button
                                        type="button"
                                        class="btn btn-ghost btn-sm"
                                        onclick={() => removeUserTeamRow(i)}
                                        >削除</button
                                    >
                                </div>
                            {/each}
                        </div>
                    {/if}
                </div>
            </div>
        {/if}
        <div class="modal-action border-t border-base-200 p-4 shrink-0">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (userModalOpen = false)}>閉じる</button
            >
            <button type="button" class="btn btn-primary" onclick={saveUser}
                >保存</button
            >
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (userModalOpen = false)}
            >close</button
        >
    </form>
</dialog>

{#snippet orgUserListBlock(
    heading: string,
    pickButtonLabel: string,
    uuids: string[],
    editable: boolean,
    searchKind: "member" | "manager" | "head" | "leader",
    onPick: () => void,
)}
    {@const listSearch =
        searchKind === "member" || searchKind === "leader"
            ? orgMemberSearch
            : orgManagerSearch}
    {@const visibleUuids = filteredListedUuids(uuids, listSearch)}
    <div class="space-y-2">
        <div class="flex items-center justify-between gap-2">
            <span class="label-text font-medium">
                {heading}
                {#if uuids.length > 0}
                    <span class="text-base-content/60 font-normal ml-1"
                        >{visibleUuids.length}人</span
                    >
                {/if}
            </span>
            {#if editable}
                <button
                    type="button"
                    class="btn btn-outline btn-xs"
                    onclick={onPick}>{pickButtonLabel}</button
                >
            {/if}
        </div>
        {#if uuids.length === 0}
            <div
                class="rounded-xl border border-base-300 bg-base-200/30 px-3 py-2 text-sm text-base-content/50"
            >
                未設定
            </div>
        {:else}
            {#if searchKind === "manager"}
                <input
                    type="search"
                    class="input input-bordered input-sm w-full placeholder-current"
                    placeholder="氏名で検索"
                    bind:value={orgManagerSearch}
                    autocomplete="off"
                />
            {:else}
                <input
                    type="search"
                    class="input input-bordered input-sm w-full placeholder-current"
                    placeholder="氏名で検索"
                    bind:value={orgMemberSearch}
                    autocomplete="off"
                />
            {/if}
            <div
                class="border border-base-300 rounded-2xl max-h-44 overflow-y-auto p-2"
            >
                {#each visibleUuids as uid (`org-${searchKind}-${uid}`)}
                    {@const u = userByUuid(uid)}
                    <div
                        class="flex items-center gap-2 px-2 py-1.5 rounded-lg text-sm"
                    >
                        <Icon
                            url={u ? userAvatar(u) : null}
                            text={u?.name ?? userName(uid)}
                            width={28}
                            height={28}
                            alt={u?.name ?? ""}
                        ></Icon>
                        <span class="truncate">{u?.name ?? userName(uid)}</span>
                    </div>
                {:else}
                    <p
                        class="px-2 py-3 text-center text-sm text-base-content/50"
                    >
                        該当するユーザがありません
                    </p>
                {/each}
            </div>
        {/if}
    </div>
{/snippet}

<dialog class="modal" class:modal-open={orgModalOpen}>
    <div class="modal-box max-w-3xl flex flex-col p-0">
        <div
            class="flex justify-between items-start gap-2 p-6 pb-2 border-b border-base-200 shrink-0"
        >
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
                        onclick={() => (orgEditMode = true)}>編集</button
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
            <div class="overflow-y-auto px-6 py-4 flex-1 space-y-4">
                {#if orgDetail?.kind === "department"}
                    <div
                        class="rounded-xl border border-base-300 bg-base-200/30 px-3 py-2 text-sm"
                    >
                        <span class="text-base-content/60">責任者</span>
                        <span class="ml-2 font-medium"
                            >{orgDetail.headUserUuid
                                ? userName(orgDetail.headUserUuid)
                                : "未設定"}</span
                        >
                    </div>
                    <div
                        class="rounded-xl border border-base-300 bg-base-200/30 px-3 py-2 text-sm"
                    >
                        <span class="text-base-content/60">管理職</span>
                        <span class="ml-2 font-medium"
                            >{managerNames(orgDetail.managerUserUuids)}</span
                        >
                    </div>
                    {@render orgUserListBlock(
                        "メンバー",
                        "メンバーを選択",
                        plainDepartmentMemberUuids(orgDetail),
                        false,
                        "member",
                        () => openOrgPick("member"),
                    )}
                {/if}
                {#if orgDetail?.kind === "team"}
                    <div
                        class="rounded-xl border border-base-300 bg-base-200/30 px-3 py-2 text-sm"
                    >
                        <span class="text-base-content/60">チームリーダー</span>
                        <span class="ml-2 font-medium"
                            >{orgDetail.leaderUuid
                                ? userName(orgDetail.leaderUuid)
                                : "未設定"}</span
                        >
                    </div>
                    {@render orgUserListBlock(
                        "メンバー",
                        "メンバーを選択",
                        plainTeamMemberUuids(orgDetail),
                        false,
                        "member",
                        () => openOrgPick("member"),
                    )}
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
                {#if orgDetail?.kind === "department"}
                    <div class="space-y-2">
                        <span class="label-text font-medium">子部署</span>
                        <input
                            class="input input-bordered input-sm w-full placeholder-current"
                            placeholder="検索"
                            bind:value={orgChildSearch}
                        />
                        <div
                            class="border border-base-300 rounded-2xl max-h-44 overflow-y-auto p-2"
                        >
                            {#each filteredOrgDepartments() as d (`c-${d.uuid}`)}
                                {#if d.uuid !== orgDetail?.uuid && !d.isRoot}
                                    <label
                                        class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer"
                                    >
                                        <input
                                            type="checkbox"
                                            class="checkbox checkbox-sm"
                                            checked={isInList(
                                                orgEdit.childDepartmentUuids,
                                                d.uuid,
                                            )}
                                            onchange={() =>
                                                toggleOrgChild(d.uuid)}
                                        />
                                        <span>{d.name}</span>
                                    </label>
                                {/if}
                            {/each}
                        </div>
                    </div>
                {/if}
                {#if orgDetail?.kind === "department"}
                    <div class="space-y-4">
                        {@render orgUserListBlock(
                            "責任者",
                            "責任者を選択",
                            uuidList(orgEdit.headUserUuid),
                            true,
                            "head",
                            () => openOrgPick("head"),
                        )}
                        {@render orgUserListBlock(
                            "管理職",
                            "管理職を選択",
                            orgEdit.managerUserUuids,
                            true,
                            "manager",
                            () => openOrgPick("manager"),
                        )}
                        {@render orgUserListBlock(
                            "メンバー",
                            "メンバーを選択",
                            orgEdit.memberUserUuids,
                            true,
                            "member",
                            () => openOrgPick("member"),
                        )}
                        <div class="space-y-2">
                            <span class="label-text font-medium">チーム</span>
                            <input
                                class="input input-bordered input-sm w-full placeholder-current"
                                placeholder="検索"
                                bind:value={orgTeamSearch}
                            />
                            <div
                                class="border border-base-300 rounded-2xl max-h-36 overflow-y-auto p-2"
                            >
                                {#each filteredOrgTeams() as t (`t-${t.uuid}`)}
                                    <label
                                        class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-base-200 cursor-pointer"
                                    >
                                        <input
                                            type="checkbox"
                                            class="checkbox checkbox-sm"
                                            checked={isInList(
                                                orgEdit.teamUuids,
                                                t.uuid,
                                            )}
                                            onchange={() =>
                                                toggleOrgTeam(t.uuid)}
                                        />
                                        <span>{t.name}</span>
                                    </label>
                                {/each}
                            </div>
                        </div>
                    </div>
                {/if}
                {#if orgDetail?.kind === "team"}
                    {@render orgUserListBlock(
                        "チームリーダー",
                        "チームリーダーを選択",
                        uuidList(orgEdit.leaderUuid),
                        true,
                        "leader",
                        () => openOrgPick("leader"),
                    )}
                    {@render orgUserListBlock(
                        "メンバー",
                        "メンバーを選択",
                        orgEdit.memberUserUuids,
                        true,
                        "member",
                        () => openOrgPick("member"),
                    )}
                {/if}
                <div class="border-t border-base-300 pt-4">
                    <button
                        type="button"
                        class="btn btn-error btn-outline btn-sm"
                        onclick={askDeleteOrg}>削除</button
                    >
                </div>
            </div>
        {/if}
        {#if !orgEditMode}
            <div class="modal-action border-t border-base-200 p-4 shrink-0">
                <button
                    type="button"
                    class="btn btn-error btn-outline btn-sm mr-auto"
                    onclick={askDeleteOrg}>削除</button
                >
                <button
                    type="button"
                    class="btn btn-ghost"
                    onclick={() => (orgModalOpen = false)}>閉じる</button
                >
            </div>
        {/if}
        {#if orgEditMode}
            <div class="modal-action border-t border-base-200 p-4 shrink-0">
                <button
                    type="button"
                    class="btn btn-ghost"
                    onclick={() => (orgModalOpen = false)}>閉じる</button
                >
                <button type="button" class="btn btn-primary" onclick={saveOrg}
                    >保存</button
                >
            </div>
        {/if}
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (orgModalOpen = false)}
            >close</button
        >
    </form>
</dialog>

<dialog class="modal" class:modal-open={createUserOpen}>
    <div class="modal-box max-w-lg flex flex-col p-0">
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
            <div class="space-y-3">
                <div class="flex items-center justify-between gap-2">
                    <span class="label-text font-medium">部署</span>
                    <button
                        type="button"
                        class="btn btn-outline btn-xs"
                        onclick={addCreateDepartmentRow}>行を追加</button
                    >
                </div>
                {#each createUserForm.departments as row, i (`new-user-dept-${i}`)}
                    <div class="flex items-center gap-2">
                        <select
                            class="select select-bordered select-sm flex-1"
                            bind:value={row.departmentUuid}
                            onchange={onCreateDepartmentChange}
                        >
                            <option value="">部署を選択</option>
                            {#each departmentChoices(createUserForm.departments, i) as d (`new-dept-opt-${i}-${d.uuid}`)}
                                <option value={d.uuid}>{d.name}</option>
                            {/each}
                        </select>
                        <select
                            class="select select-bordered select-sm flex-1"
                            bind:value={row.roleUuid}
                        >
                            <option value="">役職を選択</option>
                            {#each options.roles as r (`new-role-opt-${i}-${r.uuid}`)}
                                <option value={r.uuid}>{r.name}</option>
                            {/each}
                        </select>
                        <button
                            type="button"
                            class="btn btn-ghost btn-sm"
                            onclick={() => removeCreateDepartmentRow(i)}
                            >削除</button
                        >
                    </div>
                {/each}
            </div>
            {#if departmentsHaveTeams(createUserForm.departments)}
                <div class="space-y-3">
                    <div class="flex items-center justify-between gap-2">
                        <span class="label-text font-medium">チーム</span>
                        <button
                            type="button"
                            class="btn btn-outline btn-xs"
                            onclick={addCreateTeamRow}>行を追加</button
                        >
                    </div>
                    {#each createUserForm.teams as row, i (`new-user-team-${i}`)}
                        <div class="flex items-center gap-2">
                            <select
                                class="select select-bordered select-sm flex-1"
                                bind:value={row.teamUuid}
                            >
                                <option value={0} selected>チームを選択</option>
                                {#each teamChoices(createUserForm.departments, createUserForm.teams, i) as t (`new-team-opt-${i}-${t.uuid}`)}
                                    <option value={t.uuid}>{t.name}</option>
                                {/each}
                            </select>
                            <button
                                type="button"
                                class="btn btn-ghost btn-sm"
                                onclick={() => removeCreateTeamRow(i)}
                                >削除</button
                            >
                        </div>
                    {/each}
                </div>
            {/if}
            <label class="form-control"
                ><span class="label-text font-medium"
                    >パスワード（任意・Basic認証用）</span
                ><input
                    type="password"
                    class="input input-bordered"
                    bind:value={createUserForm.password}
                    autocomplete="new-password"
                /></label
            >
        </div>
        <div class="modal-action border-t border-base-200 p-4 shrink-0">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (createUserOpen = false)}>キャンセル</button
            >
            <button
                type="button"
                class="btn btn-primary"
                onclick={submitCreateUser}>作成</button
            >
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (createUserOpen = false)}
            >close</button
        >
    </form>
</dialog>

<dialog class="modal" class:modal-open={createDeptOpen}>
    <div class="modal-box max-w-3xl flex flex-col p-0">
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
                <select
                    class="select select-bordered"
                    bind:value={createDeptForm.parentDepartmentUuid}
                >
                    <option value="">なし（ルート）</option>
                    {#each options.departments as d (`parent-${d.uuid}`)}
                        <option value={d.uuid}>{d.name}</option>
                    {/each}
                </select>
            </label>
            {@render orgUserListBlock(
                "責任者",
                "責任者を選択",
                uuidList(createDeptForm.headUserUuid),
                true,
                "head",
                () => openOrgPick("head", "createDept"),
            )}
            {@render orgUserListBlock(
                "管理職",
                "管理職を選択",
                createDeptForm.managerUserUuids,
                true,
                "manager",
                () => openOrgPick("manager", "createDept"),
            )}
            {@render orgUserListBlock(
                "メンバー",
                "メンバーを選択",
                createDeptForm.memberUserUuids,
                true,
                "member",
                () => openOrgPick("member", "createDept"),
            )}
        </div>
        <div class="modal-action border-t border-base-200 p-4 shrink-0">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (createDeptOpen = false)}>キャンセル</button
            >
            <button
                type="button"
                class="btn btn-primary"
                onclick={submitCreateDept}>作成</button
            >
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (createDeptOpen = false)}
            >close</button
        >
    </form>
</dialog>

<dialog class="modal" class:modal-open={confirmOpen}>
    <div class="modal-box">
        <h3 class="font-bold">{confirmTitle}</h3>
        <p class="py-4 text-sm text-base-content/80">{confirmMessage}</p>
        <div class="modal-action">
            <button type="button" class="btn btn-ghost" onclick={cancelConfirm}
                >キャンセル</button
            >
            <button
                type="button"
                class="btn"
                class:btn-error={confirmDanger}
                class:btn-primary={!confirmDanger}
                onclick={runConfirm}
                >{confirmDanger ? "確定" : "変更する"}</button
            >
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={cancelConfirm}>close</button>
    </form>
</dialog>

<dialog class="modal" class:modal-open={createTeamOpen}>
    <div class="modal-box max-w-3xl flex flex-col p-0">
        <div class="p-6 pb-2 border-b border-base-200 shrink-0">
            <h3 class="font-bold text-lg">チームを追加</h3>
        </div>
        <div class="overflow-y-auto px-6 py-4 space-y-4 flex-1">
            <label class="form-control">
                <span class="label-text font-medium">部署</span>
                <select
                    class="select select-bordered"
                    value={createTeamForm.departmentUuid}
                    onchange={onCreateTeamDepartmentChange}
                >
                    <option value="" disabled>部署を選択</option>
                    {#each options.departments as d (`create-team-dept-${d.uuid}`)}
                        <option value={d.uuid}>{d.name}</option>
                    {/each}
                </select>
            </label>
            <label class="form-control">
                <span class="label-text font-medium">チーム名</span>
                <input
                    class="input input-bordered"
                    bind:value={createTeamForm.name}
                />
            </label>
            {@render orgUserListBlock(
                "チームリーダー",
                "チームリーダーを選択",
                uuidList(createTeamForm.leaderUuid),
                true,
                "leader",
                () => openOrgPick("leader", "createTeam"),
            )}
            {@render orgUserListBlock(
                "メンバー",
                "メンバーを選択",
                createTeamForm.memberUserUuids,
                true,
                "member",
                () => openOrgPick("member", "createTeam"),
            )}
        </div>
        <div class="modal-action border-t border-base-200 p-4 shrink-0">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (createTeamOpen = false)}>キャンセル</button
            >
            <button
                type="button"
                class="btn btn-primary"
                onclick={saveCreateTeam}>作成</button
            >
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (createTeamOpen = false)}
            >close</button
        >
    </form>
</dialog>

<dialog class="modal" class:modal-open={createRoleOpen}>
    <div class="modal-box max-w-md">
        <h3 class="font-bold text-lg">役職を追加</h3>
        <label class="form-control mt-4">
            <span class="label-text font-medium">役職名</span>
            <input
                class="input input-bordered"
                bind:value={createRoleForm.name}
            />
        </label>
        <div class="modal-action">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (createRoleOpen = false)}>キャンセル</button
            >
            <button
                type="button"
                class="btn btn-primary"
                onclick={saveCreateRole}>追加</button
            >
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (createRoleOpen = false)}
            >close</button
        >
    </form>
</dialog>

<dialog class="modal" class:modal-open={editRoleOpen}>
    <div class="modal-box max-w-md">
        <h3 class="font-bold text-lg">役職名の変更</h3>
        <label class="form-control mt-4">
            <span class="label-text font-medium">役職名</span>
            <input
                class="input input-bordered"
                bind:value={editRoleForm.name}
            />
        </label>
        <div class="modal-action">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (editRoleOpen = false)}>キャンセル</button
            >
            <button type="button" class="btn btn-primary" onclick={saveEditRole}
                >保存</button
            >
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (editRoleOpen = false)}
            >close</button
        >
    </form>
</dialog>

<MemberPickModal
    open={orgMemberPickOpen}
    title={orgMemberPickTitle}
    users={orgMemberPickUsers}
    multiple={orgMemberPickMultiple}
    selectedUuids={orgMemberPickSelectedUuids}
    excludeUuids={orgMemberPickExcludeUuids}
    onApply={applyOrgMemberPick}
    onClose={() => (orgMemberPickOpen = false)}
/>

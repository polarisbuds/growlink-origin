import { showToast, refreshPageIcons } from "/static/js/notice.js";

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

function toggleInList(list, id) {
  const n = Number(id);
  const i = list.indexOf(n);
  if (i >= 0) list.splice(i, 1);
  else list.push(n);
  return [...list];
}

function filterByQuery(items, query, labelKey = "name") {
  const q = (query || "").trim().toLowerCase();
  if (!q) return items;
  return items.filter((item) => String(item[labelKey] || "").toLowerCase().includes(q));
}

document.addEventListener("alpine:init", () => {
  Alpine.data("managementPage", () => ({
    tab: "users",
    users: [],
    orgTree: [],
    options: { departments: [], teams: [], roles: [] },
    selectedUserIds: [],
    userSearch: "",
    usersLoading: true,
    orgLoading: true,
    remind: { remindAutoSend: false, remindWeekday: 1, remindTime: "09:00" },
    weekdays: WEEKDAYS,
    authorityOptions: AUTHORITY_OPTIONS,
    systemRoleOptions: SYSTEM_ROLE_OPTIONS,
    operationOptions: OPERATION_OPTIONS,
    targetSearch: "",
    orgChildSearch: "",
    orgTeamSearch: "",
    orgUserSearch: "",
    orgMemberSearch: "",
    userModalOpen: false,
    userForm: null,
    orgModalOpen: false,
    orgDetail: null,
    orgEditMode: false,
    orgEdit: { name: "", childIds: [], teamIds: [], userIds: [], memberIds: [], leaderId: 0 },
    confirmOpen: false,
    confirmTitle: "",
    confirmMessage: "",
    confirmAction: null,

    get filteredUsers() {
      const q = this.userSearch.trim().toLowerCase();
      if (!q) return this.users;
      return this.users.filter(
        (u) =>
          (u.name || "").toLowerCase().includes(q) ||
          (u.email || "").toLowerCase().includes(q)
      );
    },

    get allUsersSelected() {
      const visible = this.filteredUsers;
      return visible.length > 0 && visible.every((u) => this.selectedUserIds.includes(u.id));
    },

    async init() {
      await Promise.all([this.loadUsers(), this.loadOptions(), this.loadOrgTree(), this.loadRemind()]);
      refreshPageIcons();
    },

    async loadUsers() {
      this.usersLoading = true;
      const res = await fetch("/api/v1/management/users");
      this.users = res.ok ? await res.json() : [];
      this.usersLoading = false;
    },
    async loadOptions() {
      const res = await fetch("/api/v1/management/options");
      if (res.ok) this.options = await res.json();
    },
    async loadOrgTree() {
      this.orgLoading = true;
      const res = await fetch("/api/v1/management/org/tree");
      this.orgTree = res.ok ? await res.json() : [];
      this.orgLoading = false;
      this.$nextTick(() => refreshPageIcons());
    },
    async loadRemind() {
      const res = await fetch("/api/v1/management/settings/remind");
      if (res.ok) {
        const data = await res.json();
        this.remind = { ...data, remindTime: (data.remindTime || "09:00").slice(0, 5) };
      }
    },

    authorityLabel(v) {
      return AUTHORITY_OPTIONS.find((o) => o.value === v)?.label || v;
    },
    systemRoleLabel(v) {
      return SYSTEM_ROLE_OPTIONS.find((o) => o.value === v)?.label || v;
    },
    userName(id) {
      const u = this.users.find((x) => x.id === Number(id));
      return u ? u.name : "";
    },
    userAvatar(u) {
      return u.imageUri ? `/uploads/${u.imageUri}` : null;
    },
    userInitial(u) {
      return (u.name || "?").slice(0, 1);
    },

    toggleAllUsers(e) {
      const ids = this.filteredUsers.map((u) => u.id);
      this.selectedUserIds = e.target.checked ? ids : [];
    },

    filteredDepartments() {
      return filterByQuery(this.options.departments || [], this.targetSearch);
    },
    filteredTeams() {
      return filterByQuery(this.options.teams || [], this.targetSearch);
    },
    filteredUsersForTarget() {
      return filterByQuery(this.users, this.targetSearch);
    },
    filteredOrgDepartments() {
      return filterByQuery(this.options.departments || [], this.orgChildSearch);
    },
    filteredOrgTeams() {
      return filterByQuery(this.options.teams || [], this.orgTeamSearch);
    },
    filteredOrgUsers() {
      return filterByQuery(this.users, this.orgUserSearch);
    },
    filteredOrgMembers() {
      return filterByQuery(this.users, this.orgMemberSearch);
    },

    isInList(list, id) {
      return (list || []).map(Number).includes(Number(id));
    },
    toggleUserTarget(id) {
      if (!this.userForm) return;
      this.userForm.targetIds = toggleInList(this.userForm.targetIds || [], id);
    },
    toggleOrgChild(id) {
      this.orgEdit.childIds = toggleInList(this.orgEdit.childIds || [], id);
    },
    toggleOrgTeam(id) {
      this.orgEdit.teamIds = toggleInList(this.orgEdit.teamIds || [], id);
    },
    toggleOrgUser(id) {
      this.orgEdit.userIds = toggleInList(this.orgEdit.userIds || [], id);
    },
    toggleOrgMember(id) {
      this.orgEdit.memberIds = toggleInList(this.orgEdit.memberIds || [], id);
      if (!this.orgEdit.memberIds.includes(Number(this.orgEdit.leaderId))) {
        this.orgEdit.leaderId = this.orgEdit.memberIds[0] || 0;
      }
    },

    openUserModal(u) {
      this.targetSearch = "";
      this.userForm = {
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
      this.userModalOpen = true;
      this.$nextTick(() => refreshPageIcons());
    },

    async saveUser() {
      const body = {
        authority: this.userForm.authority,
        systemRole: this.userForm.systemRole,
        operationTarget: this.userForm.operationTarget,
        targetIds: (this.userForm.targetIds || []).map(Number),
        departmentId: this.userForm.departmentId,
        roleId: this.userForm.roleId,
        suspended: this.userForm.suspended,
      };
      const res = await fetch(`/api/v1/management/users/${this.userForm.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        showToast(err.error || "保存に失敗しました", "error");
        return;
      }
      this.userModalOpen = false;
      await this.loadUsers();
      showToast("ユーザ情報を保存しました");
    },

    async deleteUserIcon() {
      const res = await fetch(`/api/v1/management/users/${this.userForm.id}/icon`, { method: "DELETE" });
      if (!res.ok) {
        showToast("アイコンの削除に失敗しました", "error");
        return;
      }
      await this.loadUsers();
      showToast("アイコンを削除しました");
    },

    toggleSuspend(unsuspend) {
      this.userForm.suspended = !unsuspend;
    },

    askSuspend() {
      this.confirmTitle = "アカウント停止";
      this.confirmMessage = `選択した ${this.selectedUserIds.length} 件のアカウントを停止しますか？`;
      this.confirmAction = async () => {
        const res = await fetch("/api/v1/management/users/suspend", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ ids: this.selectedUserIds }),
        });
        if (!res.ok) {
          showToast("停止に失敗しました", "error");
          return;
        }
        this.selectedUserIds = [];
        await this.loadUsers();
        showToast("アカウントを停止しました");
      };
      this.confirmOpen = true;
    },

    async openOrgModal(n) {
      this.orgEditMode = false;
      this.orgChildSearch = "";
      this.orgTeamSearch = "";
      this.orgUserSearch = "";
      this.orgMemberSearch = "";
      const url = n.kind === "team" ? `/api/v1/management/teams/${n.id}` : `/api/v1/management/departments/${n.id}`;
      const res = await fetch(url);
      if (!res.ok) {
        showToast("詳細の取得に失敗しました", "error");
        return;
      }
      this.orgDetail = await res.json();
      this.orgEdit = {
        name: this.orgDetail.name,
        childIds: (this.orgDetail.childIds || []).map(Number),
        teamIds: (this.orgDetail.teamIds || []).map(Number),
        userIds: (this.orgDetail.userIds || []).map(Number),
        memberIds: (this.orgDetail.memberIds || []).map(Number),
        leaderId: this.orgDetail.leaderId || 0,
      };
      this.orgModalOpen = true;
      this.$nextTick(() => refreshPageIcons());
    },

    async saveOrg() {
      let res;
      if (this.orgDetail.kind === "team") {
        res = await fetch(`/api/v1/management/teams/${this.orgDetail.id}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            name: this.orgEdit.name,
            memberIds: this.orgEdit.memberIds.map(Number),
            leaderId: this.orgEdit.leaderId,
          }),
        });
      } else {
        res = await fetch(`/api/v1/management/departments/${this.orgDetail.id}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            name: this.orgEdit.name,
            childDepartmentIds: this.orgEdit.childIds.map(Number),
            teamIds: this.orgEdit.teamIds.map(Number),
            userIds: this.orgEdit.userIds.map(Number),
          }),
        });
      }
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        showToast(err.error || "保存に失敗しました", "error");
        return;
      }
      this.orgModalOpen = false;
      await this.loadOrgTree();
      showToast("保存しました");
    },

    askDeleteOrg() {
      this.confirmTitle = "削除の確認";
      this.confirmMessage = `${this.orgDetail.name} を削除しますか？`;
      this.confirmAction = async () => {
        const url =
          this.orgDetail.kind === "team"
            ? `/api/v1/management/teams/${this.orgDetail.id}`
            : `/api/v1/management/departments/${this.orgDetail.id}`;
        const res = await fetch(url, { method: "DELETE" });
        if (!res.ok) {
          const err = await res.json().catch(() => ({}));
          showToast(err.error || "削除に失敗しました", "error");
          return;
        }
        this.orgModalOpen = false;
        await this.loadOrgTree();
        showToast("削除しました");
      };
      this.confirmOpen = true;
    },

    async saveRemind() {
      const res = await fetch("/api/v1/management/settings/remind", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(this.remind),
      });
      if (!res.ok) {
        showToast("保存に失敗しました", "error");
        return;
      }
      showToast("通知設定を保存しました");
    },

    switchTab(name) {
      this.tab = name;
      this.$nextTick(() => refreshPageIcons());
    },

    async runConfirm() {
      this.confirmOpen = false;
      if (this.confirmAction) await this.confirmAction();
      this.confirmAction = null;
    },

    closeModals() {
      if (this.confirmOpen) {
        this.confirmOpen = false;
        return;
      }
      this.userModalOpen = false;
      this.orgModalOpen = false;
    },
  }));
});

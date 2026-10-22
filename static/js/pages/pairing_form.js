import { filterItems, createPickController } from "../searchable-select.js";
import { showToast, refreshPageIcons } from "../notice.js";

document.addEventListener("alpine:init", () => {
  Alpine.data("pairingForm", (mode) => {
    let initial = { mentorId: 0, targets: [], uuid: "" };
    const el = document.getElementById("pairing-json");
    if (el) {
      const p = JSON.parse(el.textContent || "{}");
      initial = {
        uuid: p.uuid || "",
        mentorId: p.mentor?.id || 0,
        targets: (p.targets || []).map((t) => ({
          kind: t.kind,
          targetId: t.targetId,
          label: t.label || "",
        })),
      };
    }

    const form = {
      mentorId: initial.mentorId,
    };

    const mentorPick = createPickController({
      items: () => [],
      getId: (u) => u.id,
      labelOf: (u) => `${u.name} (${u.email})`,
      getSelectedId: () => form.mentorId,
      setSelectedId: (id) => {
        form.mentorId = id;
      },
    });

    const userTargetPick = createPickController({
      items: () => [],
      getId: (u) => u.id,
      labelOf: (u) => `${u.name} (${u.email})`,
      getSelectedId: () => 0,
      setSelectedId: () => {},
    });

    return {
      mode,
      form,
      targets: [...initial.targets],
      users: [],
      teams: [],
      departments: [],
      pairingUUID: initial.uuid,
      teamPickId: 0,
      departmentPickId: 0,
      mentorPick,
      userTargetPick,
      targetKey(kind, targetId) {
        return `${kind}:${targetId}`;
      },
      userLabel(id) {
        const u = this.users.find((x) => x.id === id);
        return u ? u.name : `#${id}`;
      },
      teamLabel(id) {
        const t = this.teams.find((x) => x.id === id);
        return t ? t.name : `#${id}`;
      },
      departmentLabel(id) {
        const d = this.departments.find((x) => x.id === id);
        return d ? d.name : `#${id}`;
      },
      addTarget(kind, targetId, label) {
        const key = this.targetKey(kind, targetId);
        if (this.targets.some((t) => this.targetKey(t.kind, t.targetId) === key)) return;
        this.targets.push({ kind, targetId, label });
      },
      mentorPickFiltered() {
        return filterItems(this.users, this.mentorPick.query, (u) => `${u.name} (${u.email})`);
      },
      userTargetPickFiltered() {
        const picked = new Set(this.targets.filter((t) => t.kind === "user").map((t) => t.targetId));
        const list = this.users.filter((u) => !picked.has(u.id));
        return filterItems(list, this.userTargetPick.query, (u) => `${u.name} (${u.email})`);
      },
      pickUserTarget(u) {
        this.addTarget("user", u.id, `${u.name} (${u.email})`);
        this.userTargetPick.clear();
      },
      async init() {
        const res = await fetch("/api/v1/planning/directory");
        if (res.ok) {
          const data = await res.json();
          this.users = data.users || [];
          this.teams = data.teams || [];
          this.departments = data.departments || [];
          for (const t of this.targets) {
            if (t.label) continue;
            if (t.kind === "user") t.label = this.userLabel(t.targetId);
            if (t.kind === "team") t.label = this.teamLabel(t.targetId);
            if (t.kind === "department") t.label = this.departmentLabel(t.targetId);
          }
        }
        window.addEventListener("mousedown", (ev) => {
          const t = ev.target;
          if (!(t instanceof Element)) return;
          if (!t.closest('[data-pick-root="mentor"]')) this.mentorPick.close();
          if (!t.closest('[data-pick-root="user-target"]')) this.userTargetPick.close();
        });
        refreshPageIcons();
      },
      removeTarget(index) {
        this.targets.splice(index, 1);
      },
      addTeamTarget() {
        const id = Number(this.teamPickId);
        if (!id) return;
        this.addTarget("team", id, this.teamLabel(id));
        this.teamPickId = 0;
      },
      addDepartmentTarget() {
        const id = Number(this.departmentPickId);
        if (!id) return;
        this.addTarget("department", id, this.departmentLabel(id));
        this.departmentPickId = 0;
      },
      targetKindLabel(kind) {
        if (kind === "user") return "ユーザー";
        if (kind === "team") return "チーム";
        if (kind === "department") return "部署";
        return kind;
      },
      payload() {
        return {
          mentorId: this.form.mentorId,
          targets: this.targets.map((t) => ({ kind: t.kind, targetId: t.targetId })),
        };
      },
      async doSave() {
        if (!this.form.mentorId || !this.targets.length) {
          showToast("メンターと対象を選択してください", "error");
          return;
        }
        const body = this.payload();
        const res =
          this.mode === "create"
            ? await fetch("/api/v1/pairings", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) })
            : await fetch(`/api/v1/pairings/${this.pairingUUID}`, { method: "PUT", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) });
        if (!res.ok) {
          showToast("保存に失敗しました", "error");
          return;
        }
        showToast("保存しました");
        location.href = "/planning/pairs";
      },
      requestSave() {
        this.doSave();
      },
    };
  });
});

import { createPickController } from "../searchable-select.js";
import { showToast, refreshPageIcons } from "../notice.js";

const WEEKDAYS = [
  { value: 0, label: "日曜" },
  { value: 1, label: "月曜" },
  { value: 2, label: "火曜" },
  { value: 3, label: "水曜" },
  { value: 4, label: "木曜" },
  { value: 5, label: "金曜" },
  { value: 6, label: "土曜" },
];

document.addEventListener("alpine:init", () => {
  Alpine.data("pairingForm", (mode) => {
    let initial = { mentorId: 0, menteeIds: [], recurWeek: null, recurWeekday: null, uuid: "", sharerUserIds: [] };
    const el = document.getElementById("pairing-json");
    if (el) {
      const p = JSON.parse(el.textContent || "{}");
      const menteeIds = (p.mentees || []).map((m) => m.id);
      if (!menteeIds.length && p.mentee?.id) menteeIds.push(p.mentee.id);
      initial = {
        uuid: p.uuid || "",
        mentorId: p.mentor?.id || 0,
        menteeIds,
        recurWeek: p.recurWeek,
        recurWeekday: p.recurWeekday,
        sharerUserIds: p.sharerUserIds || [],
      };
    }

    const users = { list: [] };
    const form = {
      mentorId: initial.mentorId,
      menteeIds: [...initial.menteeIds],
      recurWeek: initial.recurWeek,
      recurWeekday: initial.recurWeekday,
    };

    const mentorPick = createPickController({
      items: () => users.list,
      getId: (u) => u.id,
      labelOf: (u) => `${u.name} (${u.email})`,
      getSelectedId: () => form.mentorId,
      setSelectedId: (id) => {
        form.mentorId = id;
      },
    });
    const menteePick = createPickController({
      items: () => users.list.filter((u) => !form.menteeIds.includes(u.id)),
      getId: (u) => u.id,
      labelOf: (u) => `${u.name} (${u.email})`,
      getSelectedId: () => 0,
      setSelectedId: (id) => {
        if (id && !form.menteeIds.includes(id)) form.menteeIds.push(id);
        menteePick.clear();
      },
    });

    return {
      mode,
      weekdays: WEEKDAYS,
      teams: [],
      form,
      pairingUUID: initial.uuid,
      sharerIdsText: (initial.sharerUserIds || []).join(","),
      teamPickId: 0,
      confirmOpen: false,
      confirmTitle: "",
      confirmMessage: "",
      pendingSave: null,
      mentorPick,
      menteePick,
      get users() {
        return users.list;
      },
      userLabel(id) {
        const u = users.list.find((x) => x.id === id);
        return u ? u.name : `#${id}`;
      },
      async init() {
        const res = await fetch("/api/v1/planning/directory");
        if (res.ok) {
          const data = await res.json();
          users.list = data.users || [];
          this.teams = data.teams || [];
        }
        window.addEventListener("mousedown", (ev) => {
          const t = ev.target;
          if (!(t instanceof Element)) return;
          if (!t.closest('[data-pick-root="mentor"]')) mentorPick.close();
          if (!t.closest('[data-pick-root="mentee"]')) menteePick.close();
        });
        refreshPageIcons();
      },
      removeMentee(id) {
        form.menteeIds = form.menteeIds.filter((x) => x !== id);
      },
      addTeamMembers() {
        const team = this.teams.find((t) => t.id === Number(this.teamPickId));
        if (!team) return;
        for (const m of team.members || []) {
          if (!form.menteeIds.includes(m.id)) form.menteeIds.push(m.id);
        }
        this.teamPickId = 0;
      },
      scheduleLabel() {
        if (form.recurWeek == null) return "繰り返しなし";
        const wd = WEEKDAYS.find((d) => d.value === form.recurWeekday);
        return `第${form.recurWeek}週・${wd?.label || "—"}`;
      },
      payload() {
        const recurWeek = form.recurWeek || null;
        const recurWeekday = form.recurWeekday ?? null;
        return {
          mentorId: form.mentorId,
          menteeIds: form.menteeIds,
          recurWeek: recurWeek === 0 ? null : recurWeek,
          recurWeekday: recurWeek == null ? null : recurWeekday,
          sharerUserIds: (this.sharerIdsText || "")
            .split(",")
            .map((s) => parseInt(s.trim(), 10))
            .filter((n) => !Number.isNaN(n)),
        };
      },
      requestSave() {
        if (!form.mentorId || !form.menteeIds.length) {
          showToast("メンターとメンティを選択してください", "error");
          return;
        }
        if (form.recurWeek != null) {
          this.confirmTitle = "定期スケジュールの確認";
          this.confirmMessage = `このペアに ${this.scheduleLabel()} の繰り返しを設定します。よろしいですか？`;
          this.pendingSave = () => this.doSave();
          this.confirmOpen = true;
          return;
        }
        this.doSave();
      },
      async runConfirm() {
        this.confirmOpen = false;
        if (this.pendingSave) await this.pendingSave();
        this.pendingSave = null;
      },
      cancelConfirm() {
        this.confirmOpen = false;
        this.pendingSave = null;
      },
      async doSave() {
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
    };
  });
});

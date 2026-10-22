import { createPickController, filterItems } from "../searchable-select.js";

function todayISO() {
  const d = new Date();
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

function currentMonthISO() {
  const d = new Date();
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  return `${y}-${m}`;
}

function pairingSummary(p) {
  const mentor = p.mentor?.name || "—";
  const targetLabels = (p.targets || []).map((t) => t.label).filter(Boolean);
  const targets = targetLabels.length ? targetLabels.join("、") : p.mentee?.name || "—";
  return `${mentor} / ${targets}`;
}

document.addEventListener("alpine:init", () => {
  const pairingsEl = document.getElementById("pairings-json");
  const templatesEl = document.getElementById("templates-json");
  const usersEl = document.getElementById("users-json");
  const pairings = pairingsEl ? JSON.parse(pairingsEl.textContent || "[]") : [];
  const templates = templatesEl ? JSON.parse(templatesEl.textContent || "[]") : [];
  const users = usersEl ? JSON.parse(usersEl.textContent || "[]") : [];

  Alpine.data("meetingNewPage", () => {
    const form = {
      pairingId: 0,
      templateId: templates[0]?.id || 0,
      implementedAt: todayISO(),
      baseMonth: currentMonthISO(),
      week: 2,
      weekday: 3,
    };
    const pairPick = createPickController({
      items: () => pairings,
      getId: (p) => p.id,
      labelOf: pairingSummary,
      getSelectedId: () => form.pairingId,
      setSelectedId: (id) => {
        form.pairingId = id;
      },
    });
    const templatePick = createPickController({
      items: () => templates,
      getId: (t) => t.id,
      labelOf: (t) => t.title || "",
      getSelectedId: () => form.templateId,
      setSelectedId: (id) => {
        form.templateId = id;
      },
    });
    const sharerPick = createPickController({
      items: () => users,
      getId: (u) => u.id,
      labelOf: (u) => `${u.name} (${u.email})`,
      getSelectedId: () => 0,
      setSelectedId: () => {},
    });

    return {
      mode: "single",
      pairings,
      templates,
      users,
      sharerUserIds: [],
      submitting: false,
      error: "",
      form,
      pairPick,
      templatePick,
      sharerPick,
      init() {
        window.addEventListener("mousedown", (ev) => {
          const t = ev.target;
          if (!(t instanceof Element)) return;
          if (!t.closest('[data-pick-root="pair"]')) this.pairPick.close();
          if (!t.closest('[data-pick-root="template"]')) this.templatePick.close();
          if (!t.closest('[data-pick-root="sharer"]')) this.sharerPick.close();
        });
      },
      closePicks() {
        this.pairPick.close();
        this.templatePick.close();
        this.sharerPick.close();
      },
      pairingLabel(p) {
        return pairingSummary(p);
      },
      sharerPickFiltered() {
        const picked = new Set(this.sharerUserIds);
        const list = this.users.filter((u) => !picked.has(u.id));
        return filterItems(list, this.sharerPick.query, (u) => `${u.name} (${u.email})`);
      },
      userLabel(id) {
        const u = this.users.find((x) => x.id === id);
        return u ? u.name : `#${id}`;
      },
      addSharer(u) {
        if (!u?.id || this.sharerUserIds.includes(u.id)) return;
        this.sharerUserIds.push(u.id);
        this.sharerPick.clear();
      },
      removeSharer(id) {
        this.sharerUserIds = this.sharerUserIds.filter((x) => x !== id);
      },
      sharerPayload() {
        return this.sharerUserIds.length ? this.sharerUserIds : undefined;
      },
      submitLabel() {
        if (this.submitting) return "作成中…";
        if (this.mode === "single") return "作成";
        if (this.mode === "all") return "全ペアに作成";
        return "全ペアに定期作成";
      },
      async submit() {
        this.error = "";
        if (!this.form.templateId) {
          this.error = "テンプレートを選択してください。";
          return;
        }
        this.submitting = true;
        try {
          let url;
          let body;
          const sharerUserIds = this.sharerUserIds;
          if (this.mode === "single") {
            if (!this.form.pairingId) {
              this.error = "ペアを選択してください。";
              return;
            }
            if (!this.form.implementedAt) {
              this.error = "実施日を入力してください。";
              return;
            }
            url = "/api/v1/meetings";
            body = {
              pairingId: this.form.pairingId,
              templateId: this.form.templateId,
              implementedAt: this.form.implementedAt,
              sharerUserIds,
            };
          } else if (this.mode === "all") {
            if (!this.form.implementedAt) {
              this.error = "実施日を入力してください。";
              return;
            }
            url = "/api/v1/meetings/all";
            body = {
              templateId: this.form.templateId,
              implementedAt: this.form.implementedAt,
              sharerUserIds,
            };
          } else {
            if (!this.form.baseMonth) {
              this.error = "基準月を入力してください。";
              return;
            }
            if (this.form.week < 1 || this.form.week > 5) {
              this.error = "第N週は 1〜5 で指定してください。";
              return;
            }
            url = "/api/v1/meetings/monthly";
            body = {
              templateId: this.form.templateId,
              baseDate: `${this.form.baseMonth}-01`,
              week: this.form.week,
              weekday: this.form.weekday,
              sharerUserIds,
            };
          }
          const res = await fetch(url, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(body),
          });
          const data = await res.json().catch(() => ({}));
          if (!res.ok) {
            this.error = data.error || "作成に失敗しました。";
            return;
          }
          location.href = "/planning/meetings";
        } finally {
          this.submitting = false;
        }
      },
    };
  });
});

import { createPickController } from "../searchable-select.js";

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

document.addEventListener("alpine:init", () => {
  const pairingsEl = document.getElementById("pairings-json");
  const templatesEl = document.getElementById("templates-json");
  const pairings = pairingsEl ? JSON.parse(pairingsEl.textContent || "[]") : [];
  const templates = templatesEl ? JSON.parse(templatesEl.textContent || "[]") : [];

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
      labelOf: (p) => {
        const mentor = p.mentor?.name || "—";
        const mentee = p.mentee?.name || "—";
        return `${mentor} / ${mentee}`;
      },
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

    return {
      mode: "single",
      pairings,
      templates,
      submitting: false,
      error: "",
      form,
      pairPick,
      templatePick,
      init() {
        window.addEventListener("mousedown", (ev) => {
          const t = ev.target;
          if (!(t instanceof Element)) return;
          if (!t.closest('[data-pick-root="pair"]')) pairPick.close();
          if (!t.closest('[data-pick-root="template"]')) templatePick.close();
        });
      },
      closePicks() {
        pairPick.close();
        templatePick.close();
      },
      pairingLabel(p) {
        const mentor = p.mentor?.name || "—";
        const mentee = p.mentee?.name || "—";
        return `${mentor} / ${mentee}`;
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

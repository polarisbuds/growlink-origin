function parseSharers(text) {
  if (!text || !text.trim()) return [];
  return text.split(",").map((s) => parseInt(s.trim(), 10)).filter((n) => !Number.isNaN(n));
}

document.addEventListener("alpine:init", () => {
  Alpine.data("pairingForm", (mode) => {
    let initial = { mentorId: 0, menteeId: 0, recurWeek: null, recurWeekday: null, id: 0, sharerUserIds: [] };
    const el = document.getElementById("pairing-json");
    if (el) {
      const p = JSON.parse(el.textContent || "{}");
      initial = {
        id: p.id,
        mentorId: p.mentor?.id || 0,
        menteeId: p.mentee?.id || 0,
        recurWeek: p.recurWeek,
        recurWeekday: p.recurWeekday,
        sharerUserIds: p.sharerUserIds || [],
      };
    }
    return {
      mode,
      form: { mentorId: initial.mentorId, menteeId: initial.menteeId, recurWeek: initial.recurWeek, recurWeekday: initial.recurWeekday },
      pairingId: initial.id,
      sharerIdsText: (initial.sharerUserIds || []).join(","),
      payload() {
        const recurWeek = this.form.recurWeek || null;
        const recurWeekday = this.form.recurWeekday ?? null;
        return {
          mentorId: this.form.mentorId,
          menteeId: this.form.menteeId,
          recurWeek: recurWeek === 0 ? null : recurWeek,
          recurWeekday: recurWeek == null ? null : recurWeekday,
          sharerUserIds: parseSharers(this.sharerIdsText),
        };
      },
      async save() {
        const body = this.payload();
        if (this.mode === "create") {
          await fetch("/api/v1/pairings", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) });
        } else {
          await fetch(`/api/v1/pairings/${this.pairingId}`, { method: "PUT", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) });
        }
        location.href = "/pairing";
      },
    };
  });
});

import { showToast, refreshPageIcons } from "../notice.js";

document.addEventListener("alpine:init", () => {
  const el = document.getElementById("pairings-json");
  const pairings = el ? JSON.parse(el.textContent || "[]") : [];
  Alpine.data("pairingListPage", () => ({
    pairings,
    search: "",
    selected: null,
    deleteOpen: false,
    init() {
      refreshPageIcons();
    },
    get filteredPairings() {
      const q = this.search.trim().toLowerCase();
      if (!q) return this.pairings;
      return this.pairings.filter((p) => {
        const menteeNames = (p.mentees || []).map((m) => m.name).join(" ");
        const hay = [p.mentor?.name, p.mentee?.name, menteeNames].filter(Boolean).join(" ").toLowerCase();
        return hay.includes(q);
      });
    },
    recurLabel(p) {
      if (p.recurWeek == null) return "なし";
      const days = ["日", "月", "火", "水", "木", "金", "土"];
      const wd = p.recurWeekday != null ? days[p.recurWeekday] || p.recurWeekday : "—";
      return `第${p.recurWeek}週・${wd}`;
    },
    askDelete(p) {
      this.selected = p;
      this.deleteOpen = true;
    },
    async confirmDelete() {
      const res = await fetch(`/api/v1/pairings/${this.selected.uuid || this.selected.id}`, { method: "DELETE" });
      if (!res.ok) {
        showToast("削除に失敗しました", "error");
        return;
      }
      this.deleteOpen = false;
      showToast("ペアを削除しました");
      location.reload();
    },
  }));
});

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
      return this.pairings.filter((p) => this.searchHaystack(p).includes(q));
    },
    searchHaystack(p) {
      const targetLabels = (p.targets || []).map((t) => t.label).join(" ");
      return [p.mentor?.name, p.mentee?.name, targetLabels].filter(Boolean).join(" ").toLowerCase();
    },
    targetsLabel(p) {
      const labels = (p.targets || []).map((t) => t.label).filter(Boolean);
      if (labels.length) return labels.join("、");
      return p.mentee?.name || "—";
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

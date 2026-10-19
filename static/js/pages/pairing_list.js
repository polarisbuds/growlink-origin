document.addEventListener("alpine:init", () => {
  const el = document.getElementById("pairings-json");
  const pairings = el ? JSON.parse(el.textContent || "[]") : [];
  Alpine.data("pairingListPage", () => ({
    pairings,
    selected: null,
    deleteOpen: false,
    recurLabel(p) {
      if (p.recurWeek == null) return "なし";
      return `第${p.recurWeek}週・曜日${p.recurWeekday}`;
    },
    askDelete(p) {
      this.selected = p;
      this.deleteOpen = true;
    },
    async confirmDelete() {
      await fetch(`/api/v1/pairings/${this.selected.id}`, { method: "DELETE" });
      this.deleteOpen = false;
      location.reload();
    },
  }));
});

document.addEventListener("alpine:init", () => {
  const el = document.getElementById("templates-json");
  let templates = el ? JSON.parse(el.textContent || "[]") : [];
  templates.sort((a, b) => (a.displayNo || 0) - (b.displayNo || 0));

  async function persistOrder(list) {
    const items = list.map((t, i) => ({ id: t.id, displayNo: i + 1 }));
    await fetch("/api/v1/templates/reorder", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ items }),
    });
  }

  Alpine.data("templateListPage", () => ({
    templates,
    async copy(t) {
      await fetch(`/api/v1/templates/${t.id}/copy`, { method: "POST", headers: { "Content-Type": "application/json" }, body: "{}" });
      location.reload();
    },
    async moveUp(t) {
      const i = this.templates.findIndex((x) => x.id === t.id);
      if (i <= 0) return;
      const next = [...this.templates];
      [next[i - 1], next[i]] = [next[i], next[i - 1]];
      this.templates = next;
      await persistOrder(next);
    },
    async moveDown(t) {
      const i = this.templates.findIndex((x) => x.id === t.id);
      if (i < 0 || i >= this.templates.length - 1) return;
      const next = [...this.templates];
      [next[i], next[i + 1]] = [next[i + 1], next[i]];
      this.templates = next;
      await persistOrder(next);
    },
  }));
});

import { showToast, refreshPageIcons } from "/static/js/notice.js";

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
    modalOpen: false,
    editMode: false,
    active: null,
    editForm: { title: "", topics: [] },
    deleteOpen: false,

    init() {
      refreshPageIcons();
    },

    openView(t) {
      this.active = t;
      this.editMode = false;
      this.modalOpen = true;
      this.$nextTick(() => refreshPageIcons());
    },
    closeModal() {
      this.modalOpen = false;
      this.editMode = false;
      this.deleteOpen = false;
    },
    startEdit() {
      this.editForm = {
        title: this.active.title,
        topics: (this.active.topics || []).map((x) => ({
          question: x.question,
          isRequire: x.isRequire,
          displayNo: x.displayNo,
        })),
      };
      this.editMode = true;
    },
    addTopic() {
      this.editForm.topics.push({ question: "", isRequire: true, displayNo: this.editForm.topics.length + 1 });
    },
    removeTopic(i) {
      if (this.editForm.topics.length <= 1) return;
      this.editForm.topics.splice(i, 1);
    },
    async saveEdit() {
      const topics = this.editForm.topics.map((t, i) => ({ ...t, displayNo: i + 1 }));
      const res = await fetch(`/api/v1/templates/${this.active.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title: this.editForm.title, topics }),
      });
      if (!res.ok) {
        showToast("保存に失敗しました", "error");
        return;
      }
      showToast("保存しました");
      location.reload();
    },
    askDelete() {
      this.deleteOpen = true;
    },
    async confirmDelete() {
      const res = await fetch(`/api/v1/templates/${this.active.id}`, { method: "DELETE" });
      if (!res.ok) {
        showToast("削除に失敗しました", "error");
        return;
      }
      location.reload();
    },
    async copy(t) {
      const res = await fetch(`/api/v1/templates/${t.id}/copy`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: "{}",
      });
      if (!res.ok) {
        showToast("コピーに失敗しました", "error");
        return;
      }
      showToast("コピーしました");
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

import { showToast } from "../notice.js";

document.addEventListener("alpine:init", () => {
  Alpine.data("templateNewPage", () => ({
    title: "",
    topics: [{ question: "", isRequire: true, displayNo: 1 }],
    addTopic() {
      this.topics.push({ question: "", isRequire: true, displayNo: this.topics.length + 1 });
    },
    removeTopic(i) {
      if (this.topics.length <= 1) return;
      this.topics.splice(i, 1);
      this.topics.forEach((t, idx) => {
        t.displayNo = idx + 1;
      });
    },
    async save() {
      const res = await fetch("/api/v1/templates", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title: this.title, topics: this.topics }),
      });
      if (!res.ok) {
        showToast("保存に失敗しました", "error");
        return;
      }
      showToast("テンプレートを作成しました");
      location.href = "/template";
    },
  }));
});

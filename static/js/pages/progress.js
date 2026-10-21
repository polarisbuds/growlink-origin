import { refreshPageIcons } from "../notice.js";

document.addEventListener("alpine:init", () => {
  Alpine.data("progressPage", () => ({
    summary: { pending: 0, draft: 0, done: 0, total: 0 },
    async init() {
      const res = await fetch("/api/v1/progress/summary");
      if (res.ok) {
        this.summary = await res.json();
      }
      refreshPageIcons();
    },
  }));
});

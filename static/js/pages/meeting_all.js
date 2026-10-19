import { statusBadge, statusLabel } from "../store.js";

document.addEventListener("alpine:init", () => {
  const el = document.getElementById("meetings-json");
  const meetings = el ? JSON.parse(el.textContent || "[]") : [];
  const simple = location.search.includes("view=simple");
  Alpine.data("meetingAllPage", () => ({
    meetings,
    simple,
    statusBadge,
    statusLabel,
    formatDate(iso) {
      return iso ? new Date(iso).toLocaleDateString("ja-JP") : "—";
    },
  }));
});

import { statusBadge, statusLabel } from "../store.js";

document.addEventListener("alpine:init", () => {
  const el = document.getElementById("meetings-json");
  const meetings = el ? JSON.parse(el.textContent || "[]") : [];
  Alpine.data("meetingListPage", () => ({
    meetings,
    statusBadge,
    statusLabel,
    formatDate(iso) {
      return iso ? new Date(iso).toLocaleDateString("ja-JP") : "—";
    },
  }));
});

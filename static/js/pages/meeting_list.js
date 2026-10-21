import { statusBadge, statusLabel } from "../store.js";
import { refreshPageIcons } from "../notice.js";

document.addEventListener("alpine:init", () => {
  const el = document.getElementById("meetings-json");
  const meetings = el ? JSON.parse(el.textContent || "[]") : [];
  Alpine.data("meetingListPage", () => ({
    meetings,
    search: "",
    statusBadge,
    statusLabel,
    init() {
      refreshPageIcons();
    },
    get filteredMeetings() {
      const q = this.search.trim().toLowerCase();
      if (!q) return this.meetings;
      return this.meetings.filter((m) => {
        const hay = [
          m.template?.title,
          m.pairing?.mentor?.name,
          m.pairing?.mentee?.name,
          this.statusLabel(m.status),
        ]
          .filter(Boolean)
          .join(" ")
          .toLowerCase();
        return hay.includes(q);
      });
    },
    formatDate(iso) {
      return iso ? new Date(iso).toLocaleDateString("ja-JP") : "—";
    },
  }));
});

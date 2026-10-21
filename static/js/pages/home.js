import { statusBadge, statusLabel } from "../store.js";
import { refreshPageIcons } from "../notice.js";

function parseJSON(id) {
  const el = document.getElementById(id);
  if (!el) return [];
  try {
    return JSON.parse(el.textContent || "[]");
  } catch {
    return [];
  }
}

function splitByToday(meetings) {
  const today = new Date();
  today.setHours(23, 59, 59, 999);
  const untilToday = [];
  const fromTomorrow = [];
  for (const m of meetings) {
    const d = new Date(m.implementedAt);
    if (d <= today) untilToday.push(m);
    else fromTomorrow.push(m);
  }
  return { untilToday, fromTomorrow };
}

document.addEventListener("alpine:init", () => {
  const own = parseJSON("meetings-json");
  const sharedMeetings = parseJSON("shared-json");
  const { untilToday, fromTomorrow } = splitByToday(own);
  Alpine.data("growlinkHome", () => ({
    tab: "own",
    untilToday,
    fromTomorrow,
    sharedMeetings,
    statusBadge,
    statusLabel,
    init() {
      refreshPageIcons();
    },
    pairLabel(item) {
      const p = item.pairing;
      if (!p) return "";
      const mentees = (p.mentees || []).map((m) => m.name).filter(Boolean);
      const menteeLabel = mentees.length ? mentees.join("・") : p.mentee?.name || "";
      return `${p.mentor?.name || ""} / ${menteeLabel}`;
    },
    formatDate(iso) {
      if (!iso) return "—";
      return new Date(iso).toLocaleDateString("ja-JP");
    },
  }));
});

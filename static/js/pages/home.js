import { statusBadge, statusLabel } from "../store.js";

function parseMeetings() {
  const el = document.getElementById("meetings-json");
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
  const all = parseMeetings();
  const { untilToday, fromTomorrow } = splitByToday(all);
  Alpine.data("growlinkHome", () => ({
    untilToday,
    fromTomorrow,
    statusBadge,
    statusLabel,
    pairLabel(item) {
      const p = item.pairing;
      if (!p) return "";
      return `${p.mentor?.name || ""} / ${p.mentee?.name || ""}`;
    },
    formatDate(iso) {
      if (!iso) return "—";
      return new Date(iso).toLocaleDateString("ja-JP");
    },
  }));
});

function statusLabel(status) {
  const map = {
    1: "未着手",
    2: "下書き",
    3: "完了",
    pending: "未着手",
    draft: "下書き",
    done: "完了"
  };
  if (typeof status === "string") return map[status] || status;
  return map[status] || String(status);
}
function statusBadge(status) {
  const s = status;
  if (s === 3 || s === "done") return "badge-success";
  if (s === 2 || s === "draft") return "badge-warning";
  return "badge-ghost";
}
function formatDate(iso) {
  if (!iso) return "—";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "—";
  return d.toLocaleDateString("ja-JP");
}
function pairLabel(item) {
  const p = item.pairing;
  if (!p) return "";
  const mentees = (p.mentees || []).map((m) => m.name).filter(Boolean);
  const menteeLabel = mentees.length ? mentees.join("・") : p.mentee?.name || "";
  return `${p.mentor?.name || ""} / ${menteeLabel}`;
}
export {
  statusLabel as a,
  formatDate as f,
  pairLabel as p,
  statusBadge as s
};

import { statusBadge, statusLabel } from "../store.js";
import { showToast } from "../notice.js";

export function meetingPageData() {
  const el = document.getElementById("meeting-json");
  const payload = el ? JSON.parse(el.textContent || "{}") : {};
  return {
    uuid: payload.uuid || "",
    meeting: payload.meeting || {},
    answers: payload.answers || [],
    statusBadge,
    statusLabel,
    async api(method, path, body) {
      const res = await fetch(path, {
        method,
        headers: { "Content-Type": "application/json" },
        body: body ? JSON.stringify(body) : undefined,
      });
      if (!res.ok) throw new Error(await res.text());
      return res.json().catch(() => ({}));
    },
    toast(msg, type = "success") {
      showToast(msg, type);
    },
  };
}

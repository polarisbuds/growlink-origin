/** @typedef {'success' | 'error' | 'info'} ToastType */

function ensureToastHost() {
  let host = document.getElementById("growlink-toast-host");
  if (host) return host;
  host = document.createElement("div");
  host.id = "growlink-toast-host";
  host.className = "toast toast-end toast-bottom z-[100] gap-2";
  host.setAttribute("aria-live", "polite");
  document.body.appendChild(host);
  return host;
}

/**
 * @param {string} message
 * @param {ToastType} [type]
 */
export function showToast(message, type = "success") {
  const host = ensureToastHost();
  const alertClass =
    type === "error" ? "alert-error" : type === "info" ? "alert-info" : "alert-success";
  const el = document.createElement("div");
  el.className = `alert ${alertClass} shadow-lg text-sm max-w-sm`;
  el.textContent = message;
  host.appendChild(el);
  window.setTimeout(() => {
    el.remove();
    if (!host.childElementCount) host.remove();
  }, 3200);
}

export function refreshPageIcons() {
  if (window.lucide && typeof window.lucide.createIcons === "function") {
    window.lucide.createIcons();
  }
}

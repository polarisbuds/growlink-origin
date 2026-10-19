const THEME_KEY = "growlink-theme";
const DEFAULT_THEME = "growlink";

function currentTheme() {
    return localStorage.getItem(THEME_KEY) || DEFAULT_THEME;
}

function applyTheme(theme) {
    document.documentElement.setAttribute("data-theme", theme);
    localStorage.setItem(THEME_KEY, theme);
}

document.addEventListener("alpine:init", () => {
    Alpine.store("ui", {
        sidebarOpen: Alpine.$persist(true).as("growlink.sidebar_open"),
        theme: currentTheme(),
        toggleSidebar() {
            this.sidebarOpen = !this.sidebarOpen;
        },
        setTheme(theme) {
            this.theme = theme;
            applyTheme(theme);
        },
    });
});

export function statusLabel(status) {
    const map = { 1: "未着手", 2: "下書き", 3: "完了", pending: "未着手", draft: "下書き", done: "完了" };
    if (typeof status === "string") return map[status] || status;
    return map[status] || String(status);
}

export function statusBadge(status) {
    const s = typeof status === "number" ? status : status;
    if (s === 3 || s === "done") return "badge-success";
    if (s === 2 || s === "draft") return "badge-warning";
    return "badge-ghost";
}

function refreshIcons() {
    if (window.lucide && typeof window.lucide.createIcons === "function") {
        window.lucide.createIcons();
    }
}

document.addEventListener("DOMContentLoaded", refreshIcons);
document.addEventListener("alpine:initialized", refreshIcons);

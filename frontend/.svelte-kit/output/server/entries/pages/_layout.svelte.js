import { h as head } from "../../chunks/index.js";
import "@sveltejs/kit/internal";
import "../../chunks/exports.js";
import "../../chunks/utils2.js";
import "@sveltejs/kit/internal/server";
import "../../chunks/root.js";
import "../../chunks/state.svelte.js";
import "clsx";
const STORAGE_KEY = "growlink.sidebar_open";
function readSidebar() {
  if (typeof localStorage === "undefined") return true;
  const v = localStorage.getItem(STORAGE_KEY);
  if (v === "false") return false;
  if (v === "true") return true;
  return true;
}
class UiStore {
  sidebarOpen = readSidebar();
  toggleSidebar() {
    this.sidebarOpen = !this.sidebarOpen;
    localStorage.setItem(STORAGE_KEY, String(this.sidebarOpen));
  }
}
new UiStore();
function _layout($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let { children } = $$props;
    head("12qhfyh", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>Growlink</title>`);
      });
    });
    {
      $$renderer2.push(`<!--[0--><div class="flex min-h-screen items-center justify-center bg-base-200 text-sm text-base-content/60">読み込み中…</div>`);
    }
    $$renderer2.push(`<!--]-->`);
  });
}
export {
  _layout as default
};

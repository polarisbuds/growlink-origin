import { h as head, a as attr_class, c as attr, b as ensure_array_like, e as escape_html } from "./index.js";
import "@sveltejs/kit/internal";
import "./exports.js";
import "./utils2.js";
import "@sveltejs/kit/internal/server";
import "./root.js";
import "./state.svelte.js";
function MeetingNewView($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    function todayISO() {
      const d = /* @__PURE__ */ new Date();
      return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
    }
    function currentMonthISO() {
      const d = /* @__PURE__ */ new Date();
      return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}`;
    }
    let users = [];
    let mode = "single";
    let sharerUserIds = [];
    let submitting = false;
    let form = {
      implementedAt: todayISO(),
      baseMonth: currentMonthISO()
    };
    function userLabel(id) {
      const u = users.find((x) => x.id === id);
      return u ? u.name : `#${id}`;
    }
    function submitLabel() {
      return "作成";
    }
    head("rmpbd0", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>1on1作成 - Growlink</title>`);
      });
    });
    $$renderer2.push(`<main class="mx-auto w-full max-w-3xl min-w-0 p-4 md:p-8"><h1 class="font-display text-2xl font-semibold">1on1ミーティング作成</h1> <p class="mt-1 text-sm text-base-content/60">ペアとテンプレートを選び、実施日または定期スケジュールを指定します</p> <div class="mt-4 flex flex-wrap gap-2"><button type="button"${attr_class("btn btn-sm", void 0, {
      "btn-primary": mode === "single",
      "btn-outline": mode !== "single"
    })}>1件作成</button> <button type="button"${attr_class("btn btn-sm", void 0, { "btn-primary": mode === "all", "btn-outline": mode !== "all" })}>全ペア一括</button> <button type="button"${attr_class("btn btn-sm", void 0, {
      "btn-primary": mode === "monthly",
      "btn-outline": mode !== "monthly"
    })}>定期作成</button></div> <div class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">`);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--> `);
    {
      $$renderer2.push(`<!--[0--><div class="form-control relative" data-pick-root="pair"><span class="label-text font-medium">ペア</span> <div class="relative mt-1"><input type="text" class="input input-bordered w-full pr-10" autocomplete="off"${attr("placeholder", "名前で検索して選択")}${attr("value", "")}/> `);
      {
        $$renderer2.push("<!--[-1-->");
      }
      $$renderer2.push(`<!--]--></div> `);
      {
        $$renderer2.push("<!--[-1-->");
      }
      $$renderer2.push(`<!--]--></div>`);
    }
    $$renderer2.push(`<!--]--> <div class="form-control relative" data-pick-root="template"><span class="label-text font-medium">テンプレート</span> <div class="relative mt-1"><input type="text" class="input input-bordered w-full pr-10" autocomplete="off"${attr("placeholder", "名前で検索して選択")}${attr("value", "")}/> `);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> `);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <div class="form-control relative" data-pick-root="sharer"><span class="label-text font-medium">共有者（任意・複数）</span> <div class="mt-2 flex flex-wrap gap-2"><!--[-->`);
    const each_array_2 = ensure_array_like(sharerUserIds);
    for (let $$index_2 = 0, $$length = each_array_2.length; $$index_2 < $$length; $$index_2++) {
      let id = each_array_2[$$index_2];
      $$renderer2.push(`<span class="badge badge-lg gap-1"><span>${escape_html(userLabel(id))}</span> <button type="button" class="btn btn-ghost btn-xs btn-circle">×</button></span>`);
    }
    $$renderer2.push(`<!--]--></div> <div class="relative mt-2"><input type="text" class="input input-bordered w-full" placeholder="名前で検索して追加" autocomplete="off"${attr("value", "")}/></div> `);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> `);
    {
      $$renderer2.push(`<!--[0--><label class="form-control"><span class="label-text font-medium">実施日</span> <input class="input input-bordered" type="date"${attr("value", form.implementedAt)}/></label>`);
    }
    $$renderer2.push(`<!--]--> <div class="flex flex-wrap gap-2 pt-2"><button type="button" class="btn btn-primary"${attr("disabled", submitting, true)}>${escape_html(submitLabel())}</button> <a href="/planning/meetings" class="btn btn-ghost">戻る</a></div></div></main>`);
  });
}
export {
  MeetingNewView as M
};

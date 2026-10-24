import { h as head, b as ensure_array_like, e as escape_html, a as attr_class } from "../../../../chunks/index.js";
import { P as PlanningTabs } from "../../../../chunks/PlanningTabs.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let templates = [];
    let modalOpen = false;
    let active = null;
    let deleteOpen = false;
    head("12090mm", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>テンプレート - Growlink</title>`);
      });
    });
    PlanningTabs($$renderer2, { active: "templates" });
    $$renderer2.push(`<!----> <div class="mt-4 flex flex-wrap items-center justify-between gap-2"><div><h1 class="font-display text-2xl font-semibold">テンプレート一覧</h1> <p class="mt-1 text-sm text-base-content/60">行をクリックして内容を確認できます</p></div> <a href="/planning/templates/new" class="btn btn-primary btn-sm">新規作成</a></div> <div class="mt-6 space-y-2">`);
    if (templates.length === 0) {
      $$renderer2.push(`<!--[0--><div class="rounded-3xl border border-base-300 bg-base-100 p-12 text-center text-sm text-base-content/50">テンプレートがありません</div>`);
    } else {
      $$renderer2.push(`<!--[-1--><!--[-->`);
      const each_array = ensure_array_like(templates);
      for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
        let t = each_array[$$index];
        $$renderer2.push(`<div role="button" tabindex="0" class="group flex cursor-pointer flex-wrap items-center justify-between gap-2 rounded-2xl border border-base-300 bg-base-100 p-4 shadow-sm transition-all hover:border-primary/50 hover:shadow-md"><div class="flex min-w-0 items-center gap-2"><span class="badge badge-ghost shrink-0">#${escape_html(t.displayNo)}</span> <span class="truncate font-medium">${escape_html(t.title)}</span> `);
        if (t.isSample) {
          $$renderer2.push(`<!--[0--><span class="badge badge-outline shrink-0">サンプル</span>`);
        } else {
          $$renderer2.push("<!--[-1-->");
        }
        $$renderer2.push(`<!--]--></div> <div class="flex items-center gap-1" role="presentation"><button type="button" class="btn btn-xs btn-outline">コピー</button> <button type="button" class="btn btn-xs btn-ghost">↑</button> <button type="button" class="btn btn-xs btn-ghost">↓</button></div></div>`);
      }
      $$renderer2.push(`<!--]-->`);
    }
    $$renderer2.push(`<!--]--></div> <dialog${attr_class("modal", void 0, { "modal-open": modalOpen })}><div class="modal-box flex max-h-[85vh] max-w-2xl flex-col p-0"><div class="flex shrink-0 items-start justify-between gap-2 border-b border-base-200 p-6 pb-2"><div class="min-w-0"><h3 class="truncate text-lg font-bold">${escape_html(active?.title)}</h3> <p class="text-sm text-base-content/60">${escape_html("質問一覧")}</p></div> <div class="flex shrink-0 gap-1">`);
    {
      $$renderer2.push(`<!--[0--><button type="button" class="btn btn-sm btn-outline">編集</button>`);
    }
    $$renderer2.push(`<!--]--> <button type="button" class="btn btn-sm btn-error btn-outline">削除</button> <button type="button" class="btn btn-sm btn-circle btn-ghost" aria-label="閉じる">×</button></div></div> `);
    {
      $$renderer2.push(`<!--[0--><div class="flex-1 overflow-y-auto px-6 py-4"><!--[-->`);
      const each_array_1 = ensure_array_like([]);
      for (let $$index_1 = 0, $$length = each_array_1.length; $$index_1 < $$length; $$index_1++) {
        let topic = each_array_1[$$index_1];
        $$renderer2.push(`<div class="mb-3 rounded-2xl border border-base-300 bg-base-200/30 p-4"><div class="flex items-start gap-2"><span class="badge badge-primary badge-outline shrink-0">Q${escape_html(topic.displayNo)}</span> <div class="min-w-0 flex-1"><p class="text-sm leading-relaxed">${escape_html(topic.question)}</p> `);
        if (topic.description) {
          $$renderer2.push(`<!--[0--><p class="mt-1 whitespace-pre-wrap text-xs text-base-content/60">${escape_html(topic.description)}</p>`);
        } else {
          $$renderer2.push("<!--[-1-->");
        }
        $$renderer2.push(`<!--]--></div> `);
        if (topic.isRequire) {
          $$renderer2.push(`<!--[0--><span class="badge badge-error badge-outline badge-sm shrink-0">必須</span>`);
        } else {
          $$renderer2.push("<!--[-1-->");
        }
        $$renderer2.push(`<!--]--></div></div>`);
      }
      $$renderer2.push(`<!--]--></div>`);
    }
    $$renderer2.push(`<!--]--> <div class="modal-action shrink-0 border-t border-base-200 p-4"><button type="button" class="btn btn-ghost">閉じる</button> `);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div></div> <form method="dialog" class="modal-backdrop"><button type="button">close</button></form></dialog> <dialog${attr_class("modal", void 0, { "modal-open": deleteOpen })}><div class="modal-box max-h-[85vh]"><h3 class="font-bold">削除の確認</h3> <p class="py-4 text-sm text-base-content/80">このテンプレートを削除しますか？</p> <div class="modal-action"><button type="button" class="btn btn-ghost">キャンセル</button> <button type="button" class="btn btn-error">削除確定</button></div></div> <form method="dialog" class="modal-backdrop"><button type="button">close</button></form></dialog>`);
  });
}
export {
  _page as default
};

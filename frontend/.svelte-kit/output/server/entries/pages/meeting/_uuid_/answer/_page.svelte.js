import { h as head, b as ensure_array_like, e as escape_html } from "../../../../../chunks/index.js";
import "@sveltejs/kit/internal";
import "../../../../../chunks/exports.js";
import "../../../../../chunks/utils2.js";
import "@sveltejs/kit/internal/server";
import "../../../../../chunks/root.js";
import "../../../../../chunks/state.svelte.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let answers = [];
    let drafts = {};
    head("1dvlqgh", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>面談記入 - Growlink</title>`);
      });
    });
    $$renderer2.push(`<div class="flex max-w-2xl items-center justify-between gap-4"><div><h1 class="font-display text-2xl font-semibold">面談記入</h1> <p class="mt-1 text-sm text-base-content/60">各質問への回答を入力し、下書き保存または提出できます</p></div> <a href="/meeting" class="btn btn-ghost btn-sm shrink-0">戻る</a></div> <div class="mt-6 max-w-2xl space-y-4"><!--[-->`);
    const each_array = ensure_array_like(answers);
    for (let i = 0, $$length = each_array.length; i < $$length; i++) {
      let a = each_array[i];
      $$renderer2.push(`<div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm"><div class="flex items-start gap-2"><span class="badge badge-primary badge-outline shrink-0">Q${escape_html(a.topic?.displayNo || i + 1)}</span> <div class="min-w-0 flex-1"><p class="font-medium leading-relaxed">${escape_html(a.topic?.question)}</p> `);
      if (a.topic?.description) {
        $$renderer2.push(`<!--[0--><p class="mt-1 whitespace-pre-wrap text-sm text-base-content/60">${escape_html(a.topic.description)}</p>`);
      } else {
        $$renderer2.push("<!--[-1-->");
      }
      $$renderer2.push(`<!--]--></div> `);
      if (a.topic?.isRequire) {
        $$renderer2.push(`<!--[0--><span class="badge badge-error badge-outline badge-sm shrink-0">必須</span>`);
      } else {
        $$renderer2.push("<!--[-1-->");
      }
      $$renderer2.push(`<!--]--></div> <textarea class="textarea textarea-bordered mt-3 w-full" rows="3">`);
      const $$body = escape_html(drafts[a.id] ?? "");
      if ($$body) {
        $$renderer2.push(`${$$body}`);
      }
      $$renderer2.push(`</textarea></div>`);
    }
    $$renderer2.push(`<!--]--></div> <div class="mt-6 flex max-w-2xl flex-wrap gap-2"><button type="button" class="btn btn-outline">下書き保存</button> <button type="button" class="btn btn-primary">提出</button> <button type="button" class="btn btn-ghost">提出キャンセル</button></div>`);
  });
}
export {
  _page as default
};

import { h as head, c as attr, b as ensure_array_like, e as escape_html } from "../../../../../chunks/index.js";
import "@sveltejs/kit/internal";
import "../../../../../chunks/exports.js";
import "../../../../../chunks/utils2.js";
import "@sveltejs/kit/internal/server";
import "../../../../../chunks/root.js";
import "../../../../../chunks/state.svelte.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let title = "";
    let topics = [
      { question: "", description: "", isRequire: true, displayNo: 1 }
    ];
    head("1fe0em1", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>テンプレート作成 - Growlink</title>`);
      });
    });
    $$renderer2.push(`<main class="mx-auto w-full max-w-3xl p-4 md:p-8"><h1 class="font-display text-2xl font-semibold">テンプレート新規作成</h1> <p class="mt-1 text-sm text-base-content/60">質問を追加して面談のひな形を作ります</p> <div class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm"><label class="form-control"><span class="label-text font-medium">タイトル</span> <input class="input input-bordered"${attr("value", title)}/></label> <!--[-->`);
    const each_array = ensure_array_like(topics);
    for (let i = 0, $$length = each_array.length; i < $$length; i++) {
      let topic = each_array[i];
      $$renderer2.push(`<div class="flex items-center gap-2 rounded-2xl border border-base-300 p-2"><span class="badge badge-ghost shrink-0">Q${escape_html(i + 1)}</span> <div class="min-w-0 flex-1 space-y-1"><input class="input input-bordered input-sm w-full"${attr("value", topic.question)} placeholder="質問文"/> <textarea class="textarea textarea-bordered textarea-sm min-h-[2.5rem] w-full" rows="2" placeholder="説明（改行可能）">`);
      const $$body = escape_html(topic.description);
      if ($$body) {
        $$renderer2.push(`${$$body}`);
      }
      $$renderer2.push(`</textarea></div> <label class="label shrink-0 cursor-pointer gap-1 py-0"><input type="checkbox" class="checkbox checkbox-sm"${attr("checked", topic.isRequire, true)}/> <span class="label-text text-xs">必須</span></label> <button type="button" class="btn btn-ghost btn-xs btn-square text-error" title="削除"${attr("disabled", topics.length <= 1, true)}>×</button></div>`);
    }
    $$renderer2.push(`<!--]--> <button type="button" class="btn btn-outline btn-sm w-full">質問を追加</button> <button type="button" class="btn btn-primary">保存</button></div></main>`);
  });
}
export {
  _page as default
};

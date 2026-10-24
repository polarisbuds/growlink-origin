import { h as head, c as attr } from "../../../../../chunks/index.js";
import "../../../../../chunks/state.svelte.js";
import "@sveltejs/kit/internal";
import "../../../../../chunks/exports.js";
import "../../../../../chunks/utils2.js";
import "@sveltejs/kit/internal/server";
import "../../../../../chunks/root.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let date = "";
    head("1ukbuf5", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>ミーティング設定 - Growlink</title>`);
      });
    });
    $$renderer2.push(`<h1 class="font-display text-2xl font-semibold">実施予定日</h1> <p class="mt-1 text-sm text-base-content/60">この面談の実施予定日を変更します</p> <div class="mt-6 max-w-lg rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm"><label class="form-control"><span class="label-text font-medium">実施予定日</span> <input type="date" class="input input-bordered"${attr("value", date)}/></label> <div class="mt-6 flex flex-wrap gap-2"><button type="button" class="btn btn-primary">保存</button> <a href="/meeting" class="btn btn-ghost">戻る</a></div></div>`);
  });
}
export {
  _page as default
};

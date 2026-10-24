import { h as head, e as escape_html } from "../../../chunks/index.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let summary = { pending: 0, draft: 0, done: 0, total: 0 };
    head("tllgyw", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>進捗 - Growlink</title>`);
      });
    });
    $$renderer2.push(`<h1 class="font-display text-2xl font-semibold">進捗</h1> <p class="mt-1 text-sm text-base-content/60">面談ステータス別の件数</p> <div class="mt-8 grid gap-4 md:grid-cols-2 xl:grid-cols-4"><div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm"><p class="text-sm text-base-content/60">未着手</p> <p class="mt-2 text-3xl font-semibold text-primary">${escape_html(summary.pending)}</p></div> <div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm"><p class="text-sm text-base-content/60">下書き</p> <p class="mt-2 text-3xl font-semibold text-warning">${escape_html(summary.draft)}</p></div> <div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm"><p class="text-sm text-base-content/60">完了</p> <p class="mt-2 text-3xl font-semibold text-success">${escape_html(summary.done)}</p></div> <div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm"><p class="text-sm text-base-content/60">合計</p> <p class="mt-2 text-3xl font-semibold">${escape_html(summary.total)}</p></div></div>`);
  });
}
export {
  _page as default
};

import { h as head, c as attr, b as ensure_array_like, e as escape_html, s as stringify, a as attr_class, o as derived } from "../../../../chunks/index.js";
import { P as PlanningTabs } from "../../../../chunks/PlanningTabs.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let pairings = [];
    let search = "";
    let deleteOpen = false;
    function searchHaystack(p) {
      const targetLabels = (p.targets || []).map((t) => t.label).join(" ");
      return [p.mentor?.name, p.mentee?.name, targetLabels].filter(Boolean).join(" ").toLowerCase();
    }
    function targetsLabel(p) {
      const labels = (p.targets || []).map((t) => t.label).filter(Boolean);
      if (labels.length) return labels.join("、");
      return p.mentee?.name || "—";
    }
    let filteredPairings = derived(() => {
      const q = search.trim().toLowerCase();
      if (!q) return pairings;
      return pairings.filter((p) => searchHaystack(p).includes(q));
    });
    head("ikhw7c", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>ペア一覧 - Growlink</title>`);
      });
    });
    PlanningTabs($$renderer2, { active: "pairs" });
    $$renderer2.push(`<!----> <div class="mt-4 flex flex-wrap items-start justify-between gap-3"><div><h1 class="font-display text-2xl font-semibold">ペア一覧</h1> <p class="mt-1 text-sm text-base-content/60">メンターとメンティのペア</p></div> <a href="/planning/pairs/new" class="btn btn-primary btn-sm">+ 作成</a></div> <label class="input input-bordered mt-4 flex max-w-md items-center gap-2"><input type="search" class="grow" placeholder="メンター・対象で検索"${attr("value", search)}/></label> <div class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">`);
    if (filteredPairings().length === 0) {
      $$renderer2.push(`<!--[0--><div class="p-12 text-center text-sm text-base-content/50">該当するペアがありません</div>`);
    } else {
      $$renderer2.push(`<!--[-1--><table class="table"><thead><tr class="text-base-content/70"><th>メンター</th><th>対象</th><th></th></tr></thead><tbody><!--[-->`);
      const each_array = ensure_array_like(filteredPairings());
      for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
        let p = each_array[$$index];
        $$renderer2.push(`<tr class="hover:bg-base-200/40"><td>${escape_html(p.mentor?.name)}</td><td class="text-sm">${escape_html(targetsLabel(p))}</td><td class="space-x-1 whitespace-nowrap text-right"><a${attr("href", `/planning/pairs/${stringify(p.uuid)}/edit`)} class="btn btn-xs btn-outline">編集</a> <button type="button" class="btn btn-xs btn-error btn-outline">削除</button> <a${attr("href", `/planning/pairs/${stringify(p.uuid)}/edit`)} class="btn btn-ghost btn-xs btn-square inline-flex align-middle">›</a></td></tr>`);
      }
      $$renderer2.push(`<!--]--></tbody></table>`);
    }
    $$renderer2.push(`<!--]--></div> <dialog${attr_class("modal", void 0, { "modal-open": deleteOpen })}><div class="modal-box max-h-[85vh]"><h3 class="text-lg font-bold">削除の確認</h3> <p class="py-4 text-sm text-base-content/80">このペアを削除しますか？</p> <div class="modal-action"><button type="button" class="btn btn-ghost">キャンセル</button> <button type="button" class="btn btn-error">削除確定</button></div></div> <form method="dialog" class="modal-backdrop"><button type="button">close</button></form></dialog>`);
  });
}
export {
  _page as default
};

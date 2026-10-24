import { h as head, c as attr, e as escape_html, b as ensure_array_like, a as attr_class, s as stringify, o as derived } from "../../../../chunks/index.js";
import { p as page } from "../../../../chunks/index2.js";
import { P as PlanningTabs } from "../../../../chunks/PlanningTabs.js";
import { s as statusBadge, a as statusLabel, f as formatDate } from "../../../../chunks/meeting-helpers.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let meetings = [];
    let search = "";
    let simple = derived(() => page.url.searchParams.get("view") === "simple");
    let filteredMeetings = derived(() => {
      const q = search.trim().toLowerCase();
      if (!q) return meetings;
      return meetings.filter((m) => {
        const hay = [
          m.template?.title,
          m.pairing?.mentor?.name,
          m.pairing?.mentee?.name
        ].filter(Boolean).join(" ").toLowerCase();
        return hay.includes(q);
      });
    });
    head("oieb73", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>全ての面談 - Growlink</title>`);
      });
    });
    PlanningTabs($$renderer2, { active: "meetings" });
    $$renderer2.push(`<!----> <div class="mt-4 flex flex-wrap items-start justify-between gap-3"><div><h1 class="font-display text-2xl font-semibold">面談一覧</h1> <p class="mt-1 text-sm text-base-content/60">組織内の面談を一覧表示</p></div> <div class="flex flex-wrap items-center gap-2"><a href="/planning/meetings/new" class="btn btn-sm btn-primary">1on1を追加</a> <a${attr("href", simple() ? "/planning/meetings" : "/planning/meetings?view=simple")} class="btn btn-sm btn-outline">${escape_html(simple() ? "詳細表示" : "簡易表示")}</a></div></div> <label class="input input-bordered mt-4 flex max-w-md items-center gap-2"><input type="search" class="grow" placeholder="検索"${attr("value", search)}/></label> <div class="mt-6 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">`);
    if (filteredMeetings().length === 0) {
      $$renderer2.push(`<!--[0--><div class="p-12 text-center text-sm text-base-content/50">該当する面談がありません</div>`);
    } else {
      $$renderer2.push(`<!--[-1--><table class="table"><thead><tr class="text-base-content/70"><th>メンター</th><th>メンティ</th><th>状態</th><th>実施日</th><th></th></tr></thead><tbody><!--[-->`);
      const each_array = ensure_array_like(filteredMeetings());
      for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
        let m = each_array[$$index];
        $$renderer2.push(`<tr class="hover:bg-base-200/40"><td>${escape_html(m.pairing?.mentor?.name)}</td><td>${escape_html(m.pairing?.mentee?.name)}</td><td><span${attr_class(`badge badge-sm ${stringify(statusBadge(m.status))}`)}>${escape_html(statusLabel(m.status))}</span></td><td>${escape_html(formatDate(m.implementedAt))}</td><td><a${attr("href", `/meeting/${stringify(m.uuid)}`)} class="btn btn-xs btn-ghost gap-1">閲覧</a></td></tr> `);
        if (simple()) {
          $$renderer2.push(`<!--[0--><tr class="bg-base-200/40"><td colspan="5" class="text-sm text-base-content/70">${escape_html(m.template?.title)}</td></tr>`);
        } else {
          $$renderer2.push("<!--[-1-->");
        }
        $$renderer2.push(`<!--]-->`);
      }
      $$renderer2.push(`<!--]--></tbody></table>`);
    }
    $$renderer2.push(`<!--]--></div>`);
  });
}
export {
  _page as default
};

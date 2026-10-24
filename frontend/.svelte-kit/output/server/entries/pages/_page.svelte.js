import { h as head, a as attr_class, b as ensure_array_like, e as escape_html, s as stringify, c as attr } from "../../chunks/index.js";
import { p as pairLabel, s as statusBadge, a as statusLabel, f as formatDate } from "../../chunks/meeting-helpers.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let tab = "own";
    let untilToday = [];
    let fromTomorrow = [];
    head("1uha8ag", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>ワークスペース - Growlink</title>`);
      });
    });
    $$renderer2.push(`<h1 class="font-display text-2xl font-semibold">ワークスペース</h1> <p class="mt-1 text-sm text-base-content/60">あなたの1on1と共有された1on1</p> <div class="tabs tabs-boxed mt-6 w-fit"><button type="button"${attr_class("tab", void 0, { "tab-active": tab === "own" })}>あなたの1on1</button> <button type="button"${attr_class("tab", void 0, { "tab-active": tab === "shared" })}>共有された1on1</button></div> `);
    {
      $$renderer2.push(`<!--[0--><section class="mt-8"><h2 class="text-sm font-semibold text-base-content/70">今日まで</h2> <div class="mt-4 grid gap-4 md:grid-cols-2 xl:grid-cols-3"><!--[-->`);
      const each_array = ensure_array_like(untilToday);
      for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
        let item = each_array[$$index];
        $$renderer2.push(`<div class="group rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm hover:border-primary/40 transition-colors"><div class="flex items-start justify-between gap-3"><div class="min-w-0"><p class="text-sm text-base-content/60 truncate">${escape_html(pairLabel(item))}</p> <h3 class="mt-1 font-semibold">${escape_html(item.template?.title || "1on1")}</h3></div> <span${attr_class(`badge shrink-0 ${stringify(statusBadge(item.status))}`)}>${escape_html(statusLabel(item.status))}</span></div> <p class="mt-3 text-sm text-base-content/70">実施予定日: ${escape_html(formatDate(item.implementedAt))}</p> <div class="mt-4 flex flex-wrap gap-2 items-center"><a${attr("href", `/meeting/${stringify(item.uuid)}/answer`)} class="btn btn-sm btn-primary">記入</a> <a${attr("href", `/meeting/${stringify(item.uuid)}`)} class="btn btn-sm btn-ghost">閲覧</a> <a${attr("href", `/meeting/${stringify(item.uuid)}/setting`)} class="btn btn-sm btn-ghost">設定</a></div></div>`);
      }
      $$renderer2.push(`<!--]--> `);
      if (untilToday.length === 0) {
        $$renderer2.push(`<!--[0--><div class="rounded-3xl border border-dashed border-base-300 bg-base-100 p-10 text-center text-sm text-base-content/50 md:col-span-2 xl:col-span-3">今日までの面談はありません</div>`);
      } else {
        $$renderer2.push("<!--[-1-->");
      }
      $$renderer2.push(`<!--]--></div> <h2 class="mt-10 text-sm font-semibold text-base-content/70">明日以降</h2> <div class="mt-4 grid gap-4 md:grid-cols-2 xl:grid-cols-3"><!--[-->`);
      const each_array_1 = ensure_array_like(fromTomorrow);
      for (let $$index_1 = 0, $$length = each_array_1.length; $$index_1 < $$length; $$index_1++) {
        let item = each_array_1[$$index_1];
        $$renderer2.push(`<div class="group rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm"><p class="text-sm text-base-content/60">${escape_html(pairLabel(item))}</p> <h3 class="mt-1 font-semibold">${escape_html(item.template?.title || "1on1")}</h3> <p class="mt-3 text-sm text-base-content/70">実施予定日: ${escape_html(formatDate(item.implementedAt))}</p> <div class="mt-4 flex items-center gap-2"><a${attr("href", `/meeting/${stringify(item.uuid)}/answer`)} class="btn btn-sm btn-primary">記入</a> <a${attr("href", `/meeting/${stringify(item.uuid)}`)} class="btn btn-sm btn-ghost">閲覧</a></div></div>`);
      }
      $$renderer2.push(`<!--]--> `);
      if (fromTomorrow.length === 0) {
        $$renderer2.push(`<!--[0--><div class="rounded-3xl border border-dashed border-base-300 bg-base-100 p-10 text-center text-sm text-base-content/50 md:col-span-2 xl:col-span-3">明日以降の面談はありません</div>`);
      } else {
        $$renderer2.push("<!--[-1-->");
      }
      $$renderer2.push(`<!--]--></div></section>`);
    }
    $$renderer2.push(`<!--]-->`);
  });
}
export {
  _page as default
};

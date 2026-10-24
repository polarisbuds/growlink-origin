import { h as head, c as attr, b as ensure_array_like, e as escape_html, a as attr_class, s as stringify, o as derived } from "../../../chunks/index.js";
import { f as formatDate, s as statusBadge, a as statusLabel } from "../../../chunks/meeting-helpers.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let meetings = [];
    let search = "";
    let filtered = derived(() => {
      const q = search.trim().toLowerCase();
      if (!q) return meetings;
      return meetings.filter((m) => {
        const template = m.template;
        const pairing = m.pairing;
        const hay = [
          template?.title,
          pairing?.mentor?.name,
          pairing?.mentee?.name,
          statusLabel(m.status)
        ].filter(Boolean).join(" ").toLowerCase();
        return hay.includes(q);
      });
    });
    head("1ksakkq", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>面談一覧 - Growlink</title>`);
      });
    });
    $$renderer2.push(`<h1 class="font-display text-2xl font-semibold">面談一覧</h1> <input type="search" class="input input-bordered mt-4 w-full max-w-md" placeholder="検索"${attr("value", search)}/> <div class="mt-6 overflow-x-auto rounded-2xl border border-base-300 bg-base-100"><table class="table"><thead><tr><th>テンプレート</th><th>ペア</th><th>実施日</th><th>状態</th><th></th></tr></thead><tbody><!--[-->`);
    const each_array = ensure_array_like(filtered());
    for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
      let m = each_array[$$index];
      $$renderer2.push(`<tr><td>${escape_html(m.template?.title || "—")}</td><td>${escape_html(m.pairing?.mentor?.name || "—")} /
						${escape_html(m.pairing?.mentee?.name || "—")}</td><td>${escape_html(formatDate(m.implementedAt))}</td><td><span${attr_class(`badge ${stringify(statusBadge(m.status))}`)}>${escape_html(statusLabel(m.status))}</span></td><td class="text-right"><a${attr("href", `/meeting/${stringify(m.uuid)}`)} class="btn btn-ghost btn-xs">開く</a></td></tr>`);
    }
    $$renderer2.push(`<!--]--></tbody></table></div>`);
  });
}
export {
  _page as default
};

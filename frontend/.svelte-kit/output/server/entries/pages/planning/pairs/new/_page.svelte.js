import "clsx";
import { h as head, e as escape_html, c as attr, b as ensure_array_like } from "../../../../../chunks/index.js";
import "@sveltejs/kit/internal";
import "../../../../../chunks/exports.js";
import "../../../../../chunks/utils2.js";
import "@sveltejs/kit/internal/server";
import "../../../../../chunks/root.js";
import "../../../../../chunks/state.svelte.js";
function PairingFormView($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let targets = [];
    let teams = [];
    let departments = [];
    let teamPickId = 0;
    let departmentPickId = 0;
    function targetKindLabel(kind) {
      if (kind === "user") return "ユーザー";
      if (kind === "team") return "チーム";
      if (kind === "department") return "部署";
      return kind;
    }
    head("1e47idh", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>${escape_html("ペア作成")} - Growlink</title>`);
      });
    });
    $$renderer2.push(`<main class="mx-auto w-full max-w-3xl min-w-0 p-4 md:p-8"><h1 class="font-display text-2xl font-semibold">${escape_html("ペア新規作成")}</h1> <p class="mt-1 text-sm text-base-content/60">メンターと対象（ユーザー・チーム・部署）を${escape_html("登録")}します</p> <div class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm"><div class="form-control relative" data-pick-root="mentor"><span class="label-text font-medium">メンター</span> <div class="relative mt-1"><input type="text" class="input input-bordered w-full" autocomplete="off"${attr("placeholder", "名前で検索")}${attr("value", "")}/></div> `);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <div><span class="label-text font-medium">対象</span> <div class="mt-2 flex flex-wrap gap-2"><!--[-->`);
    const each_array_1 = ensure_array_like(targets);
    for (let i = 0, $$length = each_array_1.length; i < $$length; i++) {
      let t = each_array_1[i];
      $$renderer2.push(`<span class="badge badge-lg gap-1"><span class="text-xs opacity-60">${escape_html(targetKindLabel(t.kind))}</span> <span>${escape_html(t.label)}</span> <button type="button" class="btn btn-ghost btn-xs btn-circle">×</button></span>`);
    }
    $$renderer2.push(`<!--]--></div></div> <div class="form-control relative" data-pick-root="user-target"><span class="label-text font-medium">ユーザーを追加</span> <div class="relative mt-1"><input type="text" class="input input-bordered w-full" placeholder="名前で検索" autocomplete="off"${attr("value", "")}/></div> `);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <div class="flex flex-wrap items-end gap-2"><label class="form-control min-w-[12rem] flex-1"><span class="label-text font-medium">チームを追加</span> `);
    $$renderer2.select({ class: "select select-bordered", value: teamPickId }, ($$renderer3) => {
      $$renderer3.option({ value: 0 }, ($$renderer4) => {
        $$renderer4.push(`選択`);
      });
      $$renderer3.push(`<!--[-->`);
      const each_array_3 = ensure_array_like(teams);
      for (let $$index_3 = 0, $$length = each_array_3.length; $$index_3 < $$length; $$index_3++) {
        let t = each_array_3[$$index_3];
        $$renderer3.option({ value: t.id }, ($$renderer4) => {
          $$renderer4.push(`${escape_html(t.name)}`);
        });
      }
      $$renderer3.push(`<!--]-->`);
    });
    $$renderer2.push(`</label> <button type="button" class="btn btn-outline">追加</button></div> <div class="flex flex-wrap items-end gap-2"><label class="form-control min-w-[12rem] flex-1"><span class="label-text font-medium">部署を追加</span> `);
    $$renderer2.select({ class: "select select-bordered", value: departmentPickId }, ($$renderer3) => {
      $$renderer3.option({ value: 0 }, ($$renderer4) => {
        $$renderer4.push(`選択`);
      });
      $$renderer3.push(`<!--[-->`);
      const each_array_4 = ensure_array_like(departments);
      for (let $$index_4 = 0, $$length = each_array_4.length; $$index_4 < $$length; $$index_4++) {
        let d = each_array_4[$$index_4];
        $$renderer3.option({ value: d.id }, ($$renderer4) => {
          $$renderer4.push(`${escape_html(d.name)}`);
        });
      }
      $$renderer3.push(`<!--]-->`);
    });
    $$renderer2.push(`</label> <button type="button" class="btn btn-outline">追加</button></div> <div class="flex flex-wrap gap-2 pt-2"><button type="button" class="btn btn-primary">${escape_html("作成")}</button> <a href="/planning/pairs" class="btn btn-ghost">戻る</a></div></div></main>`);
  });
}
function _page($$renderer) {
  PairingFormView($$renderer);
}
export {
  _page as default
};

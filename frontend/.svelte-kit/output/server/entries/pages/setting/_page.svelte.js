import { h as head, c as attr, e as escape_html, a as attr_class, o as derived } from "../../../chunks/index.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let form = { name: "", email: "", lineWorksId: "", remindEnabled: false };
    let avatarUrl = derived(() => null);
    let initial = derived(() => "?".slice(0, 1));
    head("1qghcj7", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>設定 - Growlink</title>`);
      });
    });
    $$renderer2.push(`<div class="max-w-3xl mx-auto w-full"><h1 class="font-display text-2xl font-semibold">個人設定</h1> <p class="text-sm text-base-content/60 mt-1">プロフィールと通知の設定</p> <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm space-y-5"><h2 class="font-semibold">基本情報</h2> <div class="flex flex-col sm:flex-row items-center sm:items-start gap-5">`);
    if (avatarUrl()) {
      $$renderer2.push(`<!--[0--><img${attr("src", avatarUrl())} alt="" class="w-24 h-24 rounded-full object-cover border-2 border-base-300 shadow-sm"/>`);
    } else {
      $$renderer2.push(`<!--[-1--><div class="w-24 h-24 rounded-full bg-primary text-primary-content flex items-center justify-center text-3xl font-semibold shadow-sm">${escape_html(initial())}</div>`);
    }
    $$renderer2.push(`<!--]--> <div class="flex flex-col gap-2 w-full sm:flex-1"><input type="file" class="file-input file-input-bordered file-input-sm w-full" accept="image/*"/> <button type="button" class="btn btn-ghost btn-sm w-fit">アイコンを削除</button></div></div> <label class="form-control"><span class="label-text font-medium">氏名</span> <input class="input input-bordered"${attr("value", form.name)}/></label> <label class="form-control"><span class="label-text font-medium">メール</span> <input class="input input-bordered"${attr("value", form.email)} type="email"/></label> <label class="form-control"><span class="label-text font-medium">LINE WORKS ID</span> <div class="join w-full"><input${attr_class("input input-bordered join-item w-full", void 0, { "input-error": false })}${attr("value", form.lineWorksId)} placeholder="ローカル部のみ"/> <span class="btn join-item no-animation pointer-events-none bg-base-200">@cap</span></div> <span class="label-text-alt text-base-content/50 mt-1">ドメイン @cap は固定です。変更できません。</span> `);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></label> <button type="button" class="btn btn-primary w-full sm:w-auto">保存</button></div> <div class="mt-6 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm"><h2 class="font-semibold">通知</h2> <label class="label cursor-pointer justify-start gap-3 mt-3 rounded-2xl border border-base-300 px-4 py-3"><input type="checkbox" class="toggle toggle-primary"${attr("checked", form.remindEnabled, true)}/> <span class="label-text">トピック未提出のリマインドを受け取る</span></label> <p class="text-xs text-base-content/50 mt-3">送信処理はバックエンドのジョブで実行されます（設定の保存のみ）</p> <button type="button" class="btn btn-outline btn-sm mt-4">通知設定を保存</button></div></div>`);
  });
}
export {
  _page as default
};

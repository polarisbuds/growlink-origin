import "clsx";
import { d as sanitize_props, r as rest_props, f as fallback, g as attributes, i as clsx, b as ensure_array_like, j as element, k as slot, l as bind_props, m as spread_props, h as head, a as attr_class, c as attr, e as escape_html } from "../../../chunks/index.js";
/**
 * @license lucide-svelte v0.487.0 - ISC
 *
 * ISC License
 * 
 * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
 * 
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 * 
 */
const defaultAttributes = {
  xmlns: "http://www.w3.org/2000/svg",
  width: 24,
  height: 24,
  viewBox: "0 0 24 24",
  fill: "none",
  stroke: "currentColor",
  "stroke-width": 2,
  "stroke-linecap": "round",
  "stroke-linejoin": "round"
};
function Icon($$renderer, $$props) {
  const $$sanitized_props = sanitize_props($$props);
  const $$restProps = rest_props($$sanitized_props, [
    "name",
    "color",
    "size",
    "strokeWidth",
    "absoluteStrokeWidth",
    "iconNode"
  ]);
  $$renderer.component(($$renderer2) => {
    let name = fallback($$props["name"], void 0);
    let color = fallback($$props["color"], "currentColor");
    let size = fallback($$props["size"], 24);
    let strokeWidth = fallback($$props["strokeWidth"], 2);
    let absoluteStrokeWidth = fallback($$props["absoluteStrokeWidth"], false);
    let iconNode = fallback($$props["iconNode"], () => [], true);
    const mergeClasses = (...classes) => classes.filter((className, index, array) => {
      return Boolean(className) && array.indexOf(className) === index;
    }).join(" ");
    $$renderer2.push(`<svg${attributes(
      {
        ...defaultAttributes,
        ...$$restProps,
        width: size,
        height: size,
        stroke: color,
        "stroke-width": absoluteStrokeWidth ? Number(strokeWidth) * 24 / Number(size) : strokeWidth,
        class: clsx(mergeClasses("lucide-icon", "lucide", name ? `lucide-${name}` : "", $$sanitized_props.class))
      },
      void 0,
      void 0,
      void 0,
      3
    )}><!--[-->`);
    const each_array = ensure_array_like(iconNode);
    for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
      let [tag, attrs] = each_array[$$index];
      element($$renderer2, tag, () => {
        $$renderer2.push(`${attributes({ ...attrs }, void 0, void 0, void 0, 3)}`);
      });
    }
    $$renderer2.push(`<!--]--><!--[-->`);
    slot($$renderer2, $$props, "default", {});
    $$renderer2.push(`<!--]--></svg>`);
    bind_props($$props, {
      name,
      color,
      size,
      strokeWidth,
      absoluteStrokeWidth,
      iconNode
    });
  });
}
function Users($$renderer, $$props) {
  const $$sanitized_props = sanitize_props($$props);
  /**
   * @license lucide-svelte v0.487.0 - ISC
   *
   * ISC License
   *
   * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
   *
   * Permission to use, copy, modify, and/or distribute this software for any
   * purpose with or without fee is hereby granted, provided that the above
   * copyright notice and this permission notice appear in all copies.
   *
   * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
   * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
   * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
   * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
   * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
   * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
   * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
   *
   */
  const iconNode = [
    ["path", { "d": "M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" }],
    ["circle", { "cx": "9", "cy": "7", "r": "4" }],
    ["path", { "d": "M22 21v-2a4 4 0 0 0-3-3.87" }],
    ["path", { "d": "M16 3.13a4 4 0 0 1 0 7.75" }]
  ];
  Icon($$renderer, spread_props([
    { name: "users" },
    $$sanitized_props,
    {
      /**
       * @component @name Users
       * @description Lucide SVG icon component, renders SVG Element with children.
       *
       * @preview ![img](data:image/svg+xml;base64,PHN2ZyAgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIgogIHdpZHRoPSIyNCIKICBoZWlnaHQ9IjI0IgogIHZpZXdCb3g9IjAgMCAyNCAyNCIKICBmaWxsPSJub25lIgogIHN0cm9rZT0iIzAwMCIgc3R5bGU9ImJhY2tncm91bmQtY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDJweCIKICBzdHJva2Utd2lkdGg9IjIiCiAgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIgogIHN0cm9rZS1saW5lam9pbj0icm91bmQiCj4KICA8cGF0aCBkPSJNMTYgMjF2LTJhNCA0IDAgMCAwLTQtNEg2YTQgNCAwIDAgMC00IDR2MiIgLz4KICA8Y2lyY2xlIGN4PSI5IiBjeT0iNyIgcj0iNCIgLz4KICA8cGF0aCBkPSJNMjIgMjF2LTJhNCA0IDAgMCAwLTMtMy44NyIgLz4KICA8cGF0aCBkPSJNMTYgMy4xM2E0IDQgMCAwIDEgMCA3Ljc1IiAvPgo8L3N2Zz4K) - https://lucide.dev/icons/users
       * @see https://lucide.dev/guide/packages/lucide-svelte - Documentation
       *
       * @param {Object} props - Lucide icons props and any valid SVG attribute
       * @returns {FunctionalComponent} Svelte component
       *
       */
      iconNode,
      children: ($$renderer2) => {
        $$renderer2.push(`<!--[-->`);
        slot($$renderer2, $$props, "default", {});
        $$renderer2.push(`<!--]-->`);
      },
      $$slots: { default: true }
    }
  ]));
}
function Building_2($$renderer, $$props) {
  const $$sanitized_props = sanitize_props($$props);
  /**
   * @license lucide-svelte v0.487.0 - ISC
   *
   * ISC License
   *
   * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
   *
   * Permission to use, copy, modify, and/or distribute this software for any
   * purpose with or without fee is hereby granted, provided that the above
   * copyright notice and this permission notice appear in all copies.
   *
   * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
   * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
   * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
   * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
   * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
   * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
   * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
   *
   */
  const iconNode = [
    ["path", { "d": "M6 22V4a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v18Z" }],
    ["path", { "d": "M6 12H4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h2" }],
    ["path", { "d": "M18 9h2a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-2" }],
    ["path", { "d": "M10 6h4" }],
    ["path", { "d": "M10 10h4" }],
    ["path", { "d": "M10 14h4" }],
    ["path", { "d": "M10 18h4" }]
  ];
  Icon($$renderer, spread_props([
    { name: "building-2" },
    $$sanitized_props,
    {
      /**
       * @component @name Building2
       * @description Lucide SVG icon component, renders SVG Element with children.
       *
       * @preview ![img](data:image/svg+xml;base64,PHN2ZyAgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIgogIHdpZHRoPSIyNCIKICBoZWlnaHQ9IjI0IgogIHZpZXdCb3g9IjAgMCAyNCAyNCIKICBmaWxsPSJub25lIgogIHN0cm9rZT0iIzAwMCIgc3R5bGU9ImJhY2tncm91bmQtY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDJweCIKICBzdHJva2Utd2lkdGg9IjIiCiAgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIgogIHN0cm9rZS1saW5lam9pbj0icm91bmQiCj4KICA8cGF0aCBkPSJNNiAyMlY0YTIgMiAwIDAgMSAyLTJoOGEyIDIgMCAwIDEgMiAydjE4WiIgLz4KICA8cGF0aCBkPSJNNiAxMkg0YTIgMiAwIDAgMC0yIDJ2NmEyIDIgMCAwIDAgMiAyaDIiIC8+CiAgPHBhdGggZD0iTTE4IDloMmEyIDIgMCAwIDEgMiAydjlhMiAyIDAgMCAxLTIgMmgtMiIgLz4KICA8cGF0aCBkPSJNMTAgNmg0IiAvPgogIDxwYXRoIGQ9Ik0xMCAxMGg0IiAvPgogIDxwYXRoIGQ9Ik0xMCAxNGg0IiAvPgogIDxwYXRoIGQ9Ik0xMCAxOGg0IiAvPgo8L3N2Zz4K) - https://lucide.dev/icons/building-2
       * @see https://lucide.dev/guide/packages/lucide-svelte - Documentation
       *
       * @param {Object} props - Lucide icons props and any valid SVG attribute
       * @returns {FunctionalComponent} Svelte component
       *
       */
      iconNode,
      children: ($$renderer2) => {
        $$renderer2.push(`<!--[-->`);
        slot($$renderer2, $$props, "default", {});
        $$renderer2.push(`<!--]-->`);
      },
      $$slots: { default: true }
    }
  ]));
}
function Bell($$renderer, $$props) {
  const $$sanitized_props = sanitize_props($$props);
  /**
   * @license lucide-svelte v0.487.0 - ISC
   *
   * ISC License
   *
   * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
   *
   * Permission to use, copy, modify, and/or distribute this software for any
   * purpose with or without fee is hereby granted, provided that the above
   * copyright notice and this permission notice appear in all copies.
   *
   * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
   * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
   * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
   * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
   * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
   * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
   * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
   *
   */
  const iconNode = [
    ["path", { "d": "M10.268 21a2 2 0 0 0 3.464 0" }],
    [
      "path",
      {
        "d": "M3.262 15.326A1 1 0 0 0 4 17h16a1 1 0 0 0 .74-1.673C19.41 13.956 18 12.499 18 8A6 6 0 0 0 6 8c0 4.499-1.411 5.956-2.738 7.326"
      }
    ]
  ];
  Icon($$renderer, spread_props([
    { name: "bell" },
    $$sanitized_props,
    {
      /**
       * @component @name Bell
       * @description Lucide SVG icon component, renders SVG Element with children.
       *
       * @preview ![img](data:image/svg+xml;base64,PHN2ZyAgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIgogIHdpZHRoPSIyNCIKICBoZWlnaHQ9IjI0IgogIHZpZXdCb3g9IjAgMCAyNCAyNCIKICBmaWxsPSJub25lIgogIHN0cm9rZT0iIzAwMCIgc3R5bGU9ImJhY2tncm91bmQtY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDJweCIKICBzdHJva2Utd2lkdGg9IjIiCiAgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIgogIHN0cm9rZS1saW5lam9pbj0icm91bmQiCj4KICA8cGF0aCBkPSJNMTAuMjY4IDIxYTIgMiAwIDAgMCAzLjQ2NCAwIiAvPgogIDxwYXRoIGQ9Ik0zLjI2MiAxNS4zMjZBMSAxIDAgMCAwIDQgMTdoMTZhMSAxIDAgMCAwIC43NC0xLjY3M0MxOS40MSAxMy45NTYgMTggMTIuNDk5IDE4IDhBNiA2IDAgMCAwIDYgOGMwIDQuNDk5LTEuNDExIDUuOTU2LTIuNzM4IDcuMzI2IiAvPgo8L3N2Zz4K) - https://lucide.dev/icons/bell
       * @see https://lucide.dev/guide/packages/lucide-svelte - Documentation
       *
       * @param {Object} props - Lucide icons props and any valid SVG attribute
       * @returns {FunctionalComponent} Svelte component
       *
       */
      iconNode,
      children: ($$renderer2) => {
        $$renderer2.push(`<!--[-->`);
        slot($$renderer2, $$props, "default", {});
        $$renderer2.push(`<!--]-->`);
      },
      $$slots: { default: true }
    }
  ]));
}
function Search($$renderer, $$props) {
  const $$sanitized_props = sanitize_props($$props);
  /**
   * @license lucide-svelte v0.487.0 - ISC
   *
   * ISC License
   *
   * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
   *
   * Permission to use, copy, modify, and/or distribute this software for any
   * purpose with or without fee is hereby granted, provided that the above
   * copyright notice and this permission notice appear in all copies.
   *
   * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
   * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
   * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
   * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
   * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
   * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
   * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
   *
   */
  const iconNode = [
    ["circle", { "cx": "11", "cy": "11", "r": "8" }],
    ["path", { "d": "m21 21-4.3-4.3" }]
  ];
  Icon($$renderer, spread_props([
    { name: "search" },
    $$sanitized_props,
    {
      /**
       * @component @name Search
       * @description Lucide SVG icon component, renders SVG Element with children.
       *
       * @preview ![img](data:image/svg+xml;base64,PHN2ZyAgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIgogIHdpZHRoPSIyNCIKICBoZWlnaHQ9IjI0IgogIHZpZXdCb3g9IjAgMCAyNCAyNCIKICBmaWxsPSJub25lIgogIHN0cm9rZT0iIzAwMCIgc3R5bGU9ImJhY2tncm91bmQtY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDJweCIKICBzdHJva2Utd2lkdGg9IjIiCiAgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIgogIHN0cm9rZS1saW5lam9pbj0icm91bmQiCj4KICA8Y2lyY2xlIGN4PSIxMSIgY3k9IjExIiByPSI4IiAvPgogIDxwYXRoIGQ9Im0yMSAyMS00LjMtNC4zIiAvPgo8L3N2Zz4K) - https://lucide.dev/icons/search
       * @see https://lucide.dev/guide/packages/lucide-svelte - Documentation
       *
       * @param {Object} props - Lucide icons props and any valid SVG attribute
       * @returns {FunctionalComponent} Svelte component
       *
       */
      iconNode,
      children: ($$renderer2) => {
        $$renderer2.push(`<!--[-->`);
        slot($$renderer2, $$props, "default", {});
        $$renderer2.push(`<!--]-->`);
      },
      $$slots: { default: true }
    }
  ]));
}
function ManagementView($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let tab = "users";
    let options = { departments: [], roles: [] };
    let selectedUserIds = [];
    let userSearch = "";
    let userModalOpen = false;
    let userForm = null;
    let createUserOpen = false;
    let createUserForm = {
      name: "",
      email: "",
      departmentId: 0,
      roleId: 0,
      password: ""
    };
    let createDeptOpen = false;
    let createDeptForm = { name: "", parentDepartmentId: 0 };
    let orgModalOpen = false;
    let orgDetail = null;
    let confirmOpen = false;
    let confirmTitle = "";
    let confirmMessage = "";
    let confirmDanger = false;
    head("13sf17e", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>管理 - Growlink</title>`);
      });
    });
    $$renderer2.push(`<h1 class="font-display text-2xl font-semibold">管理</h1> <p class="mt-1 text-sm text-base-content/60">ユーザ・組織・通知の設定を行います</p> <div role="tablist" class="tabs tabs-bordered mt-6 max-w-3xl"><button type="button" role="tab"${attr_class("tab gap-2", void 0, { "tab-active": tab === "users" })}>`);
    Users($$renderer2, { class: "h-4 w-4" });
    $$renderer2.push(`<!---->ユーザ一覧</button> <button type="button" role="tab"${attr_class("tab gap-2", void 0, { "tab-active": tab === "departments" })}>`);
    Building_2($$renderer2, { class: "h-4 w-4" });
    $$renderer2.push(`<!---->部署一覧</button> <button type="button" role="tab"${attr_class("tab gap-2", void 0, { "tab-active": tab === "remind" })}>`);
    Bell($$renderer2, { class: "h-4 w-4" });
    $$renderer2.push(`<!---->1on1通知</button></div> `);
    {
      $$renderer2.push(`<!--[0--><div class="mt-6"><div class="flex flex-col sm:flex-row sm:items-center gap-3 mb-4"><label class="input input-bordered flex items-center gap-2 flex-1 max-w-md">`);
      Search($$renderer2, { class: "h-4 w-4 opacity-50" });
      $$renderer2.push(`<!----> <input type="search" class="grow" placeholder="氏名・メールで検索"${attr("value", userSearch)}/></label> <button type="button" class="btn btn-primary btn-sm shrink-0">ユーザを追加</button> `);
      if (selectedUserIds.length > 0) {
        $$renderer2.push(`<!--[0--><div class="flex items-center gap-2"><span class="text-sm text-base-content/60">${escape_html(selectedUserIds.length)} 件選択中</span> <button type="button" class="btn btn-error btn-sm">アカウント停止</button></div>`);
      } else {
        $$renderer2.push("<!--[-1-->");
      }
      $$renderer2.push(`<!--]--></div> <div class="overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">`);
      {
        $$renderer2.push(`<!--[0--><div class="p-6 space-y-3"><!--[-->`);
        const each_array = ensure_array_like(Array(5));
        for (let i = 0, $$length = each_array.length; i < $$length; i++) {
          each_array[i];
          $$renderer2.push(`<div class="skeleton h-10 w-full rounded-xl"></div>`);
        }
        $$renderer2.push(`<!--]--></div>`);
      }
      $$renderer2.push(`<!--]--></div></div>`);
    }
    $$renderer2.push(`<!--]--> `);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--> `);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--> <dialog${attr_class("modal", void 0, { "modal-open": userModalOpen })}><div class="modal-box max-w-2xl max-h-[85vh] flex flex-col p-0"><div class="flex justify-between items-start gap-2 p-6 pb-2 border-b border-base-200 shrink-0"><div><h3 class="font-bold text-lg">${escape_html(userForm?.name)}</h3> <p class="text-sm text-base-content/60 mt-0.5">権限・操作範囲・所属</p></div> <button type="button" class="btn btn-sm btn-circle btn-ghost" aria-label="閉じる">×</button></div> `);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--> <div class="modal-action border-t border-base-200 p-4 shrink-0"><button type="button" class="btn btn-ghost">閉じる</button> <button type="button" class="btn btn-primary">保存</button></div></div> <form method="dialog" class="modal-backdrop"><button type="button">close</button></form></dialog> <dialog${attr_class("modal", void 0, { "modal-open": orgModalOpen })}><div class="modal-box max-w-2xl max-h-[85vh] flex flex-col p-0"><div class="flex justify-between items-start gap-2 p-6 pb-2 border-b border-base-200 shrink-0"><div><h3 class="font-bold text-lg">${escape_html(orgDetail?.name)}</h3> <p class="text-sm text-base-content/60">${escape_html("部署")}</p></div> <div class="flex gap-1">`);
    {
      $$renderer2.push(`<!--[0--><button type="button" class="btn btn-sm btn-outline">編集</button>`);
    }
    $$renderer2.push(`<!--]--> <button type="button" class="btn btn-sm btn-circle btn-ghost" aria-label="閉じる">×</button></div></div> `);
    {
      $$renderer2.push(`<!--[0--><div class="overflow-y-auto px-6 py-4 flex-1">`);
      if ([].length) {
        $$renderer2.push(`<!--[0--><ul class="space-y-2"><!--[-->`);
        const each_array_13 = ensure_array_like([]);
        for (let $$index_13 = 0, $$length = each_array_13.length; $$index_13 < $$length; $$index_13++) {
          let m = each_array_13[$$index_13];
          $$renderer2.push(`<li class="flex justify-between items-center rounded-xl border border-base-300 px-3 py-2"><span class="font-medium">${escape_html(m.userName)}</span> <span class="badge badge-ghost badge-sm">${escape_html(m.title)}</span></li>`);
        }
        $$renderer2.push(`<!--]--></ul>`);
      } else {
        $$renderer2.push(`<!--[-1--><p class="text-sm text-base-content/50 text-center py-8">所属メンバーはいません</p>`);
      }
      $$renderer2.push(`<!--]--></div>`);
    }
    $$renderer2.push(`<!--]--> `);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--> `);
    {
      $$renderer2.push(`<!--[0--><div class="modal-action border-t border-base-200 p-4 shrink-0"><button type="button" class="btn btn-error btn-outline btn-sm mr-auto">削除</button> <button type="button" class="btn btn-ghost">閉じる</button></div>`);
    }
    $$renderer2.push(`<!--]--> `);
    {
      $$renderer2.push("<!--[-1-->");
    }
    $$renderer2.push(`<!--]--></div> <form method="dialog" class="modal-backdrop"><button type="button">close</button></form></dialog> <dialog${attr_class("modal", void 0, { "modal-open": createUserOpen })}><div class="modal-box max-w-lg max-h-[85vh] flex flex-col p-0"><div class="p-6 pb-2 border-b border-base-200 shrink-0"><h3 class="font-bold text-lg">ユーザを追加</h3> <p class="text-sm text-base-content/60 mt-0.5">権限はメンバー、ロールはなしで作成されます</p></div> <div class="overflow-y-auto px-6 py-4 space-y-4 flex-1"><label class="form-control"><span class="label-text font-medium">氏名</span><input class="input input-bordered"${attr("value", createUserForm.name)}/></label> <label class="form-control"><span class="label-text font-medium">メール</span><input type="email" class="input input-bordered"${attr("value", createUserForm.email)}/></label> <label class="form-control"><span class="label-text font-medium">所属部署</span> `);
    $$renderer2.select(
      {
        class: "select select-bordered",
        value: createUserForm.departmentId
      },
      ($$renderer3) => {
        $$renderer3.option({ value: 0 }, ($$renderer4) => {
          $$renderer4.push(`選択してください`);
        });
        $$renderer3.push(`<!--[-->`);
        const each_array_19 = ensure_array_like(options.departments);
        for (let $$index_19 = 0, $$length = each_array_19.length; $$index_19 < $$length; $$index_19++) {
          let d = each_array_19[$$index_19];
          $$renderer3.option({ value: d.id }, ($$renderer4) => {
            $$renderer4.push(`${escape_html(d.name)}`);
          });
        }
        $$renderer3.push(`<!--]-->`);
      }
    );
    $$renderer2.push(`</label> <label class="form-control"><span class="label-text font-medium">役職（マスタ）</span> `);
    $$renderer2.select(
      {
        class: "select select-bordered",
        value: createUserForm.roleId
      },
      ($$renderer3) => {
        $$renderer3.option({ value: 0 }, ($$renderer4) => {
          $$renderer4.push(`選択してください`);
        });
        $$renderer3.push(`<!--[-->`);
        const each_array_20 = ensure_array_like(options.roles);
        for (let $$index_20 = 0, $$length = each_array_20.length; $$index_20 < $$length; $$index_20++) {
          let r = each_array_20[$$index_20];
          $$renderer3.option({ value: r.id }, ($$renderer4) => {
            $$renderer4.push(`${escape_html(r.name)}`);
          });
        }
        $$renderer3.push(`<!--]-->`);
      }
    );
    $$renderer2.push(`</label> <label class="form-control"><span class="label-text font-medium">パスワード（任意・Basic認証用）</span><input type="password" class="input input-bordered"${attr("value", createUserForm.password)} autocomplete="new-password"/></label></div> <div class="modal-action border-t border-base-200 p-4 shrink-0"><button type="button" class="btn btn-ghost">キャンセル</button> <button type="button" class="btn btn-primary">作成</button></div></div> <form method="dialog" class="modal-backdrop"><button type="button">close</button></form></dialog> <dialog${attr_class("modal", void 0, { "modal-open": createDeptOpen })}><div class="modal-box max-w-md max-h-[85vh] flex flex-col p-0"><div class="p-6 pb-2 border-b border-base-200 shrink-0"><h3 class="font-bold text-lg">部署を追加</h3></div> <div class="overflow-y-auto px-6 py-4 space-y-4 flex-1"><label class="form-control"><span class="label-text font-medium">部署名</span><input class="input input-bordered"${attr("value", createDeptForm.name)}/></label> <label class="form-control"><span class="label-text font-medium">親部署（任意）</span> `);
    $$renderer2.select(
      {
        class: "select select-bordered",
        value: createDeptForm.parentDepartmentId
      },
      ($$renderer3) => {
        $$renderer3.option({ value: 0 }, ($$renderer4) => {
          $$renderer4.push(`なし（ルート）`);
        });
        $$renderer3.push(`<!--[-->`);
        const each_array_21 = ensure_array_like(options.departments);
        for (let $$index_21 = 0, $$length = each_array_21.length; $$index_21 < $$length; $$index_21++) {
          let d = each_array_21[$$index_21];
          $$renderer3.option({ value: d.id }, ($$renderer4) => {
            $$renderer4.push(`${escape_html(d.name)}`);
          });
        }
        $$renderer3.push(`<!--]-->`);
      }
    );
    $$renderer2.push(`</label></div> <div class="modal-action border-t border-base-200 p-4 shrink-0"><button type="button" class="btn btn-ghost">キャンセル</button> <button type="button" class="btn btn-primary">作成</button></div></div> <form method="dialog" class="modal-backdrop"><button type="button">close</button></form></dialog> <dialog${attr_class("modal", void 0, { "modal-open": confirmOpen })}><div class="modal-box max-h-[85vh]"><h3 class="font-bold">${escape_html(confirmTitle)}</h3> <p class="py-4 text-sm text-base-content/80">${escape_html(confirmMessage)}</p> <div class="modal-action"><button type="button" class="btn btn-ghost">キャンセル</button> <button type="button"${attr_class("btn", void 0, { "btn-error": confirmDanger, "btn-primary": !confirmDanger })}>${escape_html("変更する")}</button></div></div> <form method="dialog" class="modal-backdrop"><button type="button">close</button></form></dialog>`);
  });
}
function _page($$renderer) {
  ManagementView($$renderer);
}
export {
  _page as default
};

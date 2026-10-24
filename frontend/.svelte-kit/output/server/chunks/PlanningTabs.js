import { a as attr_class } from "./index.js";
function PlanningTabs($$renderer, $$props) {
  let { active } = $$props;
  $$renderer.push(`<div class="tabs tabs-boxed w-fit"><a href="/planning/meetings"${attr_class("tab", void 0, { "tab-active": active === "meetings" })}>面談</a> <a href="/planning/pairs"${attr_class("tab", void 0, { "tab-active": active === "pairs" })}>ペア</a> <a href="/planning/templates"${attr_class("tab", void 0, { "tab-active": active === "templates" })}>テンプレート</a></div>`);
}
export {
  PlanningTabs as P
};

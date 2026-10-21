/** @param {() => string} labelOf */
export function filterItems(items, query, labelOf) {
  const q = (query || "").trim().toLowerCase();
  if (!q) return items;
  return items.filter((item) => labelOf(item).toLowerCase().includes(q));
}

/** @param {object} opts */
export function createPickController(opts) {
  const { items, getId, labelOf, getSelectedId, setSelectedId } = opts;
  return {
    query: "",
    open: false,
    get filtered() {
      const list = typeof items === "function" ? items() : items;
      return filterItems(list, this.query, labelOf);
    },
    get selected() {
      const list = typeof items === "function" ? items() : items;
      const id = getSelectedId();
      return list.find((item) => getId(item) === id) || null;
    },
    displayValue() {
      if (this.open) return this.query;
      const sel = this.selected;
      return sel ? labelOf(sel) : "";
    },
    onFocus() {
      this.open = true;
      if (this.selected) this.query = "";
    },
    onInput(value) {
      this.query = value;
      this.open = true;
    },
    pick(item) {
      setSelectedId(getId(item));
      this.query = "";
      this.open = false;
    },
    clear() {
      setSelectedId(0);
      this.query = "";
      this.open = false;
    },
    close() {
      this.open = false;
    },
  };
}

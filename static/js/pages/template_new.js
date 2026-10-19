document.addEventListener("alpine:init", () => {
  Alpine.data("templateNewPage", () => ({
    title: "",
    topics: [{ question: "", isRequire: true, displayNo: 1 }],
    async save() {
      await fetch("/api/v1/templates", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title: this.title, topics: this.topics }),
      });
      location.href = "/template";
    },
  }));
});

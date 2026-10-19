document.addEventListener("alpine:init", () => {
  const el = document.getElementById("profile-json");
  const form = el ? JSON.parse(el.textContent || "{}") : {};
  Alpine.data("settingPage", () => ({
    form: { name: form.name || "", email: form.email || "", remindEnabled: !!form.remindEnabled },
    async saveProfile() {
      await fetch("/api/v1/setting", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(this.form),
      });
      alert("保存しました");
    },
    async deleteIcon() {
      await fetch("/api/v1/setting/icon", { method: "DELETE" });
      alert("アイコンを削除しました");
    },
  }));
});

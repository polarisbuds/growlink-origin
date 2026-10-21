import { showToast } from "/static/js/notice.js";

document.addEventListener("alpine:init", () => {
  const el = document.getElementById("profile-json");
  const form = el ? JSON.parse(el.textContent || "{}") : {};
  Alpine.data("settingPage", () => ({
    form: {
      name: form.name || "",
      email: form.email || "",
      lineWorksId: form.lineWorksId || "",
      remindEnabled: !!form.remindEnabled,
    },
    imageUri: form.imageUri || null,
    fieldError: "",
    get avatarUrl() {
      return this.imageUri ? `/uploads/${this.imageUri}` : null;
    },
    get initial() {
      return (this.form.name || "?").slice(0, 1);
    },
    async saveProfile() {
      this.fieldError = "";
      const res = await fetch("/api/v1/setting", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(this.form),
      });
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        this.fieldError = err.error || "保存に失敗しました";
        showToast(this.fieldError, "error");
        return;
      }
      showToast("保存しました");
    },
    async uploadIcon(e) {
      const file = e.target.files?.[0];
      if (!file) return;
      const fd = new FormData();
      fd.append("icon", file);
      const res = await fetch("/api/v1/setting/icon", { method: "POST", body: fd });
      if (!res.ok) {
        showToast("アップロードに失敗しました", "error");
        return;
      }
      const profile = await fetch("/api/v1/setting/profile");
      if (profile.ok) {
        const data = await profile.json();
        this.imageUri = data.imageUri;
      }
      showToast("アイコンをアップロードしました");
      e.target.value = "";
    },
    async deleteIcon() {
      const res = await fetch("/api/v1/setting/icon", { method: "DELETE" });
      if (!res.ok) {
        showToast("削除に失敗しました", "error");
        return;
      }
      this.imageUri = null;
      showToast("アイコンを削除しました");
    },
  }));
});

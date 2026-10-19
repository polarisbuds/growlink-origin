import { meetingPageData } from "./common.js";

document.addEventListener("alpine:init", () => {
  const base = meetingPageData();
  Alpine.data("meetingViewPage", () => ({
    ...base,
    comment: "",
    async postComment(isLimited) {
      if (!this.comment.trim()) return;
      await this.api("POST", `/api/v1/meetings/${this.uuid}/meeting-comment`, {
        content: this.comment,
        isLimited,
      });
      this.comment = "";
      this.toast("投稿しました");
    },
  }));
});

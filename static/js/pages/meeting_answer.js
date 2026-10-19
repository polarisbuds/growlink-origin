import { meetingPageData } from "./common.js";

document.addEventListener("alpine:init", () => {
  const base = meetingPageData();
  const drafts = {};
  for (const a of base.answers) drafts[a.id] = a.content || "";
  Alpine.data("meetingAnswerPage", () => ({
    ...base,
    drafts,
    answerPayload() {
      return base.answers.map((a) => ({ ...a, content: this.drafts[a.id] || "" }));
    },
    async saveDraft() {
      await this.api("POST", `/api/v1/meetings/${this.uuid}/draft`, { answers: this.answerPayload() });
      this.toast("下書きを保存しました");
    },
    async submit() {
      await this.api("POST", `/api/v1/meetings/${this.uuid}/submit`, { answers: this.answerPayload() });
      this.toast("提出しました");
      location.href = "/meeting";
    },
    async cancelSubmit() {
      await this.api("POST", `/api/v1/meetings/${this.uuid}/draft/cancel`, null);
      this.toast("提出を取り消しました");
    },
  }));
});

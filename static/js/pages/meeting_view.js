import { meetingPageData } from "./common.js";

function formatDate(iso) {
  if (!iso) return "—";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "—";
  return d.toLocaleDateString("ja-JP");
}

document.addEventListener("alpine:init", () => {
  const base = meetingPageData();
  Alpine.data("meetingViewPage", () => ({
    ...base,
    comment: "",
    formatDate,
    displayTopics() {
      const answersByTopic = new Map();
      for (const a of this.answers || []) {
        const tid = a.topic?.id || a.topicId;
        if (tid) answersByTopic.set(tid, a);
      }
      const topics = this.topics || [];
      if (topics.length) {
        return topics.map((topic, i) => {
          const answer = answersByTopic.get(topic.id);
          return {
            topic,
            answer,
            displayNo: topic.displayNo || i + 1,
          };
        });
      }
      return (this.answers || []).map((a, i) => ({
        topic: a.topic,
        answer: a,
        displayNo: a.topic?.displayNo || i + 1,
      }));
    },
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

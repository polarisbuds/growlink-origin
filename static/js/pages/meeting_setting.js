import { meetingPageData } from "./common.js";

document.addEventListener("alpine:init", () => {
  const base = meetingPageData();
  const d = base.meeting.implementedAt ? new Date(base.meeting.implementedAt) : new Date();
  const date = d.toISOString().slice(0, 10);
  Alpine.data("meetingSettingPage", () => ({
    date,
    async save() {
      alert("実施日の API は今後 PUT /api/v1/meetings/:uuid/setting で接続します（表示のみ更新可）");
    },
  }));
});

import { meetingPageData } from "./common.js";
import { showToast } from "../notice.js";

document.addEventListener("alpine:init", () => {
  const base = meetingPageData();
  const d = base.meeting.implementedAt ? new Date(base.meeting.implementedAt) : new Date();
  const date = d.toISOString().slice(0, 10);
  Alpine.data("meetingSettingPage", () => ({
    date,
    async save() {
      showToast("実施日の API は今後接続されます（表示のみ更新可）", "info");
    },
  }));
});

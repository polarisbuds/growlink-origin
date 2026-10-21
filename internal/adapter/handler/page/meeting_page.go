package page

import (
	"net/http"

	"growlink/internal/adapter/handler/authctx"
	"growlink/internal/application/converter"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func (h *Handler) MeetingListPage(c *echo.Context) error {
	claims := authctx.Claims(c)
	data := h.baseData(c, "meeting")
	meetings, err := h.deps.Meeting.List.ListInvolvedMeetings(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	data["Meetings"] = meetings
	data["MeetingsJSON"] = mustJSON(meetings)
	return c.Render(http.StatusOK, "meeting_list.html", data)
}

func (h *Handler) MeetingNewPage(c *echo.Context) error {
	data := h.baseData(c, "meeting-new")
	pairings, err := h.deps.Management.Pairing.ListPairingsByAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	templates, err := h.deps.Management.Template.ListTemplate()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	data["PairingsJSON"] = mustJSON(pairings)
	data["TemplatesJSON"] = mustJSON(templates)
	data["Breadcrumb"] = "1on1作成"
	return c.Render(http.StatusOK, "meeting_new.html", data)
}

func (h *Handler) MeetingAllPage(c *echo.Context) error {
	data := h.baseData(c, "meeting-all")
	meetings, err := h.deps.MeetingRepo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	outs := converter.MeetingEntitiesToOutputs(meetings)
	data["Meetings"] = outs
	data["Simple"] = c.QueryParam("view") == "simple"
	data["MeetingsJSON"] = mustJSON(outs)
	return c.Render(http.StatusOK, "meeting_all.html", data)
}

func (h *Handler) MeetingDetailPage(c *echo.Context) error {
	data := h.baseData(c, "meeting-detail")
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.Redirect(http.StatusFound, "/meeting")
	}
	meeting, err := h.deps.MeetingRepo.GetByUUID(id)
	if err != nil {
		return c.Redirect(http.StatusFound, "/meeting")
	}
	out := converter.MeetingEntityToOutput(*meeting)
	answers, _ := h.deps.Meeting.AnswerReader.ListAnswersByMeeting(meeting.ID)
	data["Meeting"] = out
	data["Answers"] = answers
	data["MeetingJSON"] = mustJSON(map[string]interface{}{"meeting": out, "answers": answers})
	return c.Render(http.StatusOK, "meeting_detail.html", data)
}

func (h *Handler) meetingPageData(c *echo.Context, active string, meetingUUID uuid.UUID) (map[string]interface{}, error) {
	data := h.baseData(c, active)
	meeting, err := h.deps.MeetingRepo.GetByUUID(meetingUUID)
	if err != nil {
		return nil, err
	}
	out := converter.MeetingEntityToOutput(*meeting)
	answers, _ := h.deps.Meeting.AnswerReader.ListAnswersByMeeting(meeting.ID)
	data["Meeting"] = out
	data["Answers"] = answers
	data["MeetingJSON"] = mustJSON(map[string]interface{}{"meeting": out, "answers": answers, "uuid": meetingUUID.String()})
	data["Breadcrumb"] = "面談"
	return data, nil
}

func (h *Handler) MeetingAnswerPage(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.Redirect(http.StatusFound, "/meeting")
	}
	data, err := h.meetingPageData(c, "meeting-answer", id)
	if err != nil {
		return c.Redirect(http.StatusFound, "/meeting")
	}
	data["Breadcrumb"] = "記入"
	return c.Render(http.StatusOK, "meeting_answer.html", data)
}

func (h *Handler) MeetingSettingPage(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.Redirect(http.StatusFound, "/meeting")
	}
	data, err := h.meetingPageData(c, "meeting-setting", id)
	if err != nil {
		return c.Redirect(http.StatusFound, "/meeting")
	}
	data["Breadcrumb"] = "ミーティング設定"
	return c.Render(http.StatusOK, "meeting_setting.html", data)
}

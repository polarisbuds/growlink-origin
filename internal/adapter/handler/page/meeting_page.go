package page

import (
	"net/http"

	"growlink/internal/adapter/handler/authctx"
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"

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
	users, err := h.deps.UserRepo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	userOuts := make([]outputmodel.UserOutput, 0, len(users))
	for _, u := range users {
		if u.Suspended {
			continue
		}
		userOuts = append(userOuts, outputmodel.UserOutput{ID: u.ID, Name: u.Name, Email: u.Email, ImageURI: u.ImageURI})
	}
	data["PairingsJSON"] = mustJSON(pairings)
	data["TemplatesJSON"] = mustJSON(templates)
	data["UsersJSON"] = mustJSON(userOuts)
	data["Breadcrumb"] = "1on1作成"
	return c.Render(http.StatusOK, "meeting_new.html", data)
}

func (h *Handler) MeetingAllPage(c *echo.Context) error {
	data := h.baseData(c, "planning-meetings")
	data["PlanningTab"] = "meetings"
	data["Breadcrumb"] = "面談一覧"
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
	payload, err := h.buildMeetingPagePayload(c, id)
	if err != nil {
		return c.Redirect(http.StatusFound, "/meeting")
	}
	data["Meeting"] = payload["meeting"]
	data["Answers"] = payload["answers"]
	data["MeetingJSON"] = mustJSON(payload)
	return c.Render(http.StatusOK, "meeting_detail.html", data)
}

func (h *Handler) buildMeetingPagePayload(c *echo.Context, meetingUUID uuid.UUID) (map[string]interface{}, error) {
	meeting, err := h.deps.MeetingRepo.GetByUUID(meetingUUID)
	if err != nil {
		return nil, err
	}
	out := converter.MeetingEntityToOutput(*meeting)
	answers, _ := h.deps.Meeting.AnswerReader.ListAnswersByMeeting(meeting.ID)
	topics := out.Template.Topics
	if len(topics) == 0 && meeting.Template != nil {
		topics = converter.TopicEntitiesToOutputs(meeting.Template.Topics)
	}
	canEdit := false
	if claims := authctx.Claims(c); claims != nil {
		canEdit = authctx.CanModify(claims.UserID, claims.SystemRole, meeting.CreatedByUserID)
	}
	return map[string]interface{}{
		"meeting":  out,
		"answers":  answers,
		"topics":   topics,
		"canEdit":  canEdit,
		"uuid":     meetingUUID.String(),
	}, nil
}

func (h *Handler) meetingPageData(c *echo.Context, active string, meetingUUID uuid.UUID) (map[string]interface{}, error) {
	data := h.baseData(c, active)
	payload, err := h.buildMeetingPagePayload(c, meetingUUID)
	if err != nil {
		return nil, err
	}
	data["Meeting"] = payload["meeting"]
	data["Answers"] = payload["answers"]
	data["MeetingJSON"] = mustJSON(payload)
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

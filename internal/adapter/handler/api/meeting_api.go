package api

import (
	"net/http"

	"growlink/internal/adapter/handler/authctx"
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func (h *Handler) GetMeeting(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	meeting, err := h.deps.MeetingRepo.GetByUUID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	out := converter.MeetingEntityToOutput(*meeting)
	answers, _ := h.deps.Meeting.AnswerReader.ListAnswersByMeeting(meeting.ID)
	outAnswers := answers
	return c.JSON(http.StatusOK, map[string]interface{}{
		"meeting": out,
		"answers": outAnswers,
	})
}

func (h *Handler) SaveDraft(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	meeting, err := h.deps.MeetingRepo.GetByUUID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	var req struct {
		Answers []outputmodel.AnswerOutput `json:"answers"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Meeting.DraftSaver.SaveDraft(meeting.ID, req.Answers); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) SubmitMeeting(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	meeting, err := h.deps.MeetingRepo.GetByUUID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	var req struct {
		Answers []outputmodel.AnswerOutput `json:"answers"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Meeting.Submitter.Submit(meeting.ID, req.Answers); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) CancelSubmit(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	meeting, err := h.deps.MeetingRepo.GetByUUID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if err := h.deps.Meeting.SubmitCanceller.CancelSubmit(meeting.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) PostMeetingComment(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	meeting, err := h.deps.MeetingRepo.GetByUUID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var req struct {
		Content   string `json:"content"`
		IsLimited bool   `json:"isLimited"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Meeting.MeetingCommenter.PostComment(meeting.ID, claims.UserID, req.Content, req.IsLimited); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) PostAnswerComment(c *echo.Context) error {
	var req struct {
		AnswerID uint   `json:"answerId"`
		Content  string `json:"content"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Meeting.AnswerCommenter.PostComment(req.AnswerID, req.Content); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) ListMeetingsAll(c *echo.Context) error {
	meetings, err := h.deps.MeetingRepo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, converter.MeetingEntitiesToOutputs(meetings))
}

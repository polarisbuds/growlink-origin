package api

import (
	"net/http"
	"time"

	"growlink/internal/adapter/handler/authctx"
	inputmodel "growlink/internal/application/model/input_model"

	"github.com/labstack/echo/v5"
)

func parseDateOnly(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func (h *Handler) CreateMeeting(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var req inputmodel.CreateMeetingInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if req.MentorID == 0 || req.TemplateID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "mentorId and templateId are required"})
	}
	if len(req.MenteeIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "menteeIds are required"})
	}
	implementDate, ok := parseDateOnly(req.ImplementedAt)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid implementedAt"})
	}
	outs, err := h.deps.Management.MeetingCreator.CreatePlannedMeetings(
		req.MentorID,
		req.TemplateID,
		req.MenteeIDs,
		implementDate,
		claims.UserID,
		req.SharerUserIDs,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ok": "true", "count": len(outs)})
}

func (h *Handler) CreateMeetingByAll(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var req inputmodel.CreateMeetingByAllInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if req.TemplateID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "templateId is required"})
	}
	implementDate, ok := parseDateOnly(req.ImplementedAt)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid implementedAt"})
	}
	outs, err := h.deps.Management.MeetingCreator.CreateMeetingByAll(req.TemplateID, implementDate, claims.UserID, req.SharerUserIDs, req.PairingIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ok": "true", "count": len(outs)})
}

func (h *Handler) CreateMeetingMonthly(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var req inputmodel.CreateMeetingMonthlyInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if req.TemplateID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "templateId is required"})
	}
	if req.Week < 1 || req.Week > 5 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "week must be between 1 and 5"})
	}
	if req.Weekday < 0 || req.Weekday > 6 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "weekday must be between 0 and 6"})
	}
	fromDate, ok := parseDateOnly(req.From)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid from"})
	}
	toDate, ok := parseDateOnly(req.To)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid to"})
	}
	fromMonth := time.Date(fromDate.Year(), fromDate.Month(), 1, 0, 0, 0, 0, fromDate.Location())
	toMonth := time.Date(toDate.Year(), toDate.Month(), 1, 0, 0, 0, 0, toDate.Location())
	if fromMonth.After(toMonth) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "from must not be after to"})
	}
	outs, err := h.deps.Management.MeetingCreator.CreateMeetingEveryMonth(
		req.TemplateID,
		fromDate,
		toDate,
		req.Week,
		time.Weekday(req.Weekday),
		claims.UserID,
		req.SharerUserIDs,
		req.PairingIDs,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ok": "true", "count": len(outs)})
}

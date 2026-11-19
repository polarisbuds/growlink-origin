package api

import (
	"errors"
	"net/http"
	"strings"

	"growlink/internal/adapter/handler/authctx"
	outputmodel "growlink/internal/application/model/output_model"
	meeting_usecase "growlink/internal/application/usecase/meeting_usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func (h *Handler) GetMeeting(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	claims := authctx.Claims(c)
	var viewerID uint
	var authority, systemRole interface{}
	if claims != nil {
		viewerID = claims.UserID
		authority = claims.Authority
		systemRole = claims.SystemRole
	}
	detail, err := h.deps.Meeting.API.GetMeetingDetail(
		id,
		viewerID,
		authority,
		systemRole,
		authctx.CanModify,
	)
	if err != nil {
		if errors.Is(err, meeting_usecase.ErrMeetingNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"uuid":                  detail.MeetingUUID,
		"meeting":               detail.Meeting,
		"answers":               detail.Answers,
		"topics":                detail.Topics,
		"meetingComments":       detail.Comments,
		"canEdit":               detail.Flags.CanEdit,
		"viewerRole":            detail.Flags.ViewerRole,
		"canWriteAnswers":       detail.Flags.CanWriteAnswers,
		"canWriteTopicComments": detail.Flags.CanWriteTopicComments,
		"canEditSharers":        detail.Flags.CanEditSharers,
		"sharers":               detail.Sharers,
	})
}

func meetingAPIError(c *echo.Context, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, meeting_usecase.ErrMeetingForbidden) || err.Error() == "forbidden" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
	var req struct {
		Answers []outputmodel.AnswerOutput `json:"answers"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Meeting.API.SaveDraftByUUID(id, claims.UserID, req.Answers); err != nil {
		return meetingAPIError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) SubmitMeeting(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	var req struct {
		Answers []outputmodel.AnswerOutput `json:"answers"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Meeting.API.SubmitByUUID(id, claims.UserID, req.Answers); err != nil {
		return meetingAPIError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) CompleteMentorMeeting(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	var req struct {
		Comments []meeting_usecase.MentorTopicCommentInput `json:"comments"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Meeting.API.CompleteMentorByUUID(id, claims.UserID, req.Comments); err != nil {
		if err.Error() == "pairing not found" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return meetingAPIError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) CancelSubmit(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	if err := h.deps.Meeting.API.CancelSubmitByUUID(id, claims.UserID); err != nil {
		return meetingAPIError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) UpdateMeetingSharers(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	var req struct {
		SharerUserUuids []string `json:"sharerUserUuids"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	sharers, err := h.deps.Management.MeetingUpdater.ReplaceSharersByMeetingUUID(id, claims.UserID, req.SharerUserUuids)
	if err != nil {
		return meetingAPIError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"sharers": sharers})
}

func (h *Handler) ListMeetingSharerCandidates(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	q := strings.TrimSpace(c.QueryParam("q"))
	candidates, err := h.deps.Meeting.API.ListSharerCandidates(id, claims.UserID, q)
	if err != nil {
		return meetingAPIError(c, err)
	}
	return c.JSON(http.StatusOK, candidates)
}

func (h *Handler) PostMeetingComment(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	var req struct {
		Content   string `json:"content"`
		IsLimited bool   `json:"isLimited"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Meeting.API.PostMeetingCommentByUUID(id, claims.UserID, req.Content, req.IsLimited); err != nil {
		return meetingAPIError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) PostAnswerComment(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var req struct {
		AnswerUuid string `json:"answerUuid"`
		Content    string `json:"content"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if req.AnswerUuid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "answerUuid is required"})
	}
	answerUUID, err := uuid.Parse(req.AnswerUuid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	if err := h.deps.Meeting.API.PostAnswerCommentByAnswerUUID(answerUUID, claims.UserID, req.Content); err != nil {
		if resp := apiUUIDResolveError(c, err); resp != nil {
			return resp
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) ListMeetingsAll(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	filter, page, dirDesc, err := parsePlanningMeetingsListQuery(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	out, err := h.deps.Planning.ListPlanningMeetingsPaged(
		claims.UserID,
		claims.Authority,
		authctx.SystemRole(c),
		filter,
		page,
		dirDesc,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

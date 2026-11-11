package api

import (
	"net/http"
	"strings"

	"growlink/internal/adapter/handler/authctx"
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	meeting_usecase "growlink/internal/application/usecase/meeting_usecase"
	"growlink/internal/domain/valueobject"

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
	if claims := authctx.Claims(c); claims != nil {
		ok, err := h.deps.Planning.CanViewPlanningMeeting(claims.UserID, claims.Authority, authctx.SystemRole(c), meeting.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if !ok {
			ok = canViewMeeting(meeting, claims.UserID)
		}
		if !ok {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
	}
	out := converter.MeetingEntityToOutput(*meeting)
	answers, err := h.deps.Meeting.AnswerReader.ListAnswersByMeeting(meeting.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	topics := out.Template.Topics
	if len(topics) == 0 && meeting.Template != nil {
		topics = converter.TopicEntitiesToOutputs(meeting.Template.Topics)
	}
	canEdit := false
	var viewerID uint
	viewerRole := "viewer"
	canWriteAnswers := false
	allowTopicCommentWrite := false
	canEditSharers := false
	if claims := authctx.Claims(c); claims != nil {
		viewerID = claims.UserID
		canEdit = authctx.CanModify(claims.UserID, claims.SystemRole, meeting.CreatedByUserID)
		viewerRole = meetingViewerRole(meeting, viewerID)
		canWriteAnswers = canWriteMeetingAnswers(meeting, viewerID)
		allowTopicCommentWrite = canWriteTopicComments(meeting, viewerID)
		canEditSharers = canEditMeetingSharers(meeting, viewerID)
	}
	meetingComments := filterMeetingComments(meeting, viewerID)
	sharers := meetingSharersToOutputs(meeting)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"uuid":                  id.String(),
		"meeting":               out,
		"answers":               answers,
		"topics":                topics,
		"meetingComments":       meetingComments,
		"canEdit":               canEdit,
		"viewerRole":            viewerRole,
		"canWriteAnswers":       canWriteAnswers,
		"canWriteTopicComments": allowTopicCommentWrite,
		"canEditSharers":        canEditSharers,
		"sharers":               sharers,
	})
}

func (h *Handler) lookupAnswerOutput(answerID uint) (*outputmodel.AnswerOutput, error) {
	return h.deps.Meeting.AnswerReader.GetAnswer(answerID)
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
	if !canWriteMeetingAnswers(meeting, claims.UserID) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}
	var req struct {
		Answers []outputmodel.AnswerOutput `json:"answers"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	req.Answers, err = sanitizeAnswerOutputs(meeting.ID, claims.UserID, req.Answers, meeting, h.lookupAnswerOutput)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if err := h.deps.Meeting.DraftSaver.SaveDraft(meeting.ID, req.Answers, false); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
	meeting, err := h.deps.MeetingRepo.GetByUUID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if !canWriteMeetingAnswers(meeting, claims.UserID) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}
	var req struct {
		Answers []outputmodel.AnswerOutput `json:"answers"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	req.Answers, err = sanitizeAnswerOutputs(meeting.ID, claims.UserID, req.Answers, meeting, h.lookupAnswerOutput)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if err := h.deps.Meeting.Submitter.Submit(meeting.ID, req.Answers, valueobject.MeetingStatusSubmitted); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
	meeting, err := h.deps.MeetingRepo.GetByUUID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if !canWriteTopicComments(meeting, claims.UserID) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}
	if meeting.Pairing == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "pairing not found"})
	}
	var req struct {
		Comments []meeting_usecase.MentorTopicCommentInput `json:"comments"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Meeting.MentorMeetingCompleter.Complete(
		meeting.ID,
		claims.UserID,
		meeting.Pairing.MenteeID,
		req.Comments,
	); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
	meeting, err := h.deps.MeetingRepo.GetByUUID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if !canCancelMeetingSubmit(meeting, claims.UserID) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}
	if err := h.deps.Meeting.SubmitCanceller.CancelSubmit(meeting.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
	meeting, err := h.deps.MeetingRepo.GetByUUID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if !canEditMeetingSharers(meeting, claims.UserID) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}
	var req struct {
		SharerUserUuids []string `json:"sharerUserUuids"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	sharerIDs, err := h.resolveUserIDsFromUUIDs(req.SharerUserUuids)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.deps.MeetingRepo.ReplaceSharers(meeting.ID, sharerIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	updated, err := h.deps.MeetingRepo.GetByUUID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"sharers": meetingSharersToOutputs(updated),
	})
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
	meeting, err := h.deps.MeetingRepo.GetByUUID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if !canEditMeetingSharers(meeting, claims.UserID) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}
	q := strings.TrimSpace(c.QueryParam("q"))
	users, err := h.deps.UserRepo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	type candidate struct {
		UUID  string `json:"uuid"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	out := make([]candidate, 0, 32)
	for _, u := range users {
		if u.Suspended {
			continue
		}
		if !userMatchesSharerQuery(u, q) {
			continue
		}
		out = append(out, candidate{UUID: u.UUID.String(), Name: u.Name, Email: u.Email})
		if len(out) >= 50 {
			break
		}
	}
	return c.JSON(http.StatusOK, out)
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
	if meeting.Status != valueobject.MeetingStatusDone {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}
	if err := h.deps.Meeting.MeetingCommenter.PostComment(meeting.ID, claims.UserID, req.Content, req.IsLimited); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) PostAnswerComment(c *echo.Context) error {
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
	answer, err := h.deps.AnswerRepo.GetByUUID(answerUUID)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	if err := h.deps.Meeting.AnswerCommenter.PostComment(answer.ID, claims.UserID, req.Content); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) ListMeetingsAll(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	list, err := h.deps.Planning.ListPlanningMeetings(claims.UserID, claims.Authority, authctx.SystemRole(c))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

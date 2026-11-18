package api

import (
	"net/http"

	inputmodel "growlink/internal/application/model/input_model"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func apiUUIDResolveError(c *echo.Context, err error) error {
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	return nil
}

func (h *Handler) resolveUserIDFromUUID(userUUID string) (uint, error) {
	parsed, err := uuid.Parse(userUUID)
	if err != nil {
		return 0, err
	}
	user, err := h.deps.UserRepo.GetByUUID(parsed)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (h *Handler) resolveUserIDFromUUIDs(uuids []string) ([]uint, error) {
	if len(uuids) == 0 {
		return nil, nil
	}
	ids, err := h.resolveUserIDFromUUIDs(uuids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (h *Handler) resolveTemplateIDFromUUIDString(templateUUID string) (uint, error) {
	parsed, err := uuid.Parse(templateUUID)
	if err != nil {
		return 0, err
	}
	tpl, err := h.deps.TemplateRepo.GetByUUID(parsed)
	if err != nil {
		return 0, err
	}
	return tpl.ID, nil
}

func (h *Handler) resolvePairingIDs(pairingUUIDs []string) ([]uint, error) {
	if len(pairingUUIDs) == 0 {
		return nil, nil
	}
	ids := make([]uint, 0, len(pairingUUIDs))
	pairs, err := h.deps.PairingRepo.GetByUUIDs(pairingUUIDs)
	if err != nil {
		return nil, err
	}
	for _, pair := range pairs {
		ids = append(ids, pair.ID)
	}
	return ids, nil
}

func (h *Handler) resolvePairingInput(c *echo.Context, req *inputmodel.PairingInput) error {
	if req.MentorUuid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "mentorUuid is required"})
	}
	mentorID, err := h.resolveUserIDFromUUID(req.MentorUuid)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	req.MentorID = mentorID
	if len(req.Targets) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "targets are required"})
	}
	for i := range req.Targets {
		t := &req.Targets[i]
		if t.TargetUuid == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "targetUuid is required"})
		}
		if respErr := apiUUIDResolveError(c, err); respErr != nil {
			return respErr
		}
	}
	return nil
}

func (h *Handler) resolveCreateMeetingInput(c *echo.Context, req *inputmodel.CreateMeetingInput) error {
	if req.MentorUuid == "" || req.TemplateUuid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "mentorUuid and templateUuid are required"})
	}
	if len(req.MenteeUuids) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "menteeUuids are required"})
	}
	mentorID, err := h.resolveUserIDFromUUID(req.MentorUuid)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	templateID, err := h.resolveTemplateIDFromUUIDString(req.TemplateUuid)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	menteeIDs, err := h.resolveUserIDFromUUIDs(req.MenteeUuids)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	sharerIDs, err := h.resolveUserIDFromUUIDs(req.SharerUserUuids)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	req.MentorID = mentorID
	req.TemplateID = templateID
	req.MenteeIDs = menteeIDs
	req.SharerUserIDs = sharerIDs
	return nil
}

func (h *Handler) resolveCreateMeetingByAllInput(c *echo.Context, req *inputmodel.CreateMeetingByAllInput) error {
	if req.TemplateUuid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "templateUuid is required"})
	}
	templateID, err := h.resolveTemplateIDFromUUIDString(req.TemplateUuid)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	sharerIDs, err := h.resolveUserIDFromUUIDs(req.SharerUserUuids)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	pairingIDs, err := h.resolvePairingIDs(req.PairingUuids)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	req.TemplateID = templateID
	req.SharerUserIDs = sharerIDs
	req.PairingIDs = pairingIDs
	return nil
}

func (h *Handler) resolveCreateMeetingMonthlyInput(c *echo.Context, req *inputmodel.CreateMeetingMonthlyInput) error {
	if req.MentorUuid == "" || req.TemplateUuid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "mentorUuid and templateUuid are required"})
	}
	if len(req.MenteeUuids) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "menteeUuids are required"})
	}
	mentorID, err := h.resolveUserIDFromUUID(req.MentorUuid)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	templateID, err := h.resolveTemplateIDFromUUIDString(req.TemplateUuid)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	menteeIDs, err := h.resolveUserIDFromUUIDs(req.MenteeUuids)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	sharerIDs, err := h.resolveUserIDFromUUIDs(req.SharerUserUuids)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	req.MentorID = mentorID
	req.TemplateID = templateID
	req.MenteeIDs = menteeIDs
	req.SharerUserIDs = sharerIDs
	return nil
}

func (h *Handler) resolveCreateMeetingMonthlyByPairingsInput(c *echo.Context, req *inputmodel.CreateMeetingMonthlyByPairingsInput) error {
	if req.TemplateUuid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "templateUuid is required"})
	}
	if len(req.PairingUuids) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "pairingUuids are required"})
	}
	templateID, err := h.resolveTemplateIDFromUUIDString(req.TemplateUuid)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	sharerIDs, err := h.resolveUserIDFromUUIDs(req.SharerUserUuids)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	pairingIDs, err := h.resolvePairingIDs(req.PairingUuids)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	req.TemplateID = templateID
	req.SharerUserIDs = sharerIDs
	req.PairingIDs = pairingIDs
	return nil
}

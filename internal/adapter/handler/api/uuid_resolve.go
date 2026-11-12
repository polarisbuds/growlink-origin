package api

import (
	"fmt"
	"net/http"

	inputmodel "growlink/internal/application/model/input_model"
	"growlink/internal/domain/entity"

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

func (h *Handler) resolveUserIDsFromUUIDs(uuids []string) ([]uint, error) {
	if len(uuids) == 0 {
		return nil, nil
	}
	ids := make([]uint, 0, len(uuids))
	for _, s := range uuids {
		id, err := h.resolveUserIDFromUUID(s)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
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
	for _, s := range pairingUUIDs {
		parsed, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		pair, err := h.deps.PairingRepo.GetByUUID(parsed)
		if err != nil {
			return nil, err
		}
		ids = append(ids, pair.ID)
	}
	return ids, nil
}

func (h *Handler) resolvePairingTargetID(kind, targetUUID string) (uint, error) {
	if kind != entity.PairingTargetKindUser {
		return 0, fmt.Errorf("pairing target kind %q is not supported", kind)
	}
	parsed, err := uuid.Parse(targetUUID)
	if err != nil {
		return 0, err
	}
	user, err := h.deps.UserRepo.GetByUUID(parsed)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
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
		targetID, err := h.resolvePairingTargetID(t.Kind, t.TargetUuid)
		if respErr := apiUUIDResolveError(c, err); respErr != nil {
			return respErr
		}
		t.TargetID = targetID
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
	menteeIDs, err := h.resolveUserIDsFromUUIDs(req.MenteeUuids)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	sharerIDs, err := h.resolveUserIDsFromUUIDs(req.SharerUserUuids)
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
	sharerIDs, err := h.resolveUserIDsFromUUIDs(req.SharerUserUuids)
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
	menteeIDs, err := h.resolveUserIDsFromUUIDs(req.MenteeUuids)
	if respErr := apiUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	sharerIDs, err := h.resolveUserIDsFromUUIDs(req.SharerUserUuids)
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
	sharerIDs, err := h.resolveUserIDsFromUUIDs(req.SharerUserUuids)
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

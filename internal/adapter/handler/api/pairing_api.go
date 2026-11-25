package api

import (
	"errors"
	"net/http"

	"growlink/internal/adapter/handler/authctx"
	inputmodel "growlink/internal/application/model/input_model"
	planning_usecase "growlink/internal/application/usecase/planning_usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func (h *Handler) ListPairings(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	systemRole := authctx.SystemRole(c)
	page, hasPage, err := parsePairingsListPage(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	filter, err := parsePlanningPairingsListFilter(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if hasPage {
		out, err := h.deps.Planning.ListPairingsForManagementPaged(
			claims.UserID,
			claims.Authority,
			systemRole,
			filter,
			page,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, out)
	}
	if pairingsListQueryHasFilters(c) || !planningPairingsFilterEmpty(filter) {
		list, err := h.deps.Planning.ListPairingsForManagementFiltered(
			claims.UserID,
			claims.Authority,
			systemRole,
			filter,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, list)
	}
	list, err := h.deps.Planning.ListPairingsForManagement(claims.UserID, claims.Authority, systemRole)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *Handler) ListPairingsForMeeting(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	list, err := h.deps.Planning.ListPairingsForMeeting(claims.UserID, claims.Authority, authctx.SystemRole(c))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *Handler) GetPairing(c *echo.Context) error {
	pairUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	pair, err := h.deps.Management.Pairing.GetPairingByUUID(pairUUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, pair)
}

func (h *Handler) CreatePairing(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var req inputmodel.PairingInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := validatePairingInput(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	menteeUUIDs := make([]string, len(req.Targets))
	for idx, pairingTarget := range req.Targets {
		menteeUUIDs[idx] = pairingTarget.TargetUuid
	}

	menteeIDs, err := h.deps.Management.OrgAdmin.ResolveUserIDs(menteeUUIDs)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid target uuid"})
	}

	for _, menteeID := range menteeIDs {
		ok, err := h.deps.Planning.CanManagePairingMentee(claims.UserID, claims.Authority, claims.SystemRole, menteeID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if !ok {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
	}
	if err := h.deps.Management.PairingCreator.CreatePairingFromInput(req, claims.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) ListPairingSharers(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	sharers, err := h.deps.Planning.ListPairingSharersByUUID(
		claims.UserID,
		claims.Authority,
		claims.SystemRole,
		c.Param("uuid"),
	)
	if err != nil {
		if errors.Is(err, planning_usecase.ErrPairingSharerForbidden) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		if resp := apiUUIDResolveError(c, err); resp != nil {
			return resp
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"sharers": sharers})
}

func (h *Handler) UpdatePairingSharers(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var req struct {
		SharerUserUuids []string `json:"sharerUserUuids"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	sharers, err := h.deps.Planning.ReplacePairingSharersByUUID(
		claims.UserID,
		claims.Authority,
		claims.SystemRole,
		c.Param("uuid"),
		req.SharerUserUuids,
	)
	if err != nil {
		if errors.Is(err, planning_usecase.ErrPairingSharerForbidden) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		if errors.Is(err, planning_usecase.ErrPairingSharerInvalid) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid sharer uuid"})
		}
		if resp := apiUUIDResolveError(c, err); resp != nil {
			return resp
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"sharers": sharers})
}

func (h *Handler) DeletePairing(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	if err := h.deps.Planning.BreakPairingByUUID(
		claims.UserID,
		claims.Authority,
		claims.SystemRole,
		c.Param("uuid"),
		h.deps.Management.PairingDeleter,
	); err != nil {
		if errors.Is(err, planning_usecase.ErrBreakPairingForbidden) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		if resp := apiUUIDResolveError(c, err); resp != nil {
			return resp
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func validatePairingInput(req *inputmodel.PairingInput) error {
	if req.MentorUuid == "" {
		return errors.New("mentorUuid is required")
	}
	if len(req.Targets) == 0 {
		return errors.New("targets are required")
	}
	for i := range req.Targets {
		if req.Targets[i].TargetUuid == "" {
			return errors.New("targetUuid is required")
		}
	}
	return nil
}

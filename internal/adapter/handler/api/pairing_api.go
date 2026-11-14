package api

import (
	"net/http"
	"strconv"

	"growlink/internal/adapter/handler/authctx"
	inputmodel "growlink/internal/application/model/input_model"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func (h *Handler) resolvePairingID(param string) (uint, error) {
	if id64, err := strconv.ParseUint(param, 10, 64); err == nil {
		return uint(id64), nil
	}
	pairUUID, err := uuid.Parse(param)
	if err != nil {
		return 0, err
	}
	pair, err := h.deps.PairingRepo.GetByUUID(pairUUID)
	if err != nil {
		return 0, err
	}
	return pair.ID, nil
}

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
	if err := h.resolvePairingInput(c, &req); err != nil {
		return err
	}
	for _, t := range req.ResolvedTargets() {
		ok, err := h.deps.Planning.CanManagePairingMentee(claims.UserID, claims.Authority, claims.SystemRole, t.TargetID)
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

func (h *Handler) DeletePairing(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	pairingID, err := h.resolvePairingID(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid pairing id"})
	}
	pair, err := h.deps.PairingRepo.Get(pairingID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	ok, err := h.deps.Planning.CanBreakPairing(claims.UserID, claims.Authority, claims.SystemRole, pair.CreatedByUserID, pair.MenteeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !ok {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}
	if err := h.deps.Management.PairingDeleter.BreakPairing(pairingID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

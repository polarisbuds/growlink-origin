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
	list, err := h.deps.Management.Pairing.ListPairingsByAll()
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
	if err := h.deps.Management.PairingCreator.CreatePairingFromInput(req, claims.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) UpdatePairing(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var req inputmodel.PairingInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	pairingID, err := h.resolvePairingID(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid pairing id"})
	}
	pair, err := h.deps.PairingRepo.Get(pairingID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if !authctx.CanModify(claims.UserID, claims.SystemRole, pair.CreatedByUserID) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}
	if err := h.deps.Management.PairingUpdater.UpdatePairing(pairingID, req); err != nil {
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
	if !authctx.CanModify(claims.UserID, claims.SystemRole, pair.CreatedByUserID) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}
	if err := h.deps.Management.PairingDeleter.BreakPairing(pairingID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

package api

import (
	"net/http"

	"growlink/internal/adapter/handler/authctx"

	"github.com/labstack/echo/v5"
)

func (h *Handler) ListPlanningDirectory(c *echo.Context) error {
	claims := authctx.Claims(c)
	viewerID := uint(0)
	includeManaged := false
	if claims != nil {
		viewerID = claims.UserID
		includeManaged = h.deps.Planning != nil
	}
	out, err := h.deps.Planning.ListDirectory(viewerID, includeManaged)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}
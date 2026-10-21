package api

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *Handler) ProgressSummary(c *echo.Context) error {
	summary, err := h.deps.Management.Progress.Summary()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, summary)
}

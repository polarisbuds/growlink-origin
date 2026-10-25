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

func (h *Handler) ProgressBreakdown(c *echo.Context) error {
	group := c.QueryParam("group")
	if group != "department" && group != "team" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "group は department または team を指定してください"})
	}
	rows, err := h.deps.Management.Progress.Breakdown(group)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, rows)
}

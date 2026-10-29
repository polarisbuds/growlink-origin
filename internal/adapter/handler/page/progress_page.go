package page

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *Handler) ProgressPage(c *echo.Context) error {
	data := h.baseData(c, "progress")
	data["Breadcrumb"] = "進捗"
	return c.Render(http.StatusOK, "progress.html", data)
}

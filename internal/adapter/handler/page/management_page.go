package page

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *Handler) ManagementPage(c *echo.Context) error {
	data := h.baseData(c, "management")
	data["Breadcrumb"] = "管理"
	return c.Render(http.StatusOK, "management.html", data)
}

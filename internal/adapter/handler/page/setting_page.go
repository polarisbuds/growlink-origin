package page

import (
	"net/http"

	"growlink/internal/adapter/handler/authctx"

	"github.com/labstack/echo/v5"
)

func (h *Handler) SettingPage(c *echo.Context) error {
	claims := authctx.Claims(c)
	data := h.baseData(c, "setting")
	user, err := h.deps.UserRepo.GetByID(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	data["Profile"] = user
	data["ProfileJSON"] = mustJSON(map[string]interface{}{
		"name":          user.Name,
		"email":         user.Email,
		"remindEnabled": user.RemindEnabled,
	})
	data["Breadcrumb"] = "設定"
	return c.Render(http.StatusOK, "setting.html", data)
}

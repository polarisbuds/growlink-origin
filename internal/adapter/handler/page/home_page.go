package page

import (
	"net/http"

	"growlink/internal/adapter/handler/authctx"

	"github.com/labstack/echo/v5"
)

func (h *Handler) HomePage(c *echo.Context) error {
	claims := authctx.Claims(c)
	data := h.baseData(c, "home")
	own, err := h.deps.Meeting.List.ListInvolvedMeetings(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	shared, err := h.deps.Meeting.List.ListSharedMeetings(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	data["MeetingsJSON"] = mustJSON(own)
	data["SharedJSON"] = mustJSON(shared)
	data["Breadcrumb"] = "ワークスペース"
	return c.Render(http.StatusOK, "home.html", data)
}

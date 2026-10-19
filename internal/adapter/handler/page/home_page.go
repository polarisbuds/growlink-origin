package page

import (
	"net/http"
	"time"

	"growlink/internal/adapter/handler/authctx"
	outputmodel "growlink/internal/application/model/output_model"

	"github.com/labstack/echo/v5"
)

func (h *Handler) HomePage(c *echo.Context) error {
	claims := authctx.Claims(c)
	data := h.baseData(c, "home")
	meetings, err := h.deps.Meeting.List.ListMyMeeting(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	today := time.Now()
	var untilToday, fromTomorrow []outputmodel.MeetingOutput
	for _, m := range meetings {
		if m.ImplementedAt.After(today) {
			fromTomorrow = append(fromTomorrow, m)
		} else {
			untilToday = append(untilToday, m)
		}
	}
	data["MeetingsUntilToday"] = untilToday
	data["MeetingsFromTomorrow"] = fromTomorrow
	data["MeetingsJSON"] = mustJSON(meetings)
	data["Breadcrumb"] = "あなたの1on1"
	return c.Render(http.StatusOK, "home.html", data)
}

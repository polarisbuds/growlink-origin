package api

import (
	"net/http"
	"time"

	progress_mng_usecase "growlink/internal/application/usecase/management_usecase/progress_management_usecase"
	"growlink/internal/adapter/handler/authctx"

	"github.com/labstack/echo/v5"
)

func parseProgressDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func progressFilterFromRequest(c *echo.Context) (progress_mng_usecase.ProgressFilter, error) {
	claims := authctx.Claims(c)
	filter := progress_mng_usecase.ProgressFilter{}
	if claims != nil {
		filter.ViewerID = claims.UserID
		filter.SystemRole = claims.SystemRole
	}
	from, err := parseProgressDate(c.QueryParam("from"))
	if err != nil {
		return filter, err
	}
	to, err := parseProgressDate(c.QueryParam("to"))
	if err != nil {
		return filter, err
	}
	if to != nil {
		end := to.Add(24*time.Hour).Add(-time.Nanosecond)
		to = &end
	}
	filter.From = from
	filter.To = to
	return filter, nil
}

func (h *Handler) ProgressSummary(c *echo.Context) error {
	filter, err := progressFilterFromRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "from / to は YYYY-MM-DD 形式で指定してください"})
	}
	summary, err := h.deps.Management.Progress.Summary(filter)
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
	filter, err := progressFilterFromRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "from / to は YYYY-MM-DD 形式で指定してください"})
	}
	rows, err := h.deps.Management.Progress.Breakdown(group, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, rows)
}

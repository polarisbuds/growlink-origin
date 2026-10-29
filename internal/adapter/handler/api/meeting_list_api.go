package api

import (
	"net/http"

	"growlink/internal/adapter/handler/authctx"

	"github.com/labstack/echo/v5"
)

func (h *Handler) ListInvolvedMeetings(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	list, err := h.deps.Meeting.List.ListInvolvedMeetings(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *Handler) ListSharedMeetings(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	list, err := h.deps.Meeting.List.ListSharedMeetings(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

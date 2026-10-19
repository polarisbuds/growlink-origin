package api

import (
	"net/http"

	"growlink/internal/adapter/handler/authctx"

	"github.com/labstack/echo/v5"
)

func (h *Handler) UploadIcon(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	file, err := c.FormFile("icon")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid file"})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer src.Close()
	if err := h.deps.Setting.RegistIcon(claims.UserID, file.Filename, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) DeleteIcon(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	if err := h.deps.SettingRm.RemoveIcon(claims.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) UpdateSetting(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	user, err := h.deps.UserRepo.GetByID(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	var req struct {
		Name          string `json:"name"`
		Email         string `json:"email"`
		RemindEnabled *bool  `json:"remindEnabled"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.RemindEnabled != nil {
		user.RemindEnabled = *req.RemindEnabled
	}
	if err := h.deps.UserRepo.UpdateByID(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) GetProfile(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	user, err := h.deps.UserRepo.GetByID(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":            user.ID,
		"name":          user.Name,
		"email":         user.Email,
		"remindEnabled": user.RemindEnabled,
	})
}

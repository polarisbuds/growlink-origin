package api

import (
	"errors"
	"net/http"
	"path/filepath"

	"growlink/internal/adapter/handler/authctx"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func userImageURL(imageURI *string) *string {
	if imageURI == nil || *imageURI == "" {
		return nil
	}
	s := "/uploads/" + filepath.Base(*imageURI)
	return &s
}

func (h *Handler) AvatarMe(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	user, err := h.deps.UserRepo.GetByID(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"name":          claims.Name,
		"email":         claims.Email,
		"imageUri":      user.ImageURI,
		"isSystemAdmin": authctx.IsSystemAdmin(c),
		"isAnalyst":     authctx.IsAnalyst(c),
		"canPairing":    authctx.CanPairing(c),
	})
}

func (h *Handler) AvatarInfo(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	user, err := h.deps.UserRepo.GetByID(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.deps.UserRepo.AttachImageBase64(user)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"name":        claims.Name,
		"email":       claims.Email,
		"imageBase64": user.ImageBase64,
	})
}

func (h *Handler) AvatarByUUID(c *echo.Context) error {
	userUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	user, err := h.deps.UserRepo.GetByUUID(userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"uuid":        user.UUID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"imageBase64": user.ImageBase64,
		"imageUrl":    userImageURL(user.ImageURI),
	})
}

package api

import (
	"net/http"

	"growlink/internal/adapter/handler/authctx"

	"github.com/labstack/echo/v5"
)

func (h *Handler) AvatarMe(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	user, err := h.deps.UserRepo.GetByID(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	imageURI := user.ImageURI
	return c.JSON(http.StatusOK, map[string]interface{}{
		"name":          claims.Name,
		"email":         claims.Email,
		"imageUri":      imageURI,
		"isSystemAdmin": authctx.IsSystemAdmin(c),
		"isAnalyst":     authctx.IsAnalyst(c),
		"canPairing":    authctx.CanPairing(c),
	})
}

package authctx

import (
	"growlink/internal/shared/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

const CookieNameToken = "access_token"

func Claims(c *echo.Context) *auth.CustomClaims {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok || token == nil {
		return nil
	}
	custom, ok := token.Claims.(*auth.CustomClaims)
	if !ok {
		return nil
	}
	return custom
}

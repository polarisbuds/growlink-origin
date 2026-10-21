package authctx

import (
	"growlink/internal/domain/valueobject"

	"github.com/labstack/echo/v5"
)

func SystemRole(c *echo.Context) valueobject.SystemRole {
	claims := Claims(c)
	if claims == nil {
		return valueobject.SystemRoleNone
	}
	return claims.SystemRole
}

func IsSystemAdmin(c *echo.Context) bool {
	return SystemRole(c) == valueobject.SystemRoleAdmin
}

func IsAnalyst(c *echo.Context) bool {
	return SystemRole(c) == valueobject.SystemRoleAnalyst
}

func CanViewProgress(c *echo.Context) bool {
	r := SystemRole(c)
	return r == valueobject.SystemRoleAdmin || r == valueobject.SystemRoleAnalyst
}

func CanAccessPairingPlanning(c *echo.Context) bool {
	claims := Claims(c)
	if claims == nil {
		return false
	}
	if claims.SystemRole == valueobject.SystemRoleAdmin {
		return true
	}
	return int(claims.Authority) >= int(valueobject.AuthorityBoss)
}

func CanPairing(c *echo.Context) bool {
	return CanAccessPairingPlanning(c)
}

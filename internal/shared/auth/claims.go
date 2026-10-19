package auth

import (
	"growlink/internal/domain/valueobject"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID    uint                  `json:"user_id"`
	UUID      string                `json:"uuid"`
	Email     string                `json:"email"`
	Name      string                `json:"name"`
	Authority valueobject.Authority `json:"authority"`
	jwt.RegisteredClaims
}

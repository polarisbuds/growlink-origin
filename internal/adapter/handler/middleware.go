package handler

import (
	"log/slog"
	"net/http"

	"growlink/internal/adapter/handler/authctx"
	"growlink/internal/domain/valueobject"
	"growlink/internal/shared/auth"
	"growlink/internal/shared/config"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func JWTConfig(auditLogger *slog.Logger) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims { return new(auth.CustomClaims) },
		SigningKey:    []byte(config.GetJWTKey()),
		TokenLookup:   "cookie:" + authctx.CookieNameToken,
		ErrorHandler: func(c *echo.Context, err error) error {
			auditLogger.Error("jwt-page-error", slog.String("error", err.Error()))
			return c.Redirect(http.StatusFound, "/login")
		},
	})
}

func APIJWTConfig(auditLogger *slog.Logger) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims { return new(auth.CustomClaims) },
		SigningKey:    []byte(config.GetJWTKey()),
		TokenLookup:   "cookie:" + authctx.CookieNameToken,
		ErrorHandler: func(c *echo.Context, err error) error {
			auditLogger.Error("jwt-api-error", slog.String("error", err.Error()))
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		},
	})
}

func AuthorityAtLeast(min valueobject.Authority) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			claims := authctx.Claims(c)
			if claims == nil || int(claims.Authority) < int(min) {
				return c.Redirect(http.StatusFound, "/")
			}
			return next(c)
		}
	}
}

func AuthorityAtLeastAPI(min valueobject.Authority) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			claims := authctx.Claims(c)
			if claims == nil || int(claims.Authority) < int(min) {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
			}
			return next(c)
		}
	}
}

func auditLogging(auditLogger *slog.Logger, beforeAuth bool) echo.MiddlewareFunc {
	if auditLogger == nil {
		auditLogger = slog.Default()
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()
			attrs := []any{
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.String("remote_addr", req.RemoteAddr),
			}
			if !beforeAuth {
				attrs = append(attrs, slog.Any("claims", authctx.Claims(c)))
			}
			auditLogger.Info("request", attrs...)
			return next(c)
		}
	}
}

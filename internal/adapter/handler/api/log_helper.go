package api

import (
	"context"
	"log/slog"

	logutils "growlink/internal/shared/logUtils"

	"github.com/labstack/echo/v5"
)

func requestCtx(c *echo.Context) context.Context {
	return c.Request().Context()
}

func logHandlerWarn(ctx context.Context, op, msg string, httpStatus int, attrs ...any) {
	attrs = append(attrs, slog.Int("http_status", httpStatus))
	logutils.Warn(ctx, logutils.LayerHandler, op, msg, attrs...)
}

func logHandlerInfo(ctx context.Context, op, msg string, attrs ...any) {
	logutils.Info(ctx, logutils.LayerHandler, op, msg, attrs...)
}

func handlerAttrs(c *echo.Context) []any {
	req := c.Request()
	return []any{
		slog.String("method", req.Method),
		slog.String("path", req.URL.Path),
	}
}

package logutils

import (
	"context"
	"log/slog"
	"time"
)

const (
	LayerHandler    = "handler"
	LayerUsecase    = "usecase"
	LayerRepository = "repository"
)

func logger(ctx context.Context) *slog.Logger {
	rid := RequestIDFromContext(ctx)
	if rid == "" {
		return slog.Default()
	}
	return slog.Default().With(slog.String("request_id", rid))
}

func withBase(layer, op string, attrs []any) []any {
	out := []any{slog.String("layer", layer), slog.String("op", op)}
	return append(out, attrs...)
}

func Debug(ctx context.Context, layer, op, msg string, attrs ...any) {
	logger(ctx).Debug(msg, withBase(layer, op, attrs)...)
}

func Info(ctx context.Context, layer, op, msg string, attrs ...any) {
	logger(ctx).Info(msg, withBase(layer, op, attrs)...)
}

func Warn(ctx context.Context, layer, op, msg string, attrs ...any) {
	logger(ctx).Warn(msg, withBase(layer, op, attrs)...)
}

func Error(ctx context.Context, layer, op, msg string, attrs ...any) {
	logger(ctx).Error(msg, withBase(layer, op, attrs)...)
}

// DoneFunc は Begin が返すクロージャ。defer で呼び出し、返却 err を渡す。
type DoneFunc func(err error)

// Begin は処理開始を Info で記録し、呼び出し側が返す err に応じて処理終了または処理失敗を記録する。
func Begin(ctx context.Context, layer, op string, startAttrs ...any) DoneFunc {
	start := time.Now()
	Info(ctx, layer, op, "処理開始", startAttrs...)
	return func(err error) {
		elapsed := time.Since(start)
		if err != nil {
			Error(ctx, layer, op, "処理失敗", slog.String("err", err.Error()), slog.Duration("elapsed", elapsed))
			return
		}
		Info(ctx, layer, op, "処理終了", slog.Duration("elapsed", elapsed))
	}
}

package lineworks_usecase

import (
	"context"
	"errors"
	"fmt"
	logutils "growlink/internal/shared/logUtils"
	"strings"
)

// Sender は LINE WORKS のユーザー1人へテキストを送る。
type Sender interface {
	SendUserMessage(ctx context.Context, userID, text string) error
}

type Usecase struct {
	client Sender
}

func New(client Sender) *Usecase {
	return &Usecase{client: client}
}

// SendMessage は lineworksIDs の各ユーザーへ同じ本文を送る。
// 空の ID は飛ばす。1人でも失敗したら、残りの送信のあとでそのエラーを返す。
func (uc *Usecase) SendMessage(ctx context.Context, lineworksIDs []string, message string) (err error) {
	const op = "lineworks_usecase.SendMessage"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	message = strings.TrimSpace(message)
	if message == "" {
		return fmt.Errorf("message is empty")
	}
	var errs []error
	sent := 0
	for _, id := range lineworksIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if err = uc.client.SendUserMessage(ctx, id, message); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", id, err))
			continue
		}
		sent++
	}
	if sent == 0 && len(errs) == 0 {
		return fmt.Errorf("lineworksIDs is empty")
	}
	return errors.Join(errs...)
}

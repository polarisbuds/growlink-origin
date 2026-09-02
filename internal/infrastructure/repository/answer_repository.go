package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"
	"log/slog"

	"gorm.io/gorm"
)

type AnswerRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Answer]
	unScopedBaseRepository BaseUnScopedRepository[entity.Answer]
}

func NewAnswerRepository(ctx context.Context, dp *gorm.DB) *AnswerRepository {
	return &AnswerRepository{
		baseRepository: BaseRepository[entity.Answer]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Answer]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *AnswerRepository) GetByID(id uint) (*entity.Answer, error) {
	answer, err := repo.baseRepository.First("id ASC", "id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("ID指定での回答取得に失敗しました: %w", err)
	}
	return answer, nil
}

func (repo *AnswerRepository) GetByMeetingID(meetingID uint) ([]entity.Answer, error) {
	answers, err := repo.baseRepository.Where("id ASC", "meeting_id = ?", meetingID)
	if err != nil {
		return nil, fmt.Errorf("ID指定での回答取得に失敗しました: %w", err)
	}
	return answers, nil
}

func (repo *AnswerRepository) GetByUserID(userID uint) ([]entity.Answer, error) {
	answers, err := repo.baseRepository.Where("id ASC", "user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("ID指定での回答取得に失敗しました: %w", err)
	}
	return answers, nil
}

func (repo *AnswerRepository) Submit(answer *entity.Answer) (*entity.Answer, error) {
	if answer == nil {
		slog.Warn("投稿された回答内容が空データになっています。データの登録はスキップされました。")
		return nil, nil
	}

	if err := repo.baseRepository.Insert(answer); err != nil {
		return nil, fmt.Errorf("回答情報の投稿に失敗しました: %w", err)
	}

	return answer, nil
}

func (repo *AnswerRepository) Submits(answers []entity.Answer) ([]entity.Answer, error) {
	if len(answers) == 0 || answers == nil {
		slog.Warn("投稿された回答全ての内容が空データになっています。データの登録はスキップされました。")
		return nil, nil
	}

	if err := repo.baseRepository.InsertInBatch(answers, len(answers)); err != nil {
		return nil, fmt.Errorf("複数回答の保存に失敗しました: %w", err)
	}

	return answers, nil
}

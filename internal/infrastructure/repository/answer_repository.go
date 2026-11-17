package repository

import (
	"context"
	"errors"
	"fmt"
	"growlink/internal/domain/entity"
	logutils "growlink/internal/shared/logUtils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewAnswerRepository(ctx context.Context, dp *gorm.DB) *AnswerRepository {
	return &AnswerRepository{ctx: ctx, db: dp}
}

func (repo *AnswerRepository) GetByUUID(answerUUID uuid.UUID) (ret0 *entity.Answer, err error) {
	const op = "answer_repository.GetByUUID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	answer, err := gorm.G[entity.Answer](repo.db).Where("uuid = ?", answerUUID).Order("id ASC").First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("UUID指定での回答取得に失敗しました: %w", err)
	}
	return &answer, nil
}

func (repo *AnswerRepository) GetByID(id uint) (ret0 *entity.Answer, err error) {
	const op = "answer_repository.GetByID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	answer, err := gorm.G[entity.Answer](repo.db).Where("id = ?", id).Order("id ASC").First(repo.ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("ID指定での回答取得に失敗しました: %w", err)
	}
	return &answer, nil
}

func (repo *AnswerRepository) GetByMeetingID(meetingID uint) (ret0 []entity.Answer, err error) {
	const op = "answer_repository.GetByMeetingID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	answers, err := gorm.G[entity.Answer](repo.db).
		Where("meeting_id = ?", meetingID).
		Order("id ASC").
		Preload("Topic", nil).
		Preload("User", nil).
		Preload("AnswerComment", nil).
		Preload("AnswerComment.User", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ID指定での回答取得に失敗しました: %w", err)
	}
	return answers, nil
}

func (repo *AnswerRepository) GetByUserID(userID uint) (ret0 []entity.Answer, err error) {
	const op = "answer_repository.GetByUserID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	answers, err := gorm.G[entity.Answer](repo.db).
		Where("user_id = ?", userID).
		Order("id ASC").
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ID指定での回答取得に失敗しました: %w", err)
	}
	return answers, nil
}

func (repo *AnswerRepository) Submit(answer *entity.Answer) (ret0 *entity.Answer, err error) {
	const op = "answer_repository.Submit"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if answer == nil {
		logutils.Warn(repo.ctx, logutils.LayerRepository, op, "投稿された回答内容が空データになっています。データの登録はスキップされました。")
		return nil, nil
	}

	if err = gorm.G[entity.Answer](repo.db).Create(repo.ctx, answer); err != nil {
		return nil, fmt.Errorf("回答情報の投稿に失敗しました: %w", err)
	}

	return answer, nil
}

func (repo *AnswerRepository) Submits(answers []entity.Answer) (ret0 []entity.Answer, err error) {
	const op = "answer_repository.Submits"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if len(answers) == 0 || answers == nil {
		logutils.Warn(repo.ctx, logutils.LayerRepository, op, "投稿された回答全ての内容が空データになっています。データの登録はスキップされました。")
		return nil, nil
	}

	for i := range answers {
		a := &answers[i]
		if a.ID > 0 {
			if _, err = gorm.G[entity.Answer](repo.db).Where("id = ?", a.ID).Update(repo.ctx, "content", a.Content); err != nil {
				return nil, fmt.Errorf("回答の更新に失敗しました: %w", err)
			}
			continue
		}
		if a.MeetingID > 0 && a.TopicID > 0 && a.UserID > 0 {
			existing, err := gorm.G[entity.Answer](repo.db).
				Where("meeting_id = ? AND topic_id = ? AND user_id = ?", a.MeetingID, a.TopicID, a.UserID).
				First(repo.ctx)
			if err == nil {
				if _, err = gorm.G[entity.Answer](repo.db).Where("id = ?", existing.ID).Update(repo.ctx, "content", a.Content); err != nil {
					return nil, fmt.Errorf("回答の更新に失敗しました: %w", err)
				}
				answers[i].ID = existing.ID
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("既存回答の検索に失敗しました: %w", err)
			}
		}
		if err = gorm.G[entity.Answer](repo.db).Create(repo.ctx, a); err != nil {
			return nil, fmt.Errorf("回答の保存に失敗しました: %w", err)
		}
	}

	return answers, nil
}

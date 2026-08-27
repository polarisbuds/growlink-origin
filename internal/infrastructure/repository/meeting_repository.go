package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"
	"growlink/internal/shared/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MeetingRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.Meeting]
	unScopedBaseRepository BaseUnScopedRepository[entity.Meeting]
}

func NewMeetingRepository(ctx context.Context, dp *gorm.DB) *MeetingRepository {
	return &MeetingRepository{
		baseRepository: BaseRepository[entity.Meeting]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Meeting]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func preloadNone(_ gorm.PreloadBuilder) error {
	return nil
}

func preloadAll(chain gorm.ChainInterface[entity.Meeting]) gorm.ChainInterface[entity.Meeting] {
	return chain.
		Preload("Status", preloadNone).
		Preload("Template", preloadNone).
		Preload("Template.Topics", preloadNone).
		Preload("Peering", preloadNone).
		Preload("Peering.Mentor", preloadNone).
		Preload("Peering.Mentee", preloadNone).
		Preload("Answers", preloadNone).
		Preload("Answers.AnswerComments", preloadNone).
		Preload("MeetingComments", preloadNone).
		Preload("MeetingComments.User", preloadNone)
}

func preloadBasic(chain gorm.ChainInterface[entity.Meeting]) gorm.ChainInterface[entity.Meeting] {
	return chain.
		Preload("Status", preloadNone).
		Preload("Template", preloadNone).
		Preload("Peering", preloadNone).
		Preload("Peering.Mentor", preloadNone).
		Preload("Peering.Mentee", preloadNone)
}

func (repo *MeetingRepository) GetAll() ([]entity.Meeting, error) {
	meetings, err := preloadBasic(gorm.G[entity.Meeting](repo.baseRepository.dp).Order("id DESC")).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーティング情報の全取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) GetByID(id int64) (*entity.Meeting, error) {
	meeting, err := preloadAll(gorm.G[entity.Meeting](repo.baseRepository.dp).Order("id DESC")).
		Where("id = ?", id).
		First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーディング情報のID指定取得に失敗しました: %w", err)
	}
	return &meeting, nil
}

func (repo *MeetingRepository) GetByUUID(uuid uuid.UUID) (*entity.Meeting, error) {
	meeting, err := preloadAll(gorm.G[entity.Meeting](repo.baseRepository.dp).Order("id DESC")).
		Where("uuid = ?", uuid.String()).
		First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーディング情報のID指定取得に失敗しました: %w", err)
	}
	return &meeting, nil
}

func (repo *MeetingRepository) GetUntilToday(base time.Time) ([]entity.Meeting, error) {
	endOfBaseDay := utils.GetEndOfDay(base)

	meetings, err := preloadBasic(gorm.G[entity.Meeting](repo.baseRepository.dp).Order("id DESC")).
		Where("implemented_at <= ? ", endOfBaseDay).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーティング情報の日付指定取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) GetFromSpecificDateOnward(base time.Time) ([]entity.Meeting, error) {
	startOfBaseDay := utils.GetStartOfDay(base)

	meetings, err := preloadBasic(gorm.G[entity.Meeting](repo.baseRepository.dp).Order("id DESC")).
		Where("implemented_at >= ? ", startOfBaseDay).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーティング情報の日付指定取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) GetBetweenDates(from, to time.Time) ([]entity.Meeting, error) {
	f := from
	t := to
	if from.After(to) {
		f = to
		t = from
	}

	baseFrom := utils.GetStartOfDay(f)
	baseTo := utils.GetEndOfDay(t)

	meetings, err := preloadBasic(gorm.G[entity.Meeting](repo.baseRepository.dp).Order("id DESC")).
		Where("implemented_at BETWEEN ? AND ? ", baseFrom, baseTo).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーティング情報の日付指定取得に失敗しました: %w", err)
	}
	return meetings, nil
}

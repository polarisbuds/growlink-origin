package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
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
		Preload("Template", preloadNone).
		Preload("Template.Topics", preloadNone).
		Preload("Pairing", preloadNone).
		Preload("Pairing.Mentor", preloadNone).
		Preload("Pairing.Mentee", preloadNone).
		Preload("Answers", preloadNone).
		Preload("Answers.AnswerComments", preloadNone).
		Preload("MeetingComments", preloadNone).
		Preload("MeetingComments.User", preloadNone)
}

func preloadBasic(chain gorm.ChainInterface[entity.Meeting]) gorm.ChainInterface[entity.Meeting] {
	return chain.
		Preload("Template", preloadNone).
		Preload("Pairing", preloadNone).
		Preload("Pairing.Mentor", preloadNone).
		Preload("Pairing.Mentee", preloadNone)
}

func (repo *MeetingRepository) GetAll() ([]entity.Meeting, error) {
	meetings, err := preloadBasic(gorm.G[entity.Meeting](repo.baseRepository.dp).Order("id DESC")).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーティング情報の全取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) GetByID(id uint) (*entity.Meeting, error) {
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

func (repo *MeetingRepository) GetByPairingID(pairingID uint) ([]entity.Meeting, error) {
	pairingWhereChain, err := repo.baseRepository.WhereChain("id ASC", "pairing_id = ?", pairingID)
	if err != nil {
		return nil, fmt.Errorf("ペアリングID指定でのペアリングのチェイン情報取得に失敗しました: %w", err)
	}
	if pairingWhereChain == nil {
		return []entity.Meeting{}, nil
	}
	pairings, err := pairingWhereChain.Preload("Pairing", nil).
		Preload("Pairing.Mentor", nil).
		Preload("Pairing.Mentee", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ペアリングID指定でのペアリング情報取得に失敗しました: %w", err)
	}
	return pairings, nil
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

func (repo *MeetingRepository) Create(meeting *entity.Meeting) (*entity.Meeting, error) {
	if err := repo.baseRepository.Insert(meeting); err != nil {
		return nil, fmt.Errorf("ミーティングの新規追加に失敗しました: %w", err)
	}
	return meeting, nil
}

func (repo *MeetingRepository) Update(meetingID uint, meeting *entity.Meeting) error {
	if _, err := repo.baseRepository.Updates(meeting, "id = ?", meetingID); err != nil {
		return fmt.Errorf("ミーティングデータの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) UpdateStatus(meetingID uint, status valueobject.MeetingStatus) error {
	if _, err := repo.baseRepository.Update("status", status, "id = ?", meetingID); err != nil {
		return fmt.Errorf("ミーティングステータスの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) UpdateTemplate(meetingID, templateID uint) error {
	if _, err := repo.baseRepository.Update("template_id", templateID, "id = ?", meetingID); err != nil {
		return fmt.Errorf("使用するテンプレートの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) UpdatePairing(meetingID, pairingID uint) error {
	if _, err := repo.baseRepository.Update("pairing_id", pairingID, "id = ?", meetingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) UpdateImplementDate(meetingID uint, implementDate time.Time) error {
	if _, err := repo.baseRepository.Update("implemented_at", implementDate, "id = ?", meetingID); err != nil {
		return fmt.Errorf("実施予定日の更新に失敗しました: %w", err)
	}
	return nil
}

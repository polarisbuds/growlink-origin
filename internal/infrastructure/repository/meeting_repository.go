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
	"gorm.io/gorm/clause"
)

type MeetingRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewMeetingRepository(ctx context.Context, dp *gorm.DB) *MeetingRepository {
	return &MeetingRepository{ctx: ctx, db: dp}
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
		Preload("Pairing.Mentee", preloadNone).
		Preload("Pairing.Mentee.Department", preloadNone).
		Preload("Pairing.Mentee.TeamMembers", preloadNone).
		Preload("Pairing.Mentee.TeamMembers.Team", preloadNone).
		Preload("Pairing.Mentee.TeamMembers.Team.Department", preloadNone)
}

func (repo *MeetingRepository) GetAll() ([]entity.Meeting, error) {
	meetings, err := preloadBasic(gorm.G[entity.Meeting](repo.db).Order("id DESC")).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーティング情報の全取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) GetByID(id uint) (*entity.Meeting, error) {
	meeting, err := preloadAll(gorm.G[entity.Meeting](repo.db).Order("id DESC")).
		Where("id = ?", id).
		First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーディング情報のID指定取得に失敗しました: %w", err)
	}
	return &meeting, nil
}

func (repo *MeetingRepository) GetByUUID(uuid uuid.UUID) (*entity.Meeting, error) {
	meeting, err := preloadAll(gorm.G[entity.Meeting](repo.db).Order("id DESC")).
		Where("uuid = ?", uuid.String()).
		First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーディング情報のID指定取得に失敗しました: %w", err)
	}
	return &meeting, nil
}

func (repo *MeetingRepository) GetUntilToday(base time.Time) ([]entity.Meeting, error) {
	endOfBaseDay := utils.GetEndOfDay(base)

	meetings, err := preloadBasic(gorm.G[entity.Meeting](repo.db).Order("id DESC")).
		Where("implemented_at <= ? ", endOfBaseDay).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーティング情報の日付指定取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) GetByPairingID(pairingID uint) ([]entity.Meeting, error) {
	meetings, err := gorm.G[entity.Meeting](repo.db).
		Where("pairing_id = ?", pairingID).
		Order("id ASC").
		Preload("Pairing", nil).
		Preload("Pairing.Mentor", nil).
		Preload("Pairing.Mentee", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ペアリングID指定でのペアリング情報取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) GetFromSpecificDateOnward(base time.Time) ([]entity.Meeting, error) {
	startOfBaseDay := utils.GetStartOfDay(base)

	meetings, err := preloadBasic(gorm.G[entity.Meeting](repo.db).Order("id DESC")).
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

	meetings, err := preloadBasic(gorm.G[entity.Meeting](repo.db).Order("id DESC")).
		Where("implemented_at BETWEEN ? AND ? ", baseFrom, baseTo).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーティング情報の日付指定取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) Create(meeting *entity.Meeting) (*entity.Meeting, error) {
	if err := gorm.G[entity.Meeting](repo.db).Create(repo.ctx, meeting); err != nil {
		return nil, fmt.Errorf("ミーティングの新規追加に失敗しました: %w", err)
	}
	return meeting, nil
}

func (repo *MeetingRepository) Update(meetingID uint, meeting *entity.Meeting) error {
	if _, err := gorm.G[entity.Meeting](repo.db).Where("id = ?", meetingID).Omit(clause.Associations).Updates(repo.ctx, *meeting); err != nil {
		return fmt.Errorf("ミーティングデータの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) UpdateStatus(meetingID uint, status valueobject.MeetingStatus) error {
	if _, err := gorm.G[entity.Meeting](repo.db).Where("id = ?", meetingID).Update(repo.ctx, "status", status); err != nil {
		return fmt.Errorf("ミーティングステータスの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) UpdateTemplate(meetingID, templateID uint) error {
	if _, err := gorm.G[entity.Meeting](repo.db).Where("id = ?", meetingID).Update(repo.ctx, "template_id", templateID); err != nil {
		return fmt.Errorf("使用するテンプレートの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) UpdatePairing(meetingID, pairingID uint) error {
	if _, err := gorm.G[entity.Meeting](repo.db).Where("id = ?", meetingID).Update(repo.ctx, "pairing_id", pairingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) UpdateImplementDate(meetingID uint, implementDate time.Time) error {
	if _, err := gorm.G[entity.Meeting](repo.db).Where("id = ?", meetingID).Update(repo.ctx, "implemented_at", implementDate); err != nil {
		return fmt.Errorf("実施予定日の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) GetMeetingIDsBySharerUserID(userID uint) ([]uint, error) {
	var sharers []entity.MeetingSharer
	if err := repo.db.WithContext(repo.ctx).Where("user_id = ?", userID).Find(&sharers).Error; err != nil {
		return nil, fmt.Errorf("ミーティング共有者の取得に失敗しました: %w", err)
	}
	ids := make([]uint, 0, len(sharers))
	for _, s := range sharers {
		ids = append(ids, s.MeetingID)
	}
	return ids, nil
}

func (repo *MeetingRepository) ReplaceSharers(meetingID uint, userIDs []uint) error {
	if err := repo.db.WithContext(repo.ctx).Where("meeting_id = ?", meetingID).Delete(&entity.MeetingSharer{}).Error; err != nil {
		return fmt.Errorf("ミーティング共有者の更新に失敗しました: %w", err)
	}
	for _, uid := range userIDs {
		if uid == 0 {
			continue
		}
		if err := repo.db.WithContext(repo.ctx).Create(&entity.MeetingSharer{MeetingID: meetingID, UserID: uid}).Error; err != nil {
			return fmt.Errorf("ミーティング共有者の登録に失敗しました: %w", err)
		}
	}
	return nil
}

func (repo *MeetingRepository) Delete(meetingIDs ...uint) error {
	if len(meetingIDs) == 0 {
		return nil
	}

	if _, err := gorm.G[entity.Meeting](repo.db).Where("id IN ?", meetingIDs).Delete(repo.ctx); err != nil {
		return fmt.Errorf("ミーティングの削除に失敗しました: %w", err)
	}
	return nil
}

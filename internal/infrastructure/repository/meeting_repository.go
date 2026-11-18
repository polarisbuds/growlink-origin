package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"
	"growlink/internal/shared/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MeetingRepository struct {
	ctx      context.Context
	db       *gorm.DB
	unscoped *gorm.DB
}

func NewMeetingRepository(ctx context.Context, dp *gorm.DB) *MeetingRepository {
	return &MeetingRepository{ctx: ctx, db: dp, unscoped: dp.Session(&gorm.Session{
		PropagateUnscoped: true,
	})}
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
		Preload("Answers.AnswerComment", preloadNone).
		Preload("Answers.AnswerComment.User", preloadNone).
		Preload("MeetingComments", preloadNone).
		Preload("MeetingComments.User", preloadNone).
		Preload("Sharers", preloadNone).
		Preload("Sharers.User", preloadNone)
}

func preloadBasic(chain gorm.ChainInterface[entity.Meeting]) gorm.ChainInterface[entity.Meeting] {
	return chain.
		Preload("Template", preloadNone).
		Preload("Pairing", preloadNone).
		Preload("Pairing.Mentor", preloadNone).
		Preload("Pairing.Mentee", preloadNone).
		Preload("Pairing.Mentee.DepartmentMembers", preloadNone).
		Preload("Pairing.Mentee.DepartmentMembers.Department", preloadNone).
		Preload("Pairing.Mentee.TeamMembers", preloadNone).
		Preload("Pairing.Mentee.TeamMembers.Team", preloadNone).
		Preload("Pairing.Mentee.TeamMembers.Team.Department", preloadNone)
}

func (repo *MeetingRepository) GetAll() (ret0 []entity.Meeting, err error) {
	const op = "meeting_repository.GetAll"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	meetings, err := preloadBasic(gorm.G[entity.Meeting](repo.db).Scopes(func(stmt *gorm.Statement) {
		stmt.Unscoped = true
	}).Order("id DESC")).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーティング情報の全取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) GetByID(id uint) (ret0 *entity.Meeting, err error) {
	const op = "meeting_repository.GetByID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	meeting, err := preloadAll(gorm.G[entity.Meeting](repo.unscoped).Scopes(func(stmt *gorm.Statement) { stmt.Unscoped = true }).Order("id DESC")).
		Where("id = ?", id).
		First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーディング情報のID指定取得に失敗しました: %w", err)
	}
	return &meeting, nil
}

func (repo *MeetingRepository) GetByIDs(ids []uint) (ret0 []entity.Meeting, err error) {
	const op = "meeting_repository.GetByID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	meetings, err := preloadAll(gorm.G[entity.Meeting](repo.unscoped).Scopes(func(stmt *gorm.Statement) { stmt.Unscoped = true }).Order("id DESC")).
		Where("id IN ?", ids).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーディング情報のID指定取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) GetByUUID(uuid uuid.UUID) (ret0 *entity.Meeting, err error) {
	const op = "meeting_repository.GetByUUID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	meeting, err := preloadAll(gorm.G[entity.Meeting](repo.unscoped).Scopes(func(stmt *gorm.Statement) {
		stmt.Unscoped = true
	}).Order("id DESC")).
		Where("uuid = ?", uuid.String()).
		First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーディング情報のID指定取得に失敗しました: %w", err)
	}
	return &meeting, nil
}

func (repo *MeetingRepository) GetUntilToday(base time.Time) (ret0 []entity.Meeting, err error) {
	const op = "meeting_repository.GetUntilToday"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	endOfBaseDay := utils.GetEndOfDay(base)

	meetings, err := preloadBasic(gorm.G[entity.Meeting](repo.db).Order("id DESC")).
		Where("implemented_at <= ? ", endOfBaseDay).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーティング情報の日付指定取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) GetByPairingID(pairingID uint) (ret0 []entity.Meeting, err error) {
	const op = "meeting_repository.GetByPairingID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	meetings, err := gorm.G[entity.Meeting](repo.db).
		Scopes(func(stmt *gorm.Statement) {
			stmt.Unscoped = true
		}).
		Where("pairing_id = ?", pairingID).
		Order("id ASC").
		Preload("Pairing", nil).
		Preload("Pairing.Mentor", nil).
		Preload("Pairing.Mentee", nil).
		Preload("Template", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ペアリングID指定でのペアリング情報取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) GetFromSpecificDateOnward(base time.Time) (ret0 []entity.Meeting, err error) {
	const op = "meeting_repository.GetFromSpecificDateOnward"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	startOfBaseDay := utils.GetStartOfDay(base)

	meetings, err := preloadBasic(gorm.G[entity.Meeting](repo.db).Order("id DESC")).
		Where("implemented_at >= ? ", startOfBaseDay).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーティング情報の日付指定取得に失敗しました: %w", err)
	}
	return meetings, nil
}

func (repo *MeetingRepository) GetBetweenDates(from, to time.Time) (ret0 []entity.Meeting, err error) {
	const op = "meeting_repository.GetBetweenDates"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
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

func (repo *MeetingRepository) Create(meeting *entity.Meeting) (ret0 *entity.Meeting, err error) {
	const op = "meeting_repository.Create"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if err = gorm.G[entity.Meeting](repo.db).Create(repo.ctx, meeting); err != nil {
		return nil, fmt.Errorf("ミーティングの新規追加に失敗しました: %w", err)
	}
	return meeting, nil
}

func (repo *MeetingRepository) Update(meetingID uint, meeting *entity.Meeting) (err error) {
	const op = "meeting_repository.Update"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if _, err = gorm.G[entity.Meeting](repo.db).Where("id = ?", meetingID).Omit(clause.Associations).Updates(repo.ctx, *meeting); err != nil {
		return fmt.Errorf("ミーティングデータの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) UpdateStatus(meetingID uint, status valueobject.MeetingStatus) (err error) {
	const op = "meeting_repository.UpdateStatus"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if _, err = gorm.G[entity.Meeting](repo.db).Where("id = ?", meetingID).Update(repo.ctx, "status", status); err != nil {
		return fmt.Errorf("ミーティングステータスの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) UpdateTemplate(meetingID, templateID uint) (err error) {
	const op = "meeting_repository.UpdateTemplate"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if _, err = gorm.G[entity.Meeting](repo.db).Where("id = ?", meetingID).Update(repo.ctx, "template_id", templateID); err != nil {
		return fmt.Errorf("使用するテンプレートの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) UpdatePairing(meetingID, pairingID uint) (err error) {
	const op = "meeting_repository.UpdatePairing"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if _, err = gorm.G[entity.Meeting](repo.db).Where("id = ?", meetingID).Update(repo.ctx, "pairing_id", pairingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) UpdateImplementDate(meetingID uint, implementDate time.Time) (err error) {
	const op = "meeting_repository.UpdateImplementDate"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if _, err = gorm.G[entity.Meeting](repo.db).Where("id = ?", meetingID).Update(repo.ctx, "implemented_at", implementDate); err != nil {
		return fmt.Errorf("実施予定日の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) GetMeetingIDsBySharerUserID(userID uint) (ret0 []uint, err error) {
	const op = "meeting_repository.GetMeetingIDsBySharerUserID"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	sharers, err := gorm.G[entity.MeetingSharer](repo.db).
		Where("user_id = ?", userID).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ミーティング共有者の取得に失敗しました: %w", err)
	}
	ids := make([]uint, 0, len(sharers))
	for _, s := range sharers {
		ids = append(ids, s.MeetingID)
	}
	return ids, nil
}

func (repo *MeetingRepository) ReplaceSharers(meetingID uint, userIDs []uint) (err error) {
	const op = "meeting_repository.ReplaceSharers"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	want := map[uint]struct{}{}
	for _, uid := range userIDs {
		if uid == 0 {
			continue
		}
		want[uid] = struct{}{}
	}

	err = repo.db.WithContext(repo.ctx).Transaction(func(tx *gorm.DB) error {
		existing, err := gorm.G[entity.MeetingSharer](tx).
			Where("meeting_id = ?", meetingID).
			Find(repo.ctx)
		if err != nil {
			return err
		}
		have := map[uint]struct{}{}
		for _, row := range existing {
			have[row.UserID] = struct{}{}
			if _, ok := want[row.UserID]; !ok {
				if _, err = gorm.G[entity.MeetingSharer](tx).Where("id = ?", row.ID).Delete(repo.ctx); err != nil {
					return err
				}
			}
		}
		for uid := range want {
			if _, ok := have[uid]; ok {
				continue
			}
			sharer := entity.MeetingSharer{MeetingID: meetingID, UserID: uid}
			if err = gorm.G[entity.MeetingSharer](tx).Create(repo.ctx, &sharer); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("ミーティング共有者の更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *MeetingRepository) Delete(meetingIDs ...uint) (err error) {
	const op = "meeting_repository.Delete"
	done := logutils.Begin(repo.ctx, logutils.LayerRepository, op)
	defer func() {
		done(err)
	}()
	if len(meetingIDs) == 0 {
		return nil
	}

	if _, err = gorm.G[entity.Meeting](repo.db).Where("id IN ?", meetingIDs).Delete(repo.ctx); err != nil {
		return fmt.Errorf("ミーティングの削除に失敗しました: %w", err)
	}
	return nil
}

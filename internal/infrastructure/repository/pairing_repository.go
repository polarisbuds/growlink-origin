package repository

import (
	"context"
	"fmt"

	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PairingRepository struct {
	ctx          context.Context
	db, unscoped *gorm.DB
}

func NewPairingRepository(ctx context.Context, dp *gorm.DB) *PairingRepository {
	return &PairingRepository{ctx: ctx, db: dp, unscoped: dp.Session(&gorm.Session{PropagateUnscoped: true})}
}

func (repo *PairingRepository) CreatePairing(menteeID, mentorID, createdByUserID uint) (ret0 *entity.Pairing, err error) {
	newPair := &entity.Pairing{
		MentorID:        mentorID,
		MenteeID:        menteeID,
		CreatedByUserID: createdByUserID,
	}

	if err = gorm.G[entity.Pairing](repo.db).Create(repo.ctx, newPair); err != nil {
		return nil, fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
	}

	return newPair, nil
}

func (repo *PairingRepository) CreateTemporaryPairing(menteeID, mentorID, createdByUserID uint) (ret0 *entity.Pairing, err error) {
	newPair := &entity.Pairing{
		MentorID:        mentorID,
		MenteeID:        menteeID,
		CreatedByUserID: createdByUserID,
		IsTemporary:     true,
	}
	if err = gorm.G[entity.Pairing](repo.db).Create(repo.ctx, newPair); err != nil {
		return nil, fmt.Errorf("一時ペアリングの作成に失敗しました: %w", err)
	}
	return newPair, nil
}

func (repo *PairingRepository) BreakPairing(id uint) (err error) {
	if _, err = gorm.G[entity.Pairing](repo.db).Where("id = ?", id).Delete(repo.ctx); err != nil {
		return fmt.Errorf("ペアリングの解消に失敗しました: %w", err)
	}
	return nil
}

func (repo *PairingRepository) BreakPairings(pairingID ...uint) (err error) {
	if len(pairingID) == 0 {
		return nil
	}
	if _, err = gorm.G[entity.Pairing](repo.db).Where("mentor_id IN ?", pairingID).Delete(repo.ctx); err != nil {
		return fmt.Errorf("ペアリングの解消に失敗しました: %w", err)
	}
	return nil
}

func (repo *PairingRepository) BreakAllMenteePairing(mentorID uint) (err error) {
	if _, err = gorm.G[entity.Pairing](repo.db).Where("mentor_id = ?", mentorID).Delete(repo.ctx); err != nil {
		return fmt.Errorf("ペアリングの解消に失敗しました: %w", err)
	}
	return nil
}

func (repo *PairingRepository) Get(pairingID uint) (ret0 *entity.Pairing, err error) {
	pairing, err := gorm.G[entity.Pairing](repo.db).Where("id = ?", pairingID).Order("id ASC").
		Preload("Mentor", nil).
		Preload("Mentee", nil).
		First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ペアリングID指定での取得に失敗しました: %w", err)
	}
	return &pairing, nil
}

func (repo *PairingRepository) GetByIDs(pairingIDs []uint) (ret0 []entity.Pairing, err error) {
	pairings, err := gorm.G[entity.Pairing](repo.db).Where("id IN ?", pairingIDs).Order("id ASC").
		Preload("Mentor", nil).
		Preload("Mentee", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ペアリングID指定での取得に失敗しました: %w", err)
	}
	return pairings, nil
}

func (repo *PairingRepository) GetByUUID(pairingUUID uuid.UUID) (ret0 *entity.Pairing, err error) {
	pairing, err := gorm.G[entity.Pairing](repo.db).Where("uuid = ?", pairingUUID).Order("id ASC").
		Preload("Mentor", nil).
		Preload("Mentee", nil).
		First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ペアリングUUID指定での取得に失敗しました: %w", err)
	}
	return &pairing, nil
}

func (repo *PairingRepository) GetByUUIDs(pairingUUID []string) (ret0 []entity.Pairing, err error) {
	pairings, err := gorm.G[entity.Pairing](repo.db).Where("uuid IN ?", pairingUUID).Order("id ASC").
		Preload("Mentor", nil).
		Preload("Mentee", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ペアリングUUID指定での取得に失敗しました: %w", err)
	}
	return pairings, nil
}

func (repo *PairingRepository) GetAll() (ret0 []entity.Pairing, err error) {
	pairing, err := gorm.G[entity.Pairing](repo.db).Order("id ASC").
		Preload("Mentor", nil).
		Preload("Mentor.DepartmentMembers", nil).
		Preload("Mentor.TeamMembers", nil).
		Preload("Mentee", nil).
		Preload("Mentee.DepartmentMembers", nil).
		Preload("Mentee.TeamMembers", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("全てのペアリング情報の取得に失敗しました: %w", err)
	}
	return pairing, nil
}

func (repo *PairingRepository) GetAllContainDeleted() (ret0 []entity.Pairing, err error) {
	pairing, err := gorm.G[entity.Pairing](repo.unscoped).Scopes(func(stmt *gorm.Statement) { stmt.Unscoped = true }).Order("id ASC").
		Preload("Mentor", nil).
		Preload("Mentor.DepartmentMembers", nil).
		Preload("Mentor.TeamMembers", nil).
		Preload("Mentee", nil).
		Preload("Mentee.DepartmentMembers", nil).
		Preload("Mentee.TeamMembers", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("全てのペアリング情報の取得に失敗しました: %w", err)
	}
	return pairing, nil
}

func (repo *PairingRepository) GetByMenteeID(menteeID uint) (ret0 []entity.Pairing, err error) {
	pairings, err := gorm.G[entity.Pairing](repo.db).
		Where("mentee_id = ? AND is_temporary = ?", menteeID, false).
		Order("id ASC").
		Preload("Mentee", nil).
		Preload("Mentor", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("メンティーID指定でのペアリング取得に失敗しました: %w", err)
	}
	return pairings, nil
}

func (repo *PairingRepository) GetByMenteeIDs(menteeIDs []uint) (ret0 []entity.Pairing, err error) {
	pairings, err := gorm.G[entity.Pairing](repo.db).
		Where("mentee_id IN ? AND is_temporary = ?", menteeIDs, false).
		Order("id ASC").
		Preload("Mentee", nil).
		Preload("Mentor", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("メンティーID指定でのペアリング取得に失敗しました: %w", err)
	}
	return pairings, nil
}

func (repo *PairingRepository) GetByMentorID(mentorID uint) (ret0 []entity.Pairing, err error) {
	pairings, err := gorm.G[entity.Pairing](repo.db).
		Where("mentor_id = ?", mentorID).
		Not("is_temporary = ?", true).
		Order("id ASC").
		Preload("Mentee", nil).
		Preload("Mentor", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("メンターID指定でのペアリング取得に失敗しました: %w", err)
	}
	return pairings, nil
}

func (repo *PairingRepository) GetPairingIDsBySharerUserID(userID uint) (ret0 []uint, err error) {
	var sharers []entity.PairingSharer
	if err = repo.db.WithContext(repo.ctx).Where("user_id = ?", userID).Find(&sharers).Error; err != nil {
		return nil, fmt.Errorf("共有者ペアの取得に失敗しました: %w", err)
	}
	ids := make([]uint, len(sharers))
	for idx, s := range sharers {
		ids[idx] = s.PairingID
	}
	return ids, nil
}

func (repo *PairingRepository) UpdatePairing(pairingID, mentorID, menteeID uint) (err error) {
	if err = repo.db.WithContext(repo.ctx).Model(&entity.Pairing{}).Where("id = ?", pairingID).
		Updates(map[string]interface{}{
			"mentor_id": mentorID,
			"mentee_id": menteeID,
		}).Error; err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	return nil
}

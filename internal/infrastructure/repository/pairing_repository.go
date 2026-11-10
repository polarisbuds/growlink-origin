package repository

import (
	"context"
	"errors"
	"fmt"

	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PairingRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewPairingRepository(ctx context.Context, dp *gorm.DB) *PairingRepository {
	return &PairingRepository{ctx: ctx, db: dp}
}

func (repo *PairingRepository) CreatePairing(menteeID, mentorID, createdByUserID uint) (*entity.Pairing, error) {
	newPair := &entity.Pairing{
		MentorID:        mentorID,
		MenteeID:        menteeID,
		CreatedByUserID: createdByUserID,
	}

	if err := gorm.G[entity.Pairing](repo.db).Create(repo.ctx, newPair); err != nil {
		return nil, fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
	}

	return newPair, nil
}

func (repo *PairingRepository) CreateTemporaryPairing(menteeID, mentorID, createdByUserID uint) (*entity.Pairing, error) {
	newPair := &entity.Pairing{
		MentorID:        mentorID,
		MenteeID:        menteeID,
		CreatedByUserID: createdByUserID,
		IsTemporary:     true,
	}
	if err := gorm.G[entity.Pairing](repo.db).Create(repo.ctx, newPair); err != nil {
		return nil, fmt.Errorf("一時ペアリングの作成に失敗しました: %w", err)
	}
	return newPair, nil
}

func (repo *PairingRepository) BreakPairing(id uint) error {
	if _, err := gorm.G[entity.Pairing](repo.db).Where("id = ?", id).Delete(repo.ctx); err != nil {
		return fmt.Errorf("ペアリングの解消に失敗しました: %w", err)
	}
	return nil
}

func (repo *PairingRepository) BreakPairings(pairingID ...uint) error {
	if len(pairingID) == 0 {
		return nil
	}
	if _, err := gorm.G[entity.Pairing](repo.db).Where("mentor_id IN ?", pairingID).Delete(repo.ctx); err != nil {
		return fmt.Errorf("ペアリングの解消に失敗しました: %w", err)
	}
	return nil
}

func (repo *PairingRepository) BreakAllMenteePairing(mentorID uint) error {
	if _, err := gorm.G[entity.Pairing](repo.db).Where("mentor_id = ?", mentorID).Delete(repo.ctx); err != nil {
		return fmt.Errorf("ペアリングの解消に失敗しました: %w", err)
	}
	return nil
}

func (repo *PairingRepository) Get(pairingID uint) (*entity.Pairing, error) {
	pairing, err := gorm.G[entity.Pairing](repo.db).Where("id = ?", pairingID).Order("id ASC").
		Preload("Mentor", nil).
		Preload("Mentee", nil).
		First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ペアリングID指定での取得に失敗しました: %w", err)
	}
	return &pairing, nil
}

func (repo *PairingRepository) GetByUUID(pairingUUID uuid.UUID) (*entity.Pairing, error) {
	pairing, err := gorm.G[entity.Pairing](repo.db).Where("uuid = ?", pairingUUID).Order("id ASC").
		Preload("Mentor", nil).
		Preload("Mentee", nil).
		First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ペアリングUUID指定での取得に失敗しました: %w", err)
	}
	return &pairing, nil
}

func (repo *PairingRepository) GetAll() ([]entity.Pairing, error) {
	pairing, err := gorm.G[entity.Pairing](repo.db).Order("id ASC").
		Preload("Mentor", nil).
		Preload("Mentee", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("全てのペアリング情報の取得に失敗しました: %w", err)
	}
	return pairing, nil
}

func (repo *PairingRepository) GetByMenteeID(menteeID uint) (*entity.Pairing, error) {
	pairing, err := gorm.G[entity.Pairing](repo.db).
		Where("mentee_id = ? AND is_temporary = ?", menteeID, false).
		Order("id ASC").
		First(repo.ctx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("メンティーID指定でのペアリング取得に失敗しました: %w", err)
	}
	return &pairing, nil
}

func (repo *PairingRepository) GetByMentorID(mentorID uint) ([]entity.Pairing, error) {
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

func (repo *PairingRepository) GetPairingIDsBySharerUserID(userID uint) ([]uint, error) {
	var sharers []entity.PairingSharer
	if err := repo.db.WithContext(repo.ctx).Where("user_id = ?", userID).Find(&sharers).Error; err != nil {
		return nil, fmt.Errorf("共有者ペアの取得に失敗しました: %w", err)
	}
	ids := make([]uint, 0, len(sharers))
	for _, s := range sharers {
		ids = append(ids, s.PairingID)
	}
	return ids, nil
}

func (repo *PairingRepository) UpdatePairing(pairingID, mentorID, menteeID uint) error {
	if err := repo.db.WithContext(repo.ctx).Model(&entity.Pairing{}).Where("id = ?", pairingID).
		Updates(map[string]interface{}{
			"mentor_id": mentorID,
			"mentee_id": menteeID,
		}).Error; err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	return nil
}

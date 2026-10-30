package repository

import (
	"context"
	"errors"
	"fmt"

	repository_interface "growlink/internal/domain/interface/repository"
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
	if err := repo.db.WithContext(repo.ctx).Create(&entity.PairingTarget{
		PairingID: newPair.ID,
		Kind:      entity.PairingTargetKindUser,
		TargetID:  menteeID,
	}).Error; err != nil {
		return nil, fmt.Errorf("一時ペアリング対象の登録に失敗しました: %w", err)
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
	pairing, err := gorm.G[entity.Pairing](repo.db).Where("id = ?", pairingID).Order("id ASC").First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ペアリングID指定での取得に失敗しました: %w", err)
	}
	return &pairing, nil
}

func (repo *PairingRepository) GetByUUID(pairingUUID uuid.UUID) (*entity.Pairing, error) {
	pairing, err := gorm.G[entity.Pairing](repo.db).Where("uuid = ?", pairingUUID).Order("id ASC").First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ペアリングUUID指定での取得に失敗しました: %w", err)
	}
	return &pairing, nil
}

func (repo *PairingRepository) GetAll() ([]entity.Pairing, error) {
	pairing, err := gorm.G[entity.Pairing](repo.db).Order("id ASC").
		Preload("Mentor", nil).
		Preload("Mentee", nil).
		Preload("Targets", nil).
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
	if err == nil {
		return &pairing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("メンティーID指定でのペアリング取得に失敗しました: %w", err)
	}
	var targetLink entity.PairingTarget
	if err := repo.db.WithContext(repo.ctx).Where("kind = ? AND target_id = ?", entity.PairingTargetKindUser, menteeID).First(&targetLink).Error; err == nil {
		p, err := repo.Get(targetLink.PairingID)
		if err != nil {
			return nil, err
		}
		if !p.IsTemporary {
			return p, nil
		}
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("ペアリング対象の取得に失敗しました: %w", err)
	}
	var link entity.PairingMentee
	if err := repo.db.WithContext(repo.ctx).Where("user_id = ?", menteeID).First(&link).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("ペアリングメンティー関連の取得に失敗しました: %w", err)
	}
	p, err := repo.Get(link.PairingID)
	if err != nil {
		return nil, err
	}
	if p.IsTemporary {
		return nil, nil
	}
	return p, nil
}

func (repo *PairingRepository) GetByMentorID(mentorID uint) ([]entity.Pairing, error) {
	pairings, err := gorm.G[entity.Pairing](repo.db).
		Where("mentor_id = ?", mentorID).
		Order("id ASC").
		Preload("Mentee", nil).
		Preload("Mentor", nil).
		Preload("Targets", nil).
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

func (repo *PairingRepository) UpdatePairing(pairingID, mentorID uint) error {
	if _, err := gorm.G[entity.Pairing](repo.db).Where("id = ?", pairingID).Update(repo.ctx, "mentor_id", mentorID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	return nil
}

func primaryMenteeFromTargets(targets []repository_interface.PairingTargetOption, mentorID uint) uint {
	for _, t := range targets {
		if t.Kind == entity.PairingTargetKindUser {
			return t.TargetID
		}
	}
	return mentorID
}

func (repo *PairingRepository) ApplyPairingOptions(pairingID uint, mentorID uint, targets []repository_interface.PairingTargetOption) error {
	if len(targets) == 0 {
		return fmt.Errorf("対象が指定されていません")
	}
	primaryMentee := primaryMenteeFromTargets(targets, mentorID)
	if _, err := gorm.G[entity.Pairing](repo.db).Where("id = ?", pairingID).Update(repo.ctx, "mentor_id", mentorID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	if _, err := gorm.G[entity.Pairing](repo.db).Where("id = ?", pairingID).Update(repo.ctx, "mentee_id", primaryMentee); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	if _, err := gorm.G[entity.Pairing](repo.db).Where("id = ?", pairingID).Update(repo.ctx, "recur_week", nil); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	if _, err := gorm.G[entity.Pairing](repo.db).Where("id = ?", pairingID).Update(repo.ctx, "recur_weekday", nil); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	if err := repo.db.WithContext(repo.ctx).Where("pairing_id = ?", pairingID).Delete(&entity.PairingTarget{}).Error; err != nil {
		return fmt.Errorf("対象の更新に失敗しました: %w", err)
	}
	for _, t := range targets {
		if err := repo.db.WithContext(repo.ctx).Create(&entity.PairingTarget{
			PairingID: pairingID,
			Kind:      t.Kind,
			TargetID:  t.TargetID,
		}).Error; err != nil {
			return fmt.Errorf("対象の登録に失敗しました: %w", err)
		}
	}
	return nil
}

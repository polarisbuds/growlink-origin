package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type PairingRepository struct {
	ctx                    context.Context
	db                     *gorm.DB
	baseRepository         BaseRepository[entity.Pairing]
	unScopedBaseRepository BaseUnScopedRepository[entity.Pairing]
}

func NewPairingRepository(ctx context.Context, dp *gorm.DB) *PairingRepository {
	return &PairingRepository{
		ctx: ctx,
		db:  dp,
		baseRepository: BaseRepository[entity.Pairing]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.Pairing]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *PairingRepository) CreatePairing(menteeID, mentorID uint) (*entity.Pairing, error) {
	newPair := &entity.Pairing{
		MentorID: mentorID,
		MenteeID: menteeID,
	}

	if err := repo.baseRepository.Insert(newPair); err != nil {
		return nil, fmt.Errorf("ペアリングの作成に失敗しました: %w", err)
	}

	return newPair, nil
}

func (repo *PairingRepository) BreakPairing(id uint) error {
	if _, err := repo.baseRepository.Delete(1, "id = ?", id); err != nil {
		return fmt.Errorf("ペアリングの解消に失敗しました: %w", err)
	}
	return nil
}

func (repo *PairingRepository) BreakPairings(pairingID ...uint) error {
	if _, err := repo.baseRepository.Delete(len(pairingID), "mentor_id IN ?", pairingID); err != nil {
		return fmt.Errorf("ペアリングの解消に失敗しました: %w", err)
	}
	return nil
}

func (repo *PairingRepository) BreakAllMenteePairing(mentorID uint) error {
	menties, err := repo.baseRepository.Count("mentor_id = ?", mentorID)
	if err != nil {
		return fmt.Errorf("メンターID指定でメンティー数のカウントに失敗しました: %w", err)
	}

	if _, err := repo.baseRepository.Delete(int(menties), "mentor_id = ?", mentorID); err != nil {
		return fmt.Errorf("ペアリングの解消に失敗しました: %w", err)
	}
	return nil
}

func (repo *PairingRepository) Get(pairingID uint) (*entity.Pairing, error) {
	pairing, err := repo.baseRepository.First("id ASC", "id = ?", pairingID)
	if err != nil {
		return nil, fmt.Errorf("ペアリングID指定での取得に失敗しました: %w", err)
	}
	if pairing == nil {
		return nil, fmt.Errorf("ペアリングID指定での取得に失敗しました: %w", gorm.ErrRecordNotFound)
	}
	return pairing, nil
}

func (repo *PairingRepository) GetAll() ([]entity.Pairing, error) {
	chain := repo.baseRepository.Chain("id ASC", "Mentor", "Mentee", "Sharers")
	pairing, err := chain.Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("全てのペアリング情報の取得に失敗しました: %w", err)
	}
	return pairing, nil
}

func (repo *PairingRepository) GetByMenteeID(menteeID uint) (*entity.Pairing, error) {
	pairing, err := repo.baseRepository.First("id ASC", "mentee_id = ?", menteeID)
	if err != nil {
		return nil, fmt.Errorf("メンティーID指定でのペアリング取得に失敗しました: %w", err)
	}
	return pairing, nil
}

func (repo *PairingRepository) GetByMentorID(mentorID uint) ([]entity.Pairing, error) {
	wherePairingsChain, err := repo.baseRepository.WhereChain("id ASC", "mentor_id = ?", mentorID)
	if err != nil {
		return nil, fmt.Errorf("メンターID指定でのペアリング取得用のチェイン作成に失敗しました: %w", err)
	}

	pairings, err := wherePairingsChain.Preload("Mentee", nil).Preload("Mentor", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("メンターID指定でのペアリング取得に失敗しました: %w", err)
	}
	return pairings, nil
}

func (repo *PairingRepository) UpdatePairing(pairingID, mentorID uint) error {
	if _, err := repo.baseRepository.Update("mentor_id", mentorID, "id = ?", pairingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *PairingRepository) ApplyPairingOptions(pairingID uint, mentorID, menteeID uint, recurWeek, recurWeekday *int, sharerUserIDs []uint) error {
	if _, err := repo.baseRepository.Update("mentor_id", mentorID, "id = ?", pairingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	if _, err := repo.baseRepository.Update("mentee_id", menteeID, "id = ?", pairingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	if _, err := repo.baseRepository.Update("recur_week", recurWeek, "id = ?", pairingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	if _, err := repo.baseRepository.Update("recur_weekday", recurWeekday, "id = ?", pairingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	if err := repo.db.Where("pairing_id = ?", pairingID).Delete(&entity.PairingSharer{}).Error; err != nil {
		return fmt.Errorf("共有者の更新に失敗しました: %w", err)
	}
	for _, uid := range sharerUserIDs {
		if err := repo.db.Create(&entity.PairingSharer{PairingID: pairingID, UserID: uid}).Error; err != nil {
			return fmt.Errorf("共有者の登録に失敗しました: %w", err)
		}
	}
	return nil
}

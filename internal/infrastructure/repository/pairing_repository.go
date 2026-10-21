package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
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

func (repo *PairingRepository) CreatePairing(menteeID, mentorID, createdByUserID uint) (*entity.Pairing, error) {
	newPair := &entity.Pairing{
		MentorID:        mentorID,
		MenteeID:        menteeID,
		CreatedByUserID: createdByUserID,
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

func (repo *PairingRepository) GetByUUID(pairingUUID uuid.UUID) (*entity.Pairing, error) {
	pairing, err := repo.baseRepository.First("id ASC", "uuid = ?", pairingUUID)
	if err != nil {
		return nil, fmt.Errorf("ペアリングUUID指定での取得に失敗しました: %w", err)
	}
	if pairing == nil {
		return nil, fmt.Errorf("ペアリングUUID指定での取得に失敗しました: %w", gorm.ErrRecordNotFound)
	}
	return pairing, nil
}

func (repo *PairingRepository) GetAll() ([]entity.Pairing, error) {
	chain := repo.baseRepository.Chain("id ASC", "Mentor", "Mentee", "Mentees", "Mentees.User", "Sharers")
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
	if pairing != nil {
		return pairing, nil
	}
	var link entity.PairingMentee
	if err := repo.db.Where("user_id = ?", menteeID).First(&link).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("ペアリングメンティー関連の取得に失敗しました: %w", err)
	}
	return repo.Get(link.PairingID)
}

func (repo *PairingRepository) GetByMentorID(mentorID uint) ([]entity.Pairing, error) {
	wherePairingsChain, err := repo.baseRepository.WhereChain("id ASC", "mentor_id = ?", mentorID)
	if err != nil {
		return nil, fmt.Errorf("メンターID指定でのペアリング取得用のチェイン作成に失敗しました: %w", err)
	}
	if wherePairingsChain == nil {
		return []entity.Pairing{}, nil
	}

	pairings, err := wherePairingsChain.Preload("Mentee", nil).Preload("Mentor", nil).Preload("Mentees", nil).Preload("Mentees.User", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("メンターID指定でのペアリング取得に失敗しました: %w", err)
	}
	return pairings, nil
}

func (repo *PairingRepository) GetPairingIDsBySharerUserID(userID uint) ([]uint, error) {
	var sharers []entity.PairingSharer
	if err := repo.db.Where("user_id = ?", userID).Find(&sharers).Error; err != nil {
		return nil, fmt.Errorf("共有者ペアの取得に失敗しました: %w", err)
	}
	ids := make([]uint, 0, len(sharers))
	for _, s := range sharers {
		ids = append(ids, s.PairingID)
	}
	return ids, nil
}

func (repo *PairingRepository) UpdatePairing(pairingID, mentorID uint) error {
	if _, err := repo.baseRepository.Update("mentor_id", mentorID, "id = ?", pairingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *PairingRepository) ApplyPairingOptions(pairingID uint, mentorID uint, menteeIDs []uint, recurWeek, recurWeekday *int, sharerUserIDs []uint) error {
	if len(menteeIDs) == 0 {
		return fmt.Errorf("メンティーが指定されていません")
	}
	primaryMentee := menteeIDs[0]
	if _, err := repo.baseRepository.Update("mentor_id", mentorID, "id = ?", pairingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	if _, err := repo.baseRepository.Update("mentee_id", primaryMentee, "id = ?", pairingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	if _, err := repo.baseRepository.Update("recur_week", recurWeek, "id = ?", pairingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	if _, err := repo.baseRepository.Update("recur_weekday", recurWeekday, "id = ?", pairingID); err != nil {
		return fmt.Errorf("ペアリングの更新に失敗しました: %w", err)
	}
	if err := repo.db.Where("pairing_id = ?", pairingID).Delete(&entity.PairingMentee{}).Error; err != nil {
		return fmt.Errorf("メンティー関連の更新に失敗しました: %w", err)
	}
	for _, uid := range menteeIDs {
		if err := repo.db.Create(&entity.PairingMentee{PairingID: pairingID, UserID: uid}).Error; err != nil {
			return fmt.Errorf("メンティー関連の登録に失敗しました: %w", err)
		}
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

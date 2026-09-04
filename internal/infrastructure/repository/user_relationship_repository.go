package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
	"log/slog"

	"gorm.io/gorm"
)

type UserRelationshipRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.UserRelationship]
	unScopedBaseRepository BaseUnScopedRepository[entity.UserRelationship]
}

func NewUserRelationshipRepository(ctx context.Context, dp *gorm.DB) *UserRelationshipRepository {
	return &UserRelationshipRepository{
		baseRepository: BaseRepository[entity.UserRelationship]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.UserRelationship]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *UserRelationshipRepository) Get(relationshipID uint) (*entity.UserRelationship, error) {
	relationshipWhereChain, err := repo.baseRepository.WhereChain("id ASC", "id = ?", relationshipID)
	if err != nil {
		return nil, fmt.Errorf("リレーションタイプ指定でのユーザ間関係性のチェイン作成に失敗しました: %w", err)
	}
	relationship, err := relationshipWhereChain.Preload("SourceUser", nil).Preload("TargetUser", nil).First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ID指定でのユーザ間関係性の取得に失敗しました: %w", err)
	}
	return &relationship, nil
}

func (repo *UserRelationshipRepository) GetByRelationType(relationshipType valueobject.RelationshipType) ([]entity.UserRelationship, error) {
	relationshipWhereChain, err := repo.baseRepository.WhereChain("id ASC", "relation_type = ?", relationshipType)
	if err != nil {
		return nil, fmt.Errorf("リレーションタイプ指定でのユーザ間関係性のチェイン作成に失敗しました: %w", err)
	}
	relationships, err := relationshipWhereChain.Preload("SourceUser", nil).Preload("TargetUser", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("リレーションタイプ指定でのユーザ間関係性の取得に失敗しました: %w", err)
	}
	return relationships, nil
}

func (repo *UserRelationshipRepository) GetByBossID(bossID uint) (subordinates []entity.UserRelationship, err error) {
	relationshipWhereChain, err := repo.baseRepository.WhereChain("id ASC", "relation_type = ? AND target_id = ?", valueobject.RelationshipTypeBoss, bossID)
	if err != nil {
		return nil, fmt.Errorf("上司ユーザID指定でのユーザ間関係性のチェイン作成に失敗しました: %w", err)
	}
	relationships, err := relationshipWhereChain.Preload("SourceUser", nil).Preload("TargetUser", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("上司ユーザID指定でのユーザ間関係性の取得に失敗しました: %w", err)
	}
	return relationships, nil
}

func (repo *UserRelationshipRepository) GetBySubordinateID(subordinateID uint) (bosses []entity.UserRelationship, err error) {
	relationshipWhereChain, err := repo.baseRepository.WhereChain("id ASC", "relation_type = ? AND target_id = ?", valueobject.RelationshipTypeSubordinate, subordinateID)
	if err != nil {
		return nil, fmt.Errorf("部下ユーザID指定でのユーザ間関係性のチェイン作成に失敗しました: %w", err)
	}
	relationships, err := relationshipWhereChain.Preload("SourceUser", nil).Preload("TargetUser", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("部下ユーザID指定でのユーザ間関係性の取得に失敗しました: %w", err)
	}
	return relationships, nil
}

func (repo *UserRelationshipRepository) Create(sourceUserID, targetUserID uint, relationshipType valueobject.RelationshipType) error {
	newRelation := &entity.UserRelationship{
		RelationType: relationshipType,
		SourceUserID: sourceUserID,
		TargetUserID: targetUserID,
	}

	if err := repo.baseRepository.Insert(newRelation); err != nil {
		return fmt.Errorf("ユーザリレーションシップの新規作成に失敗しました: %w", err)
	}

	return nil
}

func (repo *UserRelationshipRepository) Update(relationshipID uint, relationship *entity.UserRelationship) error {
	if relationship == nil {
		slog.Warn("ユーザリレーションシップの更新で空のUserRelationshipが送られました")
		return nil
	}

	if _, err := repo.baseRepository.Updates(relationship, "id = ?", relationshipID); err != nil {
		return fmt.Errorf("ユーザリレーションシップの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserRelationshipRepository) Delete(relationshipID uint) error {
	if _, err := repo.baseRepository.Delete(1, "id = ?", relationshipID); err != nil {
		return fmt.Errorf("ユーザリレーションシップの削除に失敗しました: %w", err)
	}
	return nil

}

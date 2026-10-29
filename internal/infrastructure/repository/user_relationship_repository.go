package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
	"log/slog"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRelationshipRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewUserRelationshipRepository(ctx context.Context, dp *gorm.DB) *UserRelationshipRepository {
	return &UserRelationshipRepository{ctx: ctx, db: dp}
}

func (repo *UserRelationshipRepository) Get(relationshipID uint) (*entity.UserRelationship, error) {
	relationship, err := gorm.G[entity.UserRelationship](repo.db).
		Where("id = ?", relationshipID).
		Order("id ASC").
		Preload("SourceUser", nil).
		Preload("TargetUser", nil).
		First(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("ID指定でのユーザ間関係性の取得に失敗しました: %w", err)
	}
	return &relationship, nil
}

func (repo *UserRelationshipRepository) GetByRelationType(relationshipType valueobject.RelationshipType) ([]entity.UserRelationship, error) {
	relationships, err := gorm.G[entity.UserRelationship](repo.db).
		Where("relation_type = ?", relationshipType).
		Order("id ASC").
		Preload("SourceUser", nil).
		Preload("TargetUser", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("リレーションタイプ指定でのユーザ間関係性の取得に失敗しました: %w", err)
	}
	return relationships, nil
}

func (repo *UserRelationshipRepository) GetByBossID(bossID uint) (subordinates []entity.UserRelationship, err error) {
	relationships, err := gorm.G[entity.UserRelationship](repo.db).
		Where("relation_type = ? AND target_user_id = ?", valueobject.RelationshipTypeBoss, bossID).
		Order("id ASC").
		Preload("SourceUser", nil).
		Preload("TargetUser", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("上司ユーザID指定でのユーザ間関係性の取得に失敗しました: %w", err)
	}
	return relationships, nil
}

func (repo *UserRelationshipRepository) GetBySubordinateID(subordinateID uint) (bosses []entity.UserRelationship, err error) {
	relationships, err := gorm.G[entity.UserRelationship](repo.db).
		Where("relation_type = ? AND target_user_id = ?", valueobject.RelationshipTypeSubordinate, subordinateID).
		Order("id ASC").
		Preload("SourceUser", nil).
		Preload("TargetUser", nil).
		Find(repo.ctx)
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

	if err := gorm.G[entity.UserRelationship](repo.db).Create(repo.ctx, newRelation); err != nil {
		return fmt.Errorf("ユーザリレーションシップの新規作成に失敗しました: %w", err)
	}

	return nil
}

func (repo *UserRelationshipRepository) Update(relationshipID uint, relationship *entity.UserRelationship) error {
	if relationship == nil {
		slog.Warn("ユーザリレーションシップの更新で空のUserRelationshipが送られました")
		return nil
	}

	if _, err := gorm.G[entity.UserRelationship](repo.db).Where("id = ?", relationshipID).Omit(clause.Associations).Updates(repo.ctx, *relationship); err != nil {
		return fmt.Errorf("ユーザリレーションシップの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *UserRelationshipRepository) Delete(relationshipID uint) error {
	if _, err := gorm.G[entity.UserRelationship](repo.db).Where("id = ?", relationshipID).Delete(repo.ctx); err != nil {
		return fmt.Errorf("ユーザリレーションシップの削除に失敗しました: %w", err)
	}
	return nil

}

func (repo *UserRelationshipRepository) DeleteByBossUserID(bossUserID uint) error {
	if _, err := gorm.G[entity.UserRelationship](repo.db).
		Where("relation_type = ? AND target_user_id = ?", valueobject.RelationshipTypeBoss, bossUserID).
		Delete(repo.ctx); err != nil {
		return fmt.Errorf("上司ユーザのリレーション削除に失敗しました: %w", err)
	}
	return nil
}

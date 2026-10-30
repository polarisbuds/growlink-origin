package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TeamMemberRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewTeamMemberRepository(ctx context.Context, dp *gorm.DB) *TeamMemberRepository {
	return &TeamMemberRepository{ctx: ctx, db: dp}
}

func (repo *TeamMemberRepository) Create(teamID, memberID uint, isReader bool) error {
	if err := gorm.G[entity.TeamMember](repo.db).Create(repo.ctx, &entity.TeamMember{
		IsReader: isReader,
		TeamID:   teamID,
		MemberID: memberID,
	}); err != nil {
		return fmt.Errorf("チームメンバーの作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamMemberRepository) Update(teamMemberID uint, member *entity.TeamMember) error {
	if _, err := gorm.G[entity.TeamMember](repo.db).Where("id = ?", teamMemberID).Omit(clause.Associations).Updates(repo.ctx, *member); err != nil {
		return fmt.Errorf("チームメンバーの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamMemberRepository) Delete(teamMemberID uint) error {
	if _, err := gorm.G[entity.TeamMember](repo.db).Where("id = ?", teamMemberID).Delete(repo.ctx); err != nil {
		return fmt.Errorf("チームメンバーの削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamMemberRepository) GetAll() ([]entity.TeamMember, error) {
	teamMembers, err := gorm.G[entity.TeamMember](repo.db).Order("id ASC").
		Preload("Team", nil).
		Preload("Team.Department", nil).
		Preload("Member", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("チームメンバーの全取得に失敗しました: %w", err)
	}
	return teamMembers, nil
}

func (repo *TeamMemberRepository) preloadWhere(query string, args ...any) ([]entity.TeamMember, error) {
	teamMembers, err := gorm.G[entity.TeamMember](repo.db).
		Where(query, args...).
		Order("id ASC").
		Preload("Team", nil).
		Preload("Team.Department", nil).
		Preload("Member", nil).
		Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("チームメンバー情報の取得に失敗しました: %w", err)
	}
	return teamMembers, nil
}

func (repo *TeamMemberRepository) GetByTeamID(teamID uint) ([]entity.TeamMember, error) {
	teamMembers, err := repo.preloadWhere("team_id = ?", teamID)
	if err != nil {
		return nil, fmt.Errorf("チームID指定でのチームメンバー取得に失敗しました: %w", err)
	}
	return teamMembers, nil
}

func (repo *TeamMemberRepository) GetByMemberID(memberID uint) ([]entity.TeamMember, error) {
	teamMembers, err := repo.preloadWhere("member_id = ?", memberID)
	if err != nil {
		return nil, fmt.Errorf("メンバーID指定でのチームメンバー取得に失敗しました: %w", err)
	}
	return teamMembers, nil
}

func (repo *TeamMemberRepository) GetReaders() ([]entity.TeamMember, error) {
	readers, err := repo.preloadWhere("is_reader = true")
	if err != nil {
		return nil, fmt.Errorf("リーダーユーザの取得に失敗しました: %w", err)
	}
	return readers, nil
}

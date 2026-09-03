package repository

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"

	"gorm.io/gorm"
)

type TeamMemberRepository struct {
	ctx                    context.Context
	baseRepository         BaseRepository[entity.TeamMember]
	unScopedBaseRepository BaseUnScopedRepository[entity.TeamMember]
}

func NewTeamMemberRepository(ctx context.Context, dp *gorm.DB) *TeamMemberRepository {
	return &TeamMemberRepository{
		baseRepository: BaseRepository[entity.TeamMember]{
			dp:   dp,
			bctx: ctx,
		},
		unScopedBaseRepository: BaseUnScopedRepository[entity.TeamMember]{
			udp:  dp.Unscoped(),
			uctx: ctx,
		},
	}
}

func (repo *TeamMemberRepository) Create(teamID, memberID uint, isReader bool) error {
	if err := repo.baseRepository.Insert(&entity.TeamMember{
		IsReader: isReader,
		TeamID:   teamID,
		MemberID: memberID,
	}); err != nil {
		return fmt.Errorf("チームメンバーの作成に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamMemberRepository) Update(teamMemberID uint, member *entity.TeamMember) error {
	if _, err := repo.baseRepository.Updates(member, "id = ?", teamMemberID); err != nil {
		return fmt.Errorf("チームメンバーの更新に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamMemberRepository) Delete(teamMemberID uint) error {
	if _, err := repo.baseRepository.Delete(1, "id = ?", teamMemberID); err != nil {
		return fmt.Errorf("チームメンバーの削除に失敗しました: %w", err)
	}
	return nil
}

func (repo *TeamMemberRepository) GetAll() ([]entity.TeamMember, error) {
	teamMembers, err := gorm.G[entity.TeamMember](repo.baseRepository.dp).Order("id ASC").
		Preload("Team", nil).
		Preload("Team.Department", nil).
		Preload("Member", nil).Find(repo.ctx)
	if err != nil {
		return nil, fmt.Errorf("チームメンバーの全取得に失敗しました: %w", err)
	}
	return teamMembers, nil
}

func (repo *TeamMemberRepository) preloadWhere(query string, args ...any) ([]entity.TeamMember, error) {
	whereChain, err := repo.baseRepository.WhereChain("id ASC", query, args...)
	if err != nil {
		return nil, fmt.Errorf("チームメンバー情報取得用のチェインの作成に失敗しました: %w", err)
	}
	if whereChain == nil {
		return []entity.TeamMember{}, nil
	}
	teamMembers, err := whereChain.Preload("Team", nil).Preload("Team.Department", nil).
		Preload("Member", nil).Find(repo.ctx)
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

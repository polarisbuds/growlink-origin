package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type TeamMemberRepository interface {
	Create(teamID, memberID uint, isReader bool) error
	Update(teamMemberID uint, member *entity.TeamMember) error
	Delete(teamMemberID uint) error

	GetReaders() ([]entity.TeamMember, error)
	GetAll() ([]entity.TeamMember, error)
	GetByTeamID(teamID uint) ([]entity.TeamMember, error)
	GetByMemberID(memberID uint) ([]entity.TeamMember, error)
}

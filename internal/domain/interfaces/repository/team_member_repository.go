package repository_interface

import "growlink/internal/domain/entity"

type TeamMemberRepository interface {
	BaseRepository[entity.TeamMember]
	BaseUnScopedRepository[entity.TeamMember]
}

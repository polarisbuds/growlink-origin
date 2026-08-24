package repository_interface

import "terms-checker/internal/domain/entity"

type TeamMemberRepository interface {
	BaseRepository[entity.TeamMember]
	BaseUnScopedRepository[entity.TeamMember]
}

package outputmodel

import "growlink/internal/domain/valueobject"

type UserRelationshipOutput struct {
	RelationType valueobject.RelationshipType `json:"relationType"`
	SourceUser   UserOutput                   `json:"sourceUser"`
	TargetUser   UserOutput                   `json:"targetUser"`
}

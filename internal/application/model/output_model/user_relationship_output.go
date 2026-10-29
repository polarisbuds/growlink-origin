package outputmodel

import "growlink/internal/domain/valueobject"

type UserRelationshipOutput struct {
	ID           uint                         `json:"id"`
	RelationType valueobject.RelationshipType `json:"relationType"`
	SourceUserID uint                         `json:"sourceUserId"`
	TargetUserID uint                         `json:"targetUserId"`
	SourceUser   UserOutput                   `json:"sourceUser"`
	TargetUser   UserOutput                   `json:"targetUser"`
}

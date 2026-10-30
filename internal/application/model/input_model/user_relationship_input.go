package inputmodel

import "growlink/internal/domain/valueobject"

type UserRelationshipInput struct {
	RelationType valueobject.RelationshipType `json:"relationType"`
	SourceUserID uint                         `json:"sourceUserId"`
	TargetUserID uint                         `json:"targetUserId"`
}

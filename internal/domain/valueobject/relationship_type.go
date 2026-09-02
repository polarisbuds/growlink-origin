package valueobject

import "fmt"

type RelationshipType string

const (
	RelationshipTypeMentor      RelationshipType = "mentor"
	RelationshipTypeMentee      RelationshipType = "mentee"
	RelationshipTypeBoss        RelationshipType = "boss"
	RelationshipTypeSubordinate RelationshipType = "subordinate"
)

func (t RelationshipType) String() string {
	return string(t)
}

func (t RelationshipType) IsValid() bool {
	switch t {
	case RelationshipTypeMentor, RelationshipTypeMentee, RelationshipTypeBoss, RelationshipTypeSubordinate:
		return true
	default:
		return false
	}
}

func ParseRelationshipType(s string) (RelationshipType, error) {
	relationType := RelationshipType(s)
	if !relationType.IsValid() {
		return "", fmt.Errorf("invalid relationship type: %s", s)
	}
	return relationType, nil
}

package valueobject

import "fmt"

type RelationshipType int

const (
	RelationshipTypeMentor      RelationshipType = iota + 1
	RelationshipTypeMentee
	RelationshipTypeBoss
	RelationshipTypeSubordinate
)

func (t RelationshipType) String() string {
	switch t {
	case RelationshipTypeMentor:
		return "mentor"
	case RelationshipTypeMentee:
		return "mentee"
	case RelationshipTypeBoss:
		return "boss"
	case RelationshipTypeSubordinate:
		return "subordinate"
	default:
		return ""
	}
}

func (t RelationshipType) Int() int {
	return int(t)
}

func (t RelationshipType) UInt() uint {
	return uint(t)
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
	switch s {
	case RelationshipTypeMentor.String():
		return RelationshipTypeMentor, nil
	case RelationshipTypeMentee.String():
		return RelationshipTypeMentee, nil
	case RelationshipTypeBoss.String():
		return RelationshipTypeBoss, nil
	case RelationshipTypeSubordinate.String():
		return RelationshipTypeSubordinate, nil
	default:
		return 0, fmt.Errorf("invalid relationship type: %s", s)
	}
}

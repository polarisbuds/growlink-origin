package valueobject

import "fmt"

type RelationshipType int

const (
	RelationshipTypeBoss RelationshipType = iota + 1
	RelationshipTypeSubordinate
)

func (t RelationshipType) String() string {
	switch t {
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
	case RelationshipTypeBoss, RelationshipTypeSubordinate:
		return true
	default:
		return false
	}
}

func ParseRelationshipType(s string) (RelationshipType, error) {
	switch s {
	case RelationshipTypeBoss.String():
		return RelationshipTypeBoss, nil
	case RelationshipTypeSubordinate.String():
		return RelationshipTypeSubordinate, nil
	default:
		return 0, fmt.Errorf("invalid relationship type: %s", s)
	}
}

package valueobject

import "fmt"

type OperationTarget int

const (
	OperationTargetNone OperationTarget = iota + 1
	OperationTargetDepartment
	OperationTargetTeam
	OperationTargetSubordinate
)

func (t OperationTarget) String() string {
	switch t {
	case OperationTargetNone:
		return "none"
	case OperationTargetDepartment:
		return "department"
	case OperationTargetTeam:
		return "team"
	case OperationTargetSubordinate:
		return "subordinate"
	default:
		return ""
	}
}

func (t OperationTarget) Int() int {
	return int(t)
}

func (t OperationTarget) IsValid() bool {
	switch t {
	case OperationTargetNone, OperationTargetDepartment, OperationTargetTeam, OperationTargetSubordinate:
		return true
	default:
		return false
	}
}

func ParseOperationTarget(s string) (OperationTarget, error) {
	switch s {
	case OperationTargetNone.String():
		return OperationTargetNone, nil
	case OperationTargetDepartment.String():
		return OperationTargetDepartment, nil
	case OperationTargetTeam.String():
		return OperationTargetTeam, nil
	case OperationTargetSubordinate.String():
		return OperationTargetSubordinate, nil
	default:
		return 0, fmt.Errorf("invalid operation target: %s", s)
	}
}

func OperationTargetLabel(t OperationTarget) string {
	switch t {
	case OperationTargetDepartment:
		return "部署"
	case OperationTargetTeam:
		return "チーム"
	case OperationTargetSubordinate:
		return "部下"
	default:
		return "なし"
	}
}

package valueobject

import "fmt"

type SystemRole int

const (
	SystemRoleNone SystemRole = iota + 1
	SystemRoleAdmin
	SystemRoleAnalyst
)

func (r SystemRole) String() string {
	switch r {
	case SystemRoleNone:
		return "none"
	case SystemRoleAdmin:
		return "system_admin"
	case SystemRoleAnalyst:
		return "analyst"
	default:
		return ""
	}
}

func (r SystemRole) Int() int {
	return int(r)
}

func (r SystemRole) IsValid() bool {
	switch r {
	case SystemRoleNone, SystemRoleAdmin, SystemRoleAnalyst:
		return true
	default:
		return false
	}
}

func ParseSystemRole(s string) (SystemRole, error) {
	switch s {
	case SystemRoleNone.String():
		return SystemRoleNone, nil
	case SystemRoleAdmin.String():
		return SystemRoleAdmin, nil
	case SystemRoleAnalyst.String():
		return SystemRoleAnalyst, nil
	default:
		return 0, fmt.Errorf("invalid system role: %s", s)
	}
}

func SystemRoleLabel(r SystemRole) string {
	switch r {
	case SystemRoleAdmin:
		return "システム管理者"
	case SystemRoleAnalyst:
		return "分析者"
	default:
		return "なし"
	}
}

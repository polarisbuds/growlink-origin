package valueobject

import "fmt"

type Authority int

const (
	AuthorityGeneral Authority = iota + 1
	AuthorityBoss
	AuthorityOwner
)

func (a Authority) String() string {
	switch a {
	case AuthorityGeneral:
		return "general"
	case AuthorityBoss:
		return "boss"
	case AuthorityOwner:
		return "owner"
	default:
		return ""
	}
}

func (a Authority) Int() int {
	return int(a)
}

func (a Authority) UInt() uint {
	return uint(a)
}

func (a Authority) IsValid() bool {
	switch a {
	case AuthorityGeneral, AuthorityBoss, AuthorityOwner:
		return true
	default:
		return false
	}
}

func ParseAuthority(s string) (Authority, error) {
	switch s {
	case AuthorityGeneral.String():
		return AuthorityGeneral, nil
	case AuthorityBoss.String():
		return AuthorityBoss, nil
	case AuthorityOwner.String():
		return AuthorityOwner, nil
	default:
		return 0, fmt.Errorf("invalid authority: %s", s)
	}
}

package valueobject

import "fmt"

type Authority string

const (
	AuthorityGeneral Authority = "general"
	AuthorityBoss    Authority = "boss"
	AuthorityOwner   Authority = "owner"
)

func (a Authority) String() string {
	return string(a)
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
	authority := Authority(s)
	if !authority.IsValid() {
		return "", fmt.Errorf("invalid authority: %s", s)
	}
	return authority, nil
}

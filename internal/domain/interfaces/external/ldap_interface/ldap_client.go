package ldap_interface

import (
	"growlink/internal/domain/entity"
)

type LDAPClient interface {
	Search(name string) (*entity.LDAPUser, error)
	Authenticate(userName, password string) (*entity.LDAPUser, error)
}

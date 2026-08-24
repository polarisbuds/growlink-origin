package ldap_repository_interface

import (
	"growlink/internal/domain/entity"
)

type LDAPRepository interface {
	FindUser(name string) (*entity.LDAPUser, error)
	Login(userName, password string) (*entity.LDAPUser, error)
}

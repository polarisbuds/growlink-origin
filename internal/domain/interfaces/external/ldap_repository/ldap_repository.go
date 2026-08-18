package ldap_repository

import (
	"terms-checker/internal/domain/entity"
)

type LDAPRepository interface {
	FindUser(name string) (*entity.LDAPUser, error)
	Login(userName, password string) (*entity.LDAPUser, error)
}

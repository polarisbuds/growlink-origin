package ldap_interface

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import ldap_entity "growlink/internal/domain/entity/ldap"

type LDAPClient interface {
	Search(name string) (*ldap_entity.LDAPUser, error)
	Authenticate(userName, password string) (*ldap_entity.LDAPUser, error)
}

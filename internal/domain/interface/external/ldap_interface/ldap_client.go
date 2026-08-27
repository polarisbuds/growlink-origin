package ldap_interface

import ldap_entity "growlink/internal/domain/entity/ldap"

type LDAPClient interface {
	Search(name string) (*ldap_entity.LDAPUser, error)
	Authenticate(userName, password string) (*ldap_entity.LDAPUser, error)
}

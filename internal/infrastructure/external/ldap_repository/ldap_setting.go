package ldap_repository

import (
	"fmt"
	"os"
)

type LDAPSetting struct {
	URL            string
	BaseDN         string
	BindDN         string
	BinddnPassword string
	SearchAttr     string
	SearchFilter   string
}

func NewLDAPSetting() (*LDAPSetting, error) {
	setting := LDAPSetting{}

	url := os.Getenv("LDAP_URL")
	if url == "" {
		return nil, ldapError("LDAP_URL")
	}
	setting.URL = url

	baseDN := os.Getenv("LDAP_BASE_DN")
	if baseDN == "" {
		return nil, ldapError("LDAP_BASE_DN")
	}
	setting.BaseDN = baseDN

	bindDN := os.Getenv("LDAP_BIND_DN")
	if bindDN == "" {
		return nil, ldapError("LDAP_BIND_DN")
	}
	setting.BindDN = bindDN

	binddnPassword := os.Getenv("LDAP_BIND_DN_PASSWORD")
	if binddnPassword == "" {
		return nil, ldapError("LDAP_BIND_DN_PASSWORD")
	}
	setting.BinddnPassword = binddnPassword

	searchAttr := os.Getenv("LDAP_SEARCH_ATTR")
	if searchAttr == "" {
		return nil, ldapError("LDAP_SEARCH_ATTR")
	}
	setting.SearchAttr = searchAttr

	searchFilter := os.Getenv("LDAP_SEARCH_FILTER")
	if searchFilter == "" {
		return nil, ldapError("LDAP_SEARCH_FILTER")
	}
	setting.SearchFilter = searchFilter

	return &setting, nil
}

const errMsgTemplate = "LDAP認証設定エラー: %s を環境変数に登録してください。"

func ldapError(env string) error {
	return fmt.Errorf(errMsgTemplate, env)
}

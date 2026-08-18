package ldap_repository

import (
	"fmt"
	"sync"
	"terms-checker/internal/domain/entity"

	"github.com/go-ldap/ldap/v3"
)

type LDAPRepository struct {
	setting *LDAPSetting
	mu      sync.RWMutex
}

func NewLDAPRepository(setting *LDAPSetting) *LDAPRepository {
	return &LDAPRepository{
		setting: setting,
	}
}

func (repo *LDAPRepository) newConnection() (*ldap.Conn, error) {
	setting := repo.setting
	conn, err := ldap.DialURL(setting.URL)
	if err != nil {
		return nil, fmt.Errorf("コネクションの作成に失敗しました: %w", err)
	}
	return conn, err
}

func (repo *LDAPRepository) bind(conn *ldap.Conn, bindDN, password string) error {
	if err := conn.Bind(bindDN, password); err != nil {
		return fmt.Errorf("コネクションのバインドに失敗しました: %w", err)
	}
	return nil
}

func (repo *LDAPRepository) newLDAPSearchRequest(filter string) *ldap.SearchRequest {
	searchReq := ldap.NewSearchRequest(
		repo.setting.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{
			"sAMAccountName",
			"displayName",
			"department",
			"mail",
			"distinguishedName",
		},
		nil,
	)
	return searchReq
}

func (repo *LDAPRepository) makeFilterString(uname string) string {
	return fmt.Sprintf(
		"(&%s(%s=%s))",
		repo.setting.SearchFilter,
		repo.setting.SearchAttr,
		ldap.EscapeFilter(uname))
}

func (repo *LDAPRepository) FindUser(userName string) (*entity.LDAPUser, error) {
	conn, err := repo.newConnection()
	if err != nil {
		return nil, fmt.Errorf("認証用接続に失敗しました: %w", err)
	}
	defer conn.Close()

	if err := repo.bind(conn, repo.setting.BindDN, repo.setting.BinddnPassword); err != nil {
		return nil, fmt.Errorf("パスワード認証に失敗しました: %w", err)
	}

	filter := repo.makeFilterString(userName)
	req := repo.newLDAPSearchRequest(filter)
	result, err := conn.Search(req)
	if err != nil {
		return nil, fmt.Errorf("LDAP認証のユーザ検索に失敗しました: %w", err)
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("LDAP認証の結果ユーザが見つかりませんでした: %w", err)
	}

	e := result.Entries[0]

	return &entity.LDAPUser{
		UserName:   e.GetAttributeValue("sAMAccountName"),
		Name:       e.GetAttributeValue("displayName"),
		Department: e.GetAttributeValue("department"),
		EmployeeID: e.GetAttributeValue("employeeID"),
		Mail:       e.GetAttributeValue("mail"),
		DN:         e.GetAttributeValue("distinguishedName"),
	}, nil
}

func (repo *LDAPRepository) Login(userName, password string) (*entity.LDAPUser, error) {
	user, err := repo.FindUser(userName)
	if err != nil {
		return nil, err
	}

	conn, err := repo.newConnection()
	if err != nil {
		return nil, fmt.Errorf("認証用接続に失敗: %w", err)
	}
	defer conn.Close()

	if err := repo.bind(conn, user.DN, password); err != nil {
		return nil, fmt.Errorf("パスワード認証に失敗: %w", err)
	}

	return user, nil
}

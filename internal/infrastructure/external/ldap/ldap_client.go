package ldap

import (
	"fmt"
	"growlink/internal/domain/entity"

	"github.com/go-ldap/ldap/v3"
)

type LDAPClient struct {
	setting *LDAPSetting
}

func NewClient(setting *LDAPSetting) *LDAPClient {
	return &LDAPClient{
		setting: setting,
	}
}

func (client *LDAPClient) newConnection() (*ldap.Conn, error) {
	setting := client.setting
	conn, err := ldap.DialURL(setting.URL)
	if err != nil {
		return nil, fmt.Errorf("コネクションの作成に失敗しました: %w", err)
	}
	return conn, err
}

func (client *LDAPClient) bind(conn *ldap.Conn, bindDN, password string) error {
	if err := conn.Bind(bindDN, password); err != nil {
		return fmt.Errorf("コネクションのバインドに失敗しました: %w", err)
	}
	return nil
}

func (client *LDAPClient) newLDAPSearchRequest(filter string) *ldap.SearchRequest {
	searchReq := ldap.NewSearchRequest(
		client.setting.BaseDN,
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
			"employeeID",
		},
		nil,
	)
	return searchReq
}

func (client *LDAPClient) makeFilterString(uname string) string {
	return fmt.Sprintf(
		"(&%s(%s=%s))",
		client.setting.SearchFilter,
		client.setting.SearchAttr,
		ldap.EscapeFilter(uname))
}

func (client *LDAPClient) Search(userName string) (*entity.LDAPUser, error) {
	conn, err := client.newConnection()
	if err != nil {
		return nil, fmt.Errorf("検索用接続に失敗しました: %w", err)
	}
	defer conn.Close()

	if err := client.bind(conn, client.setting.BindDN, client.setting.BinddnPassword); err != nil {
		return nil, fmt.Errorf("バインドに失敗しました: %w", err)
	}

	filter := client.makeFilterString(userName)
	req := client.newLDAPSearchRequest(filter)
	result, err := conn.Search(req)
	if err != nil {
		return nil, fmt.Errorf("LDAP認証のユーザ検索に失敗しました: %w", err)
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("LDAP認証の結果、ユーザが見つかりませんでした\n")
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

func (client *LDAPClient) Authenticate(userName, password string) (*entity.LDAPUser, error) {
	user, err := client.Search(userName)
	if err != nil {
		return nil, err
	}

	conn, err := client.newConnection()
	if err != nil {
		return nil, fmt.Errorf("認証用接続に失敗: %w", err)
	}
	defer conn.Close()

	if err := client.bind(conn, user.DN, password); err != nil {
		return nil, fmt.Errorf("パスワード認証に失敗: %w", err)
	}

	return user, nil
}

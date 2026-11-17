package ldap

import (
	"fmt"
	"strings"

	ldap_entity "growlink/internal/domain/entity/ldap"

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

func (client *LDAPClient) Search(userName string) (*ldap_entity.LDAPUser, error) {
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
		return nil, fmt.Errorf("LDAP検索: ユーザの検索に失敗しました: %w", err)
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("LDAP検索: 検索結果が0件でした。")
	}

	e := result.Entries[0]

	return &ldap_entity.LDAPUser{
		UserName:   e.GetAttributeValue("sAMAccountName"),
		Name:       e.GetAttributeValue("displayName"),
		Department: e.GetAttributeValue("department"),
		EmployeeID: e.GetAttributeValue("employeeID"),
		Mail:       e.GetAttributeValue("mail"),
		DN:         e.GetAttributeValue("distinguishedName"),
	}, nil
}

func (client *LDAPClient) ListProperActive() ([]ldap_entity.LDAPUser, error) {
	filter := strings.TrimSpace(client.setting.SearchFilterProper)
	if filter == "" {
		return nil, fmt.Errorf("LDAP_SEARCH_FILTER_BY_PROPER が未設定です")
	}

	conn, err := client.newConnection()
	if err != nil {
		return nil, fmt.Errorf("検索用接続に失敗しました: %w", err)
	}
	defer conn.Close()

	if err := client.bind(conn, client.setting.BindDN, client.setting.BinddnPassword); err != nil {
		return nil, fmt.Errorf("バインドに失敗しました: %w", err)
	}

	req := client.newLDAPSearchRequest(filter)
	paging := ldap.NewControlPaging(100)
	req.Controls = []ldap.Control{paging}

	users := []ldap_entity.LDAPUser{}
	for {
		result, err := conn.Search(req)
		if err != nil {
			return nil, fmt.Errorf("LDAP一覧の取得に失敗しました: %w", err)
		}
		for _, entry := range result.Entries {
			if entry == nil {
				continue
			}
			mail := strings.TrimSpace(entry.GetAttributeValue("mail"))
			if mail == "" {
				continue
			}
			dn := entry.GetAttributeValue("distinguishedName")
			if !dnIsActive(dn) {
				continue
			}
			users = append(users, ldap_entity.LDAPUser{
				UserName:   entry.GetAttributeValue("sAMAccountName"),
				Name:       entry.GetAttributeValue("displayName"),
				Department: entry.GetAttributeValue("department"),
				EmployeeID: entry.GetAttributeValue("employeeID"),
				Mail:       mail,
				DN:         dn,
			})
		}

		updated, ok := ldap.FindControl(result.Controls, ldap.ControlTypePaging).(*ldap.ControlPaging)
		if !ok || len(updated.Cookie) == 0 {
			break
		}
		paging.SetCookie(updated.Cookie)
	}
	return users, nil
}

func dnIsActive(dn string) bool {
	for _, part := range strings.Split(dn, ",") {
		part = strings.TrimSpace(part)
		if len(part) < 3 || !strings.HasPrefix(part, "OU=") {
			continue
		}
		if part[3:] == "Delete Users" {
			return false
		}
	}
	return true
}

func (client *LDAPClient) Authenticate(userName, password string) (*ldap_entity.LDAPUser, error) {
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

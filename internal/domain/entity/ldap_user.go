package entity

// LDAP認証(AD認証)の結果を格納する
type LDAPUser struct {
	UserName   string `json:"name"`
	Name       string `json:"display_name"`
	Department string `json:"department"`
	EmployeeID string `json:"employee_id"`
	Mail       string `json:"email"`
	DN         string `json:"dn"`
}

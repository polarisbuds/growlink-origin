package outputmodel

type LDAPUserOutput struct {
	UserName   string `json:"userName"`
	Name       string `json:"name"`
	Department string `json:"department"`
	EmployeeID string `json:"employeeId"`
	Mail       string `json:"mail"`
	DN         string `json:"dn"`
}

package outputmodel

type ManagementUserOutput struct {
	ID                uint     `json:"id"`
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	ImageURI          *string  `json:"imageUri,omitempty"`
	Authority         string   `json:"authority"`
	SystemRole        string   `json:"systemRole"`
	OperationTarget   string   `json:"operationTarget"`
	TargetIDs         []uint   `json:"targetIds"`
	DepartmentID      uint     `json:"departmentId"`
	DepartmentName    string   `json:"departmentName"`
	RoleID            uint     `json:"roleId"`
	RoleName          string   `json:"roleName"`
	Suspended         bool     `json:"suspended"`
}

type OrgTreeNodeOutput struct {
	Kind         string `json:"kind"`
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Depth        int    `json:"depth"`
	DepartmentID uint   `json:"departmentId,omitempty"`
	LeaderName   string `json:"leaderName,omitempty"`
	IsRoot       bool   `json:"isRoot,omitempty"`
}

type OrgMemberOutput struct {
	UserID   uint   `json:"userId"`
	UserName string `json:"userName"`
	Title    string `json:"title"`
}

type OrgDetailOutput struct {
	Kind         string            `json:"kind"`
	ID           uint              `json:"id"`
	Name         string            `json:"name"`
	IsRoot       bool              `json:"isRoot"`
	DepartmentID uint              `json:"departmentId,omitempty"`
	Members      []OrgMemberOutput `json:"members"`
	ChildIDs     []uint            `json:"childIds"`
	TeamIDs      []uint            `json:"teamIds"`
	UserIDs      []uint            `json:"userIds"`
	MemberIDs    []uint            `json:"memberIds"`
	LeaderID     uint              `json:"leaderId"`
}

type RemindSettingsOutput struct {
	RemindAutoSend bool   `json:"remindAutoSend"`
	RemindWeekday  int    `json:"remindWeekday"`
	RemindTime     string `json:"remindTime"`
}

type RoleOptionOutput struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type DepartmentOptionOutput struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type TeamOptionOutput struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	DepartmentID uint   `json:"departmentId"`
}

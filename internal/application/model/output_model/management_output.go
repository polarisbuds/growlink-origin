package outputmodel

type UserDepartmentAssignmentOutput struct {
	DepartmentUUID string `json:"departmentUuid"`
	RoleUUID       string `json:"roleUuid"`
}

type UserTeamAssignmentOutput struct {
	TeamUUID string `json:"teamUuid"`
	IsLeader bool   `json:"isLeader"`
}

type ManagementUserOutput struct {
	UUID           string                           `json:"uuid"`
	Name           string                           `json:"name"`
	Email          string                           `json:"email"`
	ImageURI       *string                          `json:"imageUri,omitempty"`
	Authority      string                           `json:"authority"`
	SystemRole     string                           `json:"systemRole"`
	DepartmentUUID string                           `json:"departmentUuid"`
	DepartmentName string                           `json:"departmentName"`
	RoleUUID       string                           `json:"roleUuid"`
	RoleName       string                           `json:"roleName"`
	Departments    []UserDepartmentAssignmentOutput `json:"departments"`
	Teams          []UserTeamAssignmentOutput       `json:"teams"`
	Suspended      bool                             `json:"suspended"`
}

type OrgTreeNodeOutput struct {
	Kind            string `json:"kind"`
	UUID            string `json:"uuid"`
	Name            string `json:"name"`
	Depth           int    `json:"depth"`
	LeaderName      string `json:"leaderName,omitempty"`
	LeaderRole      string `json:"leaderRole,omitempty"`
	HeadUserName    string `json:"headUserName,omitempty"`
	HeadUserRole    string `json:"headUserRole,omitempty"`
	IsRoot          bool   `json:"isRoot,omitempty"`
	MemberCount     int    `json:"memberCount,omitempty"`
	TeamMemberCount int    `json:"teamMemberCount,omitempty"`
	TeamCount       int    `json:"teamCount,omitempty"`
	DepartmentCount int    `json:"departmentCount,omitempty"`
}

type OrgMemberOutput struct {
	UserUUID string `json:"userUuid"`
	UserName string `json:"userName"`
	Title    string `json:"title"`
}

type OrgDetailOutput struct {
	Kind                  string            `json:"kind"`
	UUID                  string            `json:"uuid"`
	Name                  string            `json:"name"`
	IsRoot                bool              `json:"isRoot"`
	DepartmentUUID        string            `json:"departmentUuid,omitempty"`
	Members               []OrgMemberOutput `json:"members"`
	ChildDepartmentUUIDs  []string          `json:"childDepartmentUuids"`
	TeamUUIDs             []string          `json:"teamUuids"`
	UserUUIDs             []string          `json:"userUuids"`
	MemberUUIDs           []string          `json:"memberUuids"`
	LeaderUUID            string            `json:"leaderUuid"`
	HeadUserUUID          string            `json:"headUserUuid"`
	ManagerUserUUIDs      []string          `json:"managerUserUuids"`
}

type RemindSettingsOutput struct {
	RemindAutoSend   bool   `json:"remindAutoSend"`
	RemindMode       string `json:"remindMode"`
	RemindDaysBefore int    `json:"remindDaysBefore"`
	RemindWeekday    int    `json:"remindWeekday"`
	RemindTime       string `json:"remindTime"`
}

type RoleOptionOutput struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type DepartmentOptionOutput struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	IsRoot bool   `json:"isRoot"`
}

type TeamOptionOutput struct {
	UUID           string `json:"uuid"`
	Name           string `json:"name"`
	DepartmentUUID string `json:"departmentUuid"`
}

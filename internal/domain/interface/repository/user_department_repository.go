package repository_interface

type UserDepartmentRepository interface {
	GetUserIDsByDepartmentID(departmentID uint) ([]uint, error)
	SyncDepartmentMembers(departmentID uint, userIDs []uint) error
	GetDepartmentIDsByUserID(userID uint) ([]uint, error)
}

package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type DepartmentMemberRepository interface {
	GetAll() ([]entity.DepartmentMember, error)
	GetByDepartmentID(departmentID uint) ([]entity.DepartmentMember, error)
	GetByDepartmentIDs(departmentIDs []uint) (ret0 map[uint][]entity.DepartmentMember, err error)
	GetByDepartmentUUIDs(departmentUUIDs []string) (ret0 map[string][]entity.DepartmentMember, err error)
	GetByUserID(userID uint) ([]entity.DepartmentMember, error)
	GetByUserIDs(userIDs []uint) (map[uint][]entity.DepartmentMember, error)
	ReplaceByDepartment(departmentID uint, members []entity.DepartmentMember) error
	ReplaceByUser(userID uint, member entity.DepartmentMember) error
}

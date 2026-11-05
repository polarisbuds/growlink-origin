package repository_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type DepartmentMemberRepository interface {
	GetByDepartmentID(departmentID uint) ([]entity.DepartmentMember, error)
	GetByUserID(userID uint) ([]entity.DepartmentMember, error)
	ReplaceByDepartment(departmentID uint, members []entity.DepartmentMember) error
}

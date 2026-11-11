package user_admin_usecase

import (
	"fmt"

	outputmodel "growlink/internal/application/model/output_model"

	"github.com/google/uuid"
)

type deptAssignmentResolved struct {
	DepartmentID uint
	RoleID       uint
}

type teamAssignmentResolved struct {
	TeamID uint
}

func (uc *UserAdmin) resolveDepartmentAssignments(rows []DepartmentAssignmentInput) ([]deptAssignmentResolved, error) {
	out := make([]deptAssignmentResolved, 0, len(rows))
	for _, row := range rows {
		deptUUID, err := uuid.Parse(row.DepartmentUuid)
		if err != nil {
			return nil, fmt.Errorf("所属部署を選択してください")
		}
		roleUUID, err := uuid.Parse(row.RoleUuid)
		if err != nil {
			return nil, fmt.Errorf("役職を選択してください")
		}
		dept, err := uc.deptRepo.GetByUUID(deptUUID)
		if err != nil {
			return nil, fmt.Errorf("所属部署を選択してください")
		}
		role, err := uc.roleRepo.GetByUUID(roleUUID)
		if err != nil {
			return nil, fmt.Errorf("役職を選択してください")
		}
		out = append(out, deptAssignmentResolved{
			DepartmentID: dept.ID,
			RoleID:       role.ID,
		})
	}
	return out, nil
}

func (uc *UserAdmin) resolveTeamAssignments(rows []TeamAssignmentInput) ([]teamAssignmentResolved, error) {
	out := make([]teamAssignmentResolved, 0, len(rows))
	for _, row := range rows {
		teamUUID, err := uuid.Parse(row.TeamUuid)
		if err != nil {
			return nil, fmt.Errorf("チームを選択してください")
		}
		team, err := uc.teamRepo.GetByUUID(teamUUID)
		if err != nil {
			return nil, fmt.Errorf("チームを選択してください")
		}
		out = append(out, teamAssignmentResolved{
			TeamID: team.ID,
		})
	}
	return out, nil
}

func (uc *UserAdmin) departmentAssignmentsOutput(rows []deptAssignmentResolved) ([]outputmodel.UserDepartmentAssignmentOutput, error) {
	depts, err := uc.deptRepo.GetAll()
	if err != nil {
		return nil, err
	}
	deptIDToUUID := map[uint]string{}
	for _, dept := range depts {
		deptIDToUUID[dept.ID] = dept.UUID.String()
	}
	roles, err := uc.roleRepo.GetAll()
	if err != nil {
		return nil, err
	}
	roleIDToUUID := map[uint]string{}
	for _, role := range roles {
		roleIDToUUID[role.ID] = role.UUID.String()
	}
	out := make([]outputmodel.UserDepartmentAssignmentOutput, 0, len(rows))
	for _, row := range rows {
		out = append(out, outputmodel.UserDepartmentAssignmentOutput{
			DepartmentUUID: deptIDToUUID[row.DepartmentID],
			RoleUUID:       roleIDToUUID[row.RoleID],
		})
	}
	return out, nil
}

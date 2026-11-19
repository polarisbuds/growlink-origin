package org_admin_usecase

import (
	outputmodel "growlink/internal/application/model/output_model"

	"github.com/google/uuid"
)

func (uc *OrgAdmin) ResolveUserID(param string) (uint, error) {
	userUUID, err := uuid.Parse(param)
	if err != nil {
		return 0, err
	}
	user, err := uc.userRepo.GetByUUID(userUUID)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (uc *OrgAdmin) ResolveUserIDs(params []string) ([]uint, error) {
	users, err := uc.userRepo.GetByUUIDs(params)
	if err != nil {
		return nil, err
	}
	ids := make([]uint, len(users))
	for idx, user := range users {
		ids[idx] = user.ID
	}
	return ids, nil
}

func (uc *OrgAdmin) ResolveDepartmentID(param string) (uint, error) {
	deptUUID, err := uuid.Parse(param)
	if err != nil {
		return 0, err
	}
	dept, err := uc.deptRepo.GetByUUID(deptUUID)
	if err != nil {
		return 0, err
	}
	return dept.ID, nil
}

func (uc *OrgAdmin) ResolveTeamID(param string) (uint, error) {
	teamUUID, err := uuid.Parse(param)
	if err != nil {
		return 0, err
	}
	team, err := uc.teamRepo.GetByUUID(teamUUID)
	if err != nil {
		return 0, err
	}
	return team.ID, nil
}

func (uc *OrgAdmin) GetDepartmentDetailByUUID(param string) (*outputmodel.OrgDetailOutput, error) {
	id, err := uc.ResolveDepartmentID(param)
	if err != nil {
		return nil, err
	}
	return uc.GetDepartmentDetail(id)
}

func (uc *OrgAdmin) GetTeamDetailByUUID(param string) (*outputmodel.OrgDetailOutput, error) {
	id, err := uc.ResolveTeamID(param)
	if err != nil {
		return nil, err
	}
	return uc.GetTeamDetail(id)
}

func (uc *OrgAdmin) UpdateDepartmentByUUID(param string, input UpdateDepartmentInput) error {
	id, err := uc.ResolveDepartmentID(param)
	if err != nil {
		return err
	}
	return uc.UpdateDepartment(id, input)
}

func (uc *OrgAdmin) UpdateTeamByUUID(param string, input UpdateTeamInput) error {
	id, err := uc.ResolveTeamID(param)
	if err != nil {
		return err
	}
	return uc.UpdateTeam(id, input)
}

func (uc *OrgAdmin) DeleteDepartmentByUUID(param string) error {
	id, err := uc.ResolveDepartmentID(param)
	if err != nil {
		return err
	}
	return uc.DeleteDepartment(id)
}

func (uc *OrgAdmin) DeleteTeamByUUID(param string) error {
	id, err := uc.ResolveTeamID(param)
	if err != nil {
		return err
	}
	return uc.DeleteTeam(id)
}

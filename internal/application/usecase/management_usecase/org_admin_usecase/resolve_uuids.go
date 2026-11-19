package org_admin_usecase

import (
	"fmt"

	"github.com/google/uuid"
)

func (uc *OrgAdmin) resolveDepartmentUUIDs(uuids []string) ([]uint, error) {
	if uuids == nil {
		return nil, nil
	}
	out := make([]uint, len(uuids))
	depts, err := uc.deptRepo.GetByUUIDs(uuids)
	if err != nil {
		return nil, err
	}
	for idx, dept := range depts {
		out[idx] = dept.ID
	}
	return out, nil
}

func (uc *OrgAdmin) resolveTeamUUIDs(uuids []string) ([]uint, error) {
	if uuids == nil {
		return nil, nil
	}
	out := make([]uint, len(uuids))
	teams, err := uc.teamRepo.GetByUUIDs(uuids)
	if err != nil {
		return nil, fmt.Errorf("チームが見つかりません")
	}
	for idx, team := range teams {
		out[idx] = team.ID
	}
	return out, nil
}

func (uc *OrgAdmin) resolveUserUUIDs(uuids []string) ([]uint, error) {
	if uuids == nil {
		return nil, nil
	}
	out := make([]uint, len(uuids))
	users, err := uc.userRepo.GetByUUIDs(uuids)
	if err != nil {
		return nil, fmt.Errorf("invalid user uuid")
	}
	for idx, user := range users {
		out[idx] = user.ID
	}
	return out, nil
}

func (uc *OrgAdmin) resolveUserUUIDOptional(s string) (uint, error) {
	if s == "" {
		return 0, nil
	}
	parsed, err := uuid.Parse(s)
	if err != nil {
		return 0, fmt.Errorf("invalid user uuid")
	}
	user, err := uc.userRepo.GetByUUID(parsed)
	if err != nil {
		return 0, fmt.Errorf("ユーザが見つかりません")
	}
	return user.ID, nil
}

func (uc *OrgAdmin) resolveDepartmentUUIDOptional(s string) (uint, error) {
	if s == "" {
		return 0, nil
	}
	parsed, err := uuid.Parse(s)
	if err != nil {
		return 0, fmt.Errorf("invalid department uuid")
	}
	dept, err := uc.deptRepo.GetByUUID(parsed)
	if err != nil {
		return 0, fmt.Errorf("部署が見つかりません")
	}
	return dept.ID, nil
}

func (uc *OrgAdmin) userUUIDByID(id uint) string {
	if id == 0 {
		return ""
	}
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return ""
	}
	return user.UUID.String()
}

func (uc *OrgAdmin) userUUIDByIDs(ids ...uint) []string {
	searchIDs := []uint{}
	for _, id := range ids {
		if id != 0 {
			searchIDs = append(searchIDs, id)
		}
	}
	users, err := uc.userRepo.GetByIDs(searchIDs)
	if err != nil {
		return nil
	}

	result := []string{}
	for _, user := range users {
		result = append(result, user.UUID.String())
	}

	return result
}

func (uc *OrgAdmin) departmentUUIDByID(id uint) string {
	if id == 0 {
		return ""
	}
	dept, err := uc.deptRepo.GetByID(id)
	if err != nil {
		return ""
	}
	return dept.UUID.String()
}

func (uc *OrgAdmin) departmentUUIDByIDs(ids ...uint) []string {
	searchIDs := []uint{}
	for _, id := range ids {
		if id != 0 {
			searchIDs = append(searchIDs, id)
		}
	}
	depts, err := uc.deptRepo.GetByIDs(searchIDs)
	if err != nil {
		return nil
	}

	result := []string{}
	for _, dept := range depts {
		result = append(result, dept.UUID.String())
	}

	return result
}

func (uc *OrgAdmin) teamUUIDByID(id uint) string {
	if id == 0 {
		return ""
	}
	team, err := uc.teamRepo.GetByID(id)
	if err != nil {
		return ""
	}
	return team.UUID.String()
}

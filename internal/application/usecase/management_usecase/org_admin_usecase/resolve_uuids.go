package org_admin_usecase

import (
	"fmt"

	"github.com/google/uuid"
)

func (uc *OrgAdmin) resolveDepartmentUUIDs(uuids []string) ([]uint, error) {
	if uuids == nil {
		return nil, nil
	}
	out := make([]uint, 0, len(uuids))
	for _, s := range uuids {
		parsed, err := uuid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("invalid department uuid")
		}
		dept, err := uc.deptRepo.GetByUUID(parsed)
		if err != nil {
			return nil, fmt.Errorf("部署が見つかりません")
		}
		out = append(out, dept.ID)
	}
	return out, nil
}

func (uc *OrgAdmin) resolveTeamUUIDs(uuids []string) ([]uint, error) {
	if uuids == nil {
		return nil, nil
	}
	out := make([]uint, 0, len(uuids))
	for _, s := range uuids {
		parsed, err := uuid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("invalid team uuid")
		}
		team, err := uc.teamRepo.GetByUUID(parsed)
		if err != nil {
			return nil, fmt.Errorf("チームが見つかりません")
		}
		out = append(out, team.ID)
	}
	return out, nil
}

func (uc *OrgAdmin) resolveUserUUIDs(uuids []string) ([]uint, error) {
	if uuids == nil {
		return nil, nil
	}
	out := make([]uint, 0, len(uuids))
	for _, s := range uuids {
		parsed, err := uuid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("invalid user uuid")
		}
		user, err := uc.userRepo.GetByUUID(parsed)
		if err != nil {
			return nil, fmt.Errorf("ユーザが見つかりません")
		}
		out = append(out, user.ID)
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

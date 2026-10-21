package org_admin_usecase

import (
	"fmt"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
)

type OrgAdmin struct {
	deptRepo      repository_interface.DepartmentRepository
	teamRepo      repository_interface.TeamRepository
	hierarchyRepo repository_interface.DepartmentHierarchyRepository
	deptRoleRepo  repository_interface.DepartmentRoleRepository
	teamMemberRepo repository_interface.TeamMemberRepository
	userRepo      repository_interface.UserRepository
}

func NewOrgAdmin(
	deptRepo repository_interface.DepartmentRepository,
	teamRepo repository_interface.TeamRepository,
	hierarchyRepo repository_interface.DepartmentHierarchyRepository,
	deptRoleRepo repository_interface.DepartmentRoleRepository,
	teamMemberRepo repository_interface.TeamMemberRepository,
	userRepo repository_interface.UserRepository,
) *OrgAdmin {
	return &OrgAdmin{
		deptRepo:       deptRepo,
		teamRepo:       teamRepo,
		hierarchyRepo:  hierarchyRepo,
		deptRoleRepo:   deptRoleRepo,
		teamMemberRepo: teamMemberRepo,
		userRepo:       userRepo,
	}
}

func (uc *OrgAdmin) ListTree() ([]outputmodel.OrgTreeNodeOutput, error) {
	depts, err := uc.deptRepo.GetAll()
	if err != nil {
		return nil, err
	}
	hierarchies, err := uc.hierarchyRepo.GetAll()
	if err != nil {
		return nil, err
	}
	teams, err := uc.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	members, err := uc.teamMemberRepo.GetAll()
	if err != nil {
		return nil, err
	}

	childMap := map[uint][]uint{}
	childSet := map[uint]bool{}
	for _, h := range hierarchies {
		childMap[h.ParentDepartmentID] = append(childMap[h.ParentDepartmentID], h.ChildDepartmentID)
		childSet[h.ChildDepartmentID] = true
	}
	deptByID := map[uint]entity.Department{}
	for _, d := range depts {
		deptByID[d.ID] = d
	}
	teamsByDept := map[uint][]entity.Team{}
	for _, t := range teams {
		teamsByDept[t.DepartmentID] = append(teamsByDept[t.DepartmentID], t)
	}
	leaderByTeam := map[uint]string{}
	for _, m := range members {
		if m.IsReader && m.Member != nil {
			leaderByTeam[m.TeamID] = m.Member.Name
		}
	}

	var nodes []outputmodel.OrgTreeNodeOutput
	var walk func(deptID uint, depth int, isRoot bool)
	walk = func(deptID uint, depth int, isRoot bool) {
		d := deptByID[deptID]
		nodes = append(nodes, outputmodel.OrgTreeNodeOutput{
			Kind:   "department",
			ID:     d.ID,
			Name:   d.Name,
			Depth:  depth,
			IsRoot: isRoot,
		})
		for _, childID := range childMap[deptID] {
			walk(childID, depth+1, false)
		}
		for _, team := range teamsByDept[deptID] {
			nodes = append(nodes, outputmodel.OrgTreeNodeOutput{
				Kind:         "team",
				ID:           team.ID,
				Name:         team.Name,
				Depth:        depth + 1,
				DepartmentID: deptID,
				LeaderName:   leaderByTeam[team.ID],
			})
		}
	}

	for _, d := range depts {
		if !childSet[d.ID] {
			walk(d.ID, 0, true)
		}
	}
	return nodes, nil
}

func (uc *OrgAdmin) GetDepartmentDetail(id uint) (*outputmodel.OrgDetailOutput, error) {
	dept, err := uc.deptRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	isRoot, childIDs, err := uc.departmentMeta(id)
	if err != nil {
		return nil, err
	}
	users, err := uc.userRepo.GetByDepartmentID(id)
	if err != nil {
		return nil, err
	}
	roles, err := uc.deptRoleRepo.GetByDepartmentID(id)
	if err != nil {
		return nil, err
	}
	roleByUser := map[uint]string{}
	for _, r := range roles {
		roleByUser[r.UserID] = r.Name
	}
	members := make([]outputmodel.OrgMemberOutput, 0, len(users))
	for _, u := range users {
		title := roleByUser[u.ID]
		if title == "" {
			title = "メンバー"
		}
		members = append(members, outputmodel.OrgMemberOutput{
			UserID:   u.ID,
			UserName: u.Name,
			Title:    title,
		})
	}
	teams, _ := uc.teamRepo.GetByDepartmentID(id)
	teamIDs := make([]uint, 0, len(teams))
	for _, t := range teams {
		teamIDs = append(teamIDs, t.ID)
	}
	userIDs := make([]uint, 0, len(users))
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	return &outputmodel.OrgDetailOutput{
		Kind:     "department",
		ID:       dept.ID,
		Name:     dept.Name,
		IsRoot:   isRoot,
		Members:  members,
		ChildIDs: childIDs,
		TeamIDs:  teamIDs,
		UserIDs:  userIDs,
	}, nil
}

func (uc *OrgAdmin) departmentMeta(id uint) (isRoot bool, childIDs []uint, err error) {
	hierarchies, err := uc.hierarchyRepo.GetAll()
	if err != nil {
		return false, nil, err
	}
	isRoot = true
	childIDs = []uint{}
	for _, h := range hierarchies {
		if h.ChildDepartmentID == id {
			isRoot = false
		}
		if h.ParentDepartmentID == id {
			childIDs = append(childIDs, h.ChildDepartmentID)
		}
	}
	return isRoot, childIDs, nil
}

func (uc *OrgAdmin) GetTeamDetail(id uint) (*outputmodel.OrgDetailOutput, error) {
	team, err := uc.teamRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	members, err := uc.teamMemberRepo.GetByTeamID(id)
	if err != nil {
		return nil, err
	}
	outMembers := make([]outputmodel.OrgMemberOutput, 0, len(members))
	var leaderID uint
	memberIDs := make([]uint, 0, len(members))
	for _, m := range members {
		name := ""
		if m.Member != nil {
			name = m.Member.Name
		}
		title := "メンバー"
		if m.IsReader {
			title = "SM"
			leaderID = m.MemberID
		}
		memberIDs = append(memberIDs, m.MemberID)
		outMembers = append(outMembers, outputmodel.OrgMemberOutput{
			UserID:   m.MemberID,
			UserName: name,
			Title:    title,
		})
	}
	return &outputmodel.OrgDetailOutput{
		Kind:         "team",
		ID:           team.ID,
		Name:         team.Name,
		DepartmentID: team.DepartmentID,
		Members:      outMembers,
		MemberIDs:    memberIDs,
		LeaderID:     leaderID,
	}, nil
}

type UpdateDepartmentInput struct {
	Name               string `json:"name"`
	ChildDepartmentIDs []uint `json:"childDepartmentIds"`
	TeamIDs            []uint `json:"teamIds"`
	UserIDs            []uint `json:"userIds"`
}

func (uc *OrgAdmin) UpdateDepartment(id uint, input UpdateDepartmentInput) error {
	dept, err := uc.deptRepo.GetByID(id)
	if err != nil {
		return err
	}
	if input.Name != "" {
		dept.Name = input.Name
		if err := uc.deptRepo.UpdateByID(dept); err != nil {
			return err
		}
	}
	isRoot, _, err := uc.departmentMeta(id)
	if err != nil {
		return err
	}
	if isRoot && input.ChildDepartmentIDs != nil {
		if err := uc.hierarchyRepo.DeleteByParentDepartmentID(id); err != nil {
			return err
		}
		for _, childID := range input.ChildDepartmentIDs {
			if childID == id {
				return fmt.Errorf("自分自身を子部署に指定できません")
			}
			if err := uc.hierarchyRepo.Create(id, childID); err != nil {
				return err
			}
		}
	}
	if !isRoot {
		if input.TeamIDs != nil {
			for _, teamID := range input.TeamIDs {
				team, err := uc.teamRepo.GetByID(teamID)
				if err != nil {
					return err
				}
				team.DepartmentID = id
				if err := uc.teamRepo.UpdateByID(team); err != nil {
					return err
				}
			}
		}
		if input.UserIDs != nil {
			for _, userID := range input.UserIDs {
				user, err := uc.userRepo.GetByID(userID)
				if err != nil {
					return err
				}
				user.DepartmentID = id
				if err := uc.userRepo.UpdateByID(user); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type UpdateTeamInput struct {
	Name      string `json:"name"`
	MemberIDs []uint `json:"memberIds"`
	LeaderID  uint   `json:"leaderId"`
}

func (uc *OrgAdmin) UpdateTeam(id uint, input UpdateTeamInput) error {
	team, err := uc.teamRepo.GetByID(id)
	if err != nil {
		return err
	}
	if input.Name != "" {
		team.Name = input.Name
		if err := uc.teamRepo.UpdateByID(team); err != nil {
			return err
		}
	}
	if input.MemberIDs == nil {
		return nil
	}
	existing, err := uc.teamMemberRepo.GetByTeamID(id)
	if err != nil {
		return err
	}
	for _, m := range existing {
		if err := uc.teamMemberRepo.Delete(m.ID); err != nil {
			return err
		}
	}
	for _, memberID := range input.MemberIDs {
		isReader := memberID == input.LeaderID
		if err := uc.teamMemberRepo.Create(id, memberID, isReader); err != nil {
			return err
		}
	}
	return nil
}

func (uc *OrgAdmin) DeleteDepartment(id uint) error {
	isRoot, childIDs, err := uc.departmentMeta(id)
	if err != nil {
		return err
	}
	if len(childIDs) > 0 {
		return fmt.Errorf("子部署が存在するため削除できません")
	}
	teams, err := uc.teamRepo.GetByDepartmentID(id)
	if err != nil {
		return err
	}
	if len(teams) > 0 {
		return fmt.Errorf("チームが存在するため削除できません")
	}
	users, err := uc.userRepo.GetByDepartmentID(id)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return fmt.Errorf("所属ユーザが存在するため削除できません")
	}
	if isRoot {
		if err := uc.hierarchyRepo.DeleteByParentDepartmentID(id); err != nil {
			return err
		}
	} else {
		parents, err := uc.hierarchyRepo.GetByChildDepartmentID(id)
		if err != nil {
			return err
		}
		for _, p := range parents {
			if err := uc.hierarchyRepo.Delete(p.ID); err != nil {
				return err
			}
		}
	}
	return uc.deptRepo.Delete(id)
}

func (uc *OrgAdmin) DeleteTeam(id uint) error {
	members, err := uc.teamMemberRepo.GetByTeamID(id)
	if err != nil {
		return err
	}
	for _, m := range members {
		if err := uc.teamMemberRepo.Delete(m.ID); err != nil {
			return err
		}
	}
	return uc.teamRepo.Delete(id)
}

func (uc *OrgAdmin) ListDepartmentOptions() ([]outputmodel.DepartmentOptionOutput, error) {
	depts, err := uc.deptRepo.GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]outputmodel.DepartmentOptionOutput, 0, len(depts))
	for _, d := range depts {
		out = append(out, outputmodel.DepartmentOptionOutput{ID: d.ID, Name: d.Name})
	}
	return out, nil
}

func (uc *OrgAdmin) ListTeamOptions() ([]outputmodel.TeamOptionOutput, error) {
	teams, err := uc.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]outputmodel.TeamOptionOutput, 0, len(teams))
	for _, t := range teams {
		out = append(out, outputmodel.TeamOptionOutput{ID: t.ID, Name: t.Name, DepartmentID: t.DepartmentID})
	}
	return out, nil
}

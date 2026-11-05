package org_admin_usecase

import (
	"fmt"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"strings"
)

type OrgAdmin struct {
	deptRepo       repository_interface.DepartmentRepository
	teamRepo       repository_interface.TeamRepository
	hierarchyRepo  repository_interface.DepartmentHierarchyRepository
	deptRoleRepo   repository_interface.DepartmentRoleRepository
	deptMemberRepo repository_interface.DepartmentMemberRepository
	teamMemberRepo repository_interface.TeamMemberRepository
	userRepo       repository_interface.UserRepository
	userDeptRepo   repository_interface.UserDepartmentRepository
}

func NewOrgAdmin(
	deptRepo repository_interface.DepartmentRepository,
	teamRepo repository_interface.TeamRepository,
	hierarchyRepo repository_interface.DepartmentHierarchyRepository,
	deptRoleRepo repository_interface.DepartmentRoleRepository,
	deptMemberRepo repository_interface.DepartmentMemberRepository,
	teamMemberRepo repository_interface.TeamMemberRepository,
	userRepo repository_interface.UserRepository,
	userDeptRepo repository_interface.UserDepartmentRepository,
) *OrgAdmin {
	return &OrgAdmin{
		deptRepo:       deptRepo,
		teamRepo:       teamRepo,
		hierarchyRepo:  hierarchyRepo,
		deptRoleRepo:   deptRoleRepo,
		deptMemberRepo: deptMemberRepo,
		teamMemberRepo: teamMemberRepo,
		userRepo:       userRepo,
		userDeptRepo:   userDeptRepo,
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
	membersByTeam := map[uint][]uint{}
	for _, m := range members {
		membersByTeam[m.TeamID] = append(membersByTeam[m.TeamID], m.MemberID)
		if m.IsReader && m.Member != nil {
			leaderByTeam[m.TeamID] = m.Member.Name
		}
	}
	allUsers, err := uc.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	usersByDept := map[uint][]uint{}
	for _, u := range allUsers {
		usersByDept[u.DepartmentID] = append(usersByDept[u.DepartmentID], u.ID)
	}
	userByID := map[uint]entity.User{}
	for _, u := range allUsers {
		userByID[u.ID] = u
	}
	memberCountByDept := map[uint]int{}
	var subtreeDeptIDs func(deptID uint) []uint
	subtreeDeptIDs = func(deptID uint) []uint {
		ids := []uint{deptID}
		for _, childID := range childMap[deptID] {
			ids = append(ids, subtreeDeptIDs(childID)...)
		}
		return ids
	}
	for _, d := range depts {
		seen := map[uint]bool{}
		for _, subID := range subtreeDeptIDs(d.ID) {
			for _, uid := range usersByDept[subID] {
				seen[uid] = true
			}
			for _, team := range teamsByDept[subID] {
				for _, uid := range membersByTeam[team.ID] {
					seen[uid] = true
				}
			}
		}
		memberCountByDept[d.ID] = len(seen)
	}

	var nodes []outputmodel.OrgTreeNodeOutput
	var walk func(deptID uint, depth int, isRoot bool)
	walk = func(deptID uint, depth int, isRoot bool) {
		d := deptByID[deptID]
		headName := ""
		if d.HeadUserID != 0 {
			if u, ok := userByID[d.HeadUserID]; ok {
				headName = u.Name
			}
		}
		nodes = append(nodes, outputmodel.OrgTreeNodeOutput{
			Kind:            "department",
			ID:              d.ID,
			Name:            d.Name,
			Depth:           depth,
			IsRoot:          isRoot,
			MemberCount:     memberCountByDept[d.ID],
			HeadUserName:    headName,
			DepartmentID:    deptID,
			DepartmentCount: len(childMap[d.ID]),
			TeamCount:       len(teamsByDept[d.ID]),
		})
		if teams, exists := teamsByDept[d.ID]; exists {
			for _, team := range teams {
				nodes = append(nodes, outputmodel.OrgTreeNodeOutput{
					Kind:            "team",
					ID:              team.ID,
					Name:            team.Name,
					Depth:           depth + 1,
					DepartmentID:    deptID,
					LeaderName:      leaderByTeam[team.ID],
					TeamMemberCount: len(membersByTeam[team.ID]),
				})
			}
		}

		for _, childID := range childMap[deptID] {
			walk(childID, depth+1, false)
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
	deptMembers, err := uc.deptMemberRepo.GetByDepartmentID(id)
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
	members := make([]outputmodel.OrgMemberOutput, 0, len(deptMembers))
	for _, row := range deptMembers {
		name := ""
		if row.User != nil {
			name = row.User.Name
		}
		title := roleByUser[row.UserID]
		if title == "" {
			title = "メンバー"
		}
		members = append(members, outputmodel.OrgMemberOutput{
			UserID:   row.UserID,
			UserName: name,
			Title:    title,
		})
	}
	teams, _ := uc.teamRepo.GetByDepartmentID(id)
	teamIDs := make([]uint, 0, len(teams))
	for _, t := range teams {
		teamIDs = append(teamIDs, t.ID)
	}
	userIDs := make([]uint, 0, len(deptMembers))
	for _, row := range deptMembers {
		userIDs = append(userIDs, row.UserID)
	}
	return &outputmodel.OrgDetailOutput{
		Kind:       "department",
		ID:         dept.ID,
		Name:       dept.Name,
		IsRoot:     isRoot,
		Members:    members,
		ChildIDs:   childIDs,
		TeamIDs:    teamIDs,
		UserIDs:    userIDs,
		HeadUserID: dept.HeadUserID,
	}, nil
}

func departmentIsAncestor(ancestorID, deptID uint, parentOf map[uint]uint) bool {
	seen := map[uint]bool{}
	for cur := deptID; cur != 0; {
		if seen[cur] {
			break
		}
		seen[cur] = true
		parentID, ok := parentOf[cur]
		if !ok {
			break
		}
		if parentID == ancestorID {
			return true
		}
		cur = parentID
	}
	return false
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
	HeadUserID         uint   `json:"headUserId"`
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
	if input.ChildDepartmentIDs != nil {
		hierarchies, err := uc.hierarchyRepo.GetAll()
		if err != nil {
			return err
		}
		childSet := map[uint]bool{}
		parentOf := map[uint]uint{}
		for _, h := range hierarchies {
			childSet[h.ChildDepartmentID] = true
			parentOf[h.ChildDepartmentID] = h.ParentDepartmentID
		}
		want := map[uint]bool{}
		for _, childID := range input.ChildDepartmentIDs {
			if childID == id {
				return fmt.Errorf("自分自身を子部署に指定できません")
			}
			if !childSet[childID] {
				return fmt.Errorf("ルート部署を子部署に指定できません")
			}
			if departmentIsAncestor(childID, id, parentOf) {
				return fmt.Errorf("循環する子部署は指定できません")
			}
			want[childID] = true
		}
		existing, err := uc.hierarchyRepo.GetByParentDepartmentID(id)
		if err != nil {
			return err
		}
		have := map[uint]uint{}
		for _, row := range existing {
			have[row.ChildDepartmentID] = row.ID
		}
		for childID, hierarchyID := range have {
			if !want[childID] {
				if err := uc.hierarchyRepo.Delete(hierarchyID); err != nil {
					return err
				}
			}
		}
		for childID := range want {
			if have[childID] != 0 {
				continue
			}
			parents, err := uc.hierarchyRepo.GetByChildDepartmentID(childID)
			if err != nil {
				return err
			}
			for _, p := range parents {
				if err := uc.hierarchyRepo.Delete(p.ID); err != nil {
					return err
				}
			}
			if err := uc.hierarchyRepo.Create(id, childID); err != nil {
				return err
			}
		}
	}
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
		if err := uc.userDeptRepo.SyncDepartmentMembers(id, input.UserIDs); err != nil {
			return err
		}
		return uc.syncDepartmentMemberFlags(id, input.UserIDs, input.HeadUserID)
	}
	return uc.applyDepartmentHeadOnly(id, input.HeadUserID)
}

func (uc *OrgAdmin) applyDepartmentHeadOnly(departmentID uint, headUserID uint) error {
	existing, err := uc.deptMemberRepo.GetByDepartmentID(departmentID)
	if err != nil {
		return err
	}
	prevManager := map[uint]bool{}
	prevHead := map[uint]bool{}
	userIDsInDept := map[uint]struct{}{}
	for _, member := range existing {
		prevManager[member.UserID] = member.IsManager
		prevHead[member.UserID] = member.IsHead
		userIDsInDept[member.UserID] = struct{}{}
	}
	if headUserID != 0 {
		userIDsInDept[headUserID] = struct{}{}
	}
	rows := make([]entity.DepartmentMember, 0, len(userIDsInDept))
	for userID := range userIDsInDept {
		isHead := headUserID != 0 && userID == headUserID
		keptManager := prevManager[userID] && !prevHead[userID]
		rows = append(rows, entity.DepartmentMember{
			DepartmentID: departmentID,
			UserID:       userID,
			IsHead:       isHead,
			IsManager:    isHead || keptManager,
		})
	}
	if err := uc.deptMemberRepo.ReplaceByDepartment(departmentID, rows); err != nil {
		return err
	}
	return uc.deptRepo.UpdateHeadUserID(departmentID, headUserID)
}

func (uc *OrgAdmin) syncDepartmentMemberFlags(departmentID uint, userIDs []uint, headUserID uint) error {
	existing, err := uc.deptMemberRepo.GetByDepartmentID(departmentID)
	if err != nil {
		return err
	}
	prevManager := map[uint]bool{}
	prevHead := map[uint]bool{}
	for _, member := range existing {
		prevManager[member.UserID] = member.IsManager
		prevHead[member.UserID] = member.IsHead
	}
	want := map[uint]struct{}{}
	for _, userID := range userIDs {
		if userID != 0 {
			want[userID] = struct{}{}
		}
	}
	if headUserID != 0 {
		want[headUserID] = struct{}{}
	}
	rows := make([]entity.DepartmentMember, 0, len(want))
	for userID := range want {
		isHead := headUserID != 0 && userID == headUserID
		keptManager := prevManager[userID] && !prevHead[userID]
		rows = append(rows, entity.DepartmentMember{
			DepartmentID: departmentID,
			UserID:       userID,
			IsHead:       isHead,
			IsManager:    isHead || keptManager,
		})
	}
	if err := uc.deptMemberRepo.ReplaceByDepartment(departmentID, rows); err != nil {
		return err
	}
	return uc.deptRepo.UpdateHeadUserID(departmentID, headUserID)
}

type CreateTeamInput struct {
	Name         string `json:"name"`
	DepartmentID uint   `json:"departmentId"`
}

func (uc *OrgAdmin) CreateTeam(input CreateTeamInput) (*outputmodel.TeamOptionOutput, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("チーム名を入力してください")
	}
	if input.DepartmentID == 0 {
		return nil, fmt.Errorf("部署を指定してください")
	}
	team := &entity.Team{Name: name, DepartmentID: input.DepartmentID}
	if err := uc.teamRepo.Create(team); err != nil {
		return nil, err
	}
	return &outputmodel.TeamOptionOutput{ID: team.ID, Name: team.Name, DepartmentID: team.DepartmentID}, nil
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
	if input.MemberIDs != nil {
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
	return uc.applyTeamLeaderOnly(id, input.LeaderID)
}

func (uc *OrgAdmin) applyTeamLeaderOnly(teamID uint, leaderID uint) error {
	existing, err := uc.teamMemberRepo.GetByTeamID(teamID)
	if err != nil {
		return err
	}
	foundLeader := false
	for i := range existing {
		m := &existing[i]
		isReader := leaderID != 0 && m.MemberID == leaderID
		if isReader {
			foundLeader = true
		}
		if m.IsReader != isReader {
			m.IsReader = isReader
			if err := uc.teamMemberRepo.Update(m.ID, m); err != nil {
				return err
			}
		}
	}
	if leaderID != 0 && !foundLeader {
		return uc.teamMemberRepo.Create(teamID, leaderID, true)
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
	if err := uc.deptMemberRepo.ReplaceByDepartment(id, nil); err != nil {
		return err
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

type CreateDepartmentInput struct {
	Name               string `json:"name"`
	ParentDepartmentID uint   `json:"parentDepartmentId"`
}

func (uc *OrgAdmin) CreateDepartment(input CreateDepartmentInput) (*outputmodel.DepartmentOptionOutput, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("部署名を入力してください")
	}
	dept := &entity.Department{Name: name}
	if err := uc.deptRepo.Create(dept); err != nil {
		return nil, err
	}
	if input.ParentDepartmentID > 0 {
		if _, err := uc.deptRepo.GetByID(input.ParentDepartmentID); err != nil {
			return nil, fmt.Errorf("親部署が見つかりません")
		}
		if input.ParentDepartmentID == dept.ID {
			return nil, fmt.Errorf("自分自身を親部署に指定できません")
		}
		if err := uc.hierarchyRepo.Create(input.ParentDepartmentID, dept.ID); err != nil {
			return nil, err
		}
	}
	return &outputmodel.DepartmentOptionOutput{ID: dept.ID, Name: dept.Name}, nil
}

func (uc *OrgAdmin) ListDepartmentOptions() ([]outputmodel.DepartmentOptionOutput, error) {
	depts, err := uc.deptRepo.GetAll()
	if err != nil {
		return nil, err
	}
	hierarchies, err := uc.hierarchyRepo.GetAll()
	if err != nil {
		return nil, err
	}
	childSet := map[uint]bool{}
	for _, h := range hierarchies {
		childSet[h.ChildDepartmentID] = true
	}
	out := make([]outputmodel.DepartmentOptionOutput, 0, len(depts))
	for _, d := range depts {
		out = append(out, outputmodel.DepartmentOptionOutput{
			ID:     d.ID,
			Name:   d.Name,
			IsRoot: !childSet[d.ID],
		})
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

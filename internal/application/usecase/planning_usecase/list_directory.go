package planning_usecase

import (
	"fmt"
	"strings"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

type PlanningDirectoryUser struct {
	UUID            string   `json:"uuid"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	ImageURI        *string  `json:"imageUri,omitempty"`
	RoleName        string   `json:"roleName"`
	DepartmentNames string   `json:"departmentNames"`
	DepartmentUuids []string `json:"departmentUuids,omitempty"`
	TeamNames       string   `json:"teamNames"`
}

type PlanningDirectoryTeam struct {
	UUID           string                   `json:"uuid"`
	Name           string                   `json:"name"`
	DepartmentUuid string                   `json:"departmentUuid"`
	Members        []outputmodel.UserOutput `json:"members"`
}

type PlanningDirectoryResponse struct {
	Users        []PlanningDirectoryUser        `json:"users"`
	Teams        []PlanningDirectoryTeam        `json:"teams"`
	Departments  []outputmodel.DepartmentOutput `json:"departments"`
	ManagedTeams []PlanningDirectoryTeam        `json:"managedTeams"`
}

func stringsJoinUnique(items []string) string {
	seen := map[string]bool{}
	out := make([]string, 0, len(items))
	for _, s := range items {
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return strings.Join(out, "・")
}

func collectDepartmentIDsWithoutPreload(teams []entity.Team) []uint {
	seen := map[uint]struct{}{}
	ids := make([]uint, 0)
	for _, team := range teams {
		if team.Department != nil {
			continue
		}
		if team.DepartmentID == 0 {
			continue
		}
		if _, ok := seen[team.DepartmentID]; ok {
			continue
		}
		seen[team.DepartmentID] = struct{}{}
		ids = append(ids, team.DepartmentID)
	}
	return ids
}

func (uc *PlanningAccessUsecase) departmentUUIDByID(deptIDs []uint) (map[uint]string, error) {
	if len(deptIDs) == 0 {
		return map[uint]string{}, nil
	}
	depts, err := uc.deptRepo.GetByIDs(deptIDs)
	if err != nil {
		return nil, err
	}
	out := make(map[uint]string, len(depts))
	for _, d := range depts {
		out[d.ID] = d.UUID.String()
	}
	return out, nil
}

func departmentUUIDForTeam(team entity.Team, deptUUIDByID map[uint]string) string {
	if team.Department != nil {
		return team.Department.UUID.String()
	}
	if team.DepartmentID == 0 {
		return ""
	}
	return deptUUIDByID[team.DepartmentID]
}

func (uc *PlanningAccessUsecase) membersOutputsByTeamIDs(teamIDs []uint) (map[uint][]outputmodel.UserOutput, error) {
	if len(teamIDs) == 0 {
		return map[uint][]outputmodel.UserOutput{}, nil
	}
	membersByTeam, err := uc.teamMemberRepo.GetByTeamIDs(teamIDs)
	if err != nil {
		return nil, err
	}
	seenMember := map[uint]struct{}{}
	memberIDs := make([]uint, 0)
	for _, members := range membersByTeam {
		for _, m := range members {
			if m.MemberID == 0 {
				continue
			}
			if _, ok := seenMember[m.MemberID]; ok {
				continue
			}
			seenMember[m.MemberID] = struct{}{}
			memberIDs = append(memberIDs, m.MemberID)
		}
	}
	users, err := uc.userRepo.GetActiveUsersByIDs(memberIDs)
	if err != nil {
		return nil, fmt.Errorf("アクティブユーザの取得に失敗しました: %w", err)
	}
	userByID := make(map[uint]entity.User, len(users))
	for _, u := range users {
		userByID[u.ID] = u
	}
	out := make(map[uint][]outputmodel.UserOutput, len(membersByTeam))
	for teamID, members := range membersByTeam {
		memberOuts := make([]outputmodel.UserOutput, 0, len(members))
		for _, m := range members {
			if u, ok := userByID[m.MemberID]; ok {
				memberOuts = append(memberOuts, converter.UserEntityToOutput(u))
			}
		}
		out[teamID] = memberOuts
	}
	return out, nil
}

func (uc *PlanningAccessUsecase) ListDirectory(viewerUserID uint, includeManaged bool) (*PlanningDirectoryResponse, error) {
	users, err := uc.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	userOuts := make([]PlanningDirectoryUser, 0, len(users))
	for _, u := range users {
		if u.Suspended {
			continue
		}
		roleName := ""
		deptNames := make([]string, 0, len(u.DepartmentMembers))
		deptUuids := make([]string, 0, len(u.DepartmentMembers))
		for _, link := range u.DepartmentMembers {
			if link.Department != nil && link.Department.Name != "" {
				deptNames = append(deptNames, link.Department.Name)
				deptUuids = append(deptUuids, link.Department.UUID.String())
			}
			if roleName == "" && link.Role != nil {
				roleName = link.Role.Name
			}
		}
		teamNames := make([]string, 0, len(u.TeamMembers))
		for _, tm := range u.TeamMembers {
			if tm.Team != nil && tm.Team.Name != "" {
				teamNames = append(teamNames, tm.Team.Name)
			}
		}
		userOuts = append(userOuts, PlanningDirectoryUser{
			UUID:            u.UUID.String(),
			Name:            u.Name,
			Email:           u.Email,
			ImageURI:        u.ImageURI,
			RoleName:        roleName,
			DepartmentNames: stringsJoinUnique(deptNames),
			DepartmentUuids: deptUuids,
			TeamNames:       stringsJoinUnique(teamNames),
		})
	}
	teams, err := uc.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	teamIDs := make([]uint, len(teams))
	for idx, team := range teams {
		teamIDs[idx] = team.ID
	}
	membersByTeamID, err := uc.membersOutputsByTeamIDs(teamIDs)
	if err != nil {
		return nil, err
	}
	deptUUIDByID, err := uc.departmentUUIDByID(collectDepartmentIDsWithoutPreload(teams))
	if err != nil {
		return nil, err
	}
	teamOuts := make([]PlanningDirectoryTeam, len(teams))
	for idx, team := range teams {
		teamOuts[idx] = PlanningDirectoryTeam{
			UUID:           team.UUID.String(),
			Name:           team.Name,
			DepartmentUuid: departmentUUIDForTeam(team, deptUUIDByID),
			Members:        membersByTeamID[team.ID],
		}
	}
	depts, err := uc.deptRepo.GetAll()
	if err != nil {
		return nil, err
	}
	managedTeams := []PlanningDirectoryTeam{}
	if includeManaged && viewerUserID > 0 {
		leaderTeamIDs, err := uc.LeaderTeamIDs(viewerUserID)
		if err != nil {
			return nil, err
		}
		leaderTeams, err := uc.teamRepo.GetByIDs(leaderTeamIDs)
		if err != nil {
			return nil, err
		}
		for _, team := range leaderTeams {
			managedTeams = append(managedTeams, PlanningDirectoryTeam{
				UUID:           team.UUID.String(),
				Name:           team.Name,
				DepartmentUuid: departmentUUIDForTeam(team, deptUUIDByID),
				Members:        membersByTeamID[team.ID],
			})
		}
	}
	return &PlanningDirectoryResponse{
		Users:        userOuts,
		Teams:        teamOuts,
		Departments:  converter.DepartmentEntitiesToOutputs(depts),
		ManagedTeams: managedTeams,
	}, nil
}
package api

import (
	"fmt"
	"net/http"
	"strings"

	"growlink/internal/adapter/handler/authctx"
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"

	"github.com/labstack/echo/v5"
)

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

type planningDirectoryUser struct {
	UUID            string   `json:"uuid"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	ImageURI        *string  `json:"imageUri,omitempty"`
	RoleName        string   `json:"roleName"`
	DepartmentNames string   `json:"departmentNames"`
	DepartmentUuids []string `json:"departmentUuids,omitempty"`
	TeamNames       string   `json:"teamNames"`
}

type planningDirectoryResponse struct {
	Users        []planningDirectoryUser        `json:"users"`
	Teams        []planningTeamOutput           `json:"teams"`
	Departments  []outputmodel.DepartmentOutput `json:"departments"`
	ManagedTeams []planningTeamOutput           `json:"managedTeams"`
}

type planningTeamOutput struct {
	UUID           string                   `json:"uuid"`
	Name           string                   `json:"name"`
	DepartmentUuid string                   `json:"departmentUuid"`
	Members        []outputmodel.UserOutput `json:"members"`
}

func teamMembersToOutputs(h *Handler, teamID uint) ([]outputmodel.UserOutput, error) {
	members, err := h.deps.TeamMemberRepo.GetByTeamID(teamID)
	if err != nil {
		return nil, err
	}
	memberIDs := []uint{}
	for idx, m := range members {
		memberIDs[idx] = m.MemberID
	}

	users, err := h.deps.UserRepo.GetActiveUsersByIDs(memberIDs)
	if err != nil {
		return nil, fmt.Errorf("アクティブユーザの取得に失敗しました: %w", err)
	}
	memberOuts := converter.UserEntitiesToOutputs(users)
	return memberOuts, nil
}

func (h *Handler) ListPlanningDirectory(c *echo.Context) error {
	claims := authctx.Claims(c)

	users, err := h.deps.UserRepo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	userOuts := make([]planningDirectoryUser, 0, len(users))
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
		deptName := stringsJoinUnique(deptNames)
		teamNames := make([]string, 0, len(u.TeamMembers))
		for _, tm := range u.TeamMembers {
			if tm.Team != nil && tm.Team.Name != "" {
				teamNames = append(teamNames, tm.Team.Name)
			}
		}
		userOuts = append(userOuts, planningDirectoryUser{
			UUID:            u.UUID.String(),
			Name:            u.Name,
			Email:           u.Email,
			ImageURI:        u.ImageURI,
			RoleName:        roleName,
			DepartmentNames: deptName,
			DepartmentUuids: deptUuids,
			TeamNames:       stringsJoinUnique(teamNames),
		})
	}
	teams, err := h.deps.TeamRepo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	teamOuts := make([]planningTeamOutput, 0, len(teams))
	for _, team := range teams {
		memberOuts, err := teamMembersToOutputs(h, team.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		deptUUID := ""
		if team.Department != nil {
			deptUUID = team.Department.UUID.String()
		} else if team.DepartmentID > 0 {
			dept, err := h.deps.DepartmentRepo.GetByID(team.DepartmentID)
			if err == nil {
				deptUUID = dept.UUID.String()
			}
		}
		teamOuts = append(teamOuts, planningTeamOutput{
			UUID:           team.UUID.String(),
			Name:           team.Name,
			DepartmentUuid: deptUUID,
			Members:        memberOuts,
		})
	}
	depts, err := h.deps.DepartmentRepo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	managedTeams := []planningTeamOutput{}
	if claims != nil && h.deps.Planning != nil {
		leaderTeamIDs, err := h.deps.Planning.LeaderTeamIDs(claims.UserID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		teams, err := h.deps.TeamRepo.GetByIDs(leaderTeamIDs)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		for _, team := range teams {
			memberOuts, err := teamMembersToOutputs(h, team.ID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
			deptUUID := ""
			if team.Department != nil {
				deptUUID = team.Department.UUID.String()
			} else if team.DepartmentID > 0 {
				dept, err := h.deps.DepartmentRepo.GetByID(team.DepartmentID)
				if err == nil {
					deptUUID = dept.UUID.String()
				}
			}
			managedTeams = append(managedTeams, planningTeamOutput{
				UUID:           team.UUID.String(),
				Name:           team.Name,
				DepartmentUuid: deptUUID,
				Members:        memberOuts,
			})
		}
	}

	return c.JSON(http.StatusOK, planningDirectoryResponse{
		Users:        userOuts,
		Teams:        teamOuts,
		Departments:  converter.DepartmentEntitiesToOutputs(depts),
		ManagedTeams: managedTeams,
	})
}

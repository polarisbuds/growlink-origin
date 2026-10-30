package api

import (
	"net/http"
	"strings"

	"growlink/internal/adapter/handler/authctx"
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/valueobject"

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
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	ImageURI        *string `json:"imageUri,omitempty"`
	RoleName        string  `json:"roleName"`
	DepartmentNames string  `json:"departmentNames"`
	TeamNames       string  `json:"teamNames"`
	DepartmentID    uint    `json:"departmentId"`
	RoleID          uint    `json:"roleId"`
}

type planningDirectoryResponse struct {
	Users        []planningDirectoryUser        `json:"users"`
	Teams        []planningTeamOutput           `json:"teams"`
	Departments  []outputmodel.DepartmentOutput `json:"departments"`
	ManagedTeams []planningTeamOutput           `json:"managedTeams"`
}

type planningTeamOutput struct {
	ID           uint                     `json:"id"`
	Name         string                   `json:"name"`
	DepartmentID uint                     `json:"departmentId"`
	Members      []outputmodel.UserOutput `json:"members"`
}

func teamMembersToOutputs(h *Handler, teamID uint) ([]outputmodel.UserOutput, error) {
	members, err := h.deps.TeamMemberRepo.GetByTeamID(teamID)
	if err != nil {
		return nil, err
	}
	memberOuts := make([]outputmodel.UserOutput, 0, len(members))
	for _, m := range members {
		user, err := h.deps.UserRepo.GetByID(m.MemberID)
		if err != nil || user.Suspended {
			continue
		}
		memberOuts = append(memberOuts, outputmodel.UserOutput{ID: user.ID, Name: user.Name, Email: user.Email, ImageURI: user.ImageURI})
	}
	return memberOuts, nil
}

func (h *Handler) ListPlanningDirectory(c *echo.Context) error {
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
		if u.Role != nil {
			roleName = u.Role.Name
		}
		deptName := ""
		if u.Department != nil {
			deptName = u.Department.Name
		}
		teamNames := make([]string, 0, len(u.TeamMembers))
		for _, tm := range u.TeamMembers {
			if tm.Team != nil && tm.Team.Name != "" {
				teamNames = append(teamNames, tm.Team.Name)
			}
		}
		userOuts = append(userOuts, planningDirectoryUser{
			ID:              u.ID,
			Name:            u.Name,
			Email:           u.Email,
			ImageURI:        u.ImageURI,
			RoleName:        roleName,
			DepartmentNames: deptName,
			TeamNames:       stringsJoinUnique(teamNames),
			DepartmentID:    u.DepartmentID,
			RoleID:          u.RoleID,
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
		teamOuts = append(teamOuts, planningTeamOutput{
			ID:           team.ID,
			Name:         team.Name,
			DepartmentID: team.DepartmentID,
			Members:      memberOuts,
		})
	}
	depts, err := h.deps.DepartmentRepo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	managedTeams := []planningTeamOutput{}
	claims := authctx.Claims(c)
	if claims != nil && h.deps.TeamRefRepo != nil {
		user, err := h.deps.UserRepo.GetByID(claims.UserID)
		if err == nil && user.OperationTarget == valueobject.OperationTargetTeam {
			refs, err := h.deps.TeamRefRepo.GetByUserID(claims.UserID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
			for _, ref := range refs {
				team, err := h.deps.TeamRepo.GetByID(ref.ReferableTeamID)
				if err != nil {
					continue
				}
				memberOuts, err := teamMembersToOutputs(h, team.ID)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
				}
				managedTeams = append(managedTeams, planningTeamOutput{
					ID:           team.ID,
					Name:         team.Name,
					DepartmentID: team.DepartmentID,
					Members:      memberOuts,
				})
			}
		}
	}

	return c.JSON(http.StatusOK, planningDirectoryResponse{
		Users:        userOuts,
		Teams:        teamOuts,
		Departments:  converter.DepartmentEntitiesToOutputs(depts),
		ManagedTeams: managedTeams,
	})
}

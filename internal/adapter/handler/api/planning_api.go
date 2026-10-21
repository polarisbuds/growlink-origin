package api

import (
	"net/http"

	outputmodel "growlink/internal/application/model/output_model"

	"github.com/labstack/echo/v5"
)

type planningDirectoryResponse struct {
	Users []outputmodel.UserOutput `json:"users"`
	Teams []planningTeamOutput     `json:"teams"`
}

type planningTeamOutput struct {
	ID      uint                     `json:"id"`
	Name    string                   `json:"name"`
	Members []outputmodel.UserOutput `json:"members"`
}

func (h *Handler) ListPlanningDirectory(c *echo.Context) error {
	users, err := h.deps.UserRepo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	userOuts := make([]outputmodel.UserOutput, 0, len(users))
	for _, u := range users {
		if u.Suspended {
			continue
		}
		userOuts = append(userOuts, outputmodel.UserOutput{ID: u.ID, Name: u.Name, Email: u.Email, ImageURI: u.ImageURI})
	}
	teams, err := h.deps.TeamRepo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	teamOuts := make([]planningTeamOutput, 0, len(teams))
	for _, team := range teams {
		members, err := h.deps.TeamMemberRepo.GetByTeamID(team.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		memberOuts := make([]outputmodel.UserOutput, 0, len(members))
		for _, m := range members {
			user, err := h.deps.UserRepo.GetByID(m.MemberID)
			if err != nil || user.Suspended {
				continue
			}
			memberOuts = append(memberOuts, outputmodel.UserOutput{ID: user.ID, Name: user.Name, Email: user.Email, ImageURI: user.ImageURI})
		}
		teamOuts = append(teamOuts, planningTeamOutput{ID: team.ID, Name: team.Name, Members: memberOuts})
	}
	return c.JSON(http.StatusOK, planningDirectoryResponse{Users: userOuts, Teams: teamOuts})
}

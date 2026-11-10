package planning_usecase

import (
	"growlink/internal/domain/valueobject"
	"log/slog"
)

func fullPlanningAccess(systemRole valueobject.SystemRole) bool {
	return systemRole == valueobject.SystemRoleAdmin || systemRole == valueobject.SystemRoleAnalyst
}

// ReferableUserIDs returns user IDs whose meetings may count toward progress and org-scoped 1:1 access.
// nil map means unrestricted (system admin / analyst).
func (uc *PlanningAccessUsecase) ReferableUserIDs(userID uint, systemRole valueobject.SystemRole) (map[uint]bool, error) {
	if fullPlanningAccess(systemRole) {
		return nil, nil
	}
	visible := map[uint]bool{}

	deptUserIDs, err := uc.departmentScopeUserIDs(userID)
	if err != nil {
		return nil, err
	}
	slog.Debug("ReferableUserIDs-1", slog.Any("deptUserIDs", deptUserIDs))
	for _, id := range deptUserIDs {
		if id > 0 {
			visible[id] = true
		}
	}
	slog.Debug("ReferableUserIDs-2", slog.Any("visible", visible))

	teamIDs, err := uc.readerTeamIDs(userID)
	if err != nil {
		return nil, err
	}
	for _, teamID := range teamIDs {
		members, err := uc.teamMemberRepo.GetByTeamID(teamID)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			if member.MemberID > 0 {
				visible[member.MemberID] = true
			}
		}
	}

	pairings, err := uc.pairingRepo.GetAll()
	if err != nil {
		return nil, err
	}
	for _, p := range pairings {
		if p.MentorID == userID && p.MenteeID > 0 {
			visible[p.MenteeID] = true
		}
	}
	return visible, nil
}

// LeaderTeamIDs returns teams where the user is team leader (IsReader).
func (uc *PlanningAccessUsecase) LeaderTeamIDs(userID uint) ([]uint, error) {
	return uc.readerTeamIDs(userID)
}

package planning_usecase

import (
	"context"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"
)

func fullPlanningAccess(systemRole valueobject.SystemRole) bool {
	return systemRole == valueobject.SystemRoleAdmin
}

// ReferableUserIDs returns user IDs whose meetings may count toward progress and org-scoped 1:1 access.
// nil map means unrestricted (system admin / analyst).
func (uc *PlanningAccessUsecase) ReferableUserIDs(userID uint, systemRole valueobject.SystemRole) (ret0 map[uint]bool, err error) {
	const op = "planning_usecase.ReferableUserIDs"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if fullPlanningAccess(systemRole) {
		return nil, nil
	}
	visible := map[uint]bool{}

	deptUserIDs, err := uc.departmentScopeUserIDs(userID)
	if err != nil {
		return nil, err
	}

	for _, id := range deptUserIDs {
		if id > 0 {
			visible[id] = true
		}
	}

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
func (uc *PlanningAccessUsecase) LeaderTeamIDs(userID uint) (ret0 []uint, err error) {
	const op = "planning_usecase.LeaderTeamIDs"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	return uc.readerTeamIDs(userID)
}

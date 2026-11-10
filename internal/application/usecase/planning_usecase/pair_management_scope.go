package planning_usecase

import (
	"growlink/internal/domain/valueobject"
)

// PairManagementScopeUserIDs returns user IDs whose pairings may be viewed or managed
// on the planning pair list (dept head/manager scope and team leader scope only).
func (uc *PlanningAccessUsecase) PairManagementScopeUserIDs(userID uint, systemRole valueobject.SystemRole) (map[uint]bool, error) {
	if systemRole == valueobject.SystemRoleAdmin {
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
	return visible, nil
}

func (uc *PlanningAccessUsecase) CanManagePairingMentee(userID uint, systemRole valueobject.SystemRole, menteeID uint) (bool, error) {
	if systemRole == valueobject.SystemRoleAdmin {
		return true, nil
	}
	scope, err := uc.PairManagementScopeUserIDs(userID, systemRole)
	if err != nil {
		return false, err
	}
	return scope[menteeID], nil
}

func (uc *PlanningAccessUsecase) CanBreakPairing(userID uint, systemRole valueobject.SystemRole, createdByUserID, menteeID uint) (bool, error) {
	if systemRole == valueobject.SystemRoleAdmin {
		return true, nil
	}
	if userID == createdByUserID && createdByUserID > 0 {
		return true, nil
	}
	return uc.CanManagePairingMentee(userID, systemRole, menteeID)
}

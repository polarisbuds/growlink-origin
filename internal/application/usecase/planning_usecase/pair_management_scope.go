package planning_usecase

import (
	"growlink/internal/domain/valueobject"
)

// PairManagementScopeUserIDs returns user IDs whose pairings may be viewed or managed
// (boss dept scope, team leader scope, and mentor mentees for general users).
func (uc *PlanningAccessUsecase) PairManagementScopeUserIDs(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole) (map[uint]bool, error) {
	if systemRole == valueobject.SystemRoleAdmin || authority.IsOwner() {
		return nil, nil
	}
	visible := map[uint]bool{}

	if authority.IsBoss() {
		deptUserIDs, err := uc.departmentScopeUserIDs(userID)
		if err != nil {
			return nil, err
		}
		for _, id := range deptUserIDs {
			if id > 0 {
				visible[id] = true
			}
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

	if !authority.IsBoss() && uc.pairingRepo != nil {
		pairings, err := uc.pairingRepo.GetByMentorID(userID)
		if err != nil {
			return nil, err
		}
		for _, p := range pairings {
			if p.MenteeID > 0 {
				visible[p.MenteeID] = true
			}
		}
	}
	return visible, nil
}

func (uc *PlanningAccessUsecase) CanManagePairingMentee(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole, menteeID uint) (bool, error) {
	if systemRole == valueobject.SystemRoleAdmin || authority.IsOwner() {
		return true, nil
	}
	scope, err := uc.PairManagementScopeUserIDs(userID, authority, systemRole)
	if err != nil {
		return false, err
	}
	return scope[menteeID], nil
}

func (uc *PlanningAccessUsecase) CanBreakPairing(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole, createdByUserID, menteeID uint) (bool, error) {
	if systemRole == valueobject.SystemRoleAdmin || authority.IsOwner() {
		return true, nil
	}
	if userID == createdByUserID && createdByUserID > 0 {
		return true, nil
	}
	return uc.CanManagePairingMentee(userID, authority, systemRole, menteeID)
}

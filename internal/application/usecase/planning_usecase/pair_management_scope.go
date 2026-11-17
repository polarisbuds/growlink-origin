package planning_usecase

import (
	"context"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"
)

// PairManagementScopeUserIDs returns user IDs whose pairings may be viewed or managed
// (boss dept scope, team leader scope, and mentor mentees for general users).
func (uc *PlanningAccessUsecase) PairManagementScopeUserIDs(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole) (ret0 map[uint]bool, err error) {
	const op = "planning_usecase.PairManagementScopeUserIDs"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
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

	if uc.pairingRepo != nil {
		mentorPairings, err := uc.pairingRepo.GetByMentorID(userID)
		if err != nil {
			return nil, err
		}
		for _, p := range mentorPairings {
			if p.MenteeID > 0 {
				visible[p.MenteeID] = true
			}
		}

		pairingsByMentee, err := uc.pairingRepo.GetByMenteeID(userID)
		if err != nil {
			return nil, err
		}
		if len(pairingsByMentee) > 0 {
			visible[userID] = true
		}
	}
	return visible, nil
}

func (uc *PlanningAccessUsecase) CanManagePairingMentee(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole, menteeID uint) (ret0 bool, err error) {
	const op = "planning_usecase.CanManagePairingMentee"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if systemRole == valueobject.SystemRoleAdmin || authority.IsOwner() {
		return true, nil
	}
	scope, err := uc.PairManagementScopeUserIDs(userID, authority, systemRole)
	if err != nil {
		return false, err
	}
	return scope[menteeID], nil
}

func (uc *PlanningAccessUsecase) CanBreakPairing(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole, createdByUserID, menteeID uint) (ret0 bool, err error) {
	const op = "planning_usecase.CanBreakPairing"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if systemRole == valueobject.SystemRoleAdmin || authority.IsOwner() {
		return true, nil
	}
	if userID == createdByUserID && createdByUserID > 0 {
		return true, nil
	}
	return uc.CanManagePairingMentee(userID, authority, systemRole, menteeID)
}

package planning_usecase

import (
	"context"
	"fmt"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"
	"maps"

	"golang.org/x/sync/errgroup"
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

	group, _ := errgroup.WithContext(context.Background())
	deptVisibles := map[uint]bool{}
	group.Go(func() error {
		deptUserIDs, err := uc.departmentScopeUserIDs(userID)
		if err != nil {
			return err
		}

		for _, id := range deptUserIDs {
			if id > 0 {
				deptVisibles[id] = true
			}
		}
		return nil
	})

	teamVisibles := map[uint]bool{}
	group.Go(func() error {
		teamIDs, err := uc.readerTeamIDs(userID)
		if err != nil {
			return err
		}
		for _, teamID := range teamIDs {
			members, err := uc.teamMemberRepo.GetByTeamID(teamID)
			if err != nil {
				return err
			}
			for _, member := range members {
				if member.MemberID > 0 {
					teamVisibles[member.MemberID] = true
				}
			}
		}
		return nil
	})

	pairVisibles := map[uint]bool{}
	group.Go(func() error {
		pairings, err := uc.pairingRepo.GetByMentorID(userID)
		if err != nil {
			return fmt.Errorf("ペア情報の取得に失敗しました: %w", err)
		}
		for _, p := range pairings {
			pairVisibles[p.MenteeID] = true
		}
		return nil
	})

	sharedPairVisibles := map[uint]bool{}
	group.Go(func() error {
		sharedPairings, err := uc.pairingRepo.GetPairingIDsBySharerUserID(userID)
		if err != nil {
			return fmt.Errorf("共有されたペア情報の取得に失敗しました: %w", err)
		}
		pairings, err := uc.pairingRepo.GetByIDs(sharedPairings)
		if err != nil {
			return fmt.Errorf("ペア情報の取得に失敗しました: %w", err)
		}
		for _, pairing := range pairings {
			sharedPairVisibles[pairing.MenteeID] = true
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	visible := map[uint]bool{}
	maps.Copy(visible, deptVisibles)
	maps.Copy(visible, teamVisibles)
	maps.Copy(visible, pairVisibles)
	maps.Copy(visible, sharedPairVisibles)
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

package planning_usecase

import (
	"context"
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	meeting_mng_usecase "growlink/internal/application/usecase/management_usecase/meeting_management_usecase"
	pairing_mng_usecase "growlink/internal/application/usecase/management_usecase/pairing_management_usecase"
	"growlink/internal/application/usecase/meeting_usecase"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"

	"golang.org/x/sync/errgroup"
)

type PlanningAccessUsecase struct {
	meetingRepo    repository_interface.MeetingRepository
	userRepo       repository_interface.UserRepository
	deptRepo       repository_interface.DepartmentRepository
	deptMemberRepo repository_interface.DepartmentMemberRepository
	hierarchyRepo  repository_interface.DepartmentHierarchyRepository
	teamRepo       repository_interface.TeamRepository
	teamMemberRepo repository_interface.TeamMemberRepository
	meetingViewer  *meeting_mng_usecase.MeetingViewer
	meetingList    *meeting_usecase.MeetingListUsecase
	pairingViewer  *pairing_mng_usecase.PairingViewer
	pairingRepo    repository_interface.PairingRepository
}

func NewPlanningAccessUsecase(
	meetingRepo repository_interface.MeetingRepository,
	userRepo repository_interface.UserRepository,
	deptRepo repository_interface.DepartmentRepository,
	deptMemberRepo repository_interface.DepartmentMemberRepository,
	hierarchyRepo repository_interface.DepartmentHierarchyRepository,
	teamRepo repository_interface.TeamRepository,
	teamMemberRepo repository_interface.TeamMemberRepository,
	meetingViewer *meeting_mng_usecase.MeetingViewer,
	meetingList *meeting_usecase.MeetingListUsecase,
	pairingViewer *pairing_mng_usecase.PairingViewer,
	pairingRepo repository_interface.PairingRepository,
) *PlanningAccessUsecase {
	return &PlanningAccessUsecase{
		meetingRepo:    meetingRepo,
		userRepo:       userRepo,
		deptRepo:       deptRepo,
		deptMemberRepo: deptMemberRepo,
		hierarchyRepo:  hierarchyRepo,
		teamRepo:       teamRepo,
		teamMemberRepo: teamMemberRepo,
		meetingViewer:  meetingViewer,
		meetingList:    meetingList,
		pairingViewer:  pairingViewer,
		pairingRepo:    pairingRepo,
	}
}

func dedupeMeetings(meetings []outputmodel.MeetingOutput) []outputmodel.MeetingOutput {
	if len(meetings) == 0 {
		return meetings
	}
	seen := make(map[uint]struct{}, len(meetings))
	out := make([]outputmodel.MeetingOutput, 0, len(meetings))
	for _, m := range meetings {
		if _, ok := seen[m.ID]; ok {
			continue
		}
		seen[m.ID] = struct{}{}
		out = append(out, m)
	}
	return out
}

func dedupePairings(pairings []outputmodel.PairingOutput) []outputmodel.PairingOutput {
	if len(pairings) == 0 {
		return pairings
	}
	seen := make(map[uint]struct{}, len(pairings))
	out := make([]outputmodel.PairingOutput, 0, len(pairings))
	for _, p := range pairings {
		if _, ok := seen[p.ID]; ok {
			continue
		}
		seen[p.ID] = struct{}{}
		out = append(out, p)
	}
	return out
}

func applyPairingsToMeetings(meetings []outputmodel.MeetingOutput, pairings []entity.Pairing) []outputmodel.MeetingOutput {
	byID := make(map[uint]outputmodel.PairingOutput, len(pairings))
	for _, p := range pairings {
		byID[p.ID] = converter.PairingEntityToOutput(p)
	}
	for i := range meetings {
		pid := meetings[i].Pairing.ID
		if out, ok := byID[pid]; ok {
			meetings[i].Pairing = out
		}
	}
	return meetings
}

func deletedPairingsInPlanningScope(pairings []entity.Pairing, userID uint, referable map[uint]bool) []entity.Pairing {
	out := make([]entity.Pairing, 0)
	for _, p := range pairings {
		if !p.DeletedAt.Valid || p.IsTemporary {
			continue
		}
		if p.MentorID == userID || p.MenteeID == userID {
			out = append(out, p)
			continue
		}
		if referable != nil && referable[p.MenteeID] {
			out = append(out, p)
		}
	}
	return out
}

func (uc *PlanningAccessUsecase) meetingsForPairings(pairings []entity.Pairing) ([]outputmodel.MeetingOutput, error) {
	var merged []outputmodel.MeetingOutput
	for _, pairing := range pairings {
		meetings, err := uc.meetingRepo.GetByPairingID(pairing.ID)
		if err != nil {
			return nil, err
		}
		for i := range meetings {
			p := pairing
			meetings[i].Pairing = &p
		}
		merged = append(merged, converter.MeetingEntitiesToOutputs(meetings)...)
	}
	return merged, nil
}

func (uc *PlanningAccessUsecase) ListPlanningMeetings(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole) (ret0 []outputmodel.MeetingOutput, err error) {
	const op = "planning_usecase.ListPlanningMeetings"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	pairings, err := uc.pairingRepo.GetAllContainDeleted()
	if err != nil {
		return []outputmodel.MeetingOutput{}, fmt.Errorf("ペアリング一覧の取得に失敗しました: %w", err)
	}

	if fullPlanningAccess(systemRole) || authority.IsOwner() {
		meetings, err := uc.meetingRepo.GetAll()
		if err != nil {
			return []outputmodel.MeetingOutput{}, fmt.Errorf("ミーティング一覧の取得に失敗しました: %w", err)
		}
		outs := converter.MeetingEntitiesToOutputs(meetings)
		return applyPairingsToMeetings(outs, pairings), nil
	}

	merged := []outputmodel.MeetingOutput{}
	group, _ := errgroup.WithContext(context.Background())
	group.Go(func() error {
		involved, err := uc.meetingList.ListInvolvedMeetings(userID)
		if err != nil {
			return err
		}
		merged = append(merged, involved...)
		return nil
	})

	group.Go(func() error {
		referable, err := uc.ReferableUserIDs(userID, systemRole)
		if err != nil {
			return err
		}
		menteeIDs := make([]uint, len(referable))
		var menteeIdx int
		for id := range referable {
			menteeIDs[menteeIdx] = id
			menteeIdx++
		}
		scopedMeetings, err := uc.listMeetingsByMenteeIDs(menteeIDs)
		if err != nil {
			return err
		}
		merged = append(merged, scopedMeetings...)

		deletedScoped := deletedPairingsInPlanningScope(pairings, userID, referable)
		if len(deletedScoped) > 0 {
			fromDeleted, err := uc.meetingsForPairings(deletedScoped)
			if err != nil {
				return err
			}
			merged = append(merged, fromDeleted...)
		}
		return nil
	})

	group.Go(func() error {
		shared, err := uc.meetingList.ListSharedMeetings(userID)
		if err != nil {
			return err
		}
		merged = append(merged, shared...)
		return nil
	})

	if err := group.Wait(); err != nil {
		return []outputmodel.MeetingOutput{}, fmt.Errorf("面談一覧の取得に失敗しました: %w", err)
	}

	merged = dedupeMeetings(merged)
	return applyPairingsToMeetings(merged, pairings), nil
}

func (uc *PlanningAccessUsecase) CanViewPlanningMeeting(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole, meetingID uint) (ret0 bool, err error) {
	const op = "planning_usecase.CanViewPlanningMeeting"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if fullPlanningAccess(systemRole) || authority.IsOwner() {
		return true, nil
	}
	list, err := uc.ListPlanningMeetings(userID, authority, systemRole)
	if err != nil {
		return false, err
	}
	for _, m := range list {
		if m.ID == meetingID {
			return true, nil
		}
	}
	return false, nil
}

func (uc *PlanningAccessUsecase) HasPairManagementAccess(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole) (ret0 bool, err error) {
	const op = "planning_usecase.HasPairManagementAccess"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if fullPlanningAccess(systemRole) || authority.IsOwner() {
		return true, nil
	}
	scope, err := uc.PairManagementScopeUserIDs(userID, authority, systemRole)
	if err != nil {
		return false, err
	}
	return len(scope) > 0, nil
}

func (uc *PlanningAccessUsecase) ListPairingsForManagement(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole) (ret0 []outputmodel.PairingOutput, err error) {
	const op = "planning_usecase.ListPairingsForManagement"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if fullPlanningAccess(systemRole) || authority.IsOwner() {
		return uc.pairingViewer.ListPairingsByAll()
	}
	return uc.listPairManagementScopedPairings(userID, authority, systemRole)
}

func (uc *PlanningAccessUsecase) listPairManagementScopedPairings(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole) ([]outputmodel.PairingOutput, error) {
	scope, err := uc.PairManagementScopeUserIDs(userID, authority, systemRole)
	if err != nil {
		return nil, err
	}
	menteeIDs := make([]uint, len(scope))
	var menteeIdx int
	for id := range scope {
		menteeIDs[menteeIdx] = id
		menteeIdx++
	}
	return uc.listPairingsByMenteeIDs(menteeIDs)
}

func (uc *PlanningAccessUsecase) ListPairingsForMeeting(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole) (ret0 []outputmodel.PairingOutput, err error) {
	const op = "planning_usecase.ListPairingsForMeeting"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if fullPlanningAccess(systemRole) || authority.IsOwner() {
		return uc.pairingViewer.ListPairingsByAll()
	}

	merged := []outputmodel.PairingOutput{}
	mentorPairs, err := uc.pairingViewer.ListPairingsByMentor(userID)
	if err != nil {
		return nil, err
	}
	merged = append(merged, mentorPairs...)

	if menteePairings, err := uc.pairingRepo.GetByMenteeID(userID); err == nil {
		for _, p := range menteePairings {
			out, err := uc.pairingViewer.GetPairing(p.ID)
			if err == nil && out != nil {
				merged = append(merged, *out)
			}
		}
	}

	scoped, err := uc.listReferableScopedPairings(userID, systemRole)
	if err != nil {
		return nil, err
	}
	merged = append(merged, scoped...)
	return dedupePairings(merged), nil
}

func (uc *PlanningAccessUsecase) listReferableScopedPairings(userID uint, systemRole valueobject.SystemRole) ([]outputmodel.PairingOutput, error) {
	referable, err := uc.ReferableUserIDs(userID, systemRole)
	if err != nil {
		return nil, err
	}
	menteeIDs := make([]uint, len(referable))
	var menteeIdx int
	for id := range referable {
		menteeIDs[menteeIdx] = id
		menteeIdx++
	}
	return uc.listPairingsByMenteeIDs(menteeIDs)
}

func (uc *PlanningAccessUsecase) readerTeamIDs(userID uint) ([]uint, error) {
	if uc.teamMemberRepo == nil {
		return nil, nil
	}
	members, err := uc.teamMemberRepo.GetByMemberID(userID)
	if err != nil {
		return nil, err
	}
	ids := make([]uint, 0)
	for _, member := range members {
		if member.IsReader {
			ids = append(ids, member.TeamID)
		}
	}
	return ids, nil
}

func (uc *PlanningAccessUsecase) departmentScopeUserIDs(userID uint) ([]uint, error) {
	if uc.deptMemberRepo == nil || uc.hierarchyRepo == nil {
		return nil, nil
	}
	memberships, err := uc.deptMemberRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	headRoots := make([]uint, 0)
	managerRoots := make([]uint, 0)
	for _, membership := range memberships {
		if membership.IsHead {
			headRoots = append(headRoots, membership.DepartmentID)
			continue
		}
		if membership.IsManager {
			managerRoots = append(managerRoots, membership.DepartmentID)
		}
	}
	if len(headRoots) == 0 && len(managerRoots) == 0 {
		return nil, nil
	}
	hierarchies, err := uc.hierarchyRepo.GetAll()
	if err != nil {
		return nil, err
	}
	seen := map[uint]struct{}{}
	if len(headRoots) > 0 {
		if err := uc.mergeDepartmentTreeUsers(seen, headRoots, hierarchies, nil); err != nil {
			return nil, err
		}
	}
	if len(managerRoots) > 0 {
		skipHeadInDept := make(map[uint]struct{}, len(managerRoots))
		for _, deptID := range managerRoots {
			skipHeadInDept[deptID] = struct{}{}
		}
		if err := uc.mergeDepartmentTreeUsers(seen, managerRoots, hierarchies, skipHeadInDept); err != nil {
			return nil, err
		}
	}
	ids := make([]uint, len(seen))
	var idx int
	for id := range seen {
		ids[idx] = id
		idx++
	}
	return ids, nil
}

// skipHeadInDept: for listed department IDs only, members with IsHead are omitted (manager scope).
func (uc *PlanningAccessUsecase) mergeDepartmentTreeUsers(
	seen map[uint]struct{},
	roots []uint,
	hierarchies []entity.DepartmentHierarchy,
	skipHeadInDept map[uint]struct{},
) error {
	deptIDs := descendantDepartmentIDs(roots, hierarchies)
	headUserIDsInRoot := map[uint]struct{}{}
	for deptID := range skipHeadInDept {
		members, err := uc.deptMemberRepo.GetByDepartmentID(deptID)
		if err != nil {
			return err
		}
		for _, member := range members {
			if member.IsHead {
				headUserIDsInRoot[member.UserID] = struct{}{}
			}
		}
	}
	for _, deptID := range deptIDs {
		members, err := uc.deptMemberRepo.GetByDepartmentID(deptID)
		if err != nil {
			return err
		}
		_, skipHeadsHere := skipHeadInDept[deptID]
		for _, member := range members {
			if skipHeadsHere && member.IsHead {
				continue
			}
			seen[member.UserID] = struct{}{}
		}
		if uc.teamRepo == nil || uc.teamMemberRepo == nil {
			continue
		}
		teams, err := uc.teamRepo.GetByDepartmentID(deptID)
		if err != nil {
			return err
		}
		for _, team := range teams {
			teamMembers, err := uc.teamMemberRepo.GetByTeamID(team.ID)
			if err != nil {
				return err
			}
			for _, teamMember := range teamMembers {
				if _, isHead := headUserIDsInRoot[teamMember.MemberID]; isHead {
					continue
				}
				seen[teamMember.MemberID] = struct{}{}
			}
		}
	}
	return nil
}

func descendantDepartmentIDs(rootIDs []uint, hierarchies []entity.DepartmentHierarchy) []uint {
	children := map[uint][]uint{}
	for _, hierarchy := range hierarchies {
		children[hierarchy.ParentDepartmentID] = append(children[hierarchy.ParentDepartmentID], hierarchy.ChildDepartmentID)
	}
	seen := map[uint]bool{}
	var walk func(uint)
	walk = func(id uint) {
		if seen[id] {
			return
		}
		seen[id] = true
		for _, childID := range children[id] {
			walk(childID)
		}
	}
	for _, rootID := range rootIDs {
		walk(rootID)
	}
	ids := make([]uint, len(seen))
	var idx int
	for id := range seen {
		ids[idx] = id
		idx++
	}
	return ids
}

func (uc *PlanningAccessUsecase) listMeetingsByMenteeIDs(ids []uint) ([]outputmodel.MeetingOutput, error) {
	var merged []outputmodel.MeetingOutput
	for _, menteeID := range ids {
		pairings, err := uc.pairingRepo.GetByMenteeID(menteeID)
		if err != nil {
			return nil, err
		}
		for _, pairing := range pairings {
			meetings, err := uc.meetingRepo.GetByPairingID(pairing.ID)
			if err != nil {
				return nil, err
			}
			merged = append(merged, converter.MeetingEntitiesToOutputs(meetings)...)
		}
	}
	return merged, nil
}

func (uc *PlanningAccessUsecase) listPairingsByMenteeIDs(ids []uint) ([]outputmodel.PairingOutput, error) {
	pairings, err := uc.pairingRepo.GetByMenteeIDs(ids)
	if err != nil {
		return nil, err
	}

	result := make([]outputmodel.PairingOutput, len(pairings))
	for idx, pairing := range pairings {
		result[idx] = converter.PairingEntityToOutput(pairing)
	}

	return result, nil
}

func (uc *PlanningAccessUsecase) CanSpecifyPairMeeting(userID uint, systemRole valueobject.SystemRole, authority valueobject.Authority) (ret0 bool, err error) {
	const op = "planning_usecase.CanSpecifyPairMeeting"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if fullPlanningAccess(systemRole) || authority.IsOwner() {
		return true, nil
	}
	referable, err := uc.ReferableUserIDs(userID, systemRole)
	if err != nil {
		return false, err
	}
	return len(referable) > 0, nil
}

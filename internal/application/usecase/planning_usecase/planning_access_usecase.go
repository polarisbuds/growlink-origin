package planning_usecase

import (
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	meeting_mng_usecase "growlink/internal/application/usecase/management_usecase/meeting_management_usecase"
	pairing_mng_usecase "growlink/internal/application/usecase/management_usecase/pairing_management_usecase"
	"growlink/internal/application/usecase/meeting_usecase"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
)

type PlanningAccessUsecase struct {
	meetingRepo    repository_interface.MeetingRepository
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

func (uc *PlanningAccessUsecase) ListPlanningMeetings(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole) ([]outputmodel.MeetingOutput, error) {
	if fullPlanningAccess(systemRole) || authority.IsOwner() {
		meetings, err := uc.meetingRepo.GetAll()
		if err != nil {
			return []outputmodel.MeetingOutput{}, fmt.Errorf("ミーティング一覧の取得に失敗しました: %w", err)
		}
		return converter.MeetingEntitiesToOutputs(meetings), nil
	}
	merged := []outputmodel.MeetingOutput{}

	involved, err := uc.meetingList.ListInvolvedMeetings(userID)
	if err != nil {
		return []outputmodel.MeetingOutput{}, err
	}
	merged = append(merged, involved...)

	referable, err := uc.ReferableUserIDs(userID, systemRole)
	if err != nil {
		return []outputmodel.MeetingOutput{}, err
	}
	menteeIDs := make([]uint, 0, len(referable))
	for id := range referable {
		menteeIDs = append(menteeIDs, id)
	}
	scopedMeetings, err := uc.listMeetingsByMenteeIDs(menteeIDs)
	if err != nil {
		return []outputmodel.MeetingOutput{}, err
	}
	merged = append(merged, scopedMeetings...)

	shared, err := uc.meetingList.ListSharedMeetings(userID)
	if err != nil {
		return []outputmodel.MeetingOutput{}, err
	}
	merged = append(merged, shared...)

	return dedupeMeetings(merged), nil
}

func (uc *PlanningAccessUsecase) CanViewPlanningMeeting(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole, meetingID uint) (bool, error) {
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

func (uc *PlanningAccessUsecase) HasPairManagementAccess(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole) (bool, error) {
	if fullPlanningAccess(systemRole) || authority.IsOwner() {
		return true, nil
	}
	scope, err := uc.PairManagementScopeUserIDs(userID, authority, systemRole)
	if err != nil {
		return false, err
	}
	return len(scope) > 0, nil
}

func (uc *PlanningAccessUsecase) ListPairingsForManagement(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole) ([]outputmodel.PairingOutput, error) {
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
	menteeIDs := make([]uint, 0, len(scope))
	for id := range scope {
		menteeIDs = append(menteeIDs, id)
	}
	return uc.listPairingsByMenteeIDs(menteeIDs)
}

func (uc *PlanningAccessUsecase) ListPairingsForMeeting(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole) ([]outputmodel.PairingOutput, error) {
	if fullPlanningAccess(systemRole) || authority.IsOwner() {
		return uc.pairingViewer.ListPairingsByAll()
	}

	merged := []outputmodel.PairingOutput{}
	mentorPairs, err := uc.pairingViewer.ListPairingsByMentor(userID)
	if err != nil {
		return nil, err
	}
	merged = append(merged, mentorPairs...)

	if pairing, err := uc.pairingRepo.GetByMenteeID(userID); err == nil && pairing != nil {
		out, err := uc.pairingViewer.GetPairing(pairing.ID)
		if err == nil && out != nil {
			merged = append(merged, *out)
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
	menteeIDs := make([]uint, 0, len(referable))
	for id := range referable {
		menteeIDs = append(menteeIDs, id)
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
	roots := make([]uint, 0)
	for _, membership := range memberships {
		if membership.IsHead || membership.IsManager {
			roots = append(roots, membership.DepartmentID)
		}
	}
	if len(roots) == 0 {
		return nil, nil
	}
	hierarchies, err := uc.hierarchyRepo.GetAll()
	if err != nil {
		return nil, err
	}
	deptIDs := descendantDepartmentIDs(roots, hierarchies)
	seen := map[uint]struct{}{}
	for _, deptID := range deptIDs {
		members, err := uc.deptMemberRepo.GetByDepartmentID(deptID)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			seen[member.UserID] = struct{}{}
		}
		if uc.teamRepo == nil || uc.teamMemberRepo == nil {
			continue
		}
		teams, err := uc.teamRepo.GetByDepartmentID(deptID)
		if err != nil {
			return nil, err
		}
		for _, team := range teams {
			teamMembers, err := uc.teamMemberRepo.GetByTeamID(team.ID)
			if err != nil {
				return nil, err
			}
			for _, teamMember := range teamMembers {
				seen[teamMember.MemberID] = struct{}{}
			}
		}
	}
	ids := make([]uint, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids, nil
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
	ids := make([]uint, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

func (uc *PlanningAccessUsecase) listMeetingsByMenteeIDs(ids []uint) ([]outputmodel.MeetingOutput, error) {
	var merged []outputmodel.MeetingOutput
	for _, menteeID := range ids {
		pairing, err := uc.pairingRepo.GetByMenteeID(menteeID)
		if err != nil {
			return nil, err
		}
		if pairing == nil {
			continue
		}
		meetings, err := uc.meetingRepo.GetByPairingID(pairing.ID)
		if err != nil {
			return nil, err
		}
		merged = append(merged, converter.MeetingEntitiesToOutputs(meetings)...)
	}
	return merged, nil
}

func (uc *PlanningAccessUsecase) listPairingsByMenteeIDs(ids []uint) ([]outputmodel.PairingOutput, error) {
	var merged []outputmodel.PairingOutput
	for _, menteeID := range ids {
		pairing, err := uc.pairingRepo.GetByMenteeID(menteeID)
		if err != nil {
			return nil, err
		}
		if pairing == nil {
			continue
		}
		out, err := uc.pairingViewer.GetPairing(pairing.ID)
		if err != nil {
			return nil, err
		}
		if out != nil {
			merged = append(merged, *out)
		}
	}
	return merged, nil
}

func (uc *PlanningAccessUsecase) CanSpecifyPairMeeting(userID uint, systemRole valueobject.SystemRole, authority valueobject.Authority) (bool, error) {
	if fullPlanningAccess(systemRole) || authority.IsOwner() {
		return true, nil
	}
	referable, err := uc.ReferableUserIDs(userID, systemRole)
	if err != nil {
		return false, err
	}
	return len(referable) > 0, nil
}

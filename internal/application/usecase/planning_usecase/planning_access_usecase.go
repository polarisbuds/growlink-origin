package planning_usecase

import (
	"fmt"

	"growlink/internal/application/converter"
	meeting_mng_usecase "growlink/internal/application/usecase/management_usecase/meeting_management_usecase"
	pairing_mng_usecase "growlink/internal/application/usecase/management_usecase/pairing_management_usecase"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/application/usecase/meeting_usecase"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
)

type PlanningAccessUsecase struct {
	meetingRepo      repository_interface.MeetingRepository
	userRepo         repository_interface.UserRepository
	teamRefRepo      repository_interface.UserReferenceTeamRepository
	relationshipRepo repository_interface.UserRelationshipRepository
	teamMemberRepo   repository_interface.TeamMemberRepository
	meetingViewer    *meeting_mng_usecase.MeetingViewer
	meetingList      *meeting_usecase.MeetingListUsecase
	pairingViewer    *pairing_mng_usecase.PairingViewer
	pairingRepo      repository_interface.PairingRepository
}

func NewPlanningAccessUsecase(
	meetingRepo repository_interface.MeetingRepository,
	userRepo repository_interface.UserRepository,
	teamRefRepo repository_interface.UserReferenceTeamRepository,
	relationshipRepo repository_interface.UserRelationshipRepository,
	teamMemberRepo repository_interface.TeamMemberRepository,
	meetingViewer *meeting_mng_usecase.MeetingViewer,
	meetingList *meeting_usecase.MeetingListUsecase,
	pairingViewer *pairing_mng_usecase.PairingViewer,
	pairingRepo repository_interface.PairingRepository,
) *PlanningAccessUsecase {
	return &PlanningAccessUsecase{
		meetingRepo:      meetingRepo,
		userRepo:         userRepo,
		teamRefRepo:      teamRefRepo,
		relationshipRepo: relationshipRepo,
		teamMemberRepo:   teamMemberRepo,
		meetingViewer:    meetingViewer,
		meetingList:      meetingList,
		pairingViewer:    pairingViewer,
		pairingRepo:      pairingRepo,
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

func (uc *PlanningAccessUsecase) ListPlanningMeetings(userID uint, isSystemAdmin bool) ([]outputmodel.MeetingOutput, error) {
	if isSystemAdmin {
		meetings, err := uc.meetingRepo.GetAll()
		if err != nil {
			return nil, fmt.Errorf("ミーティング一覧の取得に失敗しました: %w", err)
		}
		return converter.MeetingEntitiesToOutputs(meetings), nil
	}
	var merged []outputmodel.MeetingOutput

	involved, err := uc.meetingList.ListInvolvedMeetings(userID)
	if err != nil {
		return nil, err
	}
	merged = append(merged, involved...)

	bossMeetings, err := uc.meetingViewer.ListMeetingsByBoss(userID)
	if err != nil {
		return nil, err
	}
	merged = append(merged, bossMeetings...)

	teamMeetings, err := uc.listMeetingsForManagedTeams(userID)
	if err != nil {
		return nil, err
	}
	merged = append(merged, teamMeetings...)

	shared, err := uc.meetingList.ListSharedMeetings(userID)
	if err != nil {
		return nil, err
	}
	merged = append(merged, shared...)

	return dedupeMeetings(merged), nil
}

func (uc *PlanningAccessUsecase) CanViewPlanningMeeting(userID uint, isSystemAdmin bool, meetingID uint) (bool, error) {
	if isSystemAdmin {
		return true, nil
	}
	list, err := uc.ListPlanningMeetings(userID, false)
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

func (uc *PlanningAccessUsecase) listMeetingsForManagedTeams(userID uint) ([]outputmodel.MeetingOutput, error) {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user.OperationTarget != valueobject.OperationTargetTeam || uc.teamRefRepo == nil {
		return nil, nil
	}
	refs, err := uc.teamRefRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	var merged []outputmodel.MeetingOutput
	for _, ref := range refs {
		teamMeetings, err := uc.meetingViewer.ListMeetingsByTeam(ref.ReferableTeamID)
		if err != nil {
			return nil, err
		}
		merged = append(merged, teamMeetings...)
	}
	return merged, nil
}

func (uc *PlanningAccessUsecase) HasPairManagementAccess(userID uint, isSystemAdmin bool) (bool, error) {
	if isSystemAdmin {
		return true, nil
	}
	rels, err := uc.relationshipRepo.GetByBossID(userID)
	if err != nil {
		return false, err
	}
	if len(rels) > 0 {
		return true, nil
	}
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}
	if user.OperationTarget != valueobject.OperationTargetTeam || uc.teamRefRepo == nil {
		return false, nil
	}
	refs, err := uc.teamRefRepo.GetByUserID(userID)
	if err != nil {
		return false, err
	}
	for _, ref := range refs {
		members, err := uc.teamMemberRepo.GetByTeamID(ref.ReferableTeamID)
		if err != nil {
			return false, err
		}
		if len(members) > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (uc *PlanningAccessUsecase) ListPairingsForManagement(userID uint, isSystemAdmin bool) ([]outputmodel.PairingOutput, error) {
	if isSystemAdmin {
		return uc.pairingViewer.ListPairingsByAll()
	}
	return uc.listSubordinateScopedPairings(userID)
}

func (uc *PlanningAccessUsecase) ListPairingsForMeeting(userID uint, isSystemAdmin bool) ([]outputmodel.PairingOutput, error) {
	if isSystemAdmin {
		return uc.pairingViewer.ListPairingsByAll()
	}
	var merged []outputmodel.PairingOutput
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

	subordinate, err := uc.listSubordinateScopedPairings(userID)
	if err != nil {
		return nil, err
	}
	merged = append(merged, subordinate...)
	return dedupePairings(merged), nil
}

func (uc *PlanningAccessUsecase) listSubordinateScopedPairings(userID uint) ([]outputmodel.PairingOutput, error) {
	var merged []outputmodel.PairingOutput
	bossPairs, err := uc.pairingViewer.ListPairingsByBoss(userID)
	if err != nil {
		return nil, err
	}
	merged = append(merged, bossPairs...)

	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user.OperationTarget == valueobject.OperationTargetTeam && uc.teamRefRepo != nil {
		refs, err := uc.teamRefRepo.GetByUserID(userID)
		if err != nil {
			return nil, err
		}
		for _, ref := range refs {
			teamPairs, err := uc.pairingViewer.ListPairingsByTeam(ref.ReferableTeamID)
			if err != nil {
				return nil, err
			}
			merged = append(merged, teamPairs...)
		}
	}
	return dedupePairings(merged), nil
}

func (uc *PlanningAccessUsecase) CanSpecifyPairMeeting(userID uint, systemRole valueobject.SystemRole, authority valueobject.Authority) (bool, error) {
	if systemRole == valueobject.SystemRoleAdmin {
		return true, nil
	}
	if int(authority) >= int(valueobject.AuthorityBoss) {
		return true, nil
	}
	pairs, err := uc.pairingViewer.ListPairingsByMentor(userID)
	if err != nil {
		return false, err
	}
	return len(pairs) > 0, nil
}

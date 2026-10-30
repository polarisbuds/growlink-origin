package meeting_usecase

import (
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
)

type MeetingListUsecase struct {
	meetingRepo      repository_interface.MeetingRepository
	pairingRepo      repository_interface.PairingRepository
	relationshipRepo repository_interface.UserRelationshipRepository
	userRepo         repository_interface.UserRepository
	teamMemberRepo   repository_interface.TeamMemberRepository
}

func NewMeetingListUsecase(
	meetingRepo repository_interface.MeetingRepository,
	pairingRepo repository_interface.PairingRepository,
	relationshipRepo repository_interface.UserRelationshipRepository,
	userRepo repository_interface.UserRepository,
	teamMemberRepo repository_interface.TeamMemberRepository,
) *MeetingListUsecase {
	return &MeetingListUsecase{
		meetingRepo:      meetingRepo,
		pairingRepo:      pairingRepo,
		relationshipRepo: relationshipRepo,
		userRepo:         userRepo,
		teamMemberRepo:   teamMemberRepo,
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

func userInvolvedInPairing(userID, departmentID uint, teamIDs map[uint]struct{}, pairing entity.Pairing) bool {
	if pairing.MentorID == userID {
		return true
	}
	if pairing.MenteeID == userID {
		return true
	}
	for _, t := range pairing.Targets {
		switch t.Kind {
		case entity.PairingTargetKindUser:
			if t.TargetID == userID {
				return true
			}
		case entity.PairingTargetKindTeam:
			if _, ok := teamIDs[t.TargetID]; ok {
				return true
			}
		case entity.PairingTargetKindDepartment:
			if t.TargetID == departmentID {
				return true
			}
		}
	}
	for _, m := range pairing.Mentees {
		if m.UserID == userID {
			return true
		}
	}
	return false
}

func (uc *MeetingListUsecase) ListMyMeeting(userID uint) ([]outputmodel.MeetingOutput, error) {
	pairing, err := uc.pairingRepo.GetByMenteeID(userID)
	if err != nil {
		return []outputmodel.MeetingOutput{}, nil
	}
	if pairing != nil {
		meetings, err := uc.meetingRepo.GetByPairingID(pairing.ID)
		if err != nil {
			return nil, fmt.Errorf("ミーティング一覧の取得に失敗しました: %w", err)
		}
		return converter.MeetingEntitiesToOutputs(meetings), nil
	}
	return []outputmodel.MeetingOutput{}, nil
}

func (uc *MeetingListUsecase) ListInvolvedMeetings(userID uint) ([]outputmodel.MeetingOutput, error) {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザ情報の取得に失敗しました: %w", err)
	}
	teamMembers, err := uc.teamMemberRepo.GetByMemberID(userID)
	if err != nil {
		return nil, fmt.Errorf("チーム所属の取得に失敗しました: %w", err)
	}
	teamIDs := make(map[uint]struct{}, len(teamMembers))
	for _, tm := range teamMembers {
		teamIDs[tm.TeamID] = struct{}{}
	}

	pairings, err := uc.pairingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("ペアリング一覧の取得に失敗しました: %w", err)
	}

	var all []outputmodel.MeetingOutput
	for _, p := range pairings {
		if !userInvolvedInPairing(userID, user.DepartmentID, teamIDs, p) {
			continue
		}
		meetings, err := uc.meetingRepo.GetByPairingID(p.ID)
		if err != nil {
			return nil, err
		}
		all = append(all, converter.MeetingEntitiesToOutputs(meetings)...)
	}
	return dedupeMeetings(all), nil
}

func (uc *MeetingListUsecase) ListSharedMeetings(userID uint) ([]outputmodel.MeetingOutput, error) {
	involved, err := uc.ListInvolvedMeetings(userID)
	if err != nil {
		return nil, err
	}
	exclude := make(map[uint]struct{}, len(involved))
	for _, m := range involved {
		exclude[m.ID] = struct{}{}
	}

	var shared []outputmodel.MeetingOutput

	meetingIDs, err := uc.meetingRepo.GetMeetingIDsBySharerUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("共有ミーティングの取得に失敗しました: %w", err)
	}
	for _, meetingID := range meetingIDs {
		meeting, err := uc.meetingRepo.GetByID(meetingID)
		if err != nil {
			return nil, err
		}
		shared = append(shared, converter.MeetingEntityToOutput(*meeting))
	}

	relationships, err := uc.relationshipRepo.GetByBossID(userID)
	if err != nil {
		return nil, fmt.Errorf("部下リレーションの取得に失敗しました: %w", err)
	}
	allPairings, err := uc.pairingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("ペアリング一覧の取得に失敗しました: %w", err)
	}
	for _, rel := range relationships {
		pairing, err := uc.pairingRepo.GetByMenteeID(rel.SourceUserID)
		if err != nil {
			return nil, err
		}
		if pairing != nil {
			meetings, err := uc.meetingRepo.GetByPairingID(pairing.ID)
			if err != nil {
				return nil, err
			}
			shared = append(shared, converter.MeetingEntitiesToOutputs(meetings)...)
		}
		for _, p := range allPairings {
			if !p.IsTemporary || p.MenteeID != rel.SourceUserID {
				continue
			}
			meetings, err := uc.meetingRepo.GetByPairingID(p.ID)
			if err != nil {
				return nil, err
			}
			shared = append(shared, converter.MeetingEntitiesToOutputs(meetings)...)
		}
	}

	filtered := make([]outputmodel.MeetingOutput, 0, len(shared))
	for _, m := range shared {
		if _, skip := exclude[m.ID]; skip {
			continue
		}
		filtered = append(filtered, m)
	}
	return dedupeMeetings(filtered), nil
}

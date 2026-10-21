package meeting_usecase

import (
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
)

type MeetingListUsecase struct {
	meetingRepo      repository_interface.MeetingRepository
	pairingRepo      repository_interface.PairingRepository
	relationshipRepo repository_interface.UserRelationshipRepository
}

func NewMeetingListUsecase(
	meetingRepo repository_interface.MeetingRepository,
	pairingRepo repository_interface.PairingRepository,
	relationshipRepo repository_interface.UserRelationshipRepository,
) *MeetingListUsecase {
	return &MeetingListUsecase{
		meetingRepo:      meetingRepo,
		pairingRepo:      pairingRepo,
		relationshipRepo: relationshipRepo,
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
	var all []outputmodel.MeetingOutput
	menteeMeetings, err := uc.ListMyMeeting(userID)
	if err != nil {
		return nil, err
	}
	all = append(all, menteeMeetings...)

	pairings, err := uc.pairingRepo.GetByMentorID(userID)
	if err != nil {
		return nil, fmt.Errorf("メンターとしてのペア取得に失敗しました: %w", err)
	}
	for _, p := range pairings {
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

	sharerPairingIDs, err := uc.pairingRepo.GetPairingIDsBySharerUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("共有ペアの取得に失敗しました: %w", err)
	}
	for _, pairingID := range sharerPairingIDs {
		meetings, err := uc.meetingRepo.GetByPairingID(pairingID)
		if err != nil {
			return nil, err
		}
		shared = append(shared, converter.MeetingEntitiesToOutputs(meetings)...)
	}

	relationships, err := uc.relationshipRepo.GetByBossID(userID)
	if err != nil {
		return nil, fmt.Errorf("部下リレーションの取得に失敗しました: %w", err)
	}
	for _, rel := range relationships {
		pairing, err := uc.pairingRepo.GetByMenteeID(rel.SourceUserID)
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
		shared = append(shared, converter.MeetingEntitiesToOutputs(meetings)...)
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

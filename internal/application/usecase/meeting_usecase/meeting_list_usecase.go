package meeting_usecase

import (
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
)

type MeetingListUsecase struct {
	meetingRepo repository_interface.MeetingRepository
	pairingRepo repository_interface.PairingRepository
}

func NewMeetingListUsecase(
	meetingRepo repository_interface.MeetingRepository,
	pairingRepo repository_interface.PairingRepository,
) *MeetingListUsecase {
	return &MeetingListUsecase{
		meetingRepo: meetingRepo,
		pairingRepo: pairingRepo,
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

func userIsMentorOrMenteeInPairing(userID uint, pairing entity.Pairing) bool {
	if pairing.MentorID == userID {
		return true
	}
	if pairing.MenteeID == userID {
		return true
	}
	return false
}

func (uc *MeetingListUsecase) ListMyMeeting(userID uint) ([]outputmodel.MeetingOutput, error) {
	pairings, err := uc.pairingRepo.GetByMenteeID(userID)
	if err != nil {
		return []outputmodel.MeetingOutput{}, nil
	}
	all := []outputmodel.MeetingOutput{}
	for _, pairing := range pairings {
		meetings, err := uc.meetingRepo.GetByPairingID(pairing.ID)
		if err != nil {
			return nil, fmt.Errorf("ミーティング一覧の取得に失敗しました: %w", err)
		}
		all = append(all, converter.MeetingEntitiesToOutputs(meetings)...)
	}
	return dedupeMeetings(all), nil
}

func (uc *MeetingListUsecase) ListInvolvedMeetings(userID uint) ([]outputmodel.MeetingOutput, error) {
	pairings, err := uc.pairingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("ペアリング一覧の取得に失敗しました: %w", err)
	}

	// TODO; ループでの取得は非効率。pairingRepo.GetMentorID(),pairingRepo.GetMenteeID()を取得し結果をマージする
	all := []outputmodel.MeetingOutput{}
	for _, p := range pairings {
		if !userIsMentorOrMenteeInPairing(userID, p) {
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

	filtered := make([]outputmodel.MeetingOutput, 0, len(shared))
	for _, m := range shared {
		if _, skip := exclude[m.ID]; skip {
			continue
		}
		filtered = append(filtered, m)
	}
	return dedupeMeetings(filtered), nil
}

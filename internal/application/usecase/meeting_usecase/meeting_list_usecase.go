package meeting_usecase

import (
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
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
	return all, nil
}

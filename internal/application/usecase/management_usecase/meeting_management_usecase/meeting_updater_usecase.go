package meeting_mng_usecase

import (
	"fmt"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"time"
)

type MeetingUpdater struct {
	meetingRepo repository_interface.MeetingRepository
}

func NewMeetingUpdater(meetingRepo repository_interface.MeetingRepository) *MeetingUpdater {
	return &MeetingUpdater{
		meetingRepo: meetingRepo,
	}
}

func (uc *MeetingUpdater) Update(meetindID uint, meeting *entity.Meeting) error {
	if err := uc.meetingRepo.Update(meetindID, meeting); err != nil {
		return fmt.Errorf("ミーティングの更新に失敗しました: %w", err)
	}
	return nil
}

func (uc *MeetingUpdater) ChangeImplementDate(meetingID uint, implementDate time.Time) error {
	if err := uc.meetingRepo.UpdateImplementDate(meetingID, implementDate); err != nil {
		return fmt.Errorf("ミーティングの実施予定日の変更に失敗しました: %w", err)
	}
	return nil
}

func (uc *MeetingUpdater) ChangePairing(meetingID, newPairingID uint) error {
	if err := uc.meetingRepo.UpdatePairing(meetingID, newPairingID); err != nil {
		return fmt.Errorf("ミーティングのペアリングの変更に失敗しました: %w", err)
	}
	return nil
}

func (uc *MeetingUpdater) ChangeTemplate(meetingID, newTemplateID uint) error {
	if err := uc.meetingRepo.UpdateTemplate(meetingID, newTemplateID); err != nil {
		return fmt.Errorf("ミーティングのテンプレートの変更に失敗しました: %w", err)
	}
	return nil
}

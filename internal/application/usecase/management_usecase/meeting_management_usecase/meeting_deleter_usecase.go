package meeting_mng_usecase

import (
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
)

type MeetingDeleter struct {
	meetingRepo repository_interface.MeetingRepository
}

func NewMeetingDeleter(meetingRepo repository_interface.MeetingRepository) *MeetingDeleter {
	return &MeetingDeleter{
		meetingRepo: meetingRepo,
	}
}

func (deleter *MeetingDeleter) DeleteMeeting(meetingIDs ...uint) (bool, error) {
	if err := deleter.meetingRepo.Delete(meetingIDs...); err != nil {
		return false, fmt.Errorf("ミーティングの削除に失敗しました: %w", err)
	}
	return true, nil
}

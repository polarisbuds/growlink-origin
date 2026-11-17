package meeting_mng_usecase

import (
	"context"
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"
)

type MeetingDeleter struct {
	meetingRepo repository_interface.MeetingRepository
}

func NewMeetingDeleter(meetingRepo repository_interface.MeetingRepository) *MeetingDeleter {
	return &MeetingDeleter{
		meetingRepo: meetingRepo,
	}
}

func (uc *MeetingDeleter) DeleteMeeting(meetingIDs ...uint) (ret0 bool, err error) {
	const op = "meeting_mng_usecase.DeleteMeeting"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.meetingRepo.Delete(meetingIDs...); err != nil {
		return false, fmt.Errorf("ミーティングの削除に失敗しました: %w", err)
	}
	return true, nil
}

package meeting_mng_usecase

import (
	"fmt"
	"time"

	meeting_usecase "growlink/internal/application/usecase/meeting_usecase"
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
)

func (uc *MeetingUpdater) ChangeImplementDateByUUID(
	meetingUUID uuid.UUID,
	viewerID uint,
	systemRole valueobject.SystemRole,
	implementDate time.Time,
) (time.Time, error) {
	meeting, err := uc.meetingRepo.GetByUUID(meetingUUID)
	if err != nil {
		return time.Time{}, err
	}
	if !meeting_usecase.CanEditMeetingImplementedDate(meeting, viewerID, systemRole) {
		return time.Time{}, fmt.Errorf("forbidden")
	}
	if err := uc.ChangeImplementDate(meeting.ID, implementDate); err != nil {
		return time.Time{}, err
	}
	return implementDate, nil
}
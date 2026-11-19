package meeting_mng_usecase

import (
	"fmt"

	meeting_usecase "growlink/internal/application/usecase/meeting_usecase"
	outputmodel "growlink/internal/application/model/output_model"

	"github.com/google/uuid"
)

func (uc *MeetingUpdater) ReplaceSharersByMeetingUUID(meetingUUID uuid.UUID, viewerID uint, sharerUUIDs []string) ([]outputmodel.UserOutput, error) {
	meeting, err := uc.meetingRepo.GetByUUID(meetingUUID)
	if err != nil {
		return nil, err
	}
	if !meeting_usecase.CanEditMeetingSharers(meeting, viewerID) {
		return nil, fmt.Errorf("forbidden")
	}
	sharerIDs, err := uc.resolveUserUUIDs(sharerUUIDs)
	if err != nil {
		return nil, err
	}
	if err := uc.meetingRepo.ReplaceSharers(meeting.ID, sharerIDs); err != nil {
		return nil, err
	}
	updated, err := uc.meetingRepo.GetByUUID(meetingUUID)
	if err != nil {
		return nil, err
	}
	return meeting_usecase.MeetingSharersToOutputs(updated), nil
}

func (uc *MeetingUpdater) resolveUserUUIDs(uuids []string) ([]uint, error) {
	if len(uuids) == 0 {
		return nil, nil
	}
	if uc.userRepo == nil {
		return nil, fmt.Errorf("user repository is not configured")
	}
	users, err := uc.userRepo.GetByUUIDs(uuids)
	if err != nil {
		return nil, err
	}
	ids := make([]uint, len(users))
	for idx, u := range users {
		ids[idx] = u.ID
	}
	return ids, nil
}
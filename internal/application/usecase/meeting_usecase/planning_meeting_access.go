package meeting_usecase

import "growlink/internal/domain/valueobject"

type PlanningMeetingAccess interface {
	CanViewPlanningMeeting(userID uint, authority valueobject.Authority, systemRole valueobject.SystemRole, meetingID uint) (bool, error)
}
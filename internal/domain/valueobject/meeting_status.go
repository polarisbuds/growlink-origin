package valueobject

import "fmt"

type MeetingStatus string

const (
	MeetingStatusPending   MeetingStatus = "pending"
	MeetingStatusDraft     MeetingStatus = "draft"
	MeetingStatusSubmitted MeetingStatus = "submitted"
	MeetingStatusDone      MeetingStatus = "done"
)

func (s MeetingStatus) String() string {
	return string(s)
}

func (s MeetingStatus) IsValid() bool {
	switch s {
	case MeetingStatusPending, MeetingStatusDraft, MeetingStatusSubmitted, MeetingStatusDone:
		return true
	default:
		return false
	}
}

func ParseMeetingStatus(s string) (MeetingStatus, error) {
	status := MeetingStatus(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid meeting status: %s", s)
	}
	return status, nil
}

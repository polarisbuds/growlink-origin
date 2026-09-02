package valueobject

import "fmt"

type MeetingStatus int

const (
	MeetingStatusPending   MeetingStatus = iota + 1
	MeetingStatusDraft
	MeetingStatusSubmitted
	MeetingStatusDone
)

func (s MeetingStatus) String() string {
	switch s {
	case MeetingStatusPending:
		return "pending"
	case MeetingStatusDraft:
		return "draft"
	case MeetingStatusSubmitted:
		return "submitted"
	case MeetingStatusDone:
		return "done"
	default:
		return ""
	}
}

func (s MeetingStatus) Int() int {
	return int(s)
}

func (s MeetingStatus) UInt() uint {
	return uint(s)
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
	switch s {
	case MeetingStatusPending.String():
		return MeetingStatusPending, nil
	case MeetingStatusDraft.String():
		return MeetingStatusDraft, nil
	case MeetingStatusSubmitted.String():
		return MeetingStatusSubmitted, nil
	case MeetingStatusDone.String():
		return MeetingStatusDone, nil
	default:
		return 0, fmt.Errorf("invalid meeting status: %s", s)
	}
}

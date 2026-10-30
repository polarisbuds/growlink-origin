package valueobject

import "fmt"

type MeetingStatus int

const (
	MeetingStatusPending MeetingStatus = iota + 1
	MeetingStatusDraft
	MeetingStatusDone
)

func (s MeetingStatus) String() string {
	switch s {
	case MeetingStatusPending:
		return "pending"
	case MeetingStatusDraft:
		return "draft"
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
	case MeetingStatusPending, MeetingStatusDraft, MeetingStatusDone:
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
	case MeetingStatusDone.String():
		return MeetingStatusDone, nil
	default:
		return 0, fmt.Errorf("invalid meeting status: %s", s)
	}
}

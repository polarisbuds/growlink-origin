package valueobject

import "fmt"

type EventType int

const (
	EventTypeLogin         EventType = iota + 1
	EventTypeMeetingDraft
	EventTypeMeetingSubmit
)

func (e EventType) String() string {
	switch e {
	case EventTypeLogin:
		return "login"
	case EventTypeMeetingDraft:
		return "meeting.draft"
	case EventTypeMeetingSubmit:
		return "meeting.submit"
	default:
		return ""
	}
}

func (e EventType) Int() int {
	return int(e)
}

func (e EventType) UInt() uint {
	return uint(e)
}

func (e EventType) IsValid() bool {
	switch e {
	case EventTypeLogin, EventTypeMeetingDraft, EventTypeMeetingSubmit:
		return true
	default:
		return false
	}
}

func ParseEventType(s string) (EventType, error) {
	switch s {
	case EventTypeLogin.String():
		return EventTypeLogin, nil
	case EventTypeMeetingDraft.String():
		return EventTypeMeetingDraft, nil
	case EventTypeMeetingSubmit.String():
		return EventTypeMeetingSubmit, nil
	default:
		return 0, fmt.Errorf("invalid event type: %s", s)
	}
}

package valueobject

import "fmt"

type EventType string

const (
	EventTypeLogin         EventType = "login"
	EventTypeMeetingDraft  EventType = "meeting.draft"
	EventTypeMeetingSubmit EventType = "meeting.submit"
)

func (e EventType) String() string {
	return string(e)
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
	eventType := EventType(s)
	if !eventType.IsValid() {
		return "", fmt.Errorf("invalid event type: %s", s)
	}
	return eventType, nil
}

package api

import (
	"testing"
	"time"

	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
)

func TestCanWriteMeetingAnswers(t *testing.T) {
	impl := time.Now()
	meeting := &entity.Meeting{
		Status: valueobject.MeetingStatusPending,
		Pairing: &entity.Pairing{
			MentorID: 1,
			MenteeID: 2,
		},
		ImplementedAt: &impl,
	}
	if !canWriteMeetingAnswers(meeting, 2) {
		t.Fatal("mentee should write while pending")
	}
	if canWriteMeetingAnswers(meeting, 3) {
		t.Fatal("viewer should not write")
	}
	meeting.Status = valueobject.MeetingStatusDone
	if canWriteMeetingAnswers(meeting, 2) {
		t.Fatal("mentee should not write when done")
	}
}

func TestMeetingViewerRole(t *testing.T) {
	meeting := &entity.Meeting{
		Pairing: &entity.Pairing{MentorID: 10, MenteeID: 20},
	}
	if meetingViewerRole(meeting, 10) != "mentor" {
		t.Fatalf("got %q", meetingViewerRole(meeting, 10))
	}
	if meetingViewerRole(meeting, 20) != "mentee" {
		t.Fatalf("got %q", meetingViewerRole(meeting, 20))
	}
}

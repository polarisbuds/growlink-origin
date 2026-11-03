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
	if canWriteMeetingAnswers(meeting, 1) {
		t.Fatal("mentor should not write while pending")
	}
	if canWriteMeetingAnswers(meeting, 3) {
		t.Fatal("viewer should not write")
	}
	meeting.Status = valueobject.MeetingStatusDraft
	if !canWriteMeetingAnswers(meeting, 2) {
		t.Fatal("mentee should write while draft")
	}
	if canWriteMeetingAnswers(meeting, 1) {
		t.Fatal("mentor should not write while draft")
	}
	meeting.Status = valueobject.MeetingStatusDone
	if canWriteMeetingAnswers(meeting, 2) {
		t.Fatal("mentee should not write when done")
	}
	if !canWriteMeetingAnswers(meeting, 1) {
		t.Fatal("mentor should write when done")
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

func TestCanCancelMeetingSubmit(t *testing.T) {
	meeting := &entity.Meeting{
		Status:  valueobject.MeetingStatusDone,
		Pairing: &entity.Pairing{MentorID: 1, MenteeID: 2},
	}
	if !canCancelMeetingSubmit(meeting, 2) {
		t.Fatal("mentee should cancel submit when done")
	}
	if canCancelMeetingSubmit(meeting, 1) {
		t.Fatal("mentor should not cancel submit")
	}
	meeting.Status = valueobject.MeetingStatusDraft
	if canCancelMeetingSubmit(meeting, 2) {
		t.Fatal("mentee should not cancel when not done")
	}
}

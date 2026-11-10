package api

import (
	"strings"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
)

func enrichAnswersForMeeting(meetingID, userID uint, answers []outputmodel.AnswerOutput) []outputmodel.AnswerOutput {
	out := make([]outputmodel.AnswerOutput, len(answers))
	for i, a := range answers {
		out[i] = a
		if out[i].MeetingID == 0 {
			out[i].MeetingID = meetingID
		}
		out[i].UserID = userID
	}
	return out
}

func sanitizeAnswerOutputs(
	meetingID, userID uint,
	answers []outputmodel.AnswerOutput,
	lookup func(answerID uint) (*outputmodel.AnswerOutput, error),
) ([]outputmodel.AnswerOutput, error) {
	out := enrichAnswersForMeeting(meetingID, userID, answers)
	for i := range out {
		if out[i].ID == 0 {
			continue
		}
		existing, err := lookup(out[i].ID)
		if err != nil {
			return nil, err
		}
		if existing == nil || existing.MeetingID != meetingID || existing.UserID != userID {
			out[i].ID = 0
		}
	}
	return out, nil
}

func meetingViewerRole(meeting *entity.Meeting, viewerID uint) string {
	if meeting == nil || viewerID == 0 || meeting.Pairing == nil {
		return "viewer"
	}
	if meeting.Pairing.MentorID == viewerID {
		return "mentor"
	}
	if meeting.Pairing.MenteeID == viewerID {
		return "mentee"
	}
	return "viewer"
}

func canWriteMeetingAnswers(meeting *entity.Meeting, viewerID uint) bool {
	if meeting == nil || viewerID == 0 {
		return false
	}
	role := meetingViewerRole(meeting, viewerID)
	switch role {
	case "mentee":
		return meeting.Status == valueobject.MeetingStatusPending ||
			meeting.Status == valueobject.MeetingStatusDraft
	case "mentor":
		return meeting.Status == valueobject.MeetingStatusSubmitted
	default:
		return false
	}
}

func canViewMeeting(meeting *entity.Meeting, viewerID uint) bool {
	if meeting == nil || viewerID == 0 {
		return false
	}
	switch meetingViewerRole(meeting, viewerID) {
	case "mentor", "mentee":
		return true
	}
	if meeting.Status != valueobject.MeetingStatusSubmitted &&
		meeting.Status != valueobject.MeetingStatusDone {
		return false
	}
	for _, s := range meeting.Sharers {
		if s.UserID == viewerID {
			return true
		}
	}
	return false
}

func canCancelMeetingSubmit(meeting *entity.Meeting, viewerID uint) bool {
	if meeting == nil || viewerID == 0 {
		return false
	}
	if meeting.Status != valueobject.MeetingStatusSubmitted {
		return false
	}
	return meetingViewerRole(meeting, viewerID) == "mentee"
}

func canEditMeetingSharers(meeting *entity.Meeting, viewerID uint) bool {
	if meeting == nil || viewerID == 0 {
		return false
	}
	role := meetingViewerRole(meeting, viewerID)
	return role == "mentor" || role == "mentee"
}

func meetingSharersToOutputs(meeting *entity.Meeting) []outputmodel.UserOutput {
	if meeting == nil || len(meeting.Sharers) == 0 {
		return nil
	}
	out := make([]outputmodel.UserOutput, 0, len(meeting.Sharers))
	for _, s := range meeting.Sharers {
		if s.User != nil {
			out = append(out, converter.UserEntityToOutput(*s.User))
		} else if s.UserID != 0 {
			out = append(out, outputmodel.UserOutput{ID: s.UserID})
		}
	}
	return out
}

func filterMeetingComments(meeting *entity.Meeting, viewerID uint) []outputmodel.MeetingCommentOutput {
	if meeting == nil || len(meeting.MeetingComments) == 0 {
		return nil
	}
	canViewLimited := canViewLimitedMeetingComments(meeting, viewerID)
	filtered := make([]entity.MeetingComment, 0, len(meeting.MeetingComments))
	for _, c := range meeting.MeetingComments {
		if c.IsLimited && !canViewLimited {
			continue
		}
		filtered = append(filtered, c)
	}
	return converter.MeetingCommentEntitiesToOutputs(filtered)
}

func canViewLimitedMeetingComments(meeting *entity.Meeting, viewerID uint) bool {
	if viewerID == 0 {
		return false
	}
	if meeting.CreatedByUserID == viewerID {
		return true
	}
	if meeting.Pairing == nil {
		return false
	}
	if meeting.Pairing.MentorID == viewerID || meeting.Pairing.MenteeID == viewerID {
		return true
	}
	return false
}

func userMatchesSharerQuery(u entity.User, q string) bool {
	if q == "" {
		return true
	}
	q = strings.ToLower(q)
	if strings.Contains(strings.ToLower(u.Name), q) {
		return true
	}
	if strings.Contains(strings.ToLower(u.Email), q) {
		return true
	}
	return false
}

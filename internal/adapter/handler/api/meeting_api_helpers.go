package api

import (
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func enrichAnswersForMeeting(meetingID, userID uint, answers []outputmodel.AnswerOutput) []outputmodel.AnswerOutput {
	out := make([]outputmodel.AnswerOutput, len(answers))
	for i, a := range answers {
		out[i] = a
		if out[i].MeetingID == 0 {
			out[i].MeetingID = meetingID
		}
		if out[i].UserID == 0 {
			out[i].UserID = userID
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

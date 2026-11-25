package meeting_usecase

import (
	"fmt"
	"strings"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"
)

func enrichAnswersForMeeting(meetingID, userID uint, answers []outputmodel.AnswerOutput, meeting *entity.Meeting) ([]outputmodel.AnswerOutput, error) {
	out := make([]outputmodel.AnswerOutput, len(answers))
	topicsByUUID := map[string]uint{}
	if meeting != nil && meeting.Template != nil {
		for _, topic := range meeting.Template.Topics {
			topicsByUUID[topic.UUID.String()] = topic.ID
		}
	}
	for i, a := range answers {
		out[i] = a
		if out[i].MeetingID == 0 {
			out[i].MeetingID = meetingID
		}
		if out[i].TopicID == 0 {
			topicID, ok := topicsByUUID[a.Topic.UUID]
			if !ok {
				return nil, fmt.Errorf("トピックIDが見つかりませんでした")
			}
			out[i].TopicID = topicID
		}
		out[i].UserID = userID
	}
	return out, nil
}

func validateRequiredAnswers(meeting *entity.Meeting, answers []outputmodel.AnswerOutput) error {
	if meeting == nil || meeting.Template == nil {
		return nil
	}
	contentByTopicUUID := make(map[string]string, len(answers))
	for _, a := range answers {
		topicUUID := a.Topic.UUID
		if topicUUID == "" {
			continue
		}
		contentByTopicUUID[topicUUID] = a.Content
	}
	for _, topic := range meeting.Template.Topics {
		if !topic.IsRequire {
			continue
		}
		content, ok := contentByTopicUUID[topic.UUID.String()]
		if !ok || content == "" {
			return ErrRequireAnswerNotAnswered
		}
	}
	return nil
}

func sanitizeAnswerOutputs(
	meetingID, userID uint,
	answers []outputmodel.AnswerOutput,
	meeting *entity.Meeting,
	lookup func(answerID uint) (*outputmodel.AnswerOutput, error),
) ([]outputmodel.AnswerOutput, error) {
	out, err := enrichAnswersForMeeting(meetingID, userID, answers, meeting)
	if err != nil {
		return nil, fmt.Errorf("回答内容の埋め込みに失敗しました: %w", err)
	}
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

func MeetingViewerRole(meeting *entity.Meeting, viewerID uint) string {
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

func CanWriteMeetingAnswers(meeting *entity.Meeting, viewerID uint) bool {
	if meeting == nil || viewerID == 0 {
		return false
	}
	if MeetingViewerRole(meeting, viewerID) != "mentee" {
		return false
	}
	return meeting.Status == valueobject.MeetingStatusPending ||
		meeting.Status == valueobject.MeetingStatusDraft
}

func CanWriteTopicComments(meeting *entity.Meeting, viewerID uint) bool {
	if meeting == nil || viewerID == 0 {
		return false
	}
	if MeetingViewerRole(meeting, viewerID) != "mentor" {
		return false
	}
	return meeting.Status == valueobject.MeetingStatusSubmitted
}

func CanViewMeeting(meeting *entity.Meeting, viewerID uint) bool {
	if meeting == nil || viewerID == 0 {
		return false
	}
	switch MeetingViewerRole(meeting, viewerID) {
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

func CanCancelMeetingSubmit(meeting *entity.Meeting, viewerID uint) bool {
	if meeting == nil || viewerID == 0 {
		return false
	}
	if meeting.Status != valueobject.MeetingStatusSubmitted {
		return false
	}
	return MeetingViewerRole(meeting, viewerID) == "mentee"
}

func CanEditMeetingSharers(meeting *entity.Meeting, viewerID uint) bool {
	if meeting == nil || viewerID == 0 {
		return false
	}
	role := MeetingViewerRole(meeting, viewerID)
	return role == "mentor" || role == "mentee"
}

func CanEditMeetingImplementedDate(meeting *entity.Meeting, viewerID uint, systemRole valueobject.SystemRole) bool {
	if meeting == nil || viewerID == 0 {
		return false
	}
	if CanEditMeetingSharers(meeting, viewerID) {
		return true
	}
	if systemRole == valueobject.SystemRoleAdmin {
		return true
	}
	if meeting.CreatedByUserID == 0 {
		return false
	}
	return viewerID == meeting.CreatedByUserID
}

func MeetingSharersToOutputs(meeting *entity.Meeting) []outputmodel.UserOutput {
	if meeting == nil || len(meeting.Sharers) == 0 {
		return nil
	}
	out := make([]outputmodel.UserOutput, len(meeting.Sharers))
	for idx, s := range meeting.Sharers {
		if s.User != nil {
			out[idx] = converter.UserEntityToOutput(*s.User)
		} else if s.UserID != 0 {
			out[idx] = outputmodel.UserOutput{ID: s.UserID}
		}
	}
	return out
}

func FilterMeetingComments(meeting *entity.Meeting, viewerID uint) []outputmodel.MeetingCommentOutput {
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
	for _, meetingComment := range meeting.MeetingComments {
		if meetingComment.UserID == viewerID {
			return true
		}
	}
	return false
}

func UserMatchesSharerQuery(u entity.User, q string) bool {
	if q == "" {
		return true
	}
	q = strings.ToLower(q)
	if strings.Contains(strings.ToLower(u.Name), q) {
		return true
	}
	return strings.Contains(strings.ToLower(u.Email), q)
}

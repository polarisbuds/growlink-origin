package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func MeetingCommentInputToEntity(in inputmodel.MeetingCommentInput) *entity.MeetingComment {
	return &entity.MeetingComment{
		Content:   in.Content,
		MeetingID: in.MeetingID,
		UserID:    in.UserID,
	}
}

func MeetingCommentEntityToOutput(e entity.MeetingComment) outputmodel.MeetingCommentOutput {
	out := outputmodel.MeetingCommentOutput{
		ID:        e.ID,
		Content:   e.Content,
		MeetingID: e.MeetingID,
		UserID:    e.UserID,
	}
	if e.Meeting != nil {
		out.Meeting = MeetingEntityToOutput(*e.Meeting)
	} else if e.MeetingID != 0 {
		out.Meeting = outputmodel.MeetingOutput{ID: e.MeetingID}
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	}
	return out
}

func MeetingCommentEntitiesToOutputs(es []entity.MeetingComment) []outputmodel.MeetingCommentOutput {
	outputs := make([]outputmodel.MeetingCommentOutput, len(es))
	for i, e := range es {
		outputs[i] = MeetingCommentEntityToOutput(e)
	}
	return outputs
}

package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func ReadCommentableInputToEntity(in inputmodel.ReadCommentableInput) *entity.ReadCommentable {
	return &entity.ReadCommentable{
		MeetingID: in.MeetingID,
		UserID:    in.UserID,
	}
}

func ReadCommentableEntityToOutput(e entity.ReadCommentable) outputmodel.ReadCommentableOutput {
	out := outputmodel.ReadCommentableOutput{
		ID:        e.ID,
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

func ReadCommentableEntitiesToOutputs(es []entity.ReadCommentable) []outputmodel.ReadCommentableOutput {
	outputs := make([]outputmodel.ReadCommentableOutput, len(es))
	for i, e := range es {
		outputs[i] = ReadCommentableEntityToOutput(e)
	}
	return outputs
}

package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func AnswerInputToEntity(in inputmodel.AnswerInput) *entity.Answer {
	return &entity.Answer{
		Content:   in.Content,
		MeetingID: in.MeetingID,
		TopicID:   in.TopicID,
		UserID:    in.UserID,
	}
}

func answerEntityToOutputShallow(e entity.Answer) outputmodel.AnswerOutput {
	out := outputmodel.AnswerOutput{
		ID:        e.ID,
		Content:   e.Content,
		MeetingID: e.MeetingID,
		TopicID:   e.TopicID,
		UserID:    e.UserID,
	}
	if e.Meeting != nil {
		out.Meeting = MeetingEntityToOutput(*e.Meeting)
	} else if e.MeetingID != 0 {
		out.Meeting = outputmodel.MeetingOutput{ID: e.MeetingID}
	}
	if e.Topic != nil {
		out.Topic = TopicEntityToOutput(*e.Topic)
	} else if e.TopicID != 0 {
		out.Topic = outputmodel.TopicOutput{ID: e.TopicID}
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	}
	return out
}

func AnswerEntityToOutput(e entity.Answer) outputmodel.AnswerOutput {
	out := answerEntityToOutputShallow(e)
	if len(e.AnswerComments) > 0 {
		out.AnswerComments = AnswerCommentEntitiesToOutputs(e.AnswerComments)
	}
	return out
}

func AnswerEntitiesToOutputs(es []entity.Answer) []outputmodel.AnswerOutput {
	outputs := make([]outputmodel.AnswerOutput, len(es))
	for i, e := range es {
		outputs[i] = AnswerEntityToOutput(e)
	}
	return outputs
}

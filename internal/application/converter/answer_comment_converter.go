package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func AnswerCommentInputToEntity(in inputmodel.AnswerCommentInput) *entity.AnswerComment {
	return &entity.AnswerComment{
		Content:  in.Content,
		AnswerID: in.AnswerID,
		UserID:   in.UserID,
	}
}

func AnswerCommentEntityToOutput(e entity.AnswerComment) outputmodel.AnswerCommentOutput {
	out := outputmodel.AnswerCommentOutput{
		ID:       e.ID,
		Content:  e.Content,
		AnswerID: e.AnswerID,
		UserID:   e.UserID,
	}
	if e.Answer != nil {
		out.Answer = answerEntityToOutputShallow(*e.Answer)
	} else if e.AnswerID != 0 {
		out.Answer = outputmodel.AnswerOutput{ID: e.AnswerID}
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	}
	return out
}

func AnswerCommentEntitiesToOutputs(es []entity.AnswerComment) []outputmodel.AnswerCommentOutput {
	outputs := make([]outputmodel.AnswerCommentOutput, len(es))
	for i, e := range es {
		outputs[i] = AnswerCommentEntityToOutput(e)
	}
	return outputs
}

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

func AnswerCommentEntityToOutput(e *entity.AnswerComment) *outputmodel.AnswerCommentOutput {
	if e == nil {
		return nil
	}

	out := &outputmodel.AnswerCommentOutput{
		Content: e.Content,
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

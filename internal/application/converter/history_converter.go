package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func HistoryInputToEntity(in inputmodel.HistoryInput) *entity.History {
	return &entity.History{
		UserID:    in.UserID,
		EventType: in.EventType,
	}
}

func HistoryEntityToOutput(e entity.History) outputmodel.HistoryOutput {
	out := outputmodel.HistoryOutput{
		ID:        e.ID,
		UserID:    e.UserID,
		EventType: e.EventType,
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	}
	return out
}

func HistoryEntitiesToOutputs(es []entity.History) []outputmodel.HistoryOutput {
	outputs := make([]outputmodel.HistoryOutput, len(es))
	for i, e := range es {
		outputs[i] = HistoryEntityToOutput(e)
	}
	return outputs
}

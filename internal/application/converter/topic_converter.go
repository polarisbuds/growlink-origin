package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func TopicInputToEntity(in inputmodel.TopicInput) *entity.Topic {
	return &entity.Topic{
		DisplayNo: in.DisplayNo,
		Question:  in.Question,
		IsRequire: in.IsRequire,
	}
}

func TopicEntityToOutput(e entity.Topic) outputmodel.TopicOutput {
	return outputmodel.TopicOutput{
		ID:        e.ID,
		DisplayNo: e.DisplayNo,
		Question:  e.Question,
		IsRequire: e.IsRequire,
	}
}

func TopicEntitiesToOutputs(es []entity.Topic) []outputmodel.TopicOutput {
	outputs := make([]outputmodel.TopicOutput, len(es))
	for i, e := range es {
		outputs[i] = TopicEntityToOutput(e)
	}
	return outputs
}

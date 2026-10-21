package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
)

func TemplateInputToEntity(in inputmodel.TemplateInput) *entity.Template {
	topics := make([]entity.Topic, len(in.Topics))
	for i, topic := range in.Topics {
		topics[i] = *TopicInputToEntity(topic)
	}
	return &entity.Template{
		Title:  in.Title,
		Topics: topics,
	}
}

func TemplateEntityToOutput(e entity.Template) outputmodel.TemplateOutput {
	return outputmodel.TemplateOutput{
		ID:        e.ID,
		UUID:      e.UUID.String(),
		Title:     e.Title,
		DisplayNo: e.DisplayNo,
		IsSample:  e.IsSample,
		Topics:    TopicEntitiesToOutputs(e.Topics),
	}
}

func TemplateEntitiesToOutputs(es []entity.Template) []outputmodel.TemplateOutput {
	outputs := make([]outputmodel.TemplateOutput, len(es))
	for i, e := range es {
		outputs[i] = TemplateEntityToOutput(e)
	}
	return outputs
}

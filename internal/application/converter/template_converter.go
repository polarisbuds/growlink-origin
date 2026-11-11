package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	"time"
)

func TemplateInputToEntity(in inputmodel.TemplateInput) *entity.Template {
	topics := make([]entity.Topic, len(in.Topics))
	for i, topic := range in.Topics {
		topics[i] = *TopicInputToEntity(topic)
	}
	return &entity.Template{
		Title:       in.Title,
		Description: in.Description,
		Topics:      topics,
	}
}

func TemplateEntityToOutput(e entity.Template) outputmodel.TemplateOutput {
	creatorName := ""
	if e.User != nil {
		creatorName = e.User.Name
	}
	createdAt := ""
	if !e.CreatedAt.IsZero() {
		createdAt = e.CreatedAt.UTC().Format(time.RFC3339)
	}
	return outputmodel.TemplateOutput{
		ID:          e.ID,
		UUID:        e.UUID.String(),
		Title:       e.Title,
		Description: e.Description,
		IsSample:    e.IsSample,
		IsLock:      e.IsLock,
		IsDefault:   e.IsDefault,
		CreatorName: creatorName,
		CreatedAt:   createdAt,
		Topics:      TopicEntitiesToOutputs(e.Topics),
	}
}

func TemplateEntitiesToOutputs(es []entity.Template) []outputmodel.TemplateOutput {
	outputs := make([]outputmodel.TemplateOutput, len(es))
	for i, e := range es {
		outputs[i] = TemplateEntityToOutput(e)
	}
	return outputs
}

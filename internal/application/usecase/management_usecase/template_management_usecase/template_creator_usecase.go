package template_mng_usecase

import (
	"fmt"
	"growlink/internal/application/converter"
	inputmodel "growlink/internal/application/model/input_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
)

type TemplateCreator struct {
	templateRepo repository_interface.TemplateRepository
	topicRepo    repository_interface.TopicRepository
}

func NewTemplateCreator(templateRepo repository_interface.TemplateRepository,
	topicRepo repository_interface.TopicRepository) *TemplateCreator {
	return &TemplateCreator{
		templateRepo: templateRepo,
		topicRepo:    topicRepo,
	}
}

func (uc *TemplateCreator) CreateTemplate(input inputmodel.TemplateInput) error {
	topics := []entity.Topic{}
	for _, inputTopic := range input.Topics {
		err := uc.topicRepo.Create(converter.TopicInputToEntity(inputTopic))
		if err != nil {
			return fmt.Errorf("トピックの作成に失敗しました: %w", err)
		}
		topic := converter.TopicInputToEntity(inputTopic)
		topics = append(topics, *topic)
	}

	if err := uc.templateRepo.Create(&entity.Template{
		Title:  input.Title,
		Topics: topics,
	}); err != nil {
		return fmt.Errorf("テンプレートの新規作成に失敗しました: %w", err)
	}

	return nil
}

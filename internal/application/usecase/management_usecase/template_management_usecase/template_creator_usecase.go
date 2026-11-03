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

func (uc *TemplateCreator) CreateTemplate(input inputmodel.TemplateInput, userID uint) error {
	input.Topics = normalizeTopicInputs(input.Topics)
	if err := validateTemplateInput(input); err != nil {
		return err
	}

	recordCount, err := uc.templateRepo.Count()
	if err != nil {
		return fmt.Errorf("テンプレート数のカウントに失敗しました: %w", err)
	}

	template := &entity.Template{
		Title:       input.Title,
		Description: input.Description,
		UserID:      userID,
		DisplayNo:   int(recordCount + 1),
	}
	if err := uc.templateRepo.Create(template); err != nil {
		return fmt.Errorf("テンプレートの新規作成に失敗しました: %w", err)
	}

	for _, inputTopic := range input.Topics {
		topic := converter.TopicInputToEntity(inputTopic)
		topic.TemplateID = template.ID
		if err := uc.topicRepo.Create(topic); err != nil {
			return fmt.Errorf("トピックの作成に失敗しました: %w", err)
		}
	}
	return nil
}

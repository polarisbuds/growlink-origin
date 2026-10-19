package template_mng_usecase

import (
	"fmt"
	"growlink/internal/application/converter"
	inputmodel "growlink/internal/application/model/input_model"
	repository_interface "growlink/internal/domain/interface/repository"
)

type TemplateUpdater struct {
	templateRepo repository_interface.TemplateRepository
	topicRepo    repository_interface.TopicRepository
}

func NewTemplateUpdater(templateRepo repository_interface.TemplateRepository, topicRepo repository_interface.TopicRepository) *TemplateUpdater {
	return &TemplateUpdater{
		templateRepo: templateRepo,
		topicRepo:    topicRepo,
	}
}

func (uc *TemplateUpdater) UpdateDisplayNo(templateID uint, displayNo int) error {
	if err := uc.templateRepo.UpdateDisplayNo(templateID, displayNo); err != nil {
		return fmt.Errorf("テンプレートの並び順更新に失敗しました: %w", err)
	}
	return nil
}

func (uc *TemplateUpdater) UpdateTemplate(templateID uint, input inputmodel.TemplateInput) error {
	newTemplate := converter.TemplateInputToEntity(input)
	if err := uc.templateRepo.Update(templateID, newTemplate); err != nil {
		return fmt.Errorf("テンプレートの情報更新に失敗しました: %w", err)
	}
	return nil
}

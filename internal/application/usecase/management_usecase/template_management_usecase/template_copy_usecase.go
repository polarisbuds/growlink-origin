package template_mng_usecase

import (
	"fmt"

	inputmodel "growlink/internal/application/model/input_model"
	repository_interface "growlink/internal/domain/interface/repository"
)

type TemplateCopier struct {
	templateRepo repository_interface.TemplateRepository
	topicRepo    repository_interface.TopicRepository
	creator      *TemplateCreator
}

func NewTemplateCopier(templateRepo repository_interface.TemplateRepository, topicRepo repository_interface.TopicRepository) *TemplateCopier {
	return &TemplateCopier{
		templateRepo: templateRepo,
		topicRepo:    topicRepo,
		creator:      NewTemplateCreator(templateRepo, topicRepo),
	}
}

func (uc *TemplateCopier) CopyTemplate(sourceID uint, newTitle string) error {
	src, err := uc.templateRepo.Get(sourceID)
	if err != nil {
		return fmt.Errorf("テンプレートの取得に失敗しました: %w", err)
	}
	in := inputmodel.TemplateInput{Title: newTitle}
	for _, t := range src.Topics {
		in.Topics = append(in.Topics, inputmodel.TopicInput{
			DisplayNo: t.DisplayNo,
			Question:  t.Question,
			IsRequire: t.IsRequire,
		})
	}
	return uc.creator.CreateTemplate(in)
}

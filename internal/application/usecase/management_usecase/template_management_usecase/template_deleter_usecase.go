package template_mng_usecase

import (
	"context"
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"
)

type TemplateDeleter struct {
	templateRepo repository_interface.TemplateRepository
	topicRepo    repository_interface.TopicRepository
}

func NewTemplateDeleter(templateRepo repository_interface.TemplateRepository, topicRepo repository_interface.TopicRepository) *TemplateDeleter {
	return &TemplateDeleter{
		templateRepo: templateRepo,
		topicRepo:    topicRepo,
	}
}

func (uc *TemplateDeleter) DeleteTemplate(templateID uint) (err error) {
	const op = "template_mng_usecase.DeleteTemplate"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	tpl, err := uc.templateRepo.Get(templateID)
	if err != nil {
		return fmt.Errorf("テンプレートの取得に失敗しました: %w", err)
	}
	if tpl.IsLock {
		return ErrTemplateLocked
	}
	if err = uc.templateRepo.Delete(templateID); err != nil {
		return fmt.Errorf("テンプレートの削除処理に失敗しました: %w", err)
	}
	return nil
}

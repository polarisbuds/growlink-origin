package template_mng_usecase

import (
	"context"
	"fmt"
	"growlink/internal/application/converter"
	inputmodel "growlink/internal/application/model/input_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"

	"github.com/google/uuid"
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

func (uc *TemplateUpdater) UpdateTemplate(templateID uint, input inputmodel.TemplateInput) (err error) {
	const op = "template_mng_usecase.UpdateTemplate"
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
	input.Topics = normalizeTopicInputs(input.Topics)
	if err = validateTemplateInput(input); err != nil {
		return err
	}
	newTemplate := converter.TemplateInputToEntity(input)
	if err = ensureTopicsBelongToTemplate(tpl.Topics, newTemplate.Topics); err != nil {
		return err
	}
	if err = uc.templateRepo.Update(templateID, newTemplate); err != nil {
		return fmt.Errorf("テンプレートの情報更新に失敗しました: %w", err)
	}
	return nil
}

func ensureTopicsBelongToTemplate(existing []entity.Topic, incoming []entity.Topic) error {
	known := make(map[uuid.UUID]struct{}, len(existing))
	for _, topic := range existing {
		known[topic.UUID] = struct{}{}
	}
	for i, topic := range incoming {
		if topic.UUID == uuid.Nil {
			continue
		}
		if _, ok := known[topic.UUID]; !ok {
			return fmt.Errorf("%w: Q%d は存在しない質問です", ErrInvalidTemplateInput, i+1)
		}
	}
	return nil
}

func (uc *TemplateUpdater) UpdateTemplateLock(templateID uint, isLock bool) (err error) {
	const op = "template_mng_usecase.UpdateTemplateLock"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.templateRepo.UpdateIsLock(templateID, isLock); err != nil {
		return fmt.Errorf("テンプレートのロック状態の更新に失敗しました: %w", err)
	}
	return nil
}

func (uc *TemplateUpdater) UpdateTemplateDefault(templateID uint, isDefault bool) (err error) {
	const op = "template_mng_usecase.UpdateTemplateDefault"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.templateRepo.UpdateIsDefault(templateID, isDefault); err != nil {
		return fmt.Errorf("テンプレートの初期フラグの更新に失敗しました: %w", err)
	}
	return nil
}

package template_mng_usecase

import (
	"context"
	"fmt"
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"

	"github.com/google/uuid"
)

type TemplateViewer struct {
	templateRepo repository_interface.TemplateRepository
}

func NewTemplateViewer(templateRepo repository_interface.TemplateRepository) *TemplateViewer {
	return &TemplateViewer{
		templateRepo: templateRepo,
	}
}

func (uc *TemplateViewer) ListTemplate() (ret0 []outputmodel.TemplateOutput, err error) {
	const op = "template_mng_usecase.ListTemplate"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	tempaltes, err := uc.templateRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("テンプレートの全取得に失敗しました: %w", err)
	}
	return converter.TemplateEntitiesToOutputs(tempaltes), nil
}

func (uc *TemplateViewer) GetTemplateByUUID(templateUUID uuid.UUID) (ret0 *outputmodel.TemplateOutput, err error) {
	const op = "template_mng_usecase.GetTemplateByUUID"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	template, err := uc.templateRepo.GetByUUID(templateUUID)
	if err != nil {
		return nil, fmt.Errorf("テンプレートのUUID指定取得に失敗しました: %w", err)
	}
	out := converter.TemplateEntityToOutput(*template)
	return &out, nil
}

func (uc *TemplateViewer) GetTemplate(templateID uint) (ret0 *outputmodel.TemplateOutput, err error) {
	const op = "template_mng_usecase.GetTemplate"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	template, err := uc.templateRepo.Get(templateID)
	if err != nil {
		return nil, fmt.Errorf("テンプレートのID指定取得に失敗しました: %w", err)
	}
	out := converter.TemplateEntityToOutput(*template)
	return &out, nil
}

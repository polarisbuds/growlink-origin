package template_mng_usecase



import (
	"fmt"
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
)



type TemplateViewer struct {

	templateRepo repository_interface.TemplateRepository

}



func NewTemplateViewer(templateRepo repository_interface.TemplateRepository) *TemplateViewer {

	return &TemplateViewer{

		templateRepo: templateRepo,

	}

}



func (uc *TemplateViewer) ListTemplate() ([]outputmodel.TemplateOutput, error) {

	tempaltes, err := uc.templateRepo.GetAll()

	if err != nil {

		return nil, fmt.Errorf("テンプレートの全取得に失敗しました: %w", err)

	}

	return converter.TemplateEntitiesToOutputs(tempaltes), nil

}



func (uc *TemplateViewer) GetTemplate(templateID uint) (*outputmodel.TemplateOutput, error) {

	template, err := uc.templateRepo.Get(templateID)

	if err != nil {

		return nil, fmt.Errorf("テンプレートのID指定取得に失敗しました: %w", err)

	}

	out := converter.TemplateEntityToOutput(*template)

	return &out, nil

}


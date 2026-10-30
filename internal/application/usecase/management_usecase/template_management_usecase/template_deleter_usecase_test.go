package template_mng_usecase_test

import (
	"testing"

	template_mng_usecase "growlink/internal/application/usecase/management_usecase/template_management_usecase"
	repomock "growlink/internal/domain/interface/repository/mock"

	"go.uber.org/mock/gomock"
)

func newTemplateDeleter(t *testing.T) (
	*template_mng_usecase.TemplateDeleter,
	*repomock.MockTemplateRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	templateRepo := repomock.NewMockTemplateRepository(ctrl)
	topicRepo := repomock.NewMockTopicRepository(ctrl)
	return template_mng_usecase.NewTemplateDeleter(templateRepo, topicRepo), templateRepo
}

func TestTemplateDeleter_DeleteTemplate(t *testing.T) {
	tests := []struct {
		name        string
		templateID  uint
		wantErr     bool
		errContains string
		setup       func(repo *repomock.MockTemplateRepository)
	}{
		{
			name:       "通常系_削除",
			templateID: validTemplateID,
			wantErr:    false,
			setup: func(repo *repomock.MockTemplateRepository) {
				repo.EXPECT().Delete(validTemplateID).Return(nil)
			},
		},
		{
			name:        "異常系_削除失敗",
			templateID:  validTemplateID,
			wantErr:     true,
			errContains: "テンプレートの削除処理に失敗しました",
			setup: func(repo *repomock.MockTemplateRepository) {
				repo.EXPECT().Delete(validTemplateID).Return(errDeleteTemplateFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo := newTemplateDeleter(t)
			if tt.setup != nil {
				tt.setup(repo)
			}

			gotErr := uc.DeleteTemplate(tt.templateID)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "DeleteTemplate")
		})
	}
}

package template_mng_usecase_test

import (
	"testing"

	template_mng_usecase "growlink/internal/application/usecase/management_usecase/template_management_usecase"
	inputmodel "growlink/internal/application/model/input_model"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type templateUpdaterCaptures struct {
	templateID   uint
	newTemplate  *entity.Template
}

func newTemplateUpdater(t *testing.T) (
	*template_mng_usecase.TemplateUpdater,
	*repomock.MockTemplateRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	templateRepo := repomock.NewMockTemplateRepository(ctrl)
	topicRepo := repomock.NewMockTopicRepository(ctrl)
	return template_mng_usecase.NewTemplateUpdater(templateRepo, topicRepo), templateRepo
}

func TestTemplateUpdater_UpdateTemplate(t *testing.T) {
	input := inputmodel.TemplateInput{
		Title: "updated-template",
		Topics: []inputmodel.TopicInput{
			{DisplayNo: 1, Question: "updated-question", IsRequire: true},
		},
	}

	tests := []struct {
		name        string
		templateID  uint
		input       inputmodel.TemplateInput
		wantErr     bool
		errContains string
		setup       func(repo *repomock.MockTemplateRepository, cap *templateUpdaterCaptures)
		assertFn    func(t *testing.T, cap *templateUpdaterCaptures)
	}{
		{
			name:       "通常系_更新",
			templateID: validTemplateID,
			input:      input,
			wantErr:    false,
			setup: func(repo *repomock.MockTemplateRepository, cap *templateUpdaterCaptures) {
				repo.EXPECT().Update(validTemplateID, gomock.Any()).DoAndReturn(func(id uint, template *entity.Template) error {
					cap.templateID = id
					cap.newTemplate = template
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *templateUpdaterCaptures) {
				assert.Equal(t, validTemplateID, cap.templateID)
				require.NotNil(t, cap.newTemplate)
				assert.Equal(t, "updated-template", cap.newTemplate.Title)
				require.Len(t, cap.newTemplate.Topics, 1)
				assert.Equal(t, "updated-question", cap.newTemplate.Topics[0].Question)
				assert.True(t, cap.newTemplate.Topics[0].IsRequire)
			},
		},
		{
			name:        "異常系_更新失敗",
			templateID:  validTemplateID,
			input:       input,
			wantErr:     true,
			errContains: "テンプレートの情報更新に失敗しました",
			setup: func(repo *repomock.MockTemplateRepository, _ *templateUpdaterCaptures) {
				repo.EXPECT().Update(validTemplateID, gomock.Any()).Return(errUpdateTemplateFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo := newTemplateUpdater(t)
			cap := &templateUpdaterCaptures{}
			if tt.setup != nil {
				tt.setup(repo, cap)
			}

			gotErr := uc.UpdateTemplate(tt.templateID, tt.input)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "UpdateTemplate")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}

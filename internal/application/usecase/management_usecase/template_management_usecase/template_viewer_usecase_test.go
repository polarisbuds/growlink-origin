package template_mng_usecase_test

import (
	"testing"

	template_mng_usecase "growlink/internal/application/usecase/management_usecase/template_management_usecase"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newTemplateViewer(t *testing.T) (
	*template_mng_usecase.TemplateViewer,
	*repomock.MockTemplateRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	templateRepo := repomock.NewMockTemplateRepository(ctrl)
	return template_mng_usecase.NewTemplateViewer(templateRepo), templateRepo
}

func TestTemplateViewer_ListTemplate(t *testing.T) {
	tests := []struct {
		name        string
		wantErr     bool
		errContains string
		setup       func(repo *repomock.MockTemplateRepository)
		assertFn    func(t *testing.T, got []outputmodel.TemplateOutput)
	}{
		{
			name:    "通常系_複数テンプレート取得",
			wantErr: false,
			setup: func(repo *repomock.MockTemplateRepository) {
				repo.EXPECT().GetAll().Return([]entity.Template{
					newTestTemplate(validTemplateID, "template-1", newTestTopic(topicID1, 1, "q1", true)),
					newTestTemplate(validTemplateID+1, "template-2"),
				}, nil)
			},
			assertFn: func(t *testing.T, got []outputmodel.TemplateOutput) {
				require.Len(t, got, 2)
				assert.Equal(t, validTemplateID, got[0].ID)
				assert.Equal(t, "template-1", got[0].Title)
				require.Len(t, got[0].Topics, 1)
				assert.Equal(t, "q1", got[0].Topics[0].Question)
				assert.Equal(t, validTemplateID+1, got[1].ID)
				assert.Equal(t, "template-2", got[1].Title)
			},
		},
		{
			name:    "通常系_空配列",
			wantErr: false,
			setup: func(repo *repomock.MockTemplateRepository) {
				repo.EXPECT().GetAll().Return([]entity.Template{}, nil)
			},
			assertFn: func(t *testing.T, got []outputmodel.TemplateOutput) {
				assert.Empty(t, got)
			},
		},
		{
			name:        "異常系_全取得失敗",
			wantErr:     true,
			errContains: "テンプレートの全取得に失敗しました",
			setup: func(repo *repomock.MockTemplateRepository) {
				repo.EXPECT().GetAll().Return(nil, errGetTemplatesFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo := newTemplateViewer(t)
			if tt.setup != nil {
				tt.setup(repo)
			}

			got, gotErr := uc.ListTemplate()
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "ListTemplate")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, got)
			}
		})
	}
}

func TestTemplateViewer_GetTemplate(t *testing.T) {
	tests := []struct {
		name        string
		templateID  uint
		wantErr     bool
		errContains string
		setup       func(repo *repomock.MockTemplateRepository)
		assertFn    func(t *testing.T, got *outputmodel.TemplateOutput)
	}{
		{
			name:       "通常系_ID指定取得",
			templateID: validTemplateID,
			wantErr:    false,
			setup: func(repo *repomock.MockTemplateRepository) {
				template := newTestTemplate(validTemplateID, "template-1", newTestTopic(topicID1, 1, "q1", true))
				repo.EXPECT().Get(validTemplateID).Return(&template, nil)
			},
			assertFn: func(t *testing.T, got *outputmodel.TemplateOutput) {
				require.NotNil(t, got)
				assert.Equal(t, validTemplateID, got.ID)
				assert.Equal(t, "template-1", got.Title)
				require.Len(t, got.Topics, 1)
				assert.Equal(t, topicID1, got.Topics[0].ID)
				assert.Equal(t, "q1", got.Topics[0].Question)
			},
		},
		{
			name:        "異常系_ID指定取得失敗",
			templateID:  validTemplateID,
			wantErr:     true,
			errContains: "テンプレートのID指定取得に失敗しました",
			setup: func(repo *repomock.MockTemplateRepository) {
				repo.EXPECT().Get(validTemplateID).Return(nil, errGetTemplateFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo := newTemplateViewer(t)
			if tt.setup != nil {
				tt.setup(repo)
			}

			got, gotErr := uc.GetTemplate(tt.templateID)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "GetTemplate")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, got)
			}
		})
	}
}

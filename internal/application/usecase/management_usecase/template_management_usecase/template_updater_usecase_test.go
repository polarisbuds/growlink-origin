package template_mng_usecase_test

import (
	"testing"

	inputmodel "growlink/internal/application/model/input_model"
	template_mng_usecase "growlink/internal/application/usecase/management_usecase/template_management_usecase"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

type templateUpdaterCaptures struct {
	templateID  uint
	newTemplate *entity.Template
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
	existingTopicUUID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	input := inputmodel.TemplateInput{
		Title:       "updated-template",
		Description: "updated description",
		Topics: []inputmodel.TopicInput{
			{UUID: existingTopicUUID.String(), DisplayNo: 1, Question: "updated-question", IsRequire: true},
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
			name:        "異常系_ロック中",
			templateID:  validTemplateID,
			input:       input,
			wantErr:     true,
			errContains: "ロックされているため",
			setup: func(repo *repomock.MockTemplateRepository, _ *templateUpdaterCaptures) {
				repo.EXPECT().Get(validTemplateID).Return(&entity.Template{Model: gorm.Model{ID: validTemplateID}, IsLock: true}, nil)
			},
		},
		{
			name:       "通常系_更新",
			templateID: validTemplateID,
			input:      input,
			wantErr:    false,
			setup: func(repo *repomock.MockTemplateRepository, cap *templateUpdaterCaptures) {
				repo.EXPECT().Get(validTemplateID).Return(&entity.Template{
					Model:  gorm.Model{ID: validTemplateID},
					Topics: []entity.Topic{{UUID: existingTopicUUID}},
				}, nil)
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
				assert.Equal(t, existingTopicUUID, cap.newTemplate.Topics[0].UUID)
				assert.Equal(t, "updated-question", cap.newTemplate.Topics[0].Question)
				assert.True(t, cap.newTemplate.Topics[0].IsRequire)
			},
		},
		{
			name:       "異常系_存在しない質問",
			templateID: validTemplateID,
			input: inputmodel.TemplateInput{
				Title:       "updated-template",
				Description: "updated description",
				Topics: []inputmodel.TopicInput{
					{UUID: "22222222-2222-2222-2222-222222222222", DisplayNo: 1, Question: "updated-question", IsRequire: true},
				},
			},
			wantErr:     true,
			errContains: "存在しない質問です",
			setup: func(repo *repomock.MockTemplateRepository, _ *templateUpdaterCaptures) {
				repo.EXPECT().Get(validTemplateID).Return(&entity.Template{
					Model: gorm.Model{ID: validTemplateID},
					Topics: []entity.Topic{{
						UUID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					}},
				}, nil)
			},
		},
		{
			name:        "異常系_更新失敗",
			templateID:  validTemplateID,
			input:       input,
			wantErr:     true,
			errContains: "テンプレートの情報更新に失敗しました",
			setup: func(repo *repomock.MockTemplateRepository, _ *templateUpdaterCaptures) {
				repo.EXPECT().Get(validTemplateID).Return(&entity.Template{
					Model:  gorm.Model{ID: validTemplateID},
					Topics: []entity.Topic{{UUID: existingTopicUUID}},
				}, nil)
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

func TestTemplateUpdater_UpdateTemplateLock(t *testing.T) {
	uc, repo := newTemplateUpdater(t)
	repo.EXPECT().UpdateIsLock(validTemplateID, true).Return(nil)
	gotErr := uc.UpdateTemplateLock(validTemplateID, true)
	assertUsecaseErr(t, gotErr, false, "", "UpdateTemplateLock")
}

func TestTemplateUpdater_UpdateTemplateDefault(t *testing.T) {
	uc, repo := newTemplateUpdater(t)
	repo.EXPECT().UpdateIsDefault(validTemplateID, true).Return(nil)
	gotErr := uc.UpdateTemplateDefault(validTemplateID, true)
	assertUsecaseErr(t, gotErr, false, "", "UpdateTemplateDefault")
}

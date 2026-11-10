package template_mng_usecase_test

import (
	"testing"

	inputmodel "growlink/internal/application/model/input_model"
	template_mng_usecase "growlink/internal/application/usecase/management_usecase/template_management_usecase"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type templateCreatorCaptures struct {
	createdTopics   []*entity.Topic
	createdTemplate *entity.Template
}

func newTemplateCreator(t *testing.T) (
	*template_mng_usecase.TemplateCreator,
	*repomock.MockTemplateRepository,
	*repomock.MockTopicRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	templateRepo := repomock.NewMockTemplateRepository(ctrl)
	topicRepo := repomock.NewMockTopicRepository(ctrl)
	return template_mng_usecase.NewTemplateCreator(templateRepo, topicRepo), templateRepo, topicRepo
}

func TestTemplateCreator_CreateTemplate(t *testing.T) {
	inputWithTopics := inputmodel.TemplateInput{
		Title:       "test-template",
		Description: "test description",
		Topics: []inputmodel.TopicInput{
			{DisplayNo: 1, Question: "question-1", IsRequire: true},
			{DisplayNo: 2, Question: "question-2", IsRequire: false},
		},
	}
	inputEmptyTopics := inputmodel.TemplateInput{
		Title:       "empty-topics-template",
		Description: "desc",
		Topics:      []inputmodel.TopicInput{},
	}

	tests := []struct {
		name        string
		input       inputmodel.TemplateInput
		wantErr     bool
		errContains string
		setup       func(templateRepo *repomock.MockTemplateRepository, topicRepo *repomock.MockTopicRepository, cap *templateCreatorCaptures)
		assertFn    func(t *testing.T, cap *templateCreatorCaptures)
	}{
		{
			name:    "通常系_トピック付きテンプレート作成",
			input:   inputWithTopics,
			wantErr: false,
			setup: func(templateRepo *repomock.MockTemplateRepository, topicRepo *repomock.MockTopicRepository, cap *templateCreatorCaptures) {
				templateRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(template *entity.Template) error {
					cap.createdTemplate = template
					return nil
				})
				topicRepo.EXPECT().Create(gomock.Any()).Times(2).DoAndReturn(func(topic *entity.Topic) error {
					cap.createdTopics = append(cap.createdTopics, topic)
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *templateCreatorCaptures) {
				require.Len(t, cap.createdTopics, 2)
				assert.Equal(t, uint(1), cap.createdTopics[0].DisplayNo)
				assert.Equal(t, "question-1", cap.createdTopics[0].Question)
				assert.True(t, cap.createdTopics[0].IsRequire)
				assert.Equal(t, uint(2), cap.createdTopics[1].DisplayNo)
				assert.Equal(t, "question-2", cap.createdTopics[1].Question)
				assert.False(t, cap.createdTopics[1].IsRequire)

				require.NotNil(t, cap.createdTemplate)
				assert.Equal(t, "test-template", cap.createdTemplate.Title)
				assert.Equal(t, "test description", cap.createdTemplate.Description)
				assert.Equal(t, validUserID, cap.createdTemplate.UserID)
			},
		},
		{
			name:        "異常系_トピック空",
			input:       inputEmptyTopics,
			wantErr:     true,
			errContains: "質問を1件以上",
			setup: func(_ *repomock.MockTemplateRepository, _ *repomock.MockTopicRepository, _ *templateCreatorCaptures) {
			},
		},
		{
			name:        "異常系_トピック作成失敗",
			input:       inputWithTopics,
			wantErr:     true,
			errContains: "トピックの作成に失敗しました",
			setup: func(templateRepo *repomock.MockTemplateRepository, topicRepo *repomock.MockTopicRepository, _ *templateCreatorCaptures) {
				templateRepo.EXPECT().Create(gomock.Any()).Return(nil)
				topicRepo.EXPECT().Create(gomock.Any()).Return(errCreateTopicFailed)
			},
		},
		{
			name:        "異常系_テンプレート作成失敗",
			input:       inputWithTopics,
			wantErr:     true,
			errContains: "テンプレートの新規作成に失敗しました",
			setup: func(templateRepo *repomock.MockTemplateRepository, topicRepo *repomock.MockTopicRepository, _ *templateCreatorCaptures) {
				templateRepo.EXPECT().Create(gomock.Any()).Return(errCreateTemplateFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, templateRepo, topicRepo := newTemplateCreator(t)
			cap := &templateCreatorCaptures{}
			if tt.setup != nil {
				tt.setup(templateRepo, topicRepo, cap)
			}

			gotErr := uc.CreateTemplate(tt.input, validUserID)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "CreateTemplate")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}

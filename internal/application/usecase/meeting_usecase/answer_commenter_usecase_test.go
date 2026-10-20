package meeting_usecase_test

import (
	"testing"

	"growlink/internal/application/usecase/meeting_usecase"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type answerCommenterCaptures struct {
	postedComment    *entity.AnswerComment
	updatedCommentID uint
	updatedContent   string
	deletedCommentID uint
}

type answerCommenterSetupFn func(t *testing.T, comment *repomock.MockAnswerCommentRepository, cap *answerCommenterCaptures)
type answerCommenterAssertFn func(t *testing.T, cap *answerCommenterCaptures)

func newAnswerCommenterUsecase(t *testing.T) (
	*meeting_usecase.AnswerCommenterUsecase,
	*repomock.MockAnswerCommentRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	commentRepo := repomock.NewMockAnswerCommentRepository(ctrl)
	return meeting_usecase.NewAnswerCommenterUsecase(commentRepo), commentRepo
}

func TestAnswerCommenterUsecase_GetComment(t *testing.T) {
	tests := []struct {
		name        string
		commentID   uint
		wantErr     bool
		errContains string
		setup       answerCommenterSetupFn
		assertFn    func(t *testing.T, got *entity.AnswerComment)
	}{
		{
			name:      "通常系_コメント取得",
			commentID: validCommentID,
			wantErr:   false,
			setup: func(_ *testing.T, comment *repomock.MockAnswerCommentRepository, _ *answerCommenterCaptures) {
				comment.EXPECT().GetComment(validCommentID).Return(newTestAnswerComment(validCommentID, validAnswerID, "comment"), nil)
			},
			assertFn: func(t *testing.T, got *entity.AnswerComment) {
				t.Helper()
				require.NotNil(t, got)
				assert.Equal(t, validCommentID, got.ID)
				requireStringPtrEqual(t, got.Content, strPtr("comment"))
			},
		},
		{
			name:        "異常系_リポジトリエラー",
			commentID:   validCommentID,
			wantErr:     true,
			errContains: "回答へのコメントの取得に失敗しました",
			setup: func(_ *testing.T, comment *repomock.MockAnswerCommentRepository, _ *answerCommenterCaptures) {
				comment.EXPECT().GetComment(validCommentID).Return(nil, errGetCommentFailed)
			},
		},
		{
			name:        "異常系_コメント未存在",
			commentID:   validCommentID,
			wantErr:     true,
			errContains: "指定されたコメントが存在しません",
			setup: func(_ *testing.T, comment *repomock.MockAnswerCommentRepository, _ *answerCommenterCaptures) {
				comment.EXPECT().GetComment(validCommentID).Return(nil, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, commentRepo := newAnswerCommenterUsecase(t)
			cap := &answerCommenterCaptures{}
			if tt.setup != nil {
				tt.setup(t, commentRepo, cap)
			}

			got, err := uc.GetComment(tt.commentID)
			assertUsecaseErr(t, err, tt.wantErr, tt.errContains, "GetComment")
			if !t.Failed() && tt.assertFn != nil {
				tt.assertFn(t, got)
			}
		})
	}
}

func TestAnswerCommenterUsecase_PostComment(t *testing.T) {
	content := "new comment"

	tests := []struct {
		name        string
		answerID    uint
		content     string
		wantErr     bool
		errContains string
		setup       answerCommenterSetupFn
		assertFn    answerCommenterAssertFn
	}{
		{
			name:     "通常系_コメント作成",
			answerID: validAnswerID,
			content:  content,
			wantErr:  false,
			setup: func(_ *testing.T, comment *repomock.MockAnswerCommentRepository, cap *answerCommenterCaptures) {
				comment.EXPECT().PostComment(gomock.Any()).DoAndReturn(func(c *entity.AnswerComment) error {
					cap.postedComment = c
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *answerCommenterCaptures) {
				t.Helper()
				require.NotNil(t, cap.postedComment)
				assert.Equal(t, validAnswerID, cap.postedComment.AnswerID)
				requireStringPtrEqual(t, cap.postedComment.Content, strPtr(content))
			},
		},
		{
			name:        "異常系_リポジトリエラー",
			answerID:    validAnswerID,
			content:     content,
			wantErr:     true,
			errContains: "回答へのコメントの作成に失敗しました",
			setup: func(_ *testing.T, comment *repomock.MockAnswerCommentRepository, _ *answerCommenterCaptures) {
				comment.EXPECT().PostComment(gomock.Any()).Return(errPostCommentFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, commentRepo := newAnswerCommenterUsecase(t)
			cap := &answerCommenterCaptures{}
			if tt.setup != nil {
				tt.setup(t, commentRepo, cap)
			}

			err := uc.PostComment(tt.answerID, tt.content)
			assertUsecaseErr(t, err, tt.wantErr, tt.errContains, "PostComment")
			if !t.Failed() && tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}

func TestAnswerCommenterUsecase_EditComment(t *testing.T) {
	content := "edited comment"

	tests := []struct {
		name        string
		commentID   uint
		content     string
		wantErr     bool
		errContains string
		setup       answerCommenterSetupFn
		assertFn    answerCommenterAssertFn
	}{
		{
			name:      "通常系_コメント編集",
			commentID: validCommentID,
			content:   content,
			wantErr:   false,
			setup: func(_ *testing.T, comment *repomock.MockAnswerCommentRepository, cap *answerCommenterCaptures) {
				comment.EXPECT().UpdateComment(validCommentID, content).DoAndReturn(func(commentID uint, newContent string) error {
					cap.updatedCommentID = commentID
					cap.updatedContent = newContent
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *answerCommenterCaptures) {
				t.Helper()
				assert.Equal(t, validCommentID, cap.updatedCommentID)
				assert.Equal(t, content, cap.updatedContent)
			},
		},
		{
			name:        "異常系_リポジトリエラー",
			commentID:   validCommentID,
			content:     content,
			wantErr:     true,
			errContains: "回答へのコメントの編集に失敗しました",
			setup: func(_ *testing.T, comment *repomock.MockAnswerCommentRepository, _ *answerCommenterCaptures) {
				comment.EXPECT().UpdateComment(validCommentID, content).Return(errUpdateCommentFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, commentRepo := newAnswerCommenterUsecase(t)
			cap := &answerCommenterCaptures{}
			if tt.setup != nil {
				tt.setup(t, commentRepo, cap)
			}

			err := uc.EditComment(tt.commentID, tt.content)
			assertUsecaseErr(t, err, tt.wantErr, tt.errContains, "EditComment")
			if !t.Failed() && tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}

func TestAnswerCommenterUsecase_DeleteComment(t *testing.T) {
	tests := []struct {
		name        string
		commentID   uint
		wantErr     bool
		errContains string
		setup       answerCommenterSetupFn
		assertFn    answerCommenterAssertFn
	}{
		{
			name:      "通常系_コメント削除",
			commentID: validCommentID,
			wantErr:   false,
			setup: func(_ *testing.T, comment *repomock.MockAnswerCommentRepository, cap *answerCommenterCaptures) {
				comment.EXPECT().DeleteComment(validCommentID).DoAndReturn(func(commentID uint) error {
					cap.deletedCommentID = commentID
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *answerCommenterCaptures) {
				t.Helper()
				assert.Equal(t, validCommentID, cap.deletedCommentID)
			},
		},
		{
			name:        "異常系_リポジトリエラー",
			commentID:   validCommentID,
			wantErr:     true,
			errContains: "回答へのコメントの削除に失敗しました",
			setup: func(_ *testing.T, comment *repomock.MockAnswerCommentRepository, _ *answerCommenterCaptures) {
				comment.EXPECT().DeleteComment(validCommentID).Return(errDeleteCommentFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, commentRepo := newAnswerCommenterUsecase(t)
			cap := &answerCommenterCaptures{}
			if tt.setup != nil {
				tt.setup(t, commentRepo, cap)
			}

			err := uc.DeleteComment(tt.commentID)
			assertUsecaseErr(t, err, tt.wantErr, tt.errContains, "DeleteComment")
			if !t.Failed() && tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}

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

type meetingCommenterCaptures struct {
	postedComment    *entity.MeetingComment
	updatedCommentID uint
	updatedContent   string
	deletedCommentID uint
}

type meetingCommenterSetupFn func(t *testing.T, comment *repomock.MockMeetingCommentRepository, cap *meetingCommenterCaptures)
type meetingCommenterAssertFn func(t *testing.T, cap *meetingCommenterCaptures)

func newMeetingCommenterUsecase(t *testing.T) (
	*meeting_usecase.MeetingCommenterUsecase,
	*repomock.MockMeetingCommentRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	commentRepo := repomock.NewMockMeetingCommentRepository(ctrl)
	return meeting_usecase.NewMeetingCommenterUsecase(commentRepo), commentRepo
}

func TestMeetingCommenterUsecase_GetComment(t *testing.T) {
	tests := []struct {
		name        string
		commentID   uint
		wantErr     bool
		errContains string
		setup       meetingCommenterSetupFn
		assertFn    func(t *testing.T, got *entity.MeetingComment)
	}{
		{
			name:      "通常系_コメント取得",
			commentID: validCommentID,
			wantErr:   false,
			setup: func(_ *testing.T, comment *repomock.MockMeetingCommentRepository, _ *meetingCommenterCaptures) {
				comment.EXPECT().GetComment(validCommentID).Return(newTestMeetingComment(validCommentID, validMeetingID, "comment"), nil)
			},
			assertFn: func(t *testing.T, got *entity.MeetingComment) {
				t.Helper()
				require.NotNil(t, got)
				assert.Equal(t, validCommentID, got.ID)
				assert.Equal(t, "comment", got.Content)
			},
		},
		{
			name:        "異常系_リポジトリエラー",
			commentID:   validCommentID,
			wantErr:     true,
			errContains: "ミーティングへのコメントの取得に失敗しました",
			setup: func(_ *testing.T, comment *repomock.MockMeetingCommentRepository, _ *meetingCommenterCaptures) {
				comment.EXPECT().GetComment(validCommentID).Return(nil, errGetCommentFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, commentRepo := newMeetingCommenterUsecase(t)
			cap := &meetingCommenterCaptures{}
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

func TestMeetingCommenterUsecase_PostComment(t *testing.T) {
	content := "new comment"

	tests := []struct {
		name        string
		meetingID   uint
		content     string
		wantErr     bool
		errContains string
		setup       meetingCommenterSetupFn
		assertFn    meetingCommenterAssertFn
	}{
		{
			name:      "通常系_コメント作成",
			meetingID: validMeetingID,
			content:   content,
			wantErr:   false,
			setup: func(_ *testing.T, comment *repomock.MockMeetingCommentRepository, cap *meetingCommenterCaptures) {
				comment.EXPECT().PostComment(gomock.Any()).DoAndReturn(func(c *entity.MeetingComment) error {
					cap.postedComment = c
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *meetingCommenterCaptures) {
				t.Helper()
				require.NotNil(t, cap.postedComment)
				assert.Equal(t, validMeetingID, cap.postedComment.MeetingID)
				assert.Equal(t, content, cap.postedComment.Content)
			},
		},
		{
			name:        "異常系_リポジトリエラー",
			meetingID:   validMeetingID,
			content:     content,
			wantErr:     true,
			errContains: "ミーティングへのコメントの作成に失敗しました",
			setup: func(_ *testing.T, comment *repomock.MockMeetingCommentRepository, _ *meetingCommenterCaptures) {
				comment.EXPECT().PostComment(gomock.Any()).Return(errPostCommentFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, commentRepo := newMeetingCommenterUsecase(t)
			cap := &meetingCommenterCaptures{}
			if tt.setup != nil {
				tt.setup(t, commentRepo, cap)
			}

			err := uc.PostComment(tt.meetingID, 1, tt.content, false)
			assertUsecaseErr(t, err, tt.wantErr, tt.errContains, "PostComment")
			if !t.Failed() && tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}

func TestMeetingCommenterUsecase_EditComment(t *testing.T) {
	content := "edited comment"

	tests := []struct {
		name        string
		commentID   uint
		content     string
		wantErr     bool
		errContains string
		setup       meetingCommenterSetupFn
		assertFn    meetingCommenterAssertFn
	}{
		{
			name:      "通常系_コメント編集",
			commentID: validCommentID,
			content:   content,
			wantErr:   false,
			setup: func(_ *testing.T, comment *repomock.MockMeetingCommentRepository, cap *meetingCommenterCaptures) {
				comment.EXPECT().UpdateComment(validCommentID, content).DoAndReturn(func(commentID uint, newContent string) error {
					cap.updatedCommentID = commentID
					cap.updatedContent = newContent
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *meetingCommenterCaptures) {
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
			errContains: "ミーティングへのコメントの編集に失敗しました",
			setup: func(_ *testing.T, comment *repomock.MockMeetingCommentRepository, _ *meetingCommenterCaptures) {
				comment.EXPECT().UpdateComment(validCommentID, content).Return(errUpdateCommentFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, commentRepo := newMeetingCommenterUsecase(t)
			cap := &meetingCommenterCaptures{}
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

func TestMeetingCommenterUsecase_DeleteComment(t *testing.T) {
	tests := []struct {
		name        string
		commentID   uint
		wantErr     bool
		errContains string
		setup       meetingCommenterSetupFn
		assertFn    meetingCommenterAssertFn
	}{
		{
			name:      "通常系_コメント削除",
			commentID: validCommentID,
			wantErr:   false,
			setup: func(_ *testing.T, comment *repomock.MockMeetingCommentRepository, cap *meetingCommenterCaptures) {
				comment.EXPECT().DeleteComment(validCommentID).DoAndReturn(func(commentID uint) error {
					cap.deletedCommentID = commentID
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *meetingCommenterCaptures) {
				t.Helper()
				assert.Equal(t, validCommentID, cap.deletedCommentID)
			},
		},
		{
			name:        "異常系_リポジトリエラー",
			commentID:   validCommentID,
			wantErr:     true,
			errContains: "ミーティングへのコメントの削除に失敗しました",
			setup: func(_ *testing.T, comment *repomock.MockMeetingCommentRepository, _ *meetingCommenterCaptures) {
				comment.EXPECT().DeleteComment(validCommentID).Return(errDeleteCommentFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, commentRepo := newMeetingCommenterUsecase(t)
			cap := &meetingCommenterCaptures{}
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

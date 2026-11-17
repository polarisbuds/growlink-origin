package meeting_usecase

import (
	"context"
	"fmt"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"
)

type MeetingCommenterUsecase struct {
	meetingCommentRepo repository_interface.MeetingCommentRepository
}

func NewMeetingCommenterUsecase(meetingCommentRepo repository_interface.MeetingCommentRepository) *MeetingCommenterUsecase {
	return &MeetingCommenterUsecase{
		meetingCommentRepo: meetingCommentRepo,
	}
}

func (uc *MeetingCommenterUsecase) GetComment(commentID uint) (ret0 *entity.MeetingComment, err error) {
	const op = "meeting_usecase.GetComment"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	comment, err := uc.meetingCommentRepo.GetComment(commentID)
	if err != nil {
		return nil, fmt.Errorf("ミーティングへのコメントの取得に失敗しました: %w", err)
	}
	if comment == nil {
		return nil, fmt.Errorf("ミーティングへのコメントの取得に失敗しました: 指定されたコメントが存在しません")
	}
	return comment, nil
}

func (uc *MeetingCommenterUsecase) PostComment(meetingID, userID uint, content string, isLimited bool) (err error) {
	const op = "meeting_usecase.PostComment"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.meetingCommentRepo.PostComment(&entity.MeetingComment{
		MeetingID: meetingID,
		UserID:    userID,
		Content:   content,
		IsLimited: isLimited,
	}); err != nil {
		return fmt.Errorf("ミーティングへのコメントの作成に失敗しました: %w", err)
	}
	return nil
}

func (uc *MeetingCommenterUsecase) EditComment(commentID uint, content string) (err error) {
	const op = "meeting_usecase.EditComment"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.meetingCommentRepo.UpdateComment(commentID, content); err != nil {
		return fmt.Errorf("ミーティングへのコメントの編集に失敗しました: %w", err)
	}
	return nil
}

func (uc *MeetingCommenterUsecase) DeleteComment(commentID uint) (err error) {
	const op = "meeting_usecase.DeleteComment"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.meetingCommentRepo.DeleteComment(commentID); err != nil {
		return fmt.Errorf("ミーティングへのコメントの削除に失敗しました: %w", err)
	}
	return nil
}

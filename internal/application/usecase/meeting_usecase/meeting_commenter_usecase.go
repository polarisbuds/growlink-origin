package meeting_usecase

import (
	"fmt"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
)

type MeetingCommenterUsecase struct {
	meetingCommentRepo repository_interface.MeetingCommentRepository
}

func NewMeetingCommenterUsecase(meetingCommentRepo repository_interface.MeetingCommentRepository) *MeetingCommenterUsecase {
	return &MeetingCommenterUsecase{
		meetingCommentRepo: meetingCommentRepo,
	}
}

func (uc *MeetingCommenterUsecase) GetComment(commentID uint) (*entity.MeetingComment, error) {
	comment, err := uc.meetingCommentRepo.GetComment(commentID)
	if err != nil {
		return nil, fmt.Errorf("ミーティングへのコメントの取得に失敗しました: %w", err)
	}
	return comment, nil
}

func (uc *MeetingCommenterUsecase) PostComment(meetingID uint, content string) error {
	if err := uc.meetingCommentRepo.PostComment(&entity.MeetingComment{
		MeetingID: meetingID,
		Content:   &content,
	}); err != nil {
		return fmt.Errorf("ミーティングへのコメントの作成に失敗しました: %w", err)
	}
	return nil
}

func (uc *MeetingCommenterUsecase) EditComment(commentID uint, content string) error {
	if err := uc.meetingCommentRepo.UpdateComment(commentID, content); err != nil {
		return fmt.Errorf("ミーティングへのコメントの編集に失敗しました: %w", err)
	}
	return nil
}

func (uc *MeetingCommenterUsecase) DeleteComment(commentID uint) error {
	if err := uc.meetingCommentRepo.DeleteComment(commentID); err != nil {
		return fmt.Errorf("ミーティングへのコメントの削除に失敗しました: %w", err)
	}
	return nil
}

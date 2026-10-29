package meeting_usecase

import (
	"fmt"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
)

type AnswerCommenterUsecase struct {
	answerCommentRepo repository_interface.AnswerCommentRepository
}

func NewAnswerCommenterUsecase(answerCommentRepo repository_interface.AnswerCommentRepository) *AnswerCommenterUsecase {
	return &AnswerCommenterUsecase{
		answerCommentRepo: answerCommentRepo,
	}
}

func (uc *AnswerCommenterUsecase) GetComment(commentID uint) (*entity.AnswerComment, error) {
	comment, err := uc.answerCommentRepo.GetComment(commentID)
	if err != nil {
		return nil, fmt.Errorf("回答へのコメントの取得に失敗しました: %w", err)
	}
	if comment == nil {
		return nil, fmt.Errorf("回答へのコメントの取得に失敗しました: 指定されたコメントが存在しません")
	}
	return comment, nil
}

func (uc *AnswerCommenterUsecase) PostComment(answerID uint, content string) error {
	if err := uc.answerCommentRepo.PostComment(&entity.AnswerComment{
		AnswerID: answerID,
		Content:  &content,
	}); err != nil {
		return fmt.Errorf("回答へのコメントの作成に失敗しました: %w", err)
	}
	return nil
}

func (uc *AnswerCommenterUsecase) EditComment(commentID uint, content string) error {
	if err := uc.answerCommentRepo.UpdateComment(commentID, content); err != nil {
		return fmt.Errorf("回答へのコメントの編集に失敗しました: %w", err)
	}
	return nil
}

func (uc *AnswerCommenterUsecase) DeleteComment(commentID uint) error {
	if err := uc.answerCommentRepo.DeleteComment(commentID); err != nil {
		return fmt.Errorf("回答へのコメントの削除に失敗しました: %w", err)
	}
	return nil
}

package line_works_messages_usecase

import (
	"context"
	"fmt"
	"strings"

	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"

	"github.com/google/uuid"
)

const lineWorksMessageCount = 6

func orderedTitles() []string {
	return []string{
		valueobject.LineWorksMessageKindUnsubmitted.DefaultTitle(),
		valueobject.LineWorksMessageKindSubmit.DefaultTitle(),
		valueobject.LineWorksMessageKindCancelSubmitted.DefaultTitle(),
		valueobject.LineWorksMessageKindNotConducted.DefaultTitle(),
		valueobject.LineWorksMessageKindCreateMeeting.DefaultTitle(),
		valueobject.LineWorksMessageKindSetImplementDate.DefaultTitle(),
	}
}

type LineWorksMessages struct {
	repo repository_interface.LineWorksMessageSettingRepository
}

func NewLineWorksMessages(repo repository_interface.LineWorksMessageSettingRepository) *LineWorksMessages {
	return &LineWorksMessages{repo: repo}
}

func (uc *LineWorksMessages) List() (ret0 []outputmodel.LineWorksMessageOutput, err error) {
	const op = "line_works_messages_usecase.List"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.ensureAllDefaults(); err != nil {
		return nil, err
	}
	rows, err := uc.repo.ListByTitles(orderedTitles())
	if err != nil {
		return nil, err
	}
	return rowsToOutput(rows), nil
}

type LineWorksMessageUpdateItem struct {
	UUID    string `json:"uuid"`
	Content string `json:"content"`
}

type UpdateLineWorksMessagesInput struct {
	Messages []LineWorksMessageUpdateItem `json:"messages"`
}

func (uc *LineWorksMessages) Update(updateUserID uint, input UpdateLineWorksMessagesInput) (err error) {
	const op = "line_works_messages_usecase.Update"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = validateUpdateInput(input); err != nil {
		return err
	}

	for _, item := range input.Messages {
		id, err := uuid.Parse(strings.TrimSpace(item.UUID))
		if err != nil {
			return fmt.Errorf("通知文の識別子が不正です")
		}
		row, err := uc.repo.GetByUUID(id)
		if err != nil {
			return err
		}
		if row == nil {
			return fmt.Errorf("通知文が見つかりません")
		}
		row.Content = strings.TrimSpace(item.Content)
		row.UpdateUserID = updateUserID
		if err := uc.repo.Update(row); err != nil {
			return err
		}
	}
	return nil
}

func validateUpdateInput(input UpdateLineWorksMessagesInput) error {
	if len(input.Messages) != lineWorksMessageCount {
		return fmt.Errorf("通知文は %d 件すべて指定してください", lineWorksMessageCount)
	}
	seen := make(map[string]struct{}, lineWorksMessageCount)
	for _, item := range input.Messages {
		uuidStr := strings.TrimSpace(item.UUID)
		if uuidStr == "" {
			return fmt.Errorf("通知文の識別子を指定してください")
		}
		if _, err := uuid.Parse(uuidStr); err != nil {
			return fmt.Errorf("通知文の識別子が不正です")
		}
		if _, dup := seen[uuidStr]; dup {
			return fmt.Errorf("通知文の識別子が重複しています")
		}
		seen[uuidStr] = struct{}{}
		if strings.TrimSpace(item.Content) == "" {
			return fmt.Errorf("本文を入力してください")
		}
	}
	return nil
}

func (uc *LineWorksMessages) ensureAllDefaults() error {
	titles := orderedTitles()
	rows, err := uc.repo.ListByTitles(titles)
	if err != nil {
		return err
	}
	exists := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		exists[row.Title] = struct{}{}
	}
	for _, title := range titles {
		if _, ok := exists[title]; ok {
			continue
		}
		kind, ok := valueobject.LineWorksMessageKindFromTitle(title)
		if !ok {
			return fmt.Errorf("既定の通知文がありません: %s", title)
		}
		content, ok := kind.DefaultContent()
		if !ok {
			return fmt.Errorf("既定の通知文がありません: %s", title)
		}
		if err := uc.repo.Create(&entity.LineWorksMessageSetting{
			Kind:    kind,
			Title:   title,
			Content: content,
		}); err != nil {
			return err
		}
	}
	return nil
}

func rowsToOutput(rows []entity.LineWorksMessageSetting) []outputmodel.LineWorksMessageOutput {
	out := make([]outputmodel.LineWorksMessageOutput, 0, len(rows))
	for _, row := range rows {
		out = append(out, outputmodel.LineWorksMessageOutput{
			UUID:    row.UUID.String(),
			Title:   row.Title,
			Content: row.Content,
		})
	}
	return out
}

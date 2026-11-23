package lineworks_usecase

import (
	"context"
	"errors"
	"fmt"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"
	"strings"

	"github.com/google/uuid"
)

// Sender は LINE WORKS のユーザー1人へテキストを送る。
type Sender interface {
	SendUserMessage(ctx context.Context, userID, text string) error
}

type Usecase struct {
	client                      Sender
	lineWorksMessageSettingRepo repository_interface.LineWorksMessageSettingRepository
	meetingRepo                 repository_interface.MeetingRepository
	pairingRepo                 repository_interface.PairingRepository
}

func New(client Sender, lineWorksMessageSetting repository_interface.LineWorksMessageSettingRepository, meetingRepo repository_interface.MeetingRepository, pairingRepo repository_interface.PairingRepository) *Usecase {
	return &Usecase{
		client:                      client,
		lineWorksMessageSettingRepo: lineWorksMessageSetting,
		meetingRepo:                 meetingRepo,
		pairingRepo:                 pairingRepo,
	}
}

// SendMessage は lineworksIDs の各ユーザーへ同じ本文を送る。
// 空の ID は飛ばす。1人でも失敗したら、残りの送信のあとでそのエラーを返す。
func (uc *Usecase) SendMessage(ctx context.Context, lineworksIDs []string, message string) (err error) {
	const op = "lineworks_usecase.SendMessage"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	message = strings.TrimSpace(message)
	if message == "" {
		return fmt.Errorf("message is empty")
	}
	var errs []error
	sent := 0
	for _, id := range lineworksIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if err = uc.client.SendUserMessage(ctx, id, message); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", id, err))
			continue
		}
		sent++
	}
	if sent == 0 && len(errs) == 0 {
		return fmt.Errorf("lineworksIDs is empty")
	}
	return errors.Join(errs...)
}

func (uc *Usecase) pairingByMeetingUUID(meetingUUID uuid.UUID) (mentorLineWorksID, menteeLineWorksID string, err error) {
	meeting, err := uc.meetingRepo.GetByUUID(meetingUUID)
	if err != nil {
		return "", "", fmt.Errorf("ミーティングの取得に失敗しました: %w", err)
	}

	pairing, err := uc.pairingRepo.Get(meeting.PairingID)
	if err != nil {
		return "", "", fmt.Errorf("ペアリングの取得に失敗しました: %w", err)
	}
	if pairing.Mentor == nil {
		return "", "", fmt.Errorf("メンターの取得に失敗しました")
	}
	if pairing.Mentee == nil {
		return "", "", fmt.Errorf("メンティーの取得に失敗しました")
	}
	return pairing.Mentor.LineWorksID, pairing.Mentee.LineWorksID, nil
}

func (uc *Usecase) SendSubmitMsg(ctx context.Context, meetingUUID uuid.UUID) (err error) {
	const op = "lineworks_usecase.SendSubmitMsg"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()

	mentorLineWorksID, _, err := uc.pairingByMeetingUUID(meetingUUID)
	if err != nil {
		return err
	}
	notifyMessage, err := uc.lineWorksMessageSettingRepo.GetLineWorksMessageByKind(ctx, valueobject.LineWorksMessageKindSubmit)
	if err != nil {
		return fmt.Errorf("提出通知文の取得に失敗しました: %w", err)
	}
	if err := uc.client.SendUserMessage(ctx, mentorLineWorksID, notifyMessage.Content); err != nil {
		return fmt.Errorf("提出完了通知に失敗しました: %w", err)
	}
	return nil
}

func (uc *Usecase) SendUnsubmittedMsg(ctx context.Context, meetingUUID uuid.UUID) (err error) {
	const op = "lineworks_usecase.SendUnsubmittedMsg"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()

	_, menteeLineWorksID, err := uc.pairingByMeetingUUID(meetingUUID)
	if err != nil {
		return err
	}
	notifyMessage, err := uc.lineWorksMessageSettingRepo.GetLineWorksMessageByKind(ctx, valueobject.LineWorksMessageKindUnsubmitted)
	if err != nil {
		return fmt.Errorf("未提出通知文の取得に失敗しました: %w", err)
	}
	if err := uc.client.SendUserMessage(ctx, menteeLineWorksID, notifyMessage.Content); err != nil {
		return fmt.Errorf("未提出通知に失敗しました: %w", err)
	}
	return nil
}

func (uc *Usecase) SendNotConductedMsg(ctx context.Context, meetingUUID uuid.UUID) (err error) {
	const op = "lineworks_usecase.SendNotConductedMsg"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()

	mentorLineWorksID, _, err := uc.pairingByMeetingUUID(meetingUUID)
	if err != nil {
		return err
	}
	notifyMessage, err := uc.lineWorksMessageSettingRepo.GetLineWorksMessageByKind(ctx, valueobject.LineWorksMessageKindNotConducted)
	if err != nil {
		return fmt.Errorf("未実施通知文の取得に失敗しました: %w", err)
	}
	if err := uc.client.SendUserMessage(ctx, mentorLineWorksID, notifyMessage.Content); err != nil {
		return fmt.Errorf("未実施通知に失敗しました: %w", err)
	}
	return nil
}

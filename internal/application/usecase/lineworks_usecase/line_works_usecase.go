package lineworks_usecase

import (
	"context"
	"errors"
	"fmt"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	"growlink/internal/shared/auth"
	logutils "growlink/internal/shared/logUtils"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// Sender は LINE WORKS のユーザー1人へテキストを送る。
type Sender interface {
	SendUserMessage(ctx context.Context, userID, text string) error
}

type LineWorksUsecase struct {
	base                        string
	client                      Sender
	lineWorksMessageSettingRepo repository_interface.LineWorksMessageSettingRepository
	meetingRepo                 repository_interface.MeetingRepository
	pairingRepo                 repository_interface.PairingRepository
}

func New(client Sender, lineWorksMessageSetting repository_interface.LineWorksMessageSettingRepository, meetingRepo repository_interface.MeetingRepository, pairingRepo repository_interface.PairingRepository) *LineWorksUsecase {
	return &LineWorksUsecase{
		base:                        os.Getenv("BASE"),
		client:                      client,
		lineWorksMessageSettingRepo: lineWorksMessageSetting,
		meetingRepo:                 meetingRepo,
		pairingRepo:                 pairingRepo,
	}
}

// SendMessage は lineworksIDs の各ユーザーへ同じ本文を送る。
// 空の ID は飛ばす。1人でも失敗したら、残りの送信のあとでそのエラーを返す。
func (uc *LineWorksUsecase) SendMessage(ctx context.Context, lineworksIDs []string, message string) (err error) {
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

func (uc *LineWorksUsecase) embededdingMessageTemplate(msg string, claims *auth.CustomClaims, meeting *entity.Meeting, pairing *entity.Pairing, isEmbeddedDirectLink bool) string {
	result := msg
	if strings.Contains(result, "[iam_name]") {
		result = strings.ReplaceAll(result, "[iam_name]", claims.Name)
	}
	if strings.Contains(result, "[iam_mail]") {
		result = strings.ReplaceAll(result, "[iam_mail]", claims.Email)
	}
	if strings.Contains(result, "[mentee]") {
		result = strings.ReplaceAll(result, "[mentee]", pairing.Mentee.Name)
	}
	if strings.Contains(result, "[mentor]") {
		result = strings.ReplaceAll(result, "[mentor]", pairing.Mentor.Name)
	}
	if strings.Contains(result, "[year]") {
		result = strings.ReplaceAll(result, "[year]", meeting.ImplementedAt.Format("2006"))
	}
	if strings.Contains(result, "[month]") {
		result = strings.ReplaceAll(result, "[month]", meeting.ImplementedAt.Format("01"))
	}
	if strings.Contains(result, "[day]") {
		result = strings.ReplaceAll(result, "[day]", meeting.ImplementedAt.Format("02"))
	}
	if strings.Contains(result, "[url]") {
		if isEmbeddedDirectLink {
			result = strings.ReplaceAll(result, "[url]", fmt.Sprintf("%s/meeting/%s ", uc.base, meeting.UUID.String()))
		} else {
			result = strings.ReplaceAll(result, "[url]", uc.base)
		}
	}
	return result
}

func (uc *LineWorksUsecase) getRequireInfo(meetingUUID uuid.UUID) (meeting *entity.Meeting, pairing *entity.Pairing, err error) {
	meeting, err = uc.meetingRepo.GetByUUID(meetingUUID)
	if err != nil {
		return nil, nil, fmt.Errorf("ミーティングの取得に失敗しました: %w", err)
	}

	pairing, err = uc.pairingRepo.Get(meeting.PairingID)
	if err != nil {
		return nil, nil, fmt.Errorf("ペアリングの取得に失敗しました: %w", err)
	}
	return meeting, pairing, nil
}

func (uc *LineWorksUsecase) SendSubmitMsg(ctx context.Context, claims *auth.CustomClaims, meetingUUID uuid.UUID) (err error) {
	const op = "lineworks_usecase.SendSubmitMsg"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()

	return uc.sendMentor(ctx, claims, meetingUUID, valueobject.LineWorksMessageKindSubmit, true)
}

func (uc *LineWorksUsecase) SendCancelSubmittedMsg(ctx context.Context, claims *auth.CustomClaims, meetingUUID uuid.UUID) (err error) {
	const op = "lineworks_usecase.SencCancelSubmittedMsg"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()

	return uc.sendMentor(ctx, claims, meetingUUID, valueobject.LineWorksMessageKindCancelSubmitted, true)
}

func (uc *LineWorksUsecase) SendUnsubmittedMsg(ctx context.Context, claims *auth.CustomClaims, meetingUUID uuid.UUID) (err error) {
	const op = "lineworks_usecase.SendUnsubmittedMsg"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()

	return uc.sendMentee(ctx, claims, meetingUUID, valueobject.LineWorksMessageKindUnsubmitted, true)
}

func (uc *LineWorksUsecase) SendNotConductedMsg(ctx context.Context, claims *auth.CustomClaims, meetingUUID uuid.UUID) (err error) {
	const op = "lineworks_usecase.SendNotConductedMsg"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()

	return uc.sendMentor(ctx, claims, meetingUUID, valueobject.LineWorksMessageKindNotConducted, true)
}

func (uc *LineWorksUsecase) SendCreateMeetingMsg(ctx context.Context, claims *auth.CustomClaims, meetingUUID uuid.UUID) (err error) {
	const op = "lineworks_usecase.SendCreateMeetingMsg"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()

	return uc.sendMentee(ctx, claims, meetingUUID, valueobject.LineWorksMessageKindCreateMeeting, true)
}

func (uc *LineWorksUsecase) SendCreateMeetingMsgs(ctx context.Context, claims *auth.CustomClaims, meetingUUIDs []uuid.UUID) (err error) {
	const op = "lineworks_usecase.SendCreateMeetingMsg"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()

	sem := semaphore.NewWeighted(3)
	bg := context.WithoutCancel(ctx)
	group, gctx := errgroup.WithContext(bg)
	for _, meetingUUID := range meetingUUIDs {
		sem.Acquire(gctx, 1)
		group.Go(func() error {
			defer sem.Release(1)
			if err := uc.sendMentee(gctx, claims, meetingUUID, valueobject.LineWorksMessageKindCreateMeeting, true); err != nil {
				return err
			}
			time.Sleep(time.Second)
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		logutils.Error(ctx, logutils.LayerRepository, op, "LINE WORKSへの通知処理に失敗しました", slog.String("err", err.Error()))
	}
	return nil
}

func (uc *LineWorksUsecase) SendChangeImplementedAtMsg(ctx context.Context, claims *auth.CustomClaims, meetingUUID uuid.UUID) (err error) {
	const op = "lineworks_usecase.SendChangeImplementDate"
	done := logutils.Begin(ctx, logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()

	return uc.sendBoth(ctx, claims, meetingUUID, valueobject.LineWorksMessageKindSetImplementDate, true)
}

func (uc *LineWorksUsecase) sendBoth(ctx context.Context, claims *auth.CustomClaims, meetingUUID uuid.UUID, kind valueobject.LineWorksMessageKind, isEmbeddedDirectLink bool) (err error) {
	meeting, pairing, err := uc.getRequireInfo(meetingUUID)
	if err != nil {
		return err
	}
	if pairing == nil {
		slog.Warn("ペアリングの情報が存在しませんでした", slog.String("meetings.uuid", meetingUUID.String()))
		return nil
	}

	sendUserLineWorksIDs := make([]string, 0, 2)
	if mentorLineWorksID := pairing.Mentor.LineWorksID; mentorLineWorksID != "" {
		sendUserLineWorksIDs = append(sendUserLineWorksIDs, mentorLineWorksID)
	}
	if menteeLineWorksID := pairing.Mentee.LineWorksID; menteeLineWorksID != "" {
		sendUserLineWorksIDs = append(sendUserLineWorksIDs, menteeLineWorksID)
	}
	if len(sendUserLineWorksIDs) == 0 {
		slog.Error("LINE WORKSのIDが設定されていません", slog.Any("uuids", sendUserLineWorksIDs))
		return nil
	}

	notifyMessage, err := uc.lineWorksMessageSettingRepo.GetLineWorksMessageByKind(ctx, kind)
	if err != nil {
		slog.Error(fmt.Sprintf("%s文の取得に失敗しました", kind.DefaultTitle()), slog.String("err", err.Error()))
		return err
	}
	message := uc.embededdingMessageTemplate(notifyMessage.Content, claims, meeting, pairing, isEmbeddedDirectLink)
	if err := uc.SendMessage(ctx, sendUserLineWorksIDs, message); err != nil {
		slog.Error(fmt.Sprintf("%sに失敗しました", kind.DefaultTitle()), slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (uc *LineWorksUsecase) sendMentor(ctx context.Context, claims *auth.CustomClaims, meetingUUID uuid.UUID, kind valueobject.LineWorksMessageKind, isEmbeddedDirectLink bool) (err error) {
	meeting, pairing, err := uc.getRequireInfo(meetingUUID)
	if err != nil {
		return err
	}
	if pairing == nil {
		slog.Error("ペアリングの情報が存在しません", slog.String("meetings.uuid", meetingUUID.String()))
		return nil
	}
	if pairing.Mentor.LineWorksID == "" {
		slog.Error("メンターのLINE WORKSのIDが設定されていません", slog.String("mentor.uuid", pairing.Mentor.UUID.String()))
		return nil
	}

	notifyMessage, err := uc.lineWorksMessageSettingRepo.GetLineWorksMessageByKind(ctx, kind)
	if err != nil {
		slog.Error(fmt.Sprintf("%s文の取得に失敗しました", kind.DefaultTitle()), slog.String("err", err.Error()))
		return err
	}
	message := uc.embededdingMessageTemplate(notifyMessage.Content, claims, meeting, pairing, isEmbeddedDirectLink)
	if err := uc.client.SendUserMessage(ctx, pairing.Mentor.LineWorksID, message); err != nil {
		slog.Error(fmt.Sprintf("%sに失敗しました", kind.DefaultTitle()), slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (uc *LineWorksUsecase) sendMentee(ctx context.Context, claims *auth.CustomClaims, meetingUUID uuid.UUID, kind valueobject.LineWorksMessageKind, isEmbeddedDirectLink bool) (err error) {
	meeting, pairing, err := uc.getRequireInfo(meetingUUID)
	if err != nil {
		return err
	}
	if pairing == nil {
		slog.Error("ペアリングの情報が存在しません", slog.String("meetings.uuid", meetingUUID.String()))
		return nil
	}
	if pairing.Mentee.LineWorksID == "" {
		slog.Error("メンティーのLINE WORKSのIDが設定されていません", slog.String("mentor.uuid", pairing.Mentor.UUID.String()))
		return nil
	}

	notifyMessage, err := uc.lineWorksMessageSettingRepo.GetLineWorksMessageByKind(ctx, kind)
	if err != nil {
		slog.Error(fmt.Sprintf("%s文の取得に失敗しました", kind.DefaultTitle()), slog.String("err", err.Error()))
	}
	message := uc.embededdingMessageTemplate(notifyMessage.Content, claims, meeting, pairing, isEmbeddedDirectLink)
	if err := uc.client.SendUserMessage(ctx, pairing.Mentee.LineWorksID, message); err != nil {
		slog.Error(fmt.Sprintf("%sに失敗しました", kind.DefaultTitle()), slog.String("err", err.Error()))
	}
	return nil
}

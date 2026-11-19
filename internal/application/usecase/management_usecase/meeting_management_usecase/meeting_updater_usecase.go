package meeting_mng_usecase

import (
	"context"
	"fmt"
	"growlink/internal/application/converter"
	inputmodel "growlink/internal/application/model/input_model"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"
	"time"
)

type MeetingUpdater struct {
	meetingRepo repository_interface.MeetingRepository
	userRepo    repository_interface.UserRepository
}

func NewMeetingUpdater(meetingRepo repository_interface.MeetingRepository, userRepo repository_interface.UserRepository) *MeetingUpdater {
	return &MeetingUpdater{
		meetingRepo: meetingRepo,
		userRepo:    userRepo,
	}
}

func (uc *MeetingUpdater) Update(meetindID uint, meeting inputmodel.MeetingInput) (err error) {
	const op = "meeting_mng_usecase.Update"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.meetingRepo.Update(meetindID, converter.MeetingInputToEntity(meeting)); err != nil {
		return fmt.Errorf("ミーティングの更新に失敗しました: %w", err)
	}
	return nil
}

func (uc *MeetingUpdater) ChangeImplementDate(meetingID uint, implementDate time.Time) (err error) {
	const op = "meeting_mng_usecase.ChangeImplementDate"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.meetingRepo.UpdateImplementDate(meetingID, implementDate); err != nil {
		return fmt.Errorf("ミーティングの実施予定日の変更に失敗しました: %w", err)
	}
	return nil
}

func (uc *MeetingUpdater) ChangePairing(meetingID, newPairingID uint) (err error) {
	const op = "meeting_mng_usecase.ChangePairing"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.meetingRepo.UpdatePairing(meetingID, newPairingID); err != nil {
		return fmt.Errorf("ミーティングのペアリングの変更に失敗しました: %w", err)
	}
	return nil
}

func (uc *MeetingUpdater) ChangeTemplate(meetingID, newTemplateID uint) (err error) {
	const op = "meeting_mng_usecase.ChangeTemplate"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if err = uc.meetingRepo.UpdateTemplate(meetingID, newTemplateID); err != nil {
		return fmt.Errorf("ミーティングのテンプレートの変更に失敗しました: %w", err)
	}
	return nil
}

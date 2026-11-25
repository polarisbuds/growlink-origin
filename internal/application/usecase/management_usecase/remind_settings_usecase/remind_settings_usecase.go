package remind_settings_usecase

import (
	"context"
	"fmt"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"
	"strings"
)

type RemindSettings struct {
	settingsRepo repository_interface.SystemSettingsRepository
}

func NewRemindSettings(settingsRepo repository_interface.SystemSettingsRepository) *RemindSettings {
	return &RemindSettings{settingsRepo: settingsRepo}
}

func (uc *RemindSettings) Get() (ret0 *outputmodel.RemindSettingsOutput, err error) {
	const op = "remind_settings_usecase.Get"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	s, err := uc.settingsRepo.GetSingleton()
	if err != nil {
		return nil, err
	}
	startWeek := s.NotifyStartWeek
	if startWeek < 1 {
		startWeek = 1
	}
	intervalDays := s.NotifyIntervalDays
	if intervalDays < 1 {
		intervalDays = 1
	}
	notifyTime := s.NotifyTime
	if notifyTime == "" {
		notifyTime = "09:00"
	}
	return &outputmodel.RemindSettingsOutput{
		NotifyStartWeek:    startWeek,
		NotifyIntervalDays: intervalDays,
		NotifyTime:         notifyTime,
	}, nil
}

type UpdateRemindInput struct {
	NotifyStartWeek    int    `json:"notifyStartWeek"`
	NotifyIntervalDays int    `json:"notifyIntervalDays"`
	NotifyTime         string `json:"notifyTime"`
}

func (uc *RemindSettings) Update(input UpdateRemindInput) (err error) {
	const op = "remind_settings_usecase.Update"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	if !strings.Contains(input.NotifyTime, ":") {
		return fmt.Errorf("時刻の形式が不正です")
	}
	if input.NotifyStartWeek < 1 || input.NotifyStartWeek > 5 {
		return fmt.Errorf("送信開始週は 1〜5 で指定してください")
	}
	if input.NotifyIntervalDays < 1 || input.NotifyIntervalDays > 5 {
		return fmt.Errorf("送信間隔は 1〜5 で指定してください")
	}
	s, err := uc.settingsRepo.GetSingleton()
	if err != nil {
		return err
	}
	s.NotifyStartWeek = input.NotifyStartWeek
	s.NotifyIntervalDays = input.NotifyIntervalDays
	s.NotifyTime = input.NotifyTime
	return uc.settingsRepo.Update(s)
}
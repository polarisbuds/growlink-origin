package remind_settings_usecase

import (
	"fmt"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
	"strings"
)

type RemindSettings struct {
	settingsRepo repository_interface.SystemSettingsRepository
}

func NewRemindSettings(settingsRepo repository_interface.SystemSettingsRepository) *RemindSettings {
	return &RemindSettings{settingsRepo: settingsRepo}
}

func (uc *RemindSettings) Get() (*outputmodel.RemindSettingsOutput, error) {
	s, err := uc.settingsRepo.GetSingleton()
	if err != nil {
		return nil, err
	}
	mode := s.RemindMode
	if mode == "" {
		mode = "prev_weekday"
	}
	daysBefore := s.RemindDaysBefore
	if daysBefore < 1 {
		daysBefore = 1
	}
	return &outputmodel.RemindSettingsOutput{
		RemindAutoSend:   s.RemindAutoSend,
		RemindMode:       mode,
		RemindDaysBefore: daysBefore,
		RemindWeekday:    s.RemindWeekday,
		RemindTime:       s.RemindTime,
	}, nil
}

type UpdateRemindInput struct {
	RemindAutoSend   bool   `json:"remindAutoSend"`
	RemindMode       string `json:"remindMode"`
	RemindDaysBefore int    `json:"remindDaysBefore"`
	RemindWeekday    int    `json:"remindWeekday"`
	RemindTime       string `json:"remindTime"`
}

func (uc *RemindSettings) Update(input UpdateRemindInput) error {
	if !strings.Contains(input.RemindTime, ":") {
		return fmt.Errorf("時刻の形式が不正です")
	}
	mode := input.RemindMode
	if mode == "" {
		mode = "prev_weekday"
	}
	if mode != "days_before" && mode != "prev_weekday" {
		return fmt.Errorf("通知の種類が不正です")
	}
	if mode == "days_before" {
		if input.RemindDaysBefore < 1 || input.RemindDaysBefore > 30 {
			return fmt.Errorf("日数は 1〜30 で指定してください")
		}
	} else if input.RemindWeekday < 1 || input.RemindWeekday > 7 {
		return fmt.Errorf("曜日が不正です")
	}
	s, err := uc.settingsRepo.GetSingleton()
	if err != nil {
		return err
	}
	s.RemindAutoSend = input.RemindAutoSend
	s.RemindMode = mode
	s.RemindDaysBefore = input.RemindDaysBefore
	s.RemindWeekday = input.RemindWeekday
	s.RemindTime = input.RemindTime
	return uc.settingsRepo.Update(s)
}

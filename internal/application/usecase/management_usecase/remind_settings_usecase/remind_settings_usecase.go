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
	return &outputmodel.RemindSettingsOutput{
		RemindAutoSend: s.RemindAutoSend,
		RemindWeekday:  s.RemindWeekday,
		RemindTime:     s.RemindTime,
	}, nil
}

type UpdateRemindInput struct {
	RemindAutoSend bool   `json:"remindAutoSend"`
	RemindWeekday  int    `json:"remindWeekday"`
	RemindTime     string `json:"remindTime"`
}

func (uc *RemindSettings) Update(input UpdateRemindInput) error {
	if !strings.Contains(input.RemindTime, ":") {
		return fmt.Errorf("時刻の形式が不正です")
	}
	s, err := uc.settingsRepo.GetSingleton()
	if err != nil {
		return err
	}
	s.RemindAutoSend = input.RemindAutoSend
	s.RemindWeekday = input.RemindWeekday
	s.RemindTime = input.RemindTime
	return uc.settingsRepo.Update(s)
}

package remind_settings_usecase

import "testing"

func TestUpdateRemindInput_validation(t *testing.T) {
	uc := &RemindSettings{}
	if err := uc.Update(UpdateRemindInput{
		NotifyStartWeek:    1,
		NotifyIntervalDays: 1,
		NotifyTime:         "bad",
	}); err == nil {
		t.Fatal("expected invalid time error")
	}
	if err := uc.Update(UpdateRemindInput{
		NotifyStartWeek:    0,
		NotifyIntervalDays: 1,
		NotifyTime:         "09:00",
	}); err == nil {
		t.Fatal("expected invalid start week error")
	}
	if err := uc.Update(UpdateRemindInput{
		NotifyStartWeek:    1,
		NotifyIntervalDays: 6,
		NotifyTime:         "09:00",
	}); err == nil {
		t.Fatal("expected invalid interval error")
	}
}
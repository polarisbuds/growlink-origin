package remind_settings_usecase

import "testing"

func TestUpdateRemindInput_validation(t *testing.T) {
	uc := &RemindSettings{}
	if err := uc.Update(UpdateRemindInput{RemindAutoSend: true, RemindWeekday: 3, RemindTime: "bad"}); err == nil {
		t.Fatal("expected invalid time error")
	}
}

package utils

import (
	"testing"
	"time"
)

func TestGetHoliday(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name            string
		args            args
		wantIsHoliday   bool
		wantHolidayName string
	}{
		{
			name:            "通常系: 元旦は祝日",
			args:            args{t: time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)},
			wantIsHoliday:   true,
			wantHolidayName: "元旦",
		},
		{
			name:            "通常系: こどもの日は祝日",
			args:            args{t: time.Date(2026, time.May, 5, 0, 0, 0, 0, time.UTC)},
			wantIsHoliday:   true,
			wantHolidayName: "こどもの日",
		},
		{
			name:            "通常系: 平日は祝日ではない",
			args:            args{t: time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC)},
			wantIsHoliday:   false,
			wantHolidayName: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsHoliday, gotHolidayName := GetHoliday(tt.args.t)
			if gotIsHoliday != tt.wantIsHoliday {
				t.Errorf("GetHoliday() isHoliday = %v, want %v", gotIsHoliday, tt.wantIsHoliday)
			}
			if gotHolidayName != tt.wantHolidayName {
				t.Errorf("GetHoliday() holidayName = %v, want %v", gotHolidayName, tt.wantHolidayName)
			}
		})
	}
}

func TestIsHoliday(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "通常系: 元旦は祝日",
			args: args{t: time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)},
			want: true,
		},
		{
			name: "通常系: 平日は祝日ではない",
			args: args{t: time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHoliday(tt.args.t); got != tt.want {
				t.Errorf("IsHoliday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNextWeekdays(t *testing.T) {
	type args struct {
		day time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "通常系: 平日はその日を返す",
			args: args{day: time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC)},
			want: time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "通常系: 祝日の翌日が平日ならその日を返す",
			args: args{day: time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)},
			want: time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "通常系: 連休なら祝日が終わるまで進める",
			args: args{day: time.Date(2026, time.May, 3, 9, 0, 0, 0, time.UTC)},
			want: time.Date(2026, time.May, 7, 9, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetNextWeekdays(tt.args.day)
			if !got.Equal(tt.want) {
				t.Errorf("GetNextWeekdays() = %v, want %v", got, tt.want)
			}
		})
	}
}

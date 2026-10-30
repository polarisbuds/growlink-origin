package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestGetStartOfDay(t *testing.T) {
	type args struct {
		day time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "通常系: 2006-01-02 15:04:05 が当日 00:00:00 になる",
			args: args{day: time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)},
			want: time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "異常系: ゼロ値（nil相当）は year 1 の 00:00:00 になる",
			args: args{day: time.Time{}},
			want: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "異常系: マイナスの年月日は正規化後の日付の 00:00:00 になる",
			args: args{day: time.Date(-1, -1, -1, 15, 4, 5, 0, time.UTC)},
			want: time.Date(-2, time.October, 30, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStartOfDay(tt.args.day); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStartOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEndOfDay(t *testing.T) {
	type args struct {
		day time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "通常系: 2006-01-02 15:04:05 が当日 23:59:59.999999999 になる",
			args: args{day: time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)},
			want: time.Date(2006, time.January, 2, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name: "異常系: ゼロ値（nil相当）は year 1 の 23:59:59.999999999 になる",
			args: args{day: time.Time{}},
			want: time.Date(1, time.January, 1, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name: "異常系: マイナスの年月日は正規化後の日付の 23:59:59.999999999 になる",
			args: args{day: time.Date(-1, -1, -1, 15, 4, 5, 0, time.UTC)},
			want: time.Date(-2, time.October, 30, 23, 59, 59, 999999999, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEndOfDay(tt.args.day); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEndOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSpecificWeekday(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)

	type args struct {
		base    time.Time
		weekNum int
		weekday time.Weekday
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "通常系: 2006年1月第2月曜日",
			args: args{
				base:    time.Date(2006, time.January, 15, 15, 4, 5, 0, time.UTC),
				weekNum: 2,
				weekday: time.Monday,
			},
			want: time.Date(2006, time.January, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "通常系: 月初が対象曜日なら第1週はその1日",
			args: args{
				base:    time.Date(2026, time.September, 20, 12, 0, 0, 0, time.UTC),
				weekNum: 1,
				weekday: time.Tuesday,
			},
			want: time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "通常系: タイムゾーンは基準日の Location を維持する",
			args: args{
				base:    time.Date(2026, time.September, 10, 12, 0, 0, 0, jst),
				weekNum: 1,
				weekday: time.Wednesday,
			},
			want: time.Date(2026, time.September, 2, 0, 0, 0, 0, jst),
		},
		{
			name: "異常系: 第5週が月末を超えると翌月に正規化される",
			args: args{
				base:    time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC),
				weekNum: 5,
				weekday: time.Thursday,
			},
			want: time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSpecificWeekday(tt.args.base, tt.args.weekNum, tt.args.weekday)
			if !got.Equal(tt.want) || got.Location() != tt.want.Location() {
				t.Errorf("GetSpecificWeekday() = %v (%s), want %v (%s)", got, got.Location(), tt.want, tt.want.Location())
			}
		})
	}
}

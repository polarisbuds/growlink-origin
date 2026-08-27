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

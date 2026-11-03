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
		name    string
		args    args
		want    time.Time
		wantOK  bool
	}{
		{
			name: "通常系: 2006年1月第2月曜日",
			args: args{
				base:    time.Date(2006, time.January, 15, 15, 4, 5, 0, time.UTC),
				weekNum: 2,
				weekday: time.Monday,
			},
			want:   time.Date(2006, time.January, 9, 0, 0, 0, 0, time.UTC),
			wantOK: true,
		},
		{
			name: "通常系: 月初が対象曜日なら第1週はその1日",
			args: args{
				base:    time.Date(2026, time.September, 20, 12, 0, 0, 0, time.UTC),
				weekNum: 1,
				weekday: time.Tuesday,
			},
			want:   time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC),
			wantOK: true,
		},
		{
			name: "通常系: タイムゾーンは基準日の Location を維持する",
			args: args{
				base:    time.Date(2026, time.September, 10, 12, 0, 0, 0, jst),
				weekNum: 1,
				weekday: time.Wednesday,
			},
			want:   time.Date(2026, time.September, 2, 0, 0, 0, 0, jst),
			wantOK: true,
		},
		{
			name: "異常系: 第5週が月内に無い場合は false",
			args: args{
				base:    time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC),
				weekNum: 5,
				weekday: time.Thursday,
			},
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetSpecificWeekday(tt.args.base, tt.args.weekNum, tt.args.weekday)
			if ok != tt.wantOK {
				t.Errorf("GetSpecificWeekday() ok = %v, want %v", ok, tt.wantOK)
			}
			if tt.wantOK && (!got.Equal(tt.want) || got.Location() != tt.want.Location()) {
				t.Errorf("GetSpecificWeekday() = %v (%s), want %v (%s)", got, got.Location(), tt.want, tt.want.Location())
			}
		})
	}
}

func TestListSpecificWeekdays(t *testing.T) {
	from := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.November, 1, 0, 0, 0, 0, time.UTC)

	t.Run("両端の月を含む", func(t *testing.T) {
		got := ListSpecificWeekdays(from, to, 2, time.Monday)
		requireLen := 3
		if len(got) != requireLen {
			t.Fatalf("len = %d, want %d", len(got), requireLen)
		}
		wantSep, ok := GetSpecificWeekday(from, 2, time.Monday)
		if !ok {
			t.Fatal("September 2nd Monday should exist")
		}
		if !got[0].Equal(wantSep) {
			t.Errorf("got[0] = %v, want %v", got[0], wantSep)
		}
	})

	t.Run("第5週が無い月だけ欠ける", func(t *testing.T) {
		// 2026-09 から 2026-10: 第5木曜は9月のみ存在しない想定で件数を確認
		got := ListSpecificWeekdays(
			time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
			5,
			time.Thursday,
		)
		sep, sepOK := GetSpecificWeekday(time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC), 5, time.Thursday)
		oct, octOK := GetSpecificWeekday(time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC), 5, time.Thursday)
		_ = sep
		wantLen := 0
		if octOK {
			wantLen++
		}
		if !sepOK && wantLen == 1 && len(got) == 1 && got[0].Equal(oct) {
			return
		}
		if len(got) != wantLen {
			t.Errorf("len = %d, want %d (sepOK=%v octOK=%v)", len(got), wantLen, sepOK, octOK)
		}
	})
}

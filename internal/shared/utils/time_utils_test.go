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

func TestMatchesRemindInterval(t *testing.T) {
	loc := time.FixedZone("JST", 9*60*60)
	// 2026年11月: 第2月曜は 11/9（日曜は11/1）
	beforeStart := time.Date(2026, time.November, 8, 10, 0, 0, 0, loc)
	startDay := time.Date(2026, time.November, 9, 10, 0, 0, 0, loc)
	tuesday := time.Date(2026, time.November, 10, 10, 0, 0, 0, loc)
	nextMonday := time.Date(2026, time.November, 16, 10, 0, 0, 0, loc)

	t.Run("開始週前は不一致", func(t *testing.T) {
		if MatchesRemindInterval(beforeStart, 2, 1) {
			t.Fatal("expected false before notify start")
		}
	})
	t.Run("間隔1は開始日以降毎日", func(t *testing.T) {
		for _, d := range []time.Time{startDay, tuesday} {
			if !MatchesRemindInterval(d, 2, 1) {
				t.Fatalf("expected true for daily interval on %v", d)
			}
		}
	})
	t.Run("間隔2は月水金", func(t *testing.T) {
		if !MatchesRemindInterval(startDay, 2, 2) {
			t.Fatal("Monday should match interval 2")
		}
		if MatchesRemindInterval(tuesday, 2, 2) {
			t.Fatal("Tuesday should not match interval 2")
		}
		if !MatchesRemindInterval(nextMonday, 2, 2) {
			t.Fatal("following Monday should match interval 2")
		}
	})
	t.Run("最終週は当月最終月曜から", func(t *testing.T) {
		lastMon := LastMondayOfMonth(time.Date(2026, time.November, 15, 0, 0, 0, 0, loc))
		if lastMon.Day() != 30 {
			t.Fatalf("last Monday = %v, want Nov 30", lastMon)
		}
		start := NotifyStartDateInMonth(lastMon, 5)
		if !start.Equal(lastMon) {
			t.Fatalf("notify start = %v, want %v", start, lastMon)
		}
		if !MatchesRemindInterval(lastMon, 5, 1) {
			t.Fatal("last Monday should match with start week 5")
		}
	})
}

func TestNextMonthFirstAt(t *testing.T) {
	loc := time.FixedZone("JST", 9*60*60)
	t.Run("当月1日10時前は当月1日10時", func(t *testing.T) {
		now := time.Date(2026, time.November, 1, 9, 0, 0, 0, loc)
		next := NextMonthFirstAt(now, 10, 0)
		want := time.Date(2026, time.November, 1, 10, 0, 0, 0, loc)
		if !next.Equal(want) {
			t.Fatalf("got %v, want %v", next, want)
		}
	})
	t.Run("当月1日10時過ぎは翌月1日10時", func(t *testing.T) {
		now := time.Date(2026, time.November, 1, 10, 1, 0, 0, loc)
		next := NextMonthFirstAt(now, 10, 0)
		want := time.Date(2026, time.December, 1, 10, 0, 0, 0, loc)
		if !next.Equal(want) {
			t.Fatalf("got %v, want %v", next, want)
		}
	})
}

func TestNextRemindNotifyAt(t *testing.T) {
	loc := time.FixedZone("JST", 9*60*60)
	now := time.Date(2026, time.November, 25, 8, 30, 0, 0, loc)
	next := NextRemindNotifyAt(now, "09:00")
	want := time.Date(2026, time.November, 25, 9, 0, 0, 0, loc)
	if !next.Equal(want) {
		t.Fatalf("got %v, want %v", next, want)
	}
	after := time.Date(2026, time.November, 25, 9, 5, 0, 0, loc)
	next = NextRemindNotifyAt(after, "09:00")
	want = time.Date(2026, time.November, 26, 9, 0, 0, 0, loc)
	if !next.Equal(want) {
		t.Fatalf("got %v, want %v", next, want)
	}
}

package utils

import (
	"strconv"
	"strings"
	"time"
)

func GetStartOfDay(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location()) // 指定日の 00:00:00.00....
}

func GetEndOfDay(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 999999999, day.Location()) // 指定日の23:59:59.99....
}

// 指定した基準日のある年月の、週目と曜日を指定して日付を返す。
// その週が月内に存在しない場合は ok が false になる。
func GetSpecificWeekday(base time.Time, weekNum int, weekday time.Weekday) (time.Time, bool) {
	firstOfMonth := time.Date(base.Year(), base.Month(), 1, 0, 0, 0, 0, base.Location())
	offset := (int(weekday) - int(firstOfMonth.Weekday()) + 7) % 7
	targetDay := 1 + offset + (weekNum-1)*7
	result := time.Date(base.Year(), base.Month(), targetDay, 0, 0, 0, 0, base.Location())
	if result.Month() != base.Month() || result.Year() != base.Year() {
		return time.Time{}, false
	}
	return result, true
}

// ListSpecificWeekdays は from から to まで（各月の1日基準・両端の月を含む）の各月について
// 第 weekNum 曜日を求め、月内に存在する月だけ日付を返す。
func ListSpecificWeekdays(from, to time.Time, weekNum int, weekday time.Weekday) []time.Time {
	loc := from.Location()
	cursor := time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, loc)
	end := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, loc)
	out := []time.Time{}
	for !cursor.After(end) {
		if d, ok := GetSpecificWeekday(cursor, weekNum, weekday); ok {
			out = append(out, d)
		}
		cursor = cursor.AddDate(0, 1, 0)
	}
	return out
}

func ListSpecificMonth(from, to time.Time) []time.Time {
	loc := from.Location()
	cursor := time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, loc)
	end := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, loc)
	out := []time.Time{}
	for !cursor.After(end) {
		cursor = cursor.AddDate(0, 1, 0)
	}
	return out
}

// LastMondayOfMonth は base と同じ月の最後の月曜日（00:00）を返す。
func LastMondayOfMonth(base time.Time) time.Time {
	loc := base.Location()
	y, m := base.Year(), base.Month()
	last := time.Date(y, m+1, 0, 0, 0, 0, 0, loc)
	for last.Weekday() != time.Monday {
		last = last.AddDate(0, 0, -1)
	}
	return last
}

// NotifyStartDateInMonth は当月のリマインド送信開始日（第 startWeek 月曜、5 は最終週）を
// 土日祝を除いた翌営業日に調整して返す。
func NotifyStartDateInMonth(day time.Time, startWeek int) time.Time {
	var monday time.Time
	if startWeek >= 5 {
		monday = LastMondayOfMonth(day)
	} else if d, ok := GetSpecificWeekday(day, startWeek, time.Monday); ok {
		monday = d
	} else {
		monday = LastMondayOfMonth(day)
	}
	return GetStartOfDay(GetNextWeekdays(monday))
}

// DaysSinceMonday は週の月曜を 0 とした経過日数を返す。
func DaysSinceMonday(day time.Time) int {
	wd := day.Weekday()
	return (int(wd) + 6) % 7
}

// MatchesRemindInterval は送信開始週以降かつ週初めからの経過日数が間隔に一致するかを返す。
func MatchesRemindInterval(day time.Time, startWeek, intervalDays int) bool {
	if intervalDays < 1 {
		intervalDays = 1
	}
	today := GetStartOfDay(day)
	start := NotifyStartDateInMonth(day, startWeek)
	if today.Before(start) {
		return false
	}
	return DaysSinceMonday(today)%intervalDays == 0
}

// ParseNotifyClock は "HH:MM" 形式の時刻をパースする。
func ParseNotifyClock(notifyTime string) (hour, minute int, ok bool) {
	parts := strings.Split(strings.TrimSpace(notifyTime), ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	h, err := strconv.Atoi(parts[0])
	if err != nil || h < 0 || h > 23 {
		return 0, 0, false
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil || m < 0 || m > 59 {
		return 0, 0, false
	}
	return h, m, true
}

// NextMonthFirstAt は now より後の、当月または翌月の1日 hour:minute を返す。
func NextMonthFirstAt(now time.Time, hour, minute int) time.Time {
	loc := now.Location()
	next := time.Date(now.Year(), now.Month(), 1, hour, minute, 0, 0, loc)
	if !next.After(now) {
		next = time.Date(now.Year(), now.Month()+1, 1, hour, minute, 0, 0, loc)
	}
	return next
}

// NextRemindNotifyAt は now より後の、notifyTime の次回実行時刻を返す。
func NextRemindNotifyAt(now time.Time, notifyTime string) time.Time {
	hour, minute := 9, 0
	if h, m, parsed := ParseNotifyClock(notifyTime); parsed {
		hour, minute = h, m
	}
	loc := now.Location()
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, loc)
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}

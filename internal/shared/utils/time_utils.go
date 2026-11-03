package utils

import (
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

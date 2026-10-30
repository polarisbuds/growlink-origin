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

// 指定した基準日のある年月の、週目と曜日を指定して日付データを受け取る
func GetSpecificWeekday(base time.Time, weekNum int, weekday time.Weekday) time.Time {
	firstOfMonth := time.Date(base.Year(), base.Month(), 1, 0, 0, 0, 0, base.Location())
	offset := (int(weekday) - int(firstOfMonth.Weekday()) + 7) % 7
	targetDay := 1 + offset + (weekNum-1)*7
	return time.Date(base.Year(), base.Month(), targetDay, 0, 0, 0, 0, base.Location())
}

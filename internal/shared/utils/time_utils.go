package utils

import "time"

func GetStartOfDay(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location()) // 指定日の 00:00:00.00....
}

func GetEndOfDay(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 999999999, day.Location()) // 指定日の23:59:59.99....
}

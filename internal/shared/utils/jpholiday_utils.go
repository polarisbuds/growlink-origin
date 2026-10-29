package utils

import (
	"time"

	"github.com/kokardy/jpholiday"
)

func GetHoliday(t time.Time) (isHoliday bool, holidayName string) {
	isHoliday, holiday := jpholiday.TimeToDate(t).Holiday()
	return isHoliday, holiday.String()
}

func IsHoliday(t time.Time) (isHoliday bool) {
	isHoliday, _ = jpholiday.TimeToDate(t).Holiday()
	return isHoliday
}

func GetNextWeekdays(day time.Time) (nextWeekDays time.Time) {
	nextWeekDays = day
	for IsHoliday(nextWeekDays) {
		nextWeekDays = nextWeekDays.AddDate(0, 0, 1)
	}
	return nextWeekDays
}

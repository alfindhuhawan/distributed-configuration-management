package util

import "time"

func CreateDeadline(hour time.Duration, dateTime time.Time) time.Time {
	return dateTime.Add(time.Hour * (hour * 24))
}

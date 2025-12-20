package util

import (
	"time"
)

func GetDayIndex(day string) int32 {
	dayIndex := make(map[string]int32, 0)

	dayIndex["Mon"] = 0
	dayIndex["Tue"] = 1
	dayIndex["Wed"] = 2
	dayIndex["Thu"] = 3
	dayIndex["Fri"] = 4
	dayIndex["Sat"] = 5
	dayIndex["Sun"] = 6

	return dayIndex[day]
}

func GetIsSellerClosed(day string, available bool, open string, close string) bool {
	isSellerClosed := true

	if available {
		if open == "00:00" && close == "00:00" {
			isSellerClosed = false
		} else {
			now := time.Now().In(LocationJakarta()).Format("15:04")
			timeNow, _ := time.Parse("15:04", now)
			timeOpen, _ := time.Parse("15:04", open)
			timeClose, _ := time.Parse("15:04", close)

			isTimeOpen := timeNow.After(timeOpen)
			isTimeClose := timeNow.Before(timeClose)

			if isTimeOpen && isTimeClose {
				isSellerClosed = false
			}
		}
	}

	return isSellerClosed
}

// TODO : future update, change to getTimeNow >> time.Now().In(util.LocationJakarta())
func LocationJakarta() *time.Location {
	location, _ := time.LoadLocation("Asia/Jakarta")

	return location
}

func LocalDateToUtc(localDate string) (time.Time, error) {
	result := time.Time{}

	// CONVERT STRING DATE TO LOCAL TIME
	local_loc, _ := time.LoadLocation("Asia/Jakarta")
	localTime, err := time.ParseInLocation("2006-01-02", localDate, local_loc)
	if err != nil {
		return result, err
	}

	// CONVERT LOCAL TIME TO UTC
	Utc_loc, _ := time.LoadLocation("UTC")
	result = localTime.In(Utc_loc)

	return result, nil
}

func ParseToStartOfDay(date string, timeZone string) (time.Time, error) {
	// CONVERT STRING DATE TO LOCAL TIME
	if timeZone == "" {
		timeZone = "UTC" // set default timezone
	}
	loc, _ := time.LoadLocation(timeZone)
	dateTime, err := time.ParseInLocation("2006-01-02", date, loc)
	result := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, loc)
	if err != nil {
		return time.Time{}, err
	}

	return result, nil
}

func GetStartOfDay(dateTime time.Time, timeZone string) (time.Time, error) {
	if timeZone == "" {
		timeZone = "UTC" // set default timezone
	}
	loc, _ := time.LoadLocation(timeZone)
	result := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, loc)
	return result, nil
}

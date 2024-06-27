package utils

import (
	"time"

	persian "github.com/yaa110/go-persian-calendar"
)

// ParseTime parses the given UTC date string into a time.Time object,
// converts it to Iran/Tehran time zone, and returns the Jalali calendar date as a string.
func ParseTime(utcDate string) (time.Time, time.Time, string, error) {
	// Parse the UTC date string
	utcTime, err := time.Parse(time.RFC3339, utcDate)
	if err != nil {
		return time.Time{}, time.Time{}, "", err
	}

	// Convert to Iran/Tehran time zone
	iranTime, err := convertToIranTime(utcTime)
	if err != nil {
		return time.Time{}, time.Time{}, "", err
	}

	// Convert to Jalali (Persian) calendar date
	jalaliDate := persian.New(iranTime).Format("yyyy/MM/dd")

	// The original time (Spain's local time) is just the parsed UTC time
	return utcTime, iranTime, jalaliDate, nil
}

func convertToIranTime(utcTime time.Time) (time.Time, error) {
	iranLocation, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return time.Time{}, err
	}
	return utcTime.In(iranLocation), nil
}

// IsYesterday checks if the given date was yesterday.
func IsYesterday(matchTime time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	// yesterday := time.Date(2024, 9, 1, 15, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	return matchTime.Year() == yesterday.Year() && matchTime.YearDay() == yesterday.YearDay()
}

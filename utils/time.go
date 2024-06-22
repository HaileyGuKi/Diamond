package utils

import (
	"time"
)

// ParseTime parses the given UTC date string into a time.Time object
// and converts it to Iran/Tehran time zone.
func ParseTime(utcDate string) (time.Time, error) {
	// Parse the UTC date string
	utcTime, err := time.Parse(time.RFC3339, utcDate)
	if err != nil {
		return time.Time{}, err
	}

	// Convert to Iran/Tehran time zone
	iranTime, err := convertToIranTime(utcTime)
	if err != nil {
		return time.Time{}, err
	}

	return iranTime, nil
}

// convertToIranTime converts the given UTC time to Iran/Tehran time zone.
func convertToIranTime(utcTime time.Time) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return time.Time{}, err
	}

	return utcTime.In(loc), nil
}

// IsYesterday checks if the given date was yesterday.
func IsYesterday(matchTime time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	// yesterday := time.Date(2024, 9, 1, 15, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	return matchTime.Year() == yesterday.Year() && matchTime.YearDay() == yesterday.YearDay()
}

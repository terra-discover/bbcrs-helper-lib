package lib

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
)

// CurrentTime func
func CurrentTime(format string) string {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	return time.Now().Format(format)
}

// TimeNow func
func TimeNow() *time.Time {
	t := time.Now()
	return &t
}

// StrfmtNow func
func StrfmtNow() *strfmt.DateTime {
	return DateTimeptr(strfmt.DateTime(*TimeNow()))
}

// RangeDate func
func RangeDate(date1, date2, types string) float64 {
	// date_1 > date_2
	// format_date : 2006-01-02 15:04:05
	var result float64

	format := "2006-01-02 15:04:05"
	dateOne, _ := time.Parse(format, date1)
	dateTwo, _ := time.Parse(format, date2)

	diff := dateOne.Sub(dateTwo)
	// number of Hours
	if types == "hours" {
		result = diff.Hours()

		// number of Nanoseconds
	} else if types == "nanoseconds" {
		result = float64(diff.Nanoseconds())

		// number of Minutes
	} else if types == "minutes" {
		result = diff.Minutes()

		// number of Seconds
	} else if types == "seconds" {
		result = diff.Seconds()

		// number of Days
	} else if types == "days" {
		result = float64(diff.Hours() / 24)
	}
	return result
}

// AddDate func
func AddDate(fromDate, format string, years, month, days int) string {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	date, _ := time.Parse(format, fromDate)
	t2 := date.AddDate(years, month, days)
	CurrentTimeAhead := t2.Format(format)
	return CurrentTimeAhead
}

// TimeStringToDuration will convert any time string to proper time.Duration
func TimeStringToDuration(notnormalTimeString string) (dur time.Duration) {
	splitNorNormalTime := strings.Split(notnormalTimeString, ":")

	hour := 0
	minute := 0
	second := 0

	if len(splitNorNormalTime) >= 1 {
		hour, _ = strconv.Atoi(splitNorNormalTime[0])
	}
	if len(splitNorNormalTime) >= 2 {
		minute, _ = strconv.Atoi(splitNorNormalTime[1])
	}
	if len(splitNorNormalTime) >= 3 {
		second, _ = strconv.Atoi(splitNorNormalTime[2])
	}

	seconds := (hour * 60 * 60) + (minute * 60) + second
	dur = time.Duration(seconds) * time.Second
	return
}

// UnixDurationToHumanDuration will return unix second to string "h m"
func UnixDurationToHumanDuration(unixTime int64) string {
	// Ex : 3915 (3600 + 300 + 15) -> 1h 5m
	if unixTime < 0 {
		unixTime = unixTime * -1
	}

	second := int(unixTime % 60)
	hour := int(math.Floor(float64(unixTime) / 3600))
	minute := int((int(unixTime) - (hour * 3600) - second) / 60)
	day := 0
	if hour >= 24 {
		day = hour / 24
		hour = hour % 24
	}

	if day > 0 {
		return fmt.Sprintf("%dd %dh %dm", day, hour, minute)
	}
	return fmt.Sprintf("%dh %dm", hour, minute)
}

// CalculateAgeByDate
func CalculateAgeByDate(x string, compareTime ...time.Time) int {
	now := time.Now()
	if len(compareTime) > 0 {
		now = compareTime[0]
	}

	birthDay, _ := time.Parse("2006-01-02", x)
	age := now.Year() - birthDay.Year()

	// additional age logic for specific month & day
	if now.Month() < birthDay.Month() || (now.Month() == birthDay.Month() && now.Day() < birthDay.Day()) {
		age = age - 1
	}

	if now.Before(birthDay.AddDate(age, 0, 0)) {
		age = age - 1
	}

	if age < 0 {
		age = 0
	}

	return age
}

// ElapsedTime
func ElapsedTime(startDate, endDate time.Time, types string) int {
	var result int
	elapsedTime := endDate.Sub(startDate)
	if types == "hours" {
		result = int(elapsedTime.Hours())
	} else if types == "nanoseconds" {
		result = int(elapsedTime.Nanoseconds())
	} else if types == "minutes" {
		result = int(elapsedTime.Minutes())
	} else if types == "seconds" {
		result = int(elapsedTime.Seconds())
	} else if types == "days" {
		result = int(elapsedTime.Hours() / 24)
	}
	return result
}

// ParseDateTime function to safely convert string to strfmt.DateTime
func ParseDateTime(dateTimeString string) strfmt.DateTime {
	parsedDateTime, err := strfmt.ParseDateTime(dateTimeString)
	if err != nil {
		// Handle the error, e.g., log it, set a default value, etc.
		return strfmt.DateTime(time.Time{})
	}
	return parsedDateTime
}

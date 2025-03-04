package lib

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofiber/fiber/v2/utils"
)

func TestCurrentTime(t *testing.T) {
	timeWithFormat := CurrentTime("2006-01-02 15:04:05")
	timeDefault := CurrentTime("")
	utils.AssertEqual(t, timeWithFormat, timeDefault)
}

func TestTimeNow(t *testing.T) {
	TimeNow()
}

func TestStrfmtNow(t *testing.T) {
	StrfmtNow()
}

func TestRangeDate(t *testing.T) {
	RangeDate("2022-01-02 15:04:05", "2006-01-05 15:04:05", "days")
	RangeDate("2022-01-02 15:04:05", "2006-01-05 15:04:05", "hours")
	RangeDate("2022-01-02 15:04:05", "2006-01-05 15:04:05", "nanoseconds")
	RangeDate("2022-01-02 15:04:05", "2006-01-05 15:04:05", "minutes")
	RangeDate("2022-01-02 15:04:05", "2006-01-05 15:04:05", "seconds")
}

func TestAddDate(t *testing.T) {
	AddDate("2022-01-02 15:04:05", "2022-01-02 15:04:05", 0, 1, 1)
	AddDate("2022-01-02 15:04:05", "", 0, 1, 1)
}

func TestTimeStringToDuration(t *testing.T) {
	TimeStringToDuration("12:10:20")
}

func TestUnixDurationToHumanDuration(t *testing.T) {
	resp := UnixDurationToHumanDuration(int64(3915))
	utils.AssertEqual(t, "1h 5m", resp, "Unix Duration to Human Duration")

	resp2 := UnixDurationToHumanDuration(int64(0))
	utils.AssertEqual(t, "0h 0m", resp2, "Unix Duration to Human Duration")

	resp3 := UnixDurationToHumanDuration(int64(-3915))
	utils.AssertEqual(t, "1h 5m", resp3, "Unix Duration to Human Duration")

	resp4 := UnixDurationToHumanDuration(int64(86400 + 3915))
	utils.AssertEqual(t, "1d 1h 5m", resp4, "Unix Duration to Human Duration")
}

func TestCalculateAgeByDate(t *testing.T) {
	now := time.Now()

	a := CalculateAgeByDate("1995-08-10")
	utils.AssertEqual(t, a >= 20, true, "Calculate Age By Date Normal case")

	b := CalculateAgeByDate(now.Format("2006-01-02"))
	if b != 0 {
		b = CalculateAgeByDate(now.AddDate(0, 0, -1).Format("2006-01-02"))
	}
	utils.AssertEqual(t, 0, b, "Calculate Age By Date zero case")

	c := CalculateAgeByDate(now.AddDate(0, 0, 1).Format("2006-01-02"))
	utils.AssertEqual(t, 0, c, "Calculate Age By Date greather date case")

	tt, _ := time.Parse("2006-01-02", "2018-01-02")
	d := CalculateAgeByDate("2015-01-02", tt)
	utils.AssertEqual(t, 3, d, "Calculate Age By Date by custom parameter")

	uu, _ := time.Parse("2006-01-02", "2023-05-18")
	e := CalculateAgeByDate("2012-05-17", uu)
	utils.AssertEqual(t, 11, e, "Calculate Age By Date same month only")

	vv, _ := time.Parse("2006-01-02", "2023-05-18")
	f := CalculateAgeByDate("2022-05-18", vv)
	utils.AssertEqual(t, 1, f, "Calculate Age By Date same month and date")

	ww, _ := time.Parse("2006-01-02", "2023-05-18")
	g := CalculateAgeByDate("2012-05-19", ww)
	utils.AssertEqual(t, 10, g, "Calculate Age By Date same month but not birthday yet")
}

func TestElapsedTime(t *testing.T) {
	startTime := time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 4, 1, 1, 0, 0, 0, time.UTC)

	// Test hours
	expectedHours := 1
	utils.AssertEqual(t, expectedHours, ElapsedTime(startTime, endTime, "hours"))

	// Test nanoseconds
	expectedNanoseconds := int(time.Hour / time.Nanosecond)
	utils.AssertEqual(t, expectedNanoseconds, ElapsedTime(startTime, endTime, "nanoseconds"))

	// Test minutes
	expectedMinutes := 60
	utils.AssertEqual(t, expectedMinutes, ElapsedTime(startTime, endTime, "minutes"))

	// Test seconds
	expectedSeconds := 3600
	utils.AssertEqual(t, expectedSeconds, ElapsedTime(startTime, endTime, "seconds"))

	// Test days
	endTime = time.Date(2023, 4, 3, 0, 0, 0, 0, time.UTC)
	expectedDays := 2
	utils.AssertEqual(t, expectedDays, ElapsedTime(startTime, endTime, "days"))
}

func TestParseDateTime(t *testing.T) {
	// Test case 1: Valid date/time string
	validDateTimeString := "2023-01-15T12:30:00Z"
	expectedParsedDateTime := strfmt.DateTime(time.Date(2023, 1, 15, 12, 30, 0, 0, time.UTC))
	actualParsedDateTime := ParseDateTime(validDateTimeString)
	utils.AssertEqual(t, expectedParsedDateTime, actualParsedDateTime, "Parsed date/time should match")

	// Test case 2: Invalid date/time string
	invalidDateTimeString := "invalid_datetime_string"
	expectedDefaultDateTime := strfmt.DateTime(time.Time{})
	actualDefaultDateTime := ParseDateTime(invalidDateTimeString)
	utils.AssertEqual(t, expectedDefaultDateTime, actualDefaultDateTime, "Invalid date/time should return default")
}

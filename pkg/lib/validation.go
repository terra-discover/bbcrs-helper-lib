package lib

import (
	"regexp"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IsEmptyFloat64Ptr
func IsEmptyFloat64Ptr(number *float64) (isEmpty bool) {
	if number == nil || *number <= 0 {
		isEmpty = true
	}

	return
}

// IsEmptyFloat64
func IsEmptyFloat64(number float64) (isEmpty bool) {
	if number <= 0 {
		isEmpty = true
	}

	return
}

// IsEmptyIntPtr
func IsEmptyIntPtr(number *int) (isEmpty bool) {
	if number == nil || *number <= 0 {
		isEmpty = true
	}

	return
}

// IsEmptyInt
func IsEmptyInt(number int) (isEmpty bool) {
	if number <= 0 {
		isEmpty = true
	}

	return
}

// IsEmptyInt64Ptr
func IsEmptyInt64Ptr(number *int64) (isEmpty bool) {
	if number == nil || *number <= 0 {
		isEmpty = true
	}

	return
}

// IsEmptyInt64
func IsEmptyInt64(number int64) (isEmpty bool) {
	if number <= 0 {
		isEmpty = true
	}

	return
}

// IsEmptyStrPtr
func IsEmptyStrPtr(str *string) (isEmpty bool) {
	if str == nil || len(strings.TrimSpace(*str)) == 0 {
		isEmpty = true
	}

	return
}

// IsEmptyStr
func IsEmptyStr(str string) (isEmpty bool) {
	if len(strings.TrimSpace(str)) == 0 {
		isEmpty = true
	}

	return
}

// IsFalsyBoolPtr
func IsFalsyBoolPtr(cond *bool) (isFalsy bool) {
	if cond == nil || !(*cond) {
		isFalsy = true
	}

	return
}

// IsEmptyUUIDPtr
func IsEmptyUUIDPtr(id *uuid.UUID) (isEmpty bool) {
	if id == nil || *id == uuid.Nil {
		isEmpty = true
	}

	return
}

// IsEmptyUUID
func IsEmptyUUID(id uuid.UUID) (isEmpty bool) {
	if id == uuid.Nil {
		isEmpty = true
	}

	return
}

// IsZeroTimePtr
func IsZeroTimePtr(moment *time.Time) (isZero bool) {
	if moment == nil || (*moment).IsZero() {
		isZero = true
	}

	return
}

// IsZeroTime
func IsZeroTime(moment time.Time) (isZero bool) {
	if moment.IsZero() {
		isZero = true
	}

	return
}

// IsZeroStrfmtTimePtr
func IsZeroStrfmtTimePtr(moment *strfmt.DateTime) (isZero bool) {
	if moment == nil || time.Time(*moment).IsZero() {
		isZero = true
	}

	return
}

// IsZeroStrfmtTime
func IsZeroStrfmtTime(moment strfmt.DateTime) (isZero bool) {
	if time.Time(moment).IsZero() {
		isZero = true
	}

	return
}

// IsSimilarStringPattern
func IsSimilarStringPattern(a, b string) (isSimilar bool) {
	// Create a regex pattern to match the "a" prefix
	pattern := "^" + a
	// Match the strings using the regex pattern
	match, _ := regexp.MatchString(pattern, b)
	isSimilar = match
	return
}

func MustReturnErrDB(err error) bool {
	return err != nil && err != gorm.ErrRecordNotFound
}

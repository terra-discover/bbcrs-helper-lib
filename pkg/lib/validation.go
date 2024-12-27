package lib

import (
	"strings"
	"time"

	"github.com/google/uuid"
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

// IsFalsyBool
func IsFalsyBool(cond bool) (isFalsy bool) {
	if !(cond) {
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

// ContainsDuplicatedUUID
func ContainsDuplicatedUUID(listUUID []uuid.UUID) (isDuplicated bool) {
	if len(listUUID) == 0 {
		return
	}

	tempListUUID := []uuid.UUID{}

	for _, item := range listUUID {
		for _, tempItem := range tempListUUID {
			if tempItem == item {
				isDuplicated = true
				return
			}
		}

		tempListUUID = append(tempListUUID, item)
	}

	return
}

// ContainsDuplicatedFoldStrPtr - check duplicated strptr using case insensitive like strings.EqualFold()
func ContainsDuplicatedFoldStrPtr(listStr []*string) (isDuplicated bool) {
	if len(listStr) == 0 {
		return
	}

	tempListUUID := []string{}

	for _, item := range listStr {
		if item == nil {
			continue
		}

		for _, tempItem := range tempListUUID {
			if strings.EqualFold(tempItem, *item) {
				isDuplicated = true
				return
			}
		}

		tempListUUID = append(tempListUUID, *item)
	}

	return
}

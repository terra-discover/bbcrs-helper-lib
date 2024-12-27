package lib

import (
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

// UUIDPtr to return memory address
func UUIDPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

// Intptr to return memory address
func Intptr(i int) *int {
	return &i
}

// Int64ptr to return memory address
func Int64ptr(i int64) *int64 {
	return &i
}

// Boolptr to return memory address
func Boolptr(b bool) *bool {
	return &b
}

// Strptr to return memory address
func Strptr(s string) *string {
	return &s
}

// Float64ptr to return memory address
func Float64ptr(f float64) *float64 {
	return &f
}

// Dateptr to return memory address
func Dateptr(d strfmt.Date) *strfmt.Date {
	return &d
}

// DateTimeptr to return memory address
func DateTimeptr(dt strfmt.DateTime) *strfmt.DateTime {
	return &dt
}

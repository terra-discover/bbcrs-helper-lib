package lib

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
)

func TestUUIDPtr(t *testing.T) {
	var a uuid.UUID = uuid.New()
	b := UUIDPtr(a)
	utils.AssertEqual(t, a, *b)
}

func TestIntptr(t *testing.T) {
	var a int = 1
	b := Intptr(a)
	utils.AssertEqual(t, a, *b)
}

func TestInt64ptr(t *testing.T) {
	var a int64 = 1
	b := Int64ptr(a)
	utils.AssertEqual(t, a, *b)
}

func TestStrptr(t *testing.T) {
	var a string = "1"
	b := Strptr(a)
	utils.AssertEqual(t, a, *b)
}

func TestBoolptr(t *testing.T) {
	b := Boolptr(true)
	utils.AssertEqual(t, true, *b)
}

func TestFloat64ptr(t *testing.T) {
	var a float64 = 1.2
	b := Float64ptr(a)
	utils.AssertEqual(t, a, *b)
}

func TestDateptr(t *testing.T) {
	now := time.Now()
	a := (strfmt.Date)(now)
	b := Dateptr(a)
	utils.AssertEqual(t, a, *b)
}
func TestDateTimeptr(t *testing.T) {
	now := time.Now()
	a := (strfmt.DateTime)(now)
	b := DateTimeptr(a)
	utils.AssertEqual(t, a, *b)
}

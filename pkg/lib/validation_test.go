package lib

import (
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
)

func TestIsEmptyFloat64Ptr(t *testing.T) {
	var number *float64
	res := IsEmptyFloat64Ptr(number)
	utils.AssertEqual(t, true, res)

	number = Float64ptr(1)
	res = IsEmptyFloat64Ptr(number)
	utils.AssertEqual(t, false, res)

}

func TestIsEmptyFloat64(t *testing.T) {
	var number float64
	res := IsEmptyFloat64(number)
	utils.AssertEqual(t, true, res)

	number = 1
	res = IsEmptyFloat64(number)
	utils.AssertEqual(t, false, res)
}

func TestIsEmptyIntPtr(t *testing.T) {
	var number *int
	res := IsEmptyIntPtr(number)
	utils.AssertEqual(t, true, res)

	number = Intptr(1)
	res = IsEmptyIntPtr(number)
	utils.AssertEqual(t, false, res)
}

func TestIsEmptyInt(t *testing.T) {
	var number int
	res := IsEmptyInt(number)
	utils.AssertEqual(t, true, res)

	number = 1
	res = IsEmptyInt(number)
	utils.AssertEqual(t, false, res)
}

func TestIsEmptyStrPtr(t *testing.T) {
	var str *string
	res := IsEmptyStrPtr(str)
	utils.AssertEqual(t, true, res)

	str = Strptr("1")
	res = IsEmptyStrPtr(str)
	utils.AssertEqual(t, false, res)
}

func TestIsEmptyStr(t *testing.T) {
	var str string
	res := IsEmptyStr(str)
	utils.AssertEqual(t, true, res)

	str = "1"
	res = IsEmptyStr(str)
	utils.AssertEqual(t, false, res)
}

func TestIsFalsyBoolPtr(t *testing.T) {
	var isBool *bool
	res := IsFalsyBoolPtr(isBool)
	utils.AssertEqual(t, true, res)

	isBool = Boolptr(true)
	res = IsFalsyBoolPtr(isBool)
	utils.AssertEqual(t, false, res)
}

func TestIsFalsyBool(t *testing.T) {
	var isBool bool
	res := IsFalsyBool(isBool)
	utils.AssertEqual(t, true, res)

	isBool = true
	res = IsFalsyBool(isBool)
	utils.AssertEqual(t, false, res)

}

func TestIsEmptyUUIDPtr(t *testing.T) {
	var id *uuid.UUID
	res := IsEmptyUUIDPtr(id)
	utils.AssertEqual(t, true, res)

	id = GenUUID()
	res = IsEmptyUUIDPtr(id)
	utils.AssertEqual(t, false, res)
}

func TestIsEmptyUUID(t *testing.T) {
	var id uuid.UUID
	res := IsEmptyUUID(id)
	utils.AssertEqual(t, true, res)

	id = *GenUUID()
	res = IsEmptyUUID(id)
	utils.AssertEqual(t, false, res)
}

func TestIsZeroTimePtr(t *testing.T) {
	now := time.Now()

	type args struct {
		moment *time.Time
	}
	tests := []struct {
		name       string
		args       args
		wantIsZero bool
	}{
		// Add test cases.
		{
			name: "must zero",
			args: args{
				moment: &time.Time{},
			},
			wantIsZero: true,
		},
		{
			name: "not zero",
			args: args{
				moment: &now,
			},
			wantIsZero: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsZero := IsZeroTimePtr(tt.args.moment); gotIsZero != tt.wantIsZero {
				t.Errorf("IsZeroTimePtr() = %v, want %v", gotIsZero, tt.wantIsZero)
			}
		})
	}
}

func TestContainsDuplicatedUUID(t *testing.T) {
	uuid1 := *GenUUID()
	uuid2 := *GenUUID()

	type args struct {
		listUUID []uuid.UUID
	}
	tests := []struct {
		name             string
		args             args
		wantIsDuplicated bool
	}{
		{
			name: "success, contains duplicated data",
			args: args{
				listUUID: []uuid.UUID{
					uuid1,
					uuid2,
					uuid1,
				},
			},
			wantIsDuplicated: true,
		},
		{
			name: "success, not contains duplicated data",
			args: args{
				listUUID: []uuid.UUID{
					uuid1,
					uuid2,
				},
			},
			wantIsDuplicated: false,
		},
		{
			name: "success, empty input data",
			args: args{
				listUUID: []uuid.UUID{},
			},
			wantIsDuplicated: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsDuplicated := ContainsDuplicatedUUID(tt.args.listUUID); gotIsDuplicated != tt.wantIsDuplicated {
				t.Errorf("ContainsDuplicatedUUID() = %v, want %v", gotIsDuplicated, tt.wantIsDuplicated)
			}
		})
	}
}

func TestContainsDuplicatedFoldStrPtr(t *testing.T) {
	str1 := "abc"
	str2 := "bcd"
	str3 := strings.ToUpper(str1)

	type args struct {
		listStr []*string
	}
	tests := []struct {
		name             string
		args             args
		wantIsDuplicated bool
	}{
		{
			name: "success, contains duplicated data",
			args: args{
				listStr: []*string{
					&str1,
					&str1,
					&str2,
				},
			},
			wantIsDuplicated: true,
		},
		{
			name: "success, contains duplicated data (check case-insensitive)",
			args: args{
				listStr: []*string{
					&str1,
					&str3,
				},
			},
			wantIsDuplicated: true,
		},
		{
			name: "success, not contains duplicated data",
			args: args{
				listStr: []*string{
					&str1,
					&str2,
				},
			},
			wantIsDuplicated: false,
		},
		{
			name: "success, empty input data",
			args: args{
				listStr: []*string{},
			},
			wantIsDuplicated: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsDuplicated := ContainsDuplicatedFoldStrPtr(tt.args.listStr); gotIsDuplicated != tt.wantIsDuplicated {
				t.Errorf("ContainsDuplicatedFoldStrPtr() = %v, want %v", gotIsDuplicated, tt.wantIsDuplicated)
			}
		})
	}
}

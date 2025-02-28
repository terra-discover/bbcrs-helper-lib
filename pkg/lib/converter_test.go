package lib

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
)

func TestConvertToMD5(t *testing.T) {
	value := 1
	ConvertToMD5(&value)
}

func TestConvertStrToMD5(t *testing.T) {
	value := "development usage"
	gen := ConvertStrToMD5(&value)
	gen2 := ConvertStrToMD5(&value)
	utils.AssertEqual(t, gen2, gen)
}

func TestConvertToSHA1(t *testing.T) {
	value := "development usage"
	ConvertToSHA1(value)
}

func TestConvertToSHA256(t *testing.T) {
	value := "development usage"
	ConvertToSHA256(value)
}

func TestIntToStr(t *testing.T) {
	value := 1
	res := IntToStr(value)
	utils.AssertEqual(t, "1", res)
}

func TestStrToInt(t *testing.T) {
	value := "1"
	res := StrToInt(value)
	utils.AssertEqual(t, 1, res)
}

func TestStrToInt64(t *testing.T) {
	value := "1"
	res := StrToInt64(value)
	utils.AssertEqual(t, int64(1), res)
}

func TestStrToFloat(t *testing.T) {
	value := "1"
	res := StrToFloat(value)
	utils.AssertEqual(t, float64(1), res)
}

func TestStrToBool(t *testing.T) {
	value := "true"
	res := StrToBool(value)
	utils.AssertEqual(t, true, res)
}

func TestFloatToStr(t *testing.T) {
	value := 1.2
	res := FloatToStr(value, 6)
	utils.AssertEqual(t, "1.200000", res)
}

func TestConvertJsonToStr(t *testing.T) {
	value := []interface{}{"first", "second"}
	res := ConvertJSONToStr(value)
	utils.AssertEqual(t, `["first","second"]`, res)
}

func TestConvertStrToObj(t *testing.T) {
	value := `{"index":"value"}`
	res := ConvertStrToObj(value)
	utils.AssertEqual(t, "value", res["index"])
}

func TestConvertStrToJson(t *testing.T) {
	expect := map[string]interface{}{
		"index": "value",
	}
	value := `{"index":"value"}`
	res := ConvertStrToJSON(value)
	utils.AssertEqual(t, expect, res)
}

func TestConvertStrToTime(t *testing.T) {
	value := "2021-05-19 11:56:30"
	gen := ConvertStrToTime(value)
	utils.AssertEqual(t, gen, gen)
}

func TestConvertStrToTimeWFormat(t *testing.T) {
	// Define test cases
	tests := []struct {
		input    string
		layout   string
		expected time.Time
	}{
		{
			input:    "2025-02-28 19:48:59",
			layout:   "2006-01-02 15:04:05",
			expected: time.Date(2025, 2, 28, 19, 48, 59, 0, time.UTC),
		},
		{
			input:    "28/02/2025",
			layout:   "02/01/2006",
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "02-28-2025",
			layout:   "01-02-2006",
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "19:48:59",
			layout:   "15:04:05",
			expected: time.Date(0, 1, 1, 19, 48, 59, 0, time.UTC),
		},
	}

	for _, test := range tests {
		result := ConvertStrToTimeWFormat(test.input, test.layout)
		if !result.Equal(test.expected) {
			t.Errorf("For input '%s' with layout '%s', expected %v but got %v",
				test.input, test.layout, test.expected, result)
		}
	}
}

func TestConvertSliceIntToStr(t *testing.T) {
	value := []int{1, 2, 3, 4}
	res := ConvertSliceIntToStr(value, ",")
	utils.AssertEqual(t, "1,2,3,4", res)
}

func TestConvertStrToArrObj(t *testing.T) {
	value := `[{"index":"value"}]`
	res := ConvertStrToArrObj(value)
	utils.AssertEqual(t, "value", res[0]["index"])
}

func TestConvertSliceStrToStr(t *testing.T) {
	value := []string{"active", "inactive", "suspend"}
	res := ConvertSliceStrToStr(value, ",")
	utils.AssertEqual(t, "active,inactive,suspend", res)
}

func TestConvertSliceUUIDToSliceStr(t *testing.T) {
	uuid1 := uuid.New()
	uuid2 := uuid.New()
	uuid3 := uuid.New()

	type args struct {
		listUUID []uuid.UUID
	}
	tests := []struct {
		name       string
		args       args
		wantResult []string
	}{
		{
			name: "slice uuid converted",
			args: args{
				listUUID: []uuid.UUID{
					uuid1,
					uuid2,
					uuid3,
				},
			},
			wantResult: []string{
				uuid1.String(),
				uuid2.String(),
				uuid3.String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ConvertSliceUUIDToSliceStr(tt.args.listUUID); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ConvertSliceUUIDToSliceStr() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestStrLeadingZerosRemove(t *testing.T) {
	str := "026"
	res := StrLeadingZerosRemove(str)
	utils.AssertEqual(t, true, len(res) > 0)
}

func TestForceStr(t *testing.T) {
	str1 := "abcd"

	type args struct {
		input *string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "filled input == result",
			args: args{
				input: &str1,
			},
			wantResult: str1,
		},
		{
			name: "nil input; result = empty string",
			args: args{
				input: nil,
			},
			wantResult: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ForceStr(tt.args.input); gotResult != tt.wantResult {
				t.Errorf("ForceStr() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestForceStrPtr(t *testing.T) {
	str1 := "abcd"
	emptyStr := ""

	type args struct {
		input *string
	}
	tests := []struct {
		name       string
		args       args
		wantResult *string
	}{
		{
			name: "filled input == result",
			args: args{
				input: &str1,
			},
			wantResult: &str1,
		},
		{
			name: "nil input; result = empty string",
			args: args{
				input: nil,
			},
			wantResult: &emptyStr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := ForceStrPtr(tt.args.input)
			utils.AssertEqual(t, tt.wantResult != nil, gotResult != nil, "compare pointer")
			if tt.wantResult != nil && gotResult != nil {
				utils.AssertEqual(t, *tt.wantResult, *gotResult, "compare value")
			}
		})
	}
}

func TestForceInt(t *testing.T) {
	num1 := 123

	type args struct {
		input *int
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
	}{
		{
			name: "filled input == result",
			args: args{
				input: &num1,
			},
			wantResult: num1,
		},
		{
			name: "nil input; result = 0",
			args: args{
				input: nil,
			},
			wantResult: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ForceInt(tt.args.input); gotResult != tt.wantResult {
				t.Errorf("ForceInt() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestForceIntPtr(t *testing.T) {
	num1 := 123
	zeroInt := 0

	type args struct {
		input *int
	}
	tests := []struct {
		name       string
		args       args
		wantResult *int
	}{
		{
			name: "filled input == result",
			args: args{
				input: &num1,
			},
			wantResult: &num1,
		},
		{
			name: "nil input; result = 0",
			args: args{
				input: nil,
			},
			wantResult: &zeroInt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := ForceIntPtr(tt.args.input)
			utils.AssertEqual(t, tt.wantResult != nil, gotResult != nil, "compare pointer")
			if tt.wantResult != nil && gotResult != nil {
				utils.AssertEqual(t, *tt.wantResult, *gotResult, "compare value")
			}
		})
	}
}

func TestForceInt64(t *testing.T) {
	num1 := int64(123)

	type args struct {
		input *int64
	}
	tests := []struct {
		name       string
		args       args
		wantResult int64
	}{
		{
			name: "filled input == result",
			args: args{
				input: &num1,
			},
			wantResult: num1,
		},
		{
			name: "nil input; result = 0",
			args: args{
				input: nil,
			},
			wantResult: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ForceInt64(tt.args.input); gotResult != tt.wantResult {
				t.Errorf("ForceInt64() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestForceInt64Ptr(t *testing.T) {
	num1 := int64(123)
	zeroInt64 := int64(0)

	type args struct {
		input *int64
	}
	tests := []struct {
		name       string
		args       args
		wantResult *int64
	}{
		{
			name: "filled input == result",
			args: args{
				input: &num1,
			},
			wantResult: &num1,
		},
		{
			name: "nil input; result = 0",
			args: args{
				input: nil,
			},
			wantResult: &zeroInt64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := ForceInt64Ptr(tt.args.input)
			utils.AssertEqual(t, tt.wantResult != nil, gotResult != nil, "compare pointer")
			if tt.wantResult != nil && gotResult != nil {
				utils.AssertEqual(t, *tt.wantResult, *gotResult, "compare value")
			}
		})
	}
}

func TestForceBool(t *testing.T) {
	cond1 := true
	cond2 := false

	type args struct {
		input *bool
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
	}{
		{
			name: "filled input == result",
			args: args{
				input: &cond1,
			},
			wantResult: cond1,
		},
		{
			name: "filled input == result",
			args: args{
				input: &cond2,
			},
			wantResult: cond2,
		},
		{
			name: "nil input; result = false",
			args: args{
				input: nil,
			},
			wantResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ForceBool(tt.args.input); gotResult != tt.wantResult {
				t.Errorf("ForceBool() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestForceBoolPtr(t *testing.T) {
	cond := true
	utils.AssertEqual(t, true, *ForceBoolPtr(&cond), "Test positive case")

	falseCond := false
	utils.AssertEqual(t, false, *ForceBoolPtr(&falseCond), "Test negative case")

	utils.AssertEqual(t, false, *ForceBoolPtr(nil), "Test nil case")
}

func TestForceFloat64(t *testing.T) {
	num1 := float64(123.00)

	type args struct {
		input *float64
	}
	tests := []struct {
		name       string
		args       args
		wantResult float64
	}{
		{
			name: "filled input == result",
			args: args{
				input: &num1,
			},
			wantResult: num1,
		},
		{
			name: "nil input; result = 0",
			args: args{
				input: nil,
			},
			wantResult: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ForceFloat64(tt.args.input); gotResult != tt.wantResult {
				t.Errorf("ForceFloat64() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestForceFloat64Ptr(t *testing.T) {
	num1 := float64(123.00)
	zeroFloat64 := float64(0)

	type args struct {
		input *float64
	}
	tests := []struct {
		name       string
		args       args
		wantResult *float64
	}{
		{
			name: "filled input == result",
			args: args{
				input: &num1,
			},
			wantResult: &num1,
		},
		{
			name: "nil input; result = 0",
			args: args{
				input: nil,
			},
			wantResult: &zeroFloat64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := ForceFloat64Ptr(tt.args.input)
			utils.AssertEqual(t, tt.wantResult != nil, gotResult != nil, "compare pointer")
			if tt.wantResult != nil && gotResult != nil {
				utils.AssertEqual(t, *tt.wantResult, *gotResult, "compare value")
			}
		})
	}
}

func TestForceStrfmtDateTime(t *testing.T) {
	datetime1 := strfmt.DateTime(time.Now())

	type args struct {
		input *strfmt.DateTime
	}
	tests := []struct {
		name       string
		args       args
		wantResult strfmt.DateTime
	}{
		{
			name: "filled input == result",
			args: args{
				input: &datetime1,
			},
			wantResult: datetime1,
		},
		{
			name: "nil input; result = strfmt.DateTime{}",
			args: args{
				input: nil,
			},
			wantResult: strfmt.DateTime{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ForceStrfmtDateTime(tt.args.input); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ForceStrfmtDateTime() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestForceStrfmtDateTimePtr(t *testing.T) {
	datetime1 := strfmt.DateTime(time.Now())
	datetimeZero := strfmt.DateTime{}

	type args struct {
		input *strfmt.DateTime
	}
	tests := []struct {
		name       string
		args       args
		wantResult *strfmt.DateTime
	}{
		{
			name: "filled input == result",
			args: args{
				input: &datetime1,
			},
			wantResult: &datetime1,
		},
		{
			name: "nil input; result = strfmt.DateTime{}",
			args: args{
				input: nil,
			},
			wantResult: &datetimeZero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := ForceStrfmtDateTimePtr(tt.args.input)
			utils.AssertEqual(t, tt.wantResult != nil, gotResult != nil, "compare pointer")
			if tt.wantResult != nil && gotResult != nil {
				utils.AssertEqual(t, *tt.wantResult, *gotResult, "compare value")
			}
		})
	}
}

func TestForceStrfmtDate(t *testing.T) {
	date1 := strfmt.Date(time.Now())

	type args struct {
		input *strfmt.Date
	}
	tests := []struct {
		name       string
		args       args
		wantResult strfmt.Date
	}{
		{
			name: "filled input == result",
			args: args{
				input: &date1,
			},
			wantResult: date1,
		},
		{
			name: "nil input; result = strfmt.Date{}",
			args: args{
				input: nil,
			},
			wantResult: strfmt.Date{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ForceStrfmtDate(tt.args.input); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ForceStrfmtDate() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestForceStrfmtDatePtr(t *testing.T) {
	date1 := strfmt.Date(time.Now())
	dateZero := strfmt.Date{}

	type args struct {
		input *strfmt.Date
	}
	tests := []struct {
		name       string
		args       args
		wantResult *strfmt.Date
	}{
		{
			name: "filled input == result",
			args: args{
				input: &date1,
			},
			wantResult: &date1,
		},
		{
			name: "nil input; result = strfmt.Date{}",
			args: args{
				input: nil,
			},
			wantResult: &dateZero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := ForceStrfmtDatePtr(tt.args.input)
			utils.AssertEqual(t, tt.wantResult != nil, gotResult != nil, "compare pointer")
			if tt.wantResult != nil && gotResult != nil {
				utils.AssertEqual(t, *tt.wantResult, *gotResult, "compare value")
			}
		})
	}
}

func TestForceTime(t *testing.T) {
	time1 := time.Now()

	type args struct {
		input *time.Time
	}
	tests := []struct {
		name       string
		args       args
		wantResult time.Time
	}{
		{
			name: "filled input == result",
			args: args{
				input: &time1,
			},
			wantResult: time1,
		},
		{
			name: "nil input; result = time.Time{}",
			args: args{
				input: nil,
			},
			wantResult: time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ForceTime(tt.args.input); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ForceTime() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestForceTimePtr(t *testing.T) {
	time1 := time.Now()
	timeZero := time.Time{}

	type args struct {
		input *time.Time
	}
	tests := []struct {
		name       string
		args       args
		wantResult *time.Time
	}{
		{
			name: "filled input == result",
			args: args{
				input: &time1,
			},
			wantResult: &time1,
		},
		{
			name: "nil input; result = time.Time{}",
			args: args{
				input: nil,
			},
			wantResult: &timeZero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := ForceTimePtr(tt.args.input)
			utils.AssertEqual(t, tt.wantResult != nil, gotResult != nil, "compare pointer")
			if tt.wantResult != nil && gotResult != nil {
				utils.AssertEqual(t, *tt.wantResult, *gotResult, "compare value")
			}
		})
	}
}

func TestForceUUID(t *testing.T) {
	uuid1 := uuid.New()

	type args struct {
		input *uuid.UUID
	}
	tests := []struct {
		name       string
		args       args
		wantResult uuid.UUID
	}{
		{
			name: "filled input == result",
			args: args{
				input: &uuid1,
			},
			wantResult: uuid1,
		},
		{
			name: "nil input; result = uuid.Nil",
			args: args{
				input: nil,
			},
			wantResult: uuid.Nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ForceUUID(tt.args.input); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ForceUUID() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestForceUUIDPtr(t *testing.T) {
	uuid1 := uuid.New()
	uuidNil := uuid.Nil

	type args struct {
		input *uuid.UUID
	}
	tests := []struct {
		name       string
		args       args
		wantResult *uuid.UUID
	}{
		{
			name: "filled input == result",
			args: args{
				input: &uuid1,
			},
			wantResult: &uuid1,
		},
		{
			name: "nil input; result = uuid.Nil",
			args: args{
				input: nil,
			},
			wantResult: &uuidNil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := ForceUUIDPtr(tt.args.input)
			utils.AssertEqual(t, tt.wantResult != nil, gotResult != nil, "compare pointer")
			if tt.wantResult != nil && gotResult != nil {
				utils.AssertEqual(t, *tt.wantResult, *gotResult, "compare value")
			}
		})
	}
}

func TestFloatToFormattedStr(t *testing.T) {
	mapTest := map[float64]string{
		10000:    "10,000",
		12500.12: "12,501",
		500:      "500",
		10.1:     "11",
		19191919: "19,191,919",
	}

	for inp, expected := range mapTest {
		utils.AssertEqual(t, expected, FloatToFormattedStr(inp), "Test FloatToFormattedStr")
	}
}

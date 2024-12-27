package lib

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2/utils"
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

func TestConvertSliceIntToStr(t *testing.T) {
	value := []int{1, 2, 3, 4}
	res := ConvertSliceIntToStr(value, ",")
	utils.AssertEqual(t, "1,2,3,4", res)
}

func TestConvertDatetimeToDate(t *testing.T) {
	datetimeNow := time.Now().UTC()

	year, month, day := datetimeNow.Date()

	dateNow, err := time.Parse("2006-1-2", fmt.Sprintf("%d-%d-%d", year, int(month), day))
	utils.AssertEqual(t, nil, err, "parsing date")

	type args struct {
		value time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				value: datetimeNow,
			},
			want:    dateNow,
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				value: time.Time{}.AddDate(-1000000, 0, 0),
			},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertDatetimeToDate(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertDatetimeToDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertDatetimeToDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveEmptySliceStrPtr(t *testing.T) {
	type args struct {
		values []*string
	}
	tests := []struct {
		name             string
		args             args
		wantResultLength int
	}{
		{
			name: "success, no result",
			args: args{
				values: []*string{
					nil,
					Strptr("   "),
					Strptr(""),
				},
			},
			wantResultLength: 0,
		},
		{
			name: "success, 2 result",
			args: args{
				values: []*string{
					nil,
					Strptr("   "),
					Strptr(""),
					Strptr(" abc "),
					Strptr("cde"),
				},
			},
			wantResultLength: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := RemoveEmptySliceStrPtr(tt.args.values); !reflect.DeepEqual(len(gotResult), tt.wantResultLength) {
				t.Errorf("RemoveEmptySliceStrPtr() length = %v, want length %v", len(gotResult), tt.wantResultLength)
			}
		})
	}
}

package lib

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ConvertToMD5 func
func ConvertToMD5(value *int) string {
	str := IntToStr(*value)
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

// ConvertStrToMD5 func
func ConvertStrToMD5(value *string) string {
	var str string = *value
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

// ConvertToSHA1 func
func ConvertToSHA1(value string) string {
	sha := sha1.New()
	sha.Write([]byte(value))
	encrypted := sha.Sum(nil)
	encryptedString := fmt.Sprintf("%x", encrypted)
	return encryptedString
}

// ConvertToSHA256 func
func ConvertToSHA256(value string) string {
	hash := sha256.Sum256([]byte(value))
	res := fmt.Sprintf("%x", hash)
	return res
}

// IntToStr func
func IntToStr(value int) string {
	return strconv.Itoa(value)
}

// StrToInt func
func StrToInt(value string) int {
	valueInt, _ := strconv.Atoi(value)
	return valueInt
}

// StrToInt64 func
func StrToInt64(value string) int64 {
	valueInt, _ := strconv.ParseInt(value, 10, 64)
	return valueInt
}

// StrToFloat func
func StrToFloat(value string) float64 {
	valueFloat, _ := strconv.ParseFloat(value, 32)
	return valueFloat
}

// StrToBool func
func StrToBool(value string) bool {
	valueBool, _ := strconv.ParseBool(value)
	return valueBool
}

// FloatToStr func
func FloatToStr(inputNum float64, prec ...int) string {
	precision := 0 // Default precision is 0 if not specified
	if len(prec) > 0 {
		precision = prec[0]
	}
	return strconv.FormatFloat(inputNum, 'f', precision, 64)
}

// ConvertJSONToStr func
func ConvertJSONToStr(payload interface{}) string {
	jsonData, _ := JSONMarshal(payload)
	return string(jsonData)
}

// ConvertStrToObj func
func ConvertStrToObj(value string) map[string]interface{} {
	var output map[string]interface{}
	JSONUnmarshal([]byte(value), &output)
	return output
}

// ConvertStrToJSON func
func ConvertStrToJSON(value string) interface{} {
	var output interface{}
	JSONUnmarshal([]byte(value), &output)
	return output
}

// ConvertStrToTime func
func ConvertStrToTime(value string) *time.Time {
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, value)
	return &t
}

// ConvertSliceIntToStr func
// Source: https://stackoverflow.com/a/37533144
func ConvertSliceIntToStr(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func ConvertDatetimeToDate(value time.Time) (time.Time, error) {
	year, month, day := value.Date()
	result, err := time.Parse("2006-1-2", fmt.Sprintf("%d-%d-%d", year, int(month), day))
	return result, err
}

func RemoveEmptySliceStrPtr(values []*string) (result []*string) {
	for _, val := range values {
		if IsEmptyStrPtr(val) {
			continue
		}

		result = append(result, val)
	}
	return
}

func FloatToFormattedStr(inputNum float64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.0f", Round(inputNum))
}

// ConvertJsonToStr func
func ConvertJsonToStr(payload interface{}) string {
	jsonData, _ := JSONMarshal(payload)
	return string(jsonData)
}

// ConvertStrToArrObj func
func ConvertStrToArrObj(value string) []map[string]interface{} {
	var output []map[string]interface{}
	JSONUnmarshal([]byte(value), &output)
	return output
}

// ConvertStrToJson func
func ConvertStrToJson(value string) interface{} {
	var output interface{}
	JSONUnmarshal([]byte(value), &output)
	return output
}

// ConvertSliceStrToStr func
// Source: https://stackoverflow.com/a/37533144
func ConvertSliceStrToStr(a []string, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

// ConvertSliceUUIDToSliceStr func
// Source: https://stackoverflow.com/a/37533144
func ConvertSliceUUIDToSliceStr(listUUID []uuid.UUID) (result []string) {
	for idx := range listUUID {
		item := listUUID[idx]
		itemStr := item.String()
		result = append(result, itemStr)
	}

	return
}

// LeadingZerosRemove func
func StrLeadingZerosRemove(str string) (result string) {
	if len(str) > 0 {
		result = strings.TrimLeft(str, "0")
	}
	return
}

/*
ForceStr

Return empty string if input == nil
*/
func ForceStr(input *string) (result string) {
	if input == nil {
		return
	}

	result = *input
	return
}

/*
ForceStrPtr

Return pointer empty string if input == nil
*/
func ForceStrPtr(input *string) (result *string) {
	if input == nil {
		result = Strptr("")
		return
	}

	result = input
	return
}

/*
ForceInt

Return 0 if input == nil
*/
func ForceInt(input *int) (result int) {
	if input == nil {
		return
	}

	result = *input
	return
}

/*
ForceIntPtr

Return pointer 0 if input == nil
*/
func ForceIntPtr(input *int) (result *int) {
	if input == nil {
		result = Intptr(0)
		return
	}

	result = input
	return
}

/*
ForceInt64

Return 0 if input == nil
*/
func ForceInt64(input *int64) (result int64) {
	if input == nil {
		return
	}

	result = *input
	return
}

/*
ForceInt64Ptr

Return pointer 0 if input == nil
*/
func ForceInt64Ptr(input *int64) (result *int64) {
	if input == nil {
		result = Int64ptr(0)
		return
	}

	result = input
	return
}

/*
ForceBool

Return false if input == nil
*/
func ForceBool(input *bool) (result bool) {
	if input == nil {
		return
	}

	result = *input
	return
}

// ForceBoolPtr will return either TRUE/FALSE pointer. Empty pointer will return false pointer.
func ForceBoolPtr(cond *bool) *bool {
	if cond == nil {
		return Boolptr(false)
	}
	return cond
}

/*
ForceFloat64

Return 0 if input == nil
*/
func ForceFloat64(input *float64) (result float64) {
	if input == nil {
		return
	}

	result = *input
	return
}

/*
ForceFloat64Ptr

Return pointer 0 if input == nil
*/
func ForceFloat64Ptr(input *float64) (result *float64) {
	if input == nil {
		result = Float64ptr(0)
		return
	}

	result = input
	return
}

/*
ForceStrfmtDateTime

Return zero time if input == nil
*/
func ForceStrfmtDateTime(input *strfmt.DateTime) (result strfmt.DateTime) {
	if input == nil {
		return
	}

	result = *input
	return
}

/*
ForceStrfmtDateTimePtr

Return pointer zero time if input == nil
*/
func ForceStrfmtDateTimePtr(input *strfmt.DateTime) (result *strfmt.DateTime) {
	if input == nil {
		result = new(strfmt.DateTime)
		return
	}

	result = input
	return
}

/*
ForceStrfmtDate

Return zero time if input == nil
*/
func ForceStrfmtDate(input *strfmt.Date) (result strfmt.Date) {
	if input == nil {
		return
	}

	result = *input
	return
}

/*
ForceStrfmtDatePtr

Return pointer zero time if input == nil
*/
func ForceStrfmtDatePtr(input *strfmt.Date) (result *strfmt.Date) {
	if input == nil {
		result = new(strfmt.Date)
		return
	}

	result = input
	return
}

/*
ForceTime

Return zero time if input == nil
*/
func ForceTime(input *time.Time) (result time.Time) {
	if input == nil {
		return
	}

	result = *input
	return
}

/*
ForceTimePtr

Return pointer zero time if input == nil
*/
func ForceTimePtr(input *time.Time) (result *time.Time) {
	if input == nil {
		zeroTime := time.Time{}
		result = &zeroTime
		return
	}

	result = input
	return
}

/*
ForceUUID

Return uuid.Nil if input == nil
*/
func ForceUUID(input *uuid.UUID) (result uuid.UUID) {
	if input == nil {
		return
	}

	result = *input
	return
}

/*
ForceUUIDPtr

Return pointer uuid.Nil if input == nil
*/
func ForceUUIDPtr(input *uuid.UUID) (result *uuid.UUID) {
	if input == nil {
		uuidNil := uuid.Nil
		result = &uuidNil
		return
	}

	result = input
	return
}

func ConvertStrToTimeWFormat(value, layout string) (result time.Time) {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}

	t, err := time.Parse(layout, value)
	if err != nil {
		return
	}

	result = t
	return
}
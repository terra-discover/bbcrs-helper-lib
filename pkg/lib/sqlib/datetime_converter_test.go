package sqlib

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofiber/fiber/v2/utils"
)

func TestConvertPriceGuaranteeTimeLimitStrPtrToTime(t *testing.T) {
	// valid value with Z
	// result == want
	value := "2021-02-05T13:05:10Z"
	want := time.Date(2021, 02, 05, 13, 05, 10, 0, &time.Location{})
	result := ConvertPriceGuaranteeTimeLimitStrPtrToTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == want")

	// valid value without Z
	// result == want
	value = "2021-02-05T13:05:10"
	want = time.Date(2021, 02, 05, 13, 05, 10, 0, &time.Location{})
	result = ConvertPriceGuaranteeTimeLimitStrPtrToTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == want")

	// invalid value
	// result == zero time
	value = "abcd"
	want = time.Time{}
	result = ConvertPriceGuaranteeTimeLimitStrPtrToTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == zero time")

	// nil value
	// result == zero time
	want = time.Time{}
	result = ConvertPriceGuaranteeTimeLimitStrPtrToTime(nil)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == zero time")
}

func TestConvertPriceGuaranteeTimeLimitStrPtrToStrfmtDateTime(t *testing.T) {
	// valid value with Z
	// result == want
	value := "2021-02-05T13:05:10Z"
	want := strfmt.DateTime(time.Date(2021, 02, 05, 13, 05, 10, 0, &time.Location{}))
	result := ConvertPriceGuaranteeTimeLimitStrPtrToStrfmtDateTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == want")

	// valid value without Z
	// result == want
	value = "2021-02-05T13:05:10"
	want = strfmt.DateTime(time.Date(2021, 02, 05, 13, 05, 10, 0, &time.Location{}))
	result = ConvertPriceGuaranteeTimeLimitStrPtrToStrfmtDateTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == want")

	// invalid value
	// result == zero time
	value = "abcd"
	want = strfmt.DateTime(time.Time{})
	result = ConvertPriceGuaranteeTimeLimitStrPtrToStrfmtDateTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == zero time")

	// nil value
	// result == zero time
	want = strfmt.DateTime(time.Time{})
	result = ConvertPriceGuaranteeTimeLimitStrPtrToStrfmtDateTime(nil)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == zero time")
}

func TestConvertPaymentTimeLimitStrPtrToTime(t *testing.T) {
	// valid value with .000Z
	// result == want
	value := "2021-02-05T13:05:10.000Z"
	want := time.Date(2021, 02, 05, 13, 05, 10, 0, &time.Location{})
	result := ConvertPaymentTimeLimitStrPtrToTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == want")

	// valid value without .000Z
	// result == want
	value = "2021-02-05T13:05:10"
	want = time.Time{}
	result = ConvertPaymentTimeLimitStrPtrToTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == want")

	// invalid value
	// result == zero time
	value = "abcd"
	want = time.Time{}
	result = ConvertPaymentTimeLimitStrPtrToTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == zero time")

	// nil value
	// result == zero time
	want = time.Time{}
	result = ConvertPaymentTimeLimitStrPtrToTime(nil)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == zero time")
}

func TestConvertPaymentTimeLimitStrPtrToStrfmtDateTime(t *testing.T) {
	// valid value with .000Z
	// result == want
	value := "2021-02-05T13:05:10.000Z"
	want := strfmt.DateTime(time.Date(2021, 02, 05, 13, 05, 10, 0, &time.Location{}))
	result := ConvertPaymentTimeLimitStrPtrToStrfmtDateTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == want")

	// valid value without .000Z
	// result == want
	value = "2021-02-05T13:05:10"
	want = strfmt.DateTime(time.Time{})
	result = ConvertPaymentTimeLimitStrPtrToStrfmtDateTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == want")

	// invalid value
	// result == zero time
	value = "abcd"
	want = strfmt.DateTime(time.Time{})
	result = ConvertPaymentTimeLimitStrPtrToStrfmtDateTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == zero time")

	// nil value
	// result == zero time
	want = strfmt.DateTime(time.Time{})
	result = ConvertPaymentTimeLimitStrPtrToStrfmtDateTime(nil)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == zero time")
}

func TestConvertOfferExpirationTimeLimitStrPtrToTime(t *testing.T) {
	// valid value with Z
	// result == want
	value := "2021-02-05T13:05:10Z"
	want := time.Date(2021, 02, 05, 13, 05, 10, 0, &time.Location{})
	result := ConvertOfferExpirationTimeLimitStrPtrToTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == want")

	// valid value without Z
	// result == want
	value = "2021-02-05T13:05:10"
	want = time.Date(2021, 02, 05, 13, 05, 10, 0, &time.Location{})
	result = ConvertOfferExpirationTimeLimitStrPtrToTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == want")

	// invalid value
	// result == zero time
	value = "abcd"
	want = time.Time{}
	result = ConvertOfferExpirationTimeLimitStrPtrToTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == zero time")

	// nil value
	// result == zero time
	want = time.Time{}
	result = ConvertOfferExpirationTimeLimitStrPtrToTime(nil)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == zero time")
}

func TestConvertOfferExpirationTimeLimitStrPtrToStrfmtDateTime(t *testing.T) {
	// valid value with Z
	// result == want
	value := "2021-02-05T13:05:10Z"
	want := strfmt.DateTime(time.Date(2021, 02, 05, 13, 05, 10, 0, &time.Location{}))
	result := ConvertOfferExpirationTimeLimitStrPtrToStrfmtDateTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == want")

	// valid value without Z
	// result == want
	value = "2021-02-05T13:05:10"
	want = strfmt.DateTime(time.Date(2021, 02, 05, 13, 05, 10, 0, &time.Location{}))
	result = ConvertOfferExpirationTimeLimitStrPtrToStrfmtDateTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == want")

	// invalid value
	// result == zero time
	value = "abcd"
	want = strfmt.DateTime(time.Time{})
	result = ConvertOfferExpirationTimeLimitStrPtrToStrfmtDateTime(&value)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == zero time")

	// nil value
	// result == zero time
	want = strfmt.DateTime(time.Time{})
	result = ConvertOfferExpirationTimeLimitStrPtrToStrfmtDateTime(nil)
	utils.AssertEqual(t, true, result.String() == want.String(), "result == zero time")
}

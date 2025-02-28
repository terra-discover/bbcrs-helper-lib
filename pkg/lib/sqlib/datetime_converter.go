package sqlib

import (
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/terra-discover/bbcrs-helper-lib/pkg/lib"
)

// ConvertPriceGuaranteeTimeLimitStrPtrToTime
func ConvertPriceGuaranteeTimeLimitStrPtrToTime(priceGuaranteeTimeLimitStr *string) (result time.Time) {
	const (
		formatGuaranteeWithZ    = "2006-01-02T15:04:05Z"
		formatGuaranteeWithoutZ = "2006-01-02T15:04:05"
	)

	if lib.IsEmptyStrPtr(priceGuaranteeTimeLimitStr) {
		return
	}

	formatDateUsingZ := strings.Contains(*priceGuaranteeTimeLimitStr, "Z")
	if formatDateUsingZ {
		result = lib.ConvertStrToTimeWFormat(*priceGuaranteeTimeLimitStr, formatGuaranteeWithZ)
		return
	}

	result = lib.ConvertStrToTimeWFormat(*priceGuaranteeTimeLimitStr, formatGuaranteeWithoutZ)
	return
}

// ConvertPriceGuaranteeTimeLimitStrPtrToStrfmtDateTime
func ConvertPriceGuaranteeTimeLimitStrPtrToStrfmtDateTime(priceGuaranteeTimeLimitStr *string) (result strfmt.DateTime) {
	theTime := ConvertPriceGuaranteeTimeLimitStrPtrToTime(priceGuaranteeTimeLimitStr)
	result = strfmt.DateTime(theTime)
	return
}

// ConvertPaymentTimeLimitStrPtrToTime
func ConvertPaymentTimeLimitStrPtrToTime(paymentTimeLimitStr *string) (result time.Time) {
	const (
		formatTimeLimit = "2006-01-02T15:04:05.000Z"
	)

	if lib.IsEmptyStrPtr(paymentTimeLimitStr) {
		return
	}

	result = lib.ConvertStrToTimeWFormat(*paymentTimeLimitStr, formatTimeLimit)
	return
}

// ConvertPaymentTimeLimitStrPtrToStrfmtDateTime
func ConvertPaymentTimeLimitStrPtrToStrfmtDateTime(paymentTimeLimitStr *string) (result strfmt.DateTime) {
	theTime := ConvertPaymentTimeLimitStrPtrToTime(paymentTimeLimitStr)
	result = strfmt.DateTime(theTime)
	return
}

// ConvertOfferExpirationTimeLimitStrPtrToTime
func ConvertOfferExpirationTimeLimitStrPtrToTime(offerExpirationTimeLimitStr *string) (result time.Time) {
	const (
		formatGuaranteeWithZ    = "2006-01-02T15:04:05Z"
		formatGuaranteeWithoutZ = "2006-01-02T15:04:05"
	)

	if lib.IsEmptyStrPtr(offerExpirationTimeLimitStr) {
		return
	}

	formatDateUsingZ := strings.Contains(*offerExpirationTimeLimitStr, "Z")
	if formatDateUsingZ {
		result = lib.ConvertStrToTimeWFormat(*offerExpirationTimeLimitStr, formatGuaranteeWithZ)
		return
	}

	result = lib.ConvertStrToTimeWFormat(*offerExpirationTimeLimitStr, formatGuaranteeWithoutZ)
	return
}

// ConvertOfferExpirationTimeLimitStrPtrToStrfmtDateTime
func ConvertOfferExpirationTimeLimitStrPtrToStrfmtDateTime(offerExpirationTimeLimitStr *string) (result strfmt.DateTime) {
	theTime := ConvertOfferExpirationTimeLimitStrPtrToTime(offerExpirationTimeLimitStr)
	result = strfmt.DateTime(theTime)
	return
}

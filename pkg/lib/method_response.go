package lib

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Method error response
// =====================
type errorCode int

const (
	_                     errorCode = 0
	errorCodeBadRequest   errorCode = 400
	errorCodeUnauthorized errorCode = 401
	errorCodeNotFound     errorCode = 404
	errorCodeTimeout      errorCode = 408
	errorCodeConflict     errorCode = 409
	errorCodeGone         errorCode = 410
	errorCodeInternal     errorCode = 500
)

const (
	_                        string = ""
	errorMessageBadRequest   string = "Bad request"
	errorMessageUnauthorized string = "Unauthorized"
	errorMessageNotFound     string = "Not found"
	errorMessageTimeout      string = "Timeout"
	errorMessageConflict     string = "Conflict"
	errorMessageGone         string = "Gone"
	errorMessageInternal     string = "Internal"
)

func SetErrorBadRequest(description ...string) (errResp ErrorResponse) {
	errResp.setCode(errorCodeBadRequest)
	errResp.setDescription(description...)
	return
}

func SetErrorUnauthorized(description ...string) (errResp ErrorResponse) {
	errResp.setCode(errorCodeUnauthorized)
	errResp.setDescription(description...)
	return
}

func SetErrorNotFound(description ...string) (errResp ErrorResponse) {
	errResp.setCode(errorCodeNotFound)
	errResp.setDescription(description...)
	return
}

func SetErrorTimeout(description ...string) (errResp ErrorResponse) {
	errResp.setCode(errorCodeTimeout)
	errResp.setDescription(description...)
	return
}

func SetErrorConflict(description ...string) (errResp ErrorResponse) {
	errResp.setCode(errorCodeConflict)
	errResp.setDescription(description...)
	return
}

func SetErrorGone(description ...string) (errResp ErrorResponse) {
	errResp.setCode(errorCodeGone)
	errResp.setDescription(description...)
	return
}

func SetErrorInternal(description ...string) (errResp ErrorResponse) {
	errResp.setCode(errorCodeInternal)
	errResp.setDescription(description...)
	return
}

type ErrorResponse struct {
	code        errorCode
	description string
}

type IErrorResponse interface {
	Code() (code int)
	Description() (description string)
	IsEmpty() (isEmpty bool)
	SendToContext(c *fiber.Ctx) (err error)
	setCode(code errorCode)
	setDescription(description ...string)
}

func (e ErrorResponse) Code() (code int) {
	if e.code == 0 {
		log.Printf("ErrorResponse_Code, code: %d is undefined", code)
	}
	code = int(e.code)
	return
}

func (e ErrorResponse) Description() (description string) {
	description = e.description
	return
}

func (e ErrorResponse) IsEmpty() (isEmpty bool) {
	if (e == ErrorResponse{}) {
		isEmpty = true
	}
	return
}

func (e ErrorResponse) SendToContext(c *fiber.Ctx) (err error) {
	switch e.code {
	case errorCodeBadRequest:
		{
			if len(e.description) == 0 {
				e.description = errorMessageBadRequest
			}
			err = ErrorBadRequest(c, e.description)
			break
		}
	case errorCodeUnauthorized:
		{
			if len(e.description) == 0 {
				e.description = errorMessageUnauthorized
			}
			err = ErrorUnauthorized(c, e.description)
			break
		}
	case errorCodeNotFound:
		{
			if len(e.description) == 0 {
				e.description = errorMessageNotFound
			}
			err = ErrorNotFound(c, e.description)
			break
		}
	case errorCodeTimeout:
		{
			if len(e.description) == 0 {
				e.description = errorMessageTimeout
			}
			err = ErrorTimeout(c, e.description)
			break
		}
	case errorCodeConflict:
		{
			if len(e.description) == 0 {
				e.description = errorMessageConflict
			}
			err = ErrorConflict(c, e.description)
			break
		}
	case errorCodeGone:
		{
			if len(e.description) == 0 {
				e.description = errorMessageGone
			}
			err = ErrorGone(c, e.description)
			break
		}
	case errorCodeInternal:
		{
			if len(e.description) == 0 {
				e.description = errorMessageInternal
			}
			err = ErrorInternal(c, e.description)
			break
		}
	default:
		{
			log.Printf("ErrorResponse_SendToContext, code: %d is undefined", e.code)
			err = ErrorInternal(c)
			break
		}
	}

	return
}

func (e *ErrorResponse) setCode(code errorCode) {
	e.code = code
}

func (e *ErrorResponse) setDescription(description ...string) {
	if len(description) > 0 {
		e.description = strings.TrimSpace(description[0])
	}
}

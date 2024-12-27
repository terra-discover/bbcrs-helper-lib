package lib

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func TestSetErrorBadRequest(t *testing.T) {
	type args struct {
		description []string
	}
	tests := []struct {
		name        string
		args        args
		wantErrResp ErrorResponse
	}{
		{
			name: "empty description",
			args: args{
				description: []string{},
			},
			wantErrResp: ErrorResponse{
				code: errorCodeBadRequest,
			},
		},
		{
			name: "filled description",
			args: args{
				description: []string{"Error 1234"},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeBadRequest,
				description: "Error 1234",
			},
		},
		{
			name: "filled description with space",
			args: args{
				description: []string{"  Error 1234  "},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeBadRequest,
				description: "Error 1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrResp := SetErrorBadRequest(tt.args.description...); !reflect.DeepEqual(gotErrResp, tt.wantErrResp) {
				t.Errorf("SetErrorBadRequest() = %v, want %v", gotErrResp, tt.wantErrResp)
			}
		})
	}
}

func TestSetErrorUnauthorized(t *testing.T) {
	type args struct {
		description []string
	}
	tests := []struct {
		name        string
		args        args
		wantErrResp ErrorResponse
	}{
		{
			name: "empty description",
			args: args{
				description: []string{},
			},
			wantErrResp: ErrorResponse{
				code: errorCodeUnauthorized,
			},
		},
		{
			name: "filled description",
			args: args{
				description: []string{"Error 1234"},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeUnauthorized,
				description: "Error 1234",
			},
		},
		{
			name: "filled description with space",
			args: args{
				description: []string{"  Error 1234  "},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeUnauthorized,
				description: "Error 1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrResp := SetErrorUnauthorized(tt.args.description...); !reflect.DeepEqual(gotErrResp, tt.wantErrResp) {
				t.Errorf("SetErrorUnauthorized() = %v, want %v", gotErrResp, tt.wantErrResp)
			}
		})
	}
}

func TestSetErrorNotFound(t *testing.T) {
	type args struct {
		description []string
	}
	tests := []struct {
		name        string
		args        args
		wantErrResp ErrorResponse
	}{
		{
			name: "empty description",
			args: args{
				description: []string{},
			},
			wantErrResp: ErrorResponse{
				code: errorCodeNotFound,
			},
		},
		{
			name: "filled description",
			args: args{
				description: []string{"Error 1234"},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeNotFound,
				description: "Error 1234",
			},
		},
		{
			name: "filled description with space",
			args: args{
				description: []string{"  Error 1234  "},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeNotFound,
				description: "Error 1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrResp := SetErrorNotFound(tt.args.description...); !reflect.DeepEqual(gotErrResp, tt.wantErrResp) {
				t.Errorf("SetErrorNotFound() = %v, want %v", gotErrResp, tt.wantErrResp)
			}
		})
	}
}

func TestSetErrorTimeout(t *testing.T) {
	type args struct {
		description []string
	}
	tests := []struct {
		name        string
		args        args
		wantErrResp ErrorResponse
	}{
		{
			name: "empty description",
			args: args{
				description: []string{},
			},
			wantErrResp: ErrorResponse{
				code: errorCodeTimeout,
			},
		},
		{
			name: "filled description",
			args: args{
				description: []string{"Error 1234"},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeTimeout,
				description: "Error 1234",
			},
		},
		{
			name: "filled description with space",
			args: args{
				description: []string{"  Error 1234  "},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeTimeout,
				description: "Error 1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrResp := SetErrorTimeout(tt.args.description...); !reflect.DeepEqual(gotErrResp, tt.wantErrResp) {
				t.Errorf("SetErrorTimeout() = %v, want %v", gotErrResp, tt.wantErrResp)
			}
		})
	}
}

func TestSetErrorConflict(t *testing.T) {
	type args struct {
		description []string
	}
	tests := []struct {
		name        string
		args        args
		wantErrResp ErrorResponse
	}{
		{
			name: "empty description",
			args: args{
				description: []string{},
			},
			wantErrResp: ErrorResponse{
				code: errorCodeConflict,
			},
		},
		{
			name: "filled description",
			args: args{
				description: []string{"Error 1234"},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeConflict,
				description: "Error 1234",
			},
		},
		{
			name: "filled description with space",
			args: args{
				description: []string{"  Error 1234  "},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeConflict,
				description: "Error 1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrResp := SetErrorConflict(tt.args.description...); !reflect.DeepEqual(gotErrResp, tt.wantErrResp) {
				t.Errorf("SetErrorConflict() = %v, want %v", gotErrResp, tt.wantErrResp)
			}
		})
	}
}

func TestSetErrorGone(t *testing.T) {
	type args struct {
		description []string
	}
	tests := []struct {
		name        string
		args        args
		wantErrResp ErrorResponse
	}{
		{
			name: "empty description",
			args: args{
				description: []string{},
			},
			wantErrResp: ErrorResponse{
				code: errorCodeGone,
			},
		},
		{
			name: "filled description",
			args: args{
				description: []string{"Error 1234"},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeGone,
				description: "Error 1234",
			},
		},
		{
			name: "filled description with space",
			args: args{
				description: []string{"  Error 1234  "},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeGone,
				description: "Error 1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrResp := SetErrorGone(tt.args.description...); !reflect.DeepEqual(gotErrResp, tt.wantErrResp) {
				t.Errorf("SetErrorGone() = %v, want %v", gotErrResp, tt.wantErrResp)
			}
		})
	}
}

func TestSetErrorInternal(t *testing.T) {
	type args struct {
		description []string
	}
	tests := []struct {
		name        string
		args        args
		wantErrResp ErrorResponse
	}{
		{
			name: "empty description",
			args: args{
				description: []string{},
			},
			wantErrResp: ErrorResponse{
				code: errorCodeInternal,
			},
		},
		{
			name: "filled description",
			args: args{
				description: []string{"Error 1234"},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeInternal,
				description: "Error 1234",
			},
		},
		{
			name: "filled description with space",
			args: args{
				description: []string{"  Error 1234  "},
			},
			wantErrResp: ErrorResponse{
				code:        errorCodeInternal,
				description: "Error 1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrResp := SetErrorInternal(tt.args.description...); !reflect.DeepEqual(gotErrResp, tt.wantErrResp) {
				t.Errorf("SetErrorInternal() = %v, want %v", gotErrResp, tt.wantErrResp)
			}
		})
	}
}

func TestErrorResponse_Code(t *testing.T) {
	type fields struct {
		code        errorCode
		description string
	}
	tests := []struct {
		name     string
		fields   fields
		wantCode int
	}{
		{
			name: "success",
			fields: fields{
				code: errorCodeBadRequest,
			},
			wantCode: int(errorCodeBadRequest),
		},
		{
			name: "undefined with log",
			fields: fields{
				code: 0,
			},
			wantCode: int(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrorResponse{
				code:        tt.fields.code,
				description: tt.fields.description,
			}
			if gotCode := e.Code(); gotCode != tt.wantCode {
				t.Errorf("ErrorResponse.Code() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func TestErrorResponse_Description(t *testing.T) {
	type fields struct {
		code        errorCode
		description string
	}
	tests := []struct {
		name            string
		fields          fields
		wantDescription string
	}{
		{
			name: "success",
			fields: fields{
				description: "error testing",
			},
			wantDescription: "error testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrorResponse{
				code:        tt.fields.code,
				description: tt.fields.description,
			}
			if gotDescription := e.Description(); gotDescription != tt.wantDescription {
				t.Errorf("ErrorResponse.Description() = %v, want %v", gotDescription, tt.wantDescription)
			}
		})
	}
}

func TestErrorResponse_IsEmpty(t *testing.T) {
	type fields struct {
		code        errorCode
		description string
	}
	tests := []struct {
		name        string
		fields      fields
		wantIsEmpty bool
	}{
		{
			name: "Error Response is not empty",
			fields: fields{
				code:        errorCodeBadRequest,
				description: "error testing",
			},
			wantIsEmpty: false,
		},
		{
			name:        "Error Response is empty (case 1)",
			fields:      fields{},
			wantIsEmpty: true,
		},
		{
			name: "Error Response is empty (case 2)",
			fields: fields{
				code:        0,
				description: "",
			},
			wantIsEmpty: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ErrorResponse{
				code:        tt.fields.code,
				description: tt.fields.description,
			}
			if gotIsEmpty := e.IsEmpty(); gotIsEmpty != tt.wantIsEmpty {
				t.Errorf("ErrorResponse.IsEmpty() = %v, want %v", gotIsEmpty, tt.wantIsEmpty)
			}
		})
	}
}

func TestErrorResponse_SendToContext(t *testing.T) {
	type fields struct {
		code        errorCode
		description string
	}
	tests := []struct {
		name        string
		fields      fields
		wantErrResp fields
	}{
		{
			name: errorMessageBadRequest,
			fields: fields{
				code:        errorCodeBadRequest,
				description: "",
			},
			wantErrResp: fields{
				code:        errorCodeBadRequest,
				description: errorMessageBadRequest,
			},
		},
		{
			name: errorMessageUnauthorized,
			fields: fields{
				code:        errorCodeUnauthorized,
				description: "",
			},
			wantErrResp: fields{
				code:        errorCodeUnauthorized,
				description: errorMessageUnauthorized,
			},
		},
		{
			name: errorMessageNotFound,
			fields: fields{
				code:        errorCodeNotFound,
				description: "",
			},
			wantErrResp: fields{
				code:        errorCodeNotFound,
				description: errorMessageNotFound,
			},
		},
		{
			name: errorMessageTimeout,
			fields: fields{
				code:        errorCodeTimeout,
				description: "",
			},
			wantErrResp: fields{
				code:        errorCodeTimeout,
				description: errorMessageTimeout,
			},
		},
		{
			name: errorMessageConflict,
			fields: fields{
				code:        errorCodeConflict,
				description: "",
			},
			wantErrResp: fields{
				code:        errorCodeConflict,
				description: errorMessageConflict,
			},
		},
		{
			name: errorMessageGone,
			fields: fields{
				code:        errorCodeGone,
				description: "",
			},
			wantErrResp: fields{
				code:        errorCodeGone,
				description: errorMessageGone,
			},
		},
		{
			name: errorMessageInternal,
			fields: fields{
				code:        errorCodeInternal,
				description: "",
			},
			wantErrResp: fields{
				code:        errorCodeInternal,
				description: errorMessageInternal,
			},
		},
		{
			name: "undefined",
			fields: fields{
				code:        0,
				description: "",
			},
			wantErrResp: fields{
				code:        errorCodeInternal,
				description: errorMessageInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Get("/default", func(c *fiber.Ctx) error {
				e := ErrorResponse{
					code:        tt.fields.code,
					description: tt.fields.description,
				}

				return e.SendToContext(c)
			})

			response, _ := app.Test(httptest.NewRequest("GET", "/default", nil))
			utils.AssertEqual(t, int(tt.wantErrResp.code), response.StatusCode, "validate response")
		})
	}
}

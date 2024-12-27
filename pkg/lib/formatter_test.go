package lib

import (
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func TestFormatEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				email: "    ABcdef123  @gmail .com    ",
			},
			want: "abcdef123@gmail.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatEmail(tt.args.email); got != tt.want {
				t.Errorf("FormatEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatEmailPtr(t *testing.T) {
	type args struct {
		email *string
	}
	tests := []struct {
		name       string
		args       args
		wantResult *string
	}{
		{
			name: "success",
			args: args{
				email: Strptr("    ABcdef123  @gmail .com    "),
			},
			wantResult: Strptr("abcdef123@gmail.com"),
		},
		{
			name: "success, but nil pointer",
			args: args{
				email: nil,
			},
			wantResult: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch gotResult := FormatEmailPtr(tt.args.email); {
			case gotResult == nil || tt.wantResult == nil:
				{
					utils.AssertEqual(t, gotResult == nil, tt.wantResult == nil, "validate pointer")
					break
				}
			case *gotResult != *tt.wantResult:
				{
					utils.AssertEqual(t, *gotResult, *tt.wantResult, "validate value")
					break
				}
			}
		})
	}
}

func TestFormatStr(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "success",
			args: args{
				s: " abcD  ",
			},
			wantResult: "abcD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := FormatStr(tt.args.s); gotResult != tt.wantResult {
				t.Errorf("FormatStr() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestFormatStrPtr(t *testing.T) {
	type args struct {
		s *string
	}
	tests := []struct {
		name       string
		args       args
		wantResult *string
	}{
		{
			name: "success, with not nil input",
			args: args{
				s: Strptr(" abcD  "),
			},
			wantResult: Strptr("abcD"),
		},
		{
			name: "success, with nil input",
			args: args{
				s: nil,
			},
			wantResult: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := FormatStrPtr(tt.args.s)
			utils.AssertEqual(t, tt.wantResult == nil, gotResult == nil, "validate pointer")
			if tt.wantResult != nil {
				utils.AssertEqual(t, *tt.wantResult, *gotResult, "validate value")
			}
		})
	}
}

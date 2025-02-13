package lib

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// VALIDATOR validate request body
var VALIDATOR *validator.Validate = validator.New()

const (
	specialcharacter            = `^[A-z0-9 ]+$`
	alphaNumUniCodeSpacePattern = `^[\p{L}\p{N} ]+$`
	customPhoneNumber           = `^\+[1-9]\d{1,31}$` // using validator e164 with custom length 32 characters (default e164 is 15 characters)
)

func init() {
	VALIDATOR.RegisterValidation("specialcharacter", customValidator(specialcharacter))
	VALIDATOR.RegisterValidation("alphanumunicodespace", customValidator(alphaNumUniCodeSpacePattern))
	VALIDATOR.RegisterValidation("customphonenumber", customValidator(customPhoneNumber))
}

func customValidator(pattern string) validator.Func {
	return func(f validator.FieldLevel) bool {
		str := f.Field().String()
		re := regexp.MustCompile(pattern)
		return re.MatchString(str)
	}
}

// HTTPClient http client interface
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// GetXUserID provide user id by http headers
func GetXUserID(c *fiber.Ctx) *uuid.UUID {
	id := string(c.Request().Header.Peek("x-user-id"))
	if id != "" {
		if current, err := uuid.Parse(id); nil == err {
			return &current
		}
	}

	return nil
}

// GetXAgentID provide user id by http headers
func GetXAgentID(c *fiber.Ctx) *uuid.UUID {
	id := string(c.Request().Header.Peek("x-agent-id"))
	if id != "" {
		if current, err := uuid.Parse(id); nil == err {
			return &current
		}
	} else if id = viper.GetString("AGENT_ID"); id != "" {
		if current, err := uuid.Parse(id); nil == err {
			return &current
		}
	}

	return nil
}

// GetXCorporateID provide corporate id by http headers
func GetXCorporateID(c *fiber.Ctx) *uuid.UUID {
	id := string(c.Request().Header.Peek("x-corporate-id"))
	if id != "" {
		if current, err := uuid.Parse(id); nil == err {
			return &current
		}
	}

	return nil
}

// GetLanguage get language by http header Accept-Language
func GetLanguage(c *fiber.Ctx) string {
	lang := viper.GetString("LANGUAGE")
	acceptLanguage := string(c.Request().Header.Peek("accept-language"))
	if acceptLanguage != "" && len(acceptLanguage) >= 2 {
		lang = acceptLanguage[0:2]
		// TODO: check to database if database exists, if not return to fallback language ...
	}

	lang = strings.ToLower(lang)
	if ok, _ := regexp.Match("^[a-z]", []byte(lang)); !ok || len(lang) < 2 {
		lang = "en"
	}

	return lang
}

// BodyParser with validation
func BodyParser(c *fiber.Ctx, payload interface{}) error {
	if err := c.BodyParser(payload); nil != err {
		return err
	}

	return VALIDATOR.Struct(payload)
}

// QueryParser with validation
func QueryParser(c *fiber.Ctx, payload interface{}) error {
	if err := c.QueryParser(payload); nil != err {
		return err
	}

	return VALIDATOR.Struct(payload)
}
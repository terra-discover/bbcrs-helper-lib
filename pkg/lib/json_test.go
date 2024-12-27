package lib

import (
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func TestSetJSONEngine(t *testing.T) {
	SetJSONEngine()
	engine := reflect.TypeOf(JSONMarshal)
	utils.AssertEqual(t, "utils.JSONMarshal", engine.String())

	SetJSONEngine("sonic")
	engine = reflect.TypeOf(JSONMarshal)
	utils.AssertEqual(t, "utils.JSONMarshal", engine.String())

	SetJSONEngine("go-json")
	engine = reflect.TypeOf(JSONMarshal)
	utils.AssertEqual(t, "utils.JSONMarshal", engine.String())
}

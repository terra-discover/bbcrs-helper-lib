package lib

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func TestLogStruct(t *testing.T) {
	var prefix = "Prefix"
	x := map[string]string{
		"test": "123",
	}
	byteX, _ := json.Marshal(x)

	loggedString := LogStruct(x, "Prefix")
	expectedCase := fmt.Sprintf("%s%s", prefix+" : ", byteX)
	utils.AssertEqual(t, expectedCase, loggedString, "TEST LOG STRUCT WITH PREFIX")

	loggedStringNoPrefix := LogStruct(x)
	expectedCaseNoPrefix := fmt.Sprintf("%s", byteX)
	utils.AssertEqual(t, expectedCaseNoPrefix, loggedStringNoPrefix, "TEST LOG STRUCT WITH NO PREFIX")
}

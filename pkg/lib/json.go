package lib

import (
	"encoding/json"

	"github.com/bytedance/sonic"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/utils"
)

// JSONMarshal default Encoder for JSON
var JSONMarshal utils.JSONMarshal = json.Marshal

// JSONUnmarshal default Decoder for JSON
var JSONUnmarshal utils.JSONUnmarshal = json.Unmarshal

// JSONEncoder get default json encoder
func JSONEncoder(engine ...string) utils.JSONMarshal {
	defaultEngine := "encoding/json"
	if len(engine) == 1 && engine[0] != "" {
		defaultEngine = engine[0]
	}

	switch defaultEngine {
	case "sonic":
		return sonic.Marshal
	case "go-json":
		return gojson.Marshal
	}

	return json.Marshal
}

// JSONDecoder get default json decoder
func JSONDecoder(engine ...string) utils.JSONUnmarshal {
	defaultEngine := "encoding/json"
	if len(engine) == 1 && engine[0] != "" {
		defaultEngine = engine[0]
	}

	switch defaultEngine {
	case "sonic":
		return sonic.Unmarshal
	case "go-json":
		return gojson.Unmarshal
	}

	return json.Unmarshal
}

// SetJSONEngine set default JSON encoder / decoder
func SetJSONEngine(engine ...string) {
	JSONMarshal = JSONEncoder(engine...)
	JSONUnmarshal = JSONDecoder(engine...)
}

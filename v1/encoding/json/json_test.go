package json

import (
	"testing"
)

type test struct {
	ABC string `json:"adc"`
	No  int
}

var jsonStr = "{\"adc\":\"123\",\"No\":1234}"
var jsonStrReverse = "{\"No\":1234,\"adc\":\"123\"}"

// TestEncode Json编码
func TestEncode(t *testing.T) {
	result := Encode(map[string]interface{}{"adc": "123", "No": 1234})
	if result != jsonStrReverse {
		panic("encoding json.Encode map error")
	}

	result = Encode(test{ABC: "123", No: 1234})
	if result != jsonStr {
		panic("encoding json.Encode struct error")
	}
}

// TestDecode Json解码
func TestDecode(t *testing.T) {
	var data interface{}
	err := Decode(jsonStr, &data)
	if err != nil {
		panic("encoding json.Decode map " + err.Error())
	}

	var testOne test
	err = Decode(jsonStr, &testOne)
	if err != nil {
		panic("encoding json.Decode struct " + err.Error())
	}
}

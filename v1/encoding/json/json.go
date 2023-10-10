package json

import (
	"encoding/json"
)

// Encode Json编码
func Encode(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}

// Decode Json解码
func Decode(str string, v interface{}) error {
	return json.Unmarshal([]byte(str), &v)
}

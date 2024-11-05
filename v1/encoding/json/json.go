package json

import (
	"encoding/json"
)

// Encode Json编码
func Encode(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}

// EncodeToByte Json编码
func EncodeToByte(v interface{}) []byte {
	bytes, _ := json.Marshal(v)
	return bytes
}

// Decode Json解码
func Decode(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), &v)
}

// DecodeFromByte Json解码
func DecodeFromByte(data []byte, v interface{}) error {
	return json.Unmarshal(data, &v)
}

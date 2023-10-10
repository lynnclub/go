package algorithm

import (
	"encoding/base64"
)

// Base64Encode base64编码
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64解码
func Base64Decode(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	return string(data), err
}

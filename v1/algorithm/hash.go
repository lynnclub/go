package algorithm

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash/crc32"
)

// MD5 返回十六进制字符串
func MD5(str string) string {
	hash := md5.Sum([]byte(str))

	return hex.EncodeToString(hash[:])
}

// SHA1 返回十六进制字符串
func SHA1(str string) string {
	hash := sha1.Sum([]byte(str))

	return hex.EncodeToString(hash[:])
}

// Crc32 crc32
func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

// HmacSHA256 返回十六进制字符串
func HmacSHA256(message string, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(message))

	return hex.EncodeToString(hash.Sum(nil))
}

// HmacSHA1 返回十六进制字符串
func HmacSHA1(s string, secret string) string {
	hash := hmac.New(sha1.New, []byte(secret))
	hash.Write([]byte(s))

	return hex.EncodeToString(hash.Sum(nil))
}

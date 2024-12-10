package algorithm

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
)

func MD5v2(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func MD5v3(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))

	return hex.EncodeToString(hash.Sum(nil))
}

func TestMD5(t *testing.T) {
	result := MD5("123")
	if result != "202cb962ac59075b964b07152d234b70" {
		panic("md5 incorrect")
	}

	resultV2 := MD5v2("123")
	if result != resultV2 {
		panic("md5v2 inconsistent")
	}

	resultV3 := MD5v3("123")
	if result != resultV3 {
		panic("md5v3 inconsistent")
	}
}

func BenchmarkMD5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5("123")
	}
}

func BenchmarkMD5V2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5v2("123")
	}
}

func BenchmarkMD5V3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5v3("123")
	}
}

func TestCRC32Hash(t *testing.T) {
	expectedHashValue := uint32(3964322768)
	hashValue := Crc32("Hello, World!")

	if hashValue != expectedHashValue {
		t.Errorf("Expected CRC32 hash: %d, got: %d", expectedHashValue, hashValue)
	}
}

package aes

import (
	"strings"
	"testing"
)

func TestEncryptDecryptECB(t *testing.T) {
	plaintext := "Hello, World!"
	key := "0123456789abcdef" // 16-byte key for AES-128

	ciphertext, err := EncryptECB(plaintext, key)
	if err != nil {
		t.Errorf("Error during encryption: %v", err)
	}

	decryptedText, err := DecryptECB(ciphertext, key)
	if err != nil {
		t.Errorf("Error during decryption: %v", err)
	}

	if decryptedText != plaintext {
		t.Errorf("Decrypted text doesn't match original plaintext. Expected: %s, Got: %s", plaintext, decryptedText)
	}
}

// TestEncryptECBWithDifferentKeyLengths 测试不同密钥长度
func TestEncryptECBWithDifferentKeyLengths(t *testing.T) {
	plaintext := "Hello, AES Encryption!"

	tests := []struct {
		name      string
		key       string
		keyLength int
	}{
		{"AES-128", "0123456789abcdef", 16},                 // 16 bytes
		{"AES-192", "0123456789abcdef01234567", 24},         // 24 bytes
		{"AES-256", "0123456789abcdef0123456789abcdef", 32}, // 32 bytes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := EncryptECB(plaintext, tt.key)
			if err != nil {
				t.Fatalf("Encryption failed for %s: %v", tt.name, err)
			}

			decrypted, err := DecryptECB(ciphertext, tt.key)
			if err != nil {
				t.Fatalf("Decryption failed for %s: %v", tt.name, err)
			}

			if decrypted != plaintext {
				t.Errorf("Decrypted text doesn't match for %s. Expected: %s, Got: %s", tt.name, plaintext, decrypted)
			}
		})
	}
}

// TestEncryptECBInvalidKey 测试无效密钥长度
func TestEncryptECBInvalidKey(t *testing.T) {
	plaintext := "Hello, World!"

	tests := []struct {
		name string
		key  string
	}{
		{"Too short", "short"},
		{"15 bytes", "012345678901234"},
		{"17 bytes", "01234567890123456"},
		{"23 bytes", "01234567890123456789012"},
		{"25 bytes", "0123456789012345678901234"},
		{"Empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EncryptECB(plaintext, tt.key)
			if err == nil {
				t.Errorf("Expected error for key length %d, but got none", len(tt.key))
			}
			if !strings.Contains(err.Error(), "key validation failed") && !strings.Contains(err.Error(), "invalid key length") {
				t.Errorf("Expected key validation error, got: %v", err)
			}
		})
	}
}

// TestEncryptECBEmptyPlaintext 测试空明文
func TestEncryptECBEmptyPlaintext(t *testing.T) {
	key := "0123456789abcdef"
	_, err := EncryptECB("", key)
	if err == nil {
		t.Error("Expected error for empty plaintext, but got none")
	}
	if err != ErrEmptyPlaintext {
		t.Errorf("Expected ErrEmptyPlaintext, got: %v", err)
	}
}

// TestDecryptECBInvalidKey 测试解密时使用无效密钥
func TestDecryptECBInvalidKey(t *testing.T) {
	ciphertext := "dGVzdA==" // 任意base64字符串

	tests := []struct {
		name string
		key  string
	}{
		{"Too short", "short"},
		{"15 bytes", "012345678901234"},
		{"Empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecryptECB(ciphertext, tt.key)
			if err == nil {
				t.Errorf("Expected error for key length %d, but got none", len(tt.key))
			}
		})
	}
}

// TestDecryptECBInvalidBase64 测试解密无效的base64字符串
func TestDecryptECBInvalidBase64(t *testing.T) {
	key := "0123456789abcdef"
	invalidBase64 := "This is not valid base64!!!"

	_, err := DecryptECB(invalidBase64, key)
	if err == nil {
		t.Error("Expected error for invalid base64, but got none")
	}
	if !strings.Contains(err.Error(), "base64 decode failed") {
		t.Errorf("Expected base64 decode error, got: %v", err)
	}
}

// TestDecryptECBEmptyCiphertext 测试解密空密文
func TestDecryptECBEmptyCiphertext(t *testing.T) {
	key := "0123456789abcdef"
	emptyCiphertext := "" // 空的base64字符串会解码为空字节数组

	_, err := DecryptECB(emptyCiphertext, key)
	if err == nil {
		t.Error("Expected error for empty ciphertext, but got none")
	}
}

// TestDecryptECBInvalidBlockSize 测试密文长度不是块大小的倍数
func TestDecryptECBInvalidBlockSize(t *testing.T) {
	key := "0123456789abcdef"
	// 创建一个长度不是16倍数的密文（5字节）
	invalidCiphertext := "AAAAA" // base64解码后是4字节（不是16的倍数）

	_, err := DecryptECB(invalidCiphertext, key)
	if err == nil {
		t.Error("Expected error for invalid block size, but got none")
	}
}

// TestPad 测试填充函数
func TestPad(t *testing.T) {
	tests := []struct {
		name         string
		data         []byte
		blockSize    int
		expectedLen  int
		expectedLast byte
	}{
		{"Empty data", []byte{}, 16, 16, 16},
		{"1 byte", []byte{1}, 16, 16, 15},
		{"15 bytes", make([]byte, 15), 16, 16, 1},
		{"16 bytes", make([]byte, 16), 16, 32, 16},
		{"17 bytes", make([]byte, 17), 16, 32, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			padded := Pad(tt.data, tt.blockSize)
			if len(padded) != tt.expectedLen {
				t.Errorf("Expected padded length %d, got %d", tt.expectedLen, len(padded))
			}
			if padded[len(padded)-1] != tt.expectedLast {
				t.Errorf("Expected last byte %d, got %d", tt.expectedLast, padded[len(padded)-1])
			}
			// 验证填充字节都是相同的
			paddingSize := int(padded[len(padded)-1])
			for i := len(padded) - paddingSize; i < len(padded); i++ {
				if padded[i] != byte(paddingSize) {
					t.Errorf("Padding byte at position %d is %d, expected %d", i, padded[i], paddingSize)
				}
			}
		})
	}
}

// TestUnpad 测试去填充函数
func TestUnpad(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
		expectedLen int
	}{
		{"Valid padding 1", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 1}, false, 15},
		{"Valid padding 2", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 2, 2}, false, 14},
		{"Valid padding 16", []byte{16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16}, false, 0},
		{"Empty data", []byte{}, true, 0},
		{"Invalid padding size", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 17}, true, 0},
		{"Invalid padding bytes", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 2}, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Unpad(tt.data)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) != tt.expectedLen {
					t.Errorf("Expected result length %d, got %d", tt.expectedLen, len(result))
				}
			}
		})
	}
}

// TestPadUnpadRoundTrip 测试填充和去填充的往返
func TestPadUnpadRoundTrip(t *testing.T) {
	blockSize := 16
	testData := [][]byte{
		[]byte("Hello"),
		[]byte("Hello, World!"),
		[]byte("This is exactly 16 bytes data!!!"),
		[]byte("A"),
		make([]byte, 15),
		make([]byte, 16),
		make([]byte, 17),
		make([]byte, 32),
		make([]byte, 100),
	}

	for i, data := range testData {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			padded := Pad(data, blockSize)
			unpadded, err := Unpad(padded)
			if err != nil {
				t.Fatalf("Unpad failed: %v", err)
			}
			if string(unpadded) != string(data) {
				t.Errorf("Round trip failed. Original: %q, After pad/unpad: %q", data, unpadded)
			}
		})
	}
}

// TestEncryptDecryptLongText 测试长文本
func TestEncryptDecryptLongText(t *testing.T) {
	longText := strings.Repeat("This is a long text for testing AES ECB encryption. ", 100)
	key := "0123456789abcdef"

	ciphertext, err := EncryptECB(longText, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := DecryptECB(ciphertext, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decrypted != longText {
		t.Error("Decrypted long text doesn't match original")
	}
}

// TestEncryptDecryptSpecialCharacters 测试特殊字符
func TestEncryptDecryptSpecialCharacters(t *testing.T) {
	tests := []string{
		"Hello, 世界!",
		"Special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
		"Newlines\nand\ttabs",
		"Emoji: 😀🎉🚀",
		"\x00\x01\x02\x03", // 控制字符
	}

	key := "0123456789abcdef"

	for _, plaintext := range tests {
		t.Run(plaintext, func(t *testing.T) {
			ciphertext, err := EncryptECB(plaintext, key)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			decrypted, err := DecryptECB(ciphertext, key)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			if decrypted != plaintext {
				t.Errorf("Mismatch. Expected: %q, Got: %q", plaintext, decrypted)
			}
		})
	}
}

// TestDecryptWithWrongKey 测试使用错误密钥解密
func TestDecryptWithWrongKey(t *testing.T) {
	plaintext := "Secret Message"
	correctKey := "0123456789abcdef"
	wrongKey := "fedcba9876543210"

	ciphertext, err := EncryptECB(plaintext, correctKey)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := DecryptECB(ciphertext, wrongKey)
	// 解密可能成功但结果是乱码，或者在unpad时失败
	if err == nil && decrypted == plaintext {
		t.Error("Should not decrypt correctly with wrong key")
	}
}

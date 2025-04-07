package aes

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"errors"
	"fmt"
)

var (
	ErrInvalidKeyLength  = errors.New("invalid key length (must be 16, 24 or 32 bytes)")
	ErrEmptyPlaintext    = errors.New("plaintext cannot be empty")
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
)

func validateKey(key []byte) error {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return ErrInvalidKeyLength
	}
	return nil
}

func EncryptECB(plainText string, key string) (string, error) {
	if len(plainText) == 0 {
		return "", ErrEmptyPlaintext
	}

	keyBytes := []byte(key)
	if err := validateKey(keyBytes); err != nil {
		return "", fmt.Errorf("key validation failed: %w", err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("cipher initialization failed: %w", err)
	}

	plainTextBytes := Pad([]byte(plainText), block.BlockSize())
	cipherText := make([]byte, len(plainTextBytes))

	for i := 0; i < len(plainTextBytes); i += block.BlockSize() {
		block.Encrypt(cipherText[i:], plainTextBytes[i:i+block.BlockSize()])
	}

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptECB(cipherTextString string, key string) (string, error) {
	keyBytes := []byte(key)
	if err := validateKey(keyBytes); err != nil {
		return "", fmt.Errorf("key validation failed: %w", err)
	}

	cipherText, err := base64.StdEncoding.DecodeString(cipherTextString)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	if len(cipherText) == 0 {
		return "", ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("cipher initialization failed: %w", err)
	}

	bs := block.BlockSize()
	if len(cipherText)%bs != 0 {
		return "", fmt.Errorf("ciphertext length %d not multiple of block size %d", len(cipherText), bs)
	}

	plainText := make([]byte, len(cipherText))
	for i := 0; i < len(cipherText); i += bs {
		block.Decrypt(plainText[i:], cipherText[i:i+bs])
	}

	unpadded, err := Unpad(plainText)
	if err != nil {
		return "", fmt.Errorf("unpad failed: %w", err)
	}

	return string(unpadded), nil
}

func Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	if padding == 0 {
		padding = blockSize
	}
	return append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data for unpadding")
	}

	padding := int(data[len(data)-1])
	if padding <= 0 || padding > len(data) {
		return nil, fmt.Errorf("invalid padding size: %d", padding)
	}

	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			return nil, errors.New("invalid padding bytes")
		}
	}
	return data[:len(data)-padding], nil
}

package aes

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"errors"
	"fmt"
)

func EncryptECB(plainText string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plainTextBytes := []byte(plainText)

	// 补全明文
	plainTextBytes = Pad(plainTextBytes, block.BlockSize())

	cipherText := make([]byte, len(plainTextBytes))

	// 使用 ECB 模式进行加密
	for i := 0; i < len(plainTextBytes); i += block.BlockSize() {
		block.Encrypt(cipherText[i:i+block.BlockSize()], plainTextBytes[i:i+block.BlockSize()])
	}

	// 将 []byte 转换为 base64 编码的字符串
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptECB(cipherTextString string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 将 base64 编码的字符串解码为 []byte
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextString)
	if err != nil {
		return "", err
	}

	if len(cipherText) == 0 {
		return "", fmt.Errorf("invalid ciphertext: %w", err)
	}

	if len(cipherText)%block.BlockSize() != 0 {
		return "", errors.New("cipherText is not a multiple of the block size")
	}

	// 使用 ECB 模式进行解密
	for i := 0; i < len(cipherText); i += block.BlockSize() {
		block.Decrypt(cipherText[i:i+block.BlockSize()], cipherText[i:i+block.BlockSize()])
	}

	// 去除填充
	cipherText = Unpad(cipherText)

	return string(cipherText), nil
}

// Pad 使用 PKCS7 填充方式对数据进行补全
func Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// Unpad 去除填充
func Unpad(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return nil
	}

	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

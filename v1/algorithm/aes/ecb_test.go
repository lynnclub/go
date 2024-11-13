package aes

import (
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

package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func EncryptPassword(plainText, keyString string) (string, error) {
	key := []byte(keyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("UNABLE TO CREATE CIPHER BLOCK : %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("UNABLE TO CREATE GCM CIPHER : %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("UNABLE TO GENERATE NONCE : %w", err)
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(cipherText), nil
}

func DecryptPassword(encryptedHex, keyString string) (string, error) {
	key := []byte(keyString)
	encryptedBytes, _ := hex.DecodeString(encryptedHex)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("UNABLE TO CREATE CIPHER BLOCK : %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("UNABLE TO CREATE GCM CIPHER : %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	nonce, cipherText := encryptedBytes[:nonceSize], encryptedBytes[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("UNABLE TO DECRYPT : %w", err)
	}

	return string(plainText), nil
}

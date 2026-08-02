package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func EncryptPassword(plainText, keyString string) (string, error) {
	key := []byte(keyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(cipherText), nil
}

func DecryptPassword(encryptedHex, keyString string) (string, error) {
	key := []byte(keyString)
	encryptedBytes, _ := hex.DecodeString(encryptedHex)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, cipherText := encryptedBytes[:nonceSize], encryptedBytes[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func Trial(str string) {
	key := os.Getenv("AES_KEY")
	encrypted, err := EncryptPassword(str, key)
	if err != nil {
		fmt.Println("Error encrypting:", err)
	}
	fmt.Println("Encrypted:", encrypted)

	decrypted, err := DecryptPassword(encrypted, key)
	if err != nil {
		fmt.Println("Error decrypting:", err)
	}
	fmt.Println("Decrypted:", decrypted)
}

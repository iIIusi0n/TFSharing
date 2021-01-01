package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func DecryptID(encryptedID string) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedID)
	if err != nil {
		return "", err
	}
	key := "0123456789012345"
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(encryptedBytes) % aes.BlockSize != 0 {
		return "", err
	}
	iv := encryptedBytes[:aes.BlockSize]
	encryptedBytes = encryptedBytes[aes.BlockSize:]

	plaintext := make([]byte, len(encryptedBytes))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, encryptedBytes)
	return string(plaintext[:56]), nil
}

func EncryptID(id string) (string, error) {
	idBytes := []byte(id)
	key := "0123456789012345"
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if mod := len(idBytes) % aes.BlockSize; mod != 0 {
		padding := make([]byte, aes.BlockSize-mod)
		idBytes = append(idBytes, padding...)
	}

	ciphertext := make([]byte, aes.BlockSize+len(idBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], idBytes)

	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encoded, nil
}

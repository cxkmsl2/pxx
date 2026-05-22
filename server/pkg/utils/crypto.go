package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

var aesKey = []byte("pxx-2026-sub-key-32bytes!!")

func AESEncrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil { return "", err }
	aesGCM, err := cipher.NewGCM(block)
	if err != nil { return "", err }
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil { return "", err }
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func AESDecrypt(cipherB64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cipherB64)
	if err != nil { return "", err }
	block, err := aes.NewCipher(aesKey)
	if err != nil { return "", err }
	aesGCM, err := cipher.NewGCM(block)
	if err != nil { return "", err }
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil { return "", err }
	return string(plaintext), nil
}

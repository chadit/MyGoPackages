package gocrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// EncryptAES - encrypt a key with AES key
func EncryptAES(key, text string) (string, error) {
	var (
		block cipher.Block
		err   error
	)

	if block, err = aes.NewCipher([]byte(key)); err != nil {
		return "", err
	}
	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES - encrypt a key with AES key
func DecryptAES(key, text string) (string, error) {
	var (
		block      cipher.Block
		err        error
		ciphertext []byte
	)
	if block, err = aes.NewCipher([]byte(key)); err != nil {
		return "", err
	}
	if ciphertext, err = base64.StdEncoding.DecodeString(text); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext), nil
}

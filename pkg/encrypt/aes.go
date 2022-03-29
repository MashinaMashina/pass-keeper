package encrypt

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"github.com/Djarvur/go-aescrypt"
	"strings"
)

func EncryptAES(key []byte, plaintext string) (string, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	bytes, err := aescrypt.EncryptAESCBCPadded([]byte(plaintext), key, iv)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(iv) + "." + hex.EncodeToString(bytes), nil
}

func DecryptAES(key []byte, ct string) (string, error) {
	parts := strings.SplitN(ct, ".", 2)

	iv, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	str, err := aescrypt.DecryptAESCBCPadded(ciphertext, key, iv)
	if err != nil {
		return "", err
	}

	return string(str), nil
}

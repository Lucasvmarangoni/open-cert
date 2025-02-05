package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"

	"errors"

	errs "github.com/Lucasvmarangoni/logella/err"
)

func Encrypt(plaintext string, key string) (string, error) {
	keyDecoded, err := hex.DecodeString(key)
	if err != nil {
		return "", errs.Wrap(err, "hex.DecodeString", http.StatusInternalServerError)
	}

	block, err := aes.NewCipher(keyDecoded)
	if err != nil {
		return "", errs.Wrap(err, "aes.NewCipher", http.StatusInternalServerError)
	}

	ciphertext := make([]byte, aes.BlockSize+len([]byte(plaintext)))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errs.Wrap(err, "io.ReadFull", http.StatusInternalServerError)
	}
	mode := cipher.NewCFBEncrypter(block, iv)
	mode.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext string, key string) (string, error) {

	keyDecoded, err := hex.DecodeString(key)
	if err != nil {
		return "", errs.Wrap(err, "hex.DecodeString", http.StatusInternalServerError)
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", errs.Wrap(err, "base64.StdEncoding.DecodeString", http.StatusInternalServerError)
	}

	block, err := aes.NewCipher(keyDecoded)
	if err != nil {
		return "", errs.Wrap(err, "aes.NewCipher", http.StatusInternalServerError)
	}

	if len(data) < aes.BlockSize {
		err = errors.New("invalid ciphertext block size")
		return "", errs.Wrap(err, "len(data) < aes.BlockSize", http.StatusInternalServerError)
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	mode := cipher.NewCFBDecrypter(block, iv)
	mode.XORKeyStream(data, data)

	return string(data), nil
}

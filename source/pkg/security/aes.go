package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"

	"errors"

	"github.com/Lucasvmarangoni/logella/err"
)

func Encrypt(plaintext string, key string) (string, error) {
	keyDecoded, err := hex.DecodeString(key)
	if err != nil {
		return "", errs.Ctx(err, "hex.DecodeString")
	}

	block, err := aes.NewCipher(keyDecoded)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len([]byte(plaintext)))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errs.Ctx(err, "io.ReadFull")
	}

	mode := cipher.NewCFBEncrypter(block, iv)
	mode.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}


func Decrypt(ciphertext string, key string) (string, error) {

	keyDecoded, err := hex.DecodeString(key)
	if err != nil {
		return "", errs.Ctx(err, "hex.DecodeString")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", errs.Ctx(err, "base64.StdEncoding.DecodeString")
	}

	block, err := aes.NewCipher(keyDecoded)
	if err != nil {
		return "", errs.Ctx(err, "aes.NewCipher")
	}

	if len(data) < aes.BlockSize {
		err = errors.New("invalid ciphertext block size")
		return "", errs.Ctx(err, "len(data) < aes.BlockSize")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	mode := cipher.NewCFBDecrypter(block, iv)
	mode.XORKeyStream(data, data)

	return string(data), nil
}


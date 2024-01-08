package aeslib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type Aes struct {
	key256 []byte // 32 bytes = 256 bits key
}

// New creates a new AES cipher with a 256 bits key
func New(key256 []byte) *Aes {
	return &Aes{
		key256: key256,
	}
}

func (c *Aes) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key256)
	if err != nil {
		return []byte{}, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}

	return aesGCM.Seal(nonce, nonce, data, nil), nil
}

func (c *Aes) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key256)
	if err != nil {
		return []byte{}, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	result, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

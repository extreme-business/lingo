package aeslib

import (
	"testing"
)

func TestNewAes(t *testing.T) {

}

func TestAes_Encrypt_Decrypt(t *testing.T) {
	t.Run("should encrypt and decrypt", func(t *testing.T) {
		c := Aes{
			key256: []byte("12345678901234567890123456789012"),
		}

		r, err := c.Encrypt([]byte("hello world"))
		if err != nil {
			t.Error(err)
		}

		r2, err := c.Decrypt(r)
		if err != nil {
			t.Error(err)
		}

		if string(r2) != "hello world" {
			t.Error("invalid decrypt")
		}
	})

	t.Run("should fail to decrypt with invalid key", func(t *testing.T) {
		c := Aes{
			key256: []byte("12345678901234567890123456789012"),
		}

		r, err := c.Encrypt([]byte("hello world"))
		if err != nil {
			t.Error(err)
		}

		c.key256 = []byte("12345678901234567890123456789013")

		_, err = c.Decrypt(r)
		if err == nil {
			t.Error("expected error")
		}
	})
}

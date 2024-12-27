package lib

import (
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func TestCipherEncryptDecrypt(t *testing.T) {
	plaintext := "password"
	key := "CIPHER_SECRETKEY_MUST_HAVE_32BIT"

	_, err := CipherEncrypt(plaintext, key[:28])
	// Invalid Key just have 28 byte in Encrypt
	utils.AssertEqual(t, "crypto/aes: invalid key size 28", err.Error())

	cipherEncrypt, err := CipherEncrypt(plaintext, key)
	utils.AssertEqual(t, nil, err)

	cipherDecrypt, err := CipherDecrypt(cipherEncrypt, key)
	utils.AssertEqual(t, nil, err)

	// Success Decrypt
	utils.AssertEqual(t, plaintext, string(cipherDecrypt))

	_, err = CipherDecrypt(cipherEncrypt, key[:28])
	// Invalid Key just have 28 byte in Decrypt
	utils.AssertEqual(t, "crypto/aes: invalid key size 28", err.Error())

	_, err = CipherDecrypt([]byte(string(cipherEncrypt)[:5]), key)
	// Len byte is different
	utils.AssertEqual(t, "illegal base64 data at input byte 5", err.Error())
}

func TestGeneratePassword(t *testing.T) {
	password := GeneratePassword(20, 6, 6, 6)

	utils.AssertEqual(t, 20, len(password))
}

func TestRandomChars(t *testing.T) {
	utils.AssertEqual(t, 10, len(RandomChars(10)))
}

func TestRandomString(t *testing.T) {
	utils.AssertEqual(t, 10, len(RandomString(10, "")))
}

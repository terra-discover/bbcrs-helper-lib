package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	mathRand "math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

// On using math rand, you must call mathRand.New(mathRand.NewSource(time.Now().UnixNano())) or call globally mathRand.Seed(time.Now().UnixNano())
// It needs for avoid returning same result
var seededRand *mathRand.Rand = mathRand.New(
	mathRand.NewSource(time.Now().UTC().UnixNano()))

// GeneratePassword for generate random password
func GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
	var (
		lowerCharSet   = "abcdefghijklmnopqrstuvwxyz"
		upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		specialCharSet = "!@#$%&*"
		numberSet      = "0123456789"
		allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
	)
	var password strings.Builder

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random, _ := rand.Int(rand.Reader, big.NewInt(int64(len(specialCharSet))))
		password.WriteString(string(specialCharSet[random.BitLen()]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random, _ := rand.Int(rand.Reader, big.NewInt(int64(len(numberSet))))
		password.WriteString(string(numberSet[random.BitLen()]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random, _ := rand.Int(rand.Reader, big.NewInt(int64(len(upperCharSet))))
		password.WriteString(string(upperCharSet[random.BitLen()]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random, _ := rand.Int(rand.Reader, big.NewInt(int64(len(allCharSet))))
		password.WriteString(string(allCharSet[random.BitLen()]))
	}
	inRune := []rune(password.String())
	seededRand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

// CipherEncrypt for encrypt data with AES algorithm
func CipherEncrypt(plaintext, key string) ([]byte, error) {
	c, err := aes.NewCipher([]byte(key))
	if err == nil {
		gcm, err := cipher.NewGCM(c)
		if err == nil {
			nonce := make([]byte, gcm.NonceSize())
			if _, err = io.ReadFull(rand.Reader, nonce); err == nil {
				return gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
			}
		}
	}

	return nil, err
}

// CipherDecrypt for decrypt data with AES algorithm
func CipherDecrypt(ciphertext []byte, key string) ([]byte, error) {
	c, err := aes.NewCipher([]byte(key))
	if err == nil {
		gcm, err := cipher.NewGCM(c)
		if err == nil {
			nonceSize := gcm.NonceSize()
			if len(ciphertext) < nonceSize {
				return nil, errors.New("ciphertext too short")
			}
			nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
			return gcm.Open(nil, nonce, ciphertext, nil)
		}
	}
	return nil, err
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomChars func
func RandomChars(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[seededRand.Intn(len(letters))]
	}
	return string(b)
}

// RandomString func
func RandomString(length int, charset string) string {
	if charset == "" {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomCode func
func RandomCode(length int) string {
	var codes = []rune("1234567890")
	b := make([]rune, length)
	for i := range b {
		b[i] = codes[mathRand.Intn(len(codes))]
	}
	return string(b)
}

// ToUUID Convert everything To UUID v5
func ToUUID(s interface{}) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("%v", s)))
}

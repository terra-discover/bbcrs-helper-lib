package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"math/big"
	mathRand "math/rand"
	"strings"
	"time"
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
func CipherEncrypt(text, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	padding := aes.BlockSize - len(text)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	msg := append([]byte(text), padtext...)

	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
	finalMsg := strings.Replace(base64.URLEncoding.EncodeToString(ciphertext), "=", "", -1)
	return []byte(finalMsg), nil
}

// CipherDecrypt for decrypt data with AES algorithm
func CipherDecrypt(ciphertext []byte, key string) ([]byte, error) {
	text := string(ciphertext)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	if m := len(text) % 4; m != 0 {
		text += strings.Repeat("=", 4-m)
	}

	decodedMsg, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return nil, err
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return nil, errors.New("blocksize must be multipe of decoded message length")
	}

	if len(decodedMsg) < aes.BlockSize {
		return nil, errors.New("invalid block size")
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	length := len(msg)
	unpadding := int(msg[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	unpadMsg := msg[:(length - unpadding)]
	if err != nil {
		return nil, err
	}

	return unpadMsg, nil
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

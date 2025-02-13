package lib

import (
	"bytes"
	"compress/zlib"
	"strings"
)

// Compress will translate raw string to compressed string
func Compress(s string) string {
	buff := bytes.Buffer{}
	compressEnc := zlib.NewWriter(&buff)
	compressEnc.Write([]byte(s))
	compressEnc.Close()
	return string(buff.Bytes())
}

// Decompress will translate compressed string to raw string
func Decompress(s string) (string, error) {
	decompressor, err := zlib.NewReader(strings.NewReader(s))
	if err != nil {
		return "", err
	}
	var decompressedBuff bytes.Buffer
	decompressedBuff.ReadFrom(decompressor)
	return decompressedBuff.String(), nil
}

// CompressBytes will translate raw []byte to compressed []byte
func CompressBytes(s []byte) []byte {
	buff := bytes.Buffer{}
	compressEnc := zlib.NewWriter(&buff)
	compressEnc.Write(s)
	compressEnc.Close()

	return buff.Bytes()
}

// DecompressBytes will translate compressed []byte to raw []byte
func DecompressBytes(s string) ([]byte, error) {
	decompressor, err := zlib.NewReader(strings.NewReader(s))
	if err != nil {
		return nil, err
	}

	var decompressedBuff bytes.Buffer
	decompressedBuff.ReadFrom(decompressor)
	return decompressedBuff.Bytes(), nil
}

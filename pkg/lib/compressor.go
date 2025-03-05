package lib

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
)

type GzipHeader struct {
	Name    string    // auto uuid.New() if empty on compressing
	ModTime time.Time // auto time.Now().UTC() if empty on compressing
	Comment string
}

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

func CompressGzipString(strVal string, header GzipHeader) (result string, err error) {
	// Configure params
	name := uuid.New().String()
	modTime := time.Now().UTC()
	comment := ""

	if !IsEmptyStr(header.Name) {
		name = header.Name
	}

	if !IsZeroTime(header.ModTime) {
		modTime = header.ModTime
	}

	if !IsEmptyStr(header.Comment) {
		comment = header.Comment
	}

	// Start compress
	var buf bytes.Buffer

	zw := gzip.NewWriter(&buf)
	zw.Name = name
	zw.ModTime = modTime
	zw.Comment = comment

	_, errWrite := zw.Write([]byte(strVal))
	if errWrite != nil {
		err = errWrite
		return
	}

	errClose := zw.Close()
	if errClose != nil {
		err = errClose
		return
	}

	// Set result
	result = buf.String()
	return
}

func DecompressGzipString(strVal string) (resultString string, header GzipHeader, err error) {
	// Start Decompress
	bufRead := new(bytes.Buffer)
	bufRead.WriteString(strVal)

	zr, errRead := gzip.NewReader(bufRead)
	if errRead != nil {
		err = errRead
		return
	}

	// fmt.Printf("Name: %s\nComment: %s\nModTime: %s\n\n", zr.Name, zr.Comment, zr.ModTime.UTC())

	strBuild := new(strings.Builder)

	_, errCopy := io.Copy(strBuild, zr)
	if errCopy != nil {
		err = errCopy
		return
	}

	errClose := zr.Close()
	if errClose != nil {
		err = errClose
		return
	}

	// Set result
	resultString = strBuild.String()

	header = GzipHeader{
		Name:    zr.Name,
		ModTime: zr.ModTime,
		Comment: zr.Comment,
	}

	return
}

package lib

import (
	"log"
	"strings"
)

func FormatEmail(email string) (result string) {
	// Note:
	// strings.TrimSpace is only remove all leading or trailing whitespace defined by utf like " ", "\t", or "\n"
	// so, we need to strings.ReplaceAll to remove whitespace between characters
	strLowercase := strings.ToLower(email)
	strTrim := strings.TrimSpace(strLowercase)
	result = strings.ReplaceAll(strTrim, " ", "")
	return
}

func FormatEmailPtr(email *string) (result *string) {
	// Note:
	// strings.TrimSpace is only remove all leading or trailing whitespace defined by utf like " ", "\t", or "\n"
	// so, we need to strings.ReplaceAll to remove whitespace between characters
	if email == nil {
		log.Println("WARNING FormatEmailPtr: formatting nil pointer email")
		return
	}
	strLowercase := strings.ToLower(*email)
	strTrim := strings.TrimSpace(strLowercase)
	result = Strptr(strings.ReplaceAll(strTrim, " ", ""))
	return
}

func FormatStr(s string) (result string) {
	result = strings.TrimSpace(s)
	return
}

func FormatStrPtr(s *string) (result *string) {
	if s == nil {
		return
	}
	result = Strptr(strings.TrimSpace(*s))
	return
}

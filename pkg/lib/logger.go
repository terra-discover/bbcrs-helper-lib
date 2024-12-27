package lib

import (
	"fmt"
	"log"
)

// LogStruct func
func LogStruct(data interface{}, message ...string) string {
	byteData, _ := JSONMarshal(data)
	var prefix string
	if len(message) > 0 {
		prefix = message[0] + " : "
	}

	printedLog := fmt.Sprintf("%s%s", prefix, byteData)
	log.Println(printedLog)
	return printedLog
}

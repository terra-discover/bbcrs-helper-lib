package lib

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	guuid "github.com/google/uuid"
)

// GenUUIDString func
func GenUUIDString() string {
	id := guuid.New().String()
	return id
}

// StringToUUID func
func StringToUUID(s string) *guuid.UUID {
	res, _ := guuid.Parse(s)
	return &res
}

// GenUUID func
func GenUUID() *guuid.UUID {
	id := guuid.New()
	return &id
}

// ConvertSlicePtrUUIDToStr func
// format: example -> '%s'
func ConvertSlicePtrUUIDToStr(listUuid []*guuid.UUID, separator string, format string) (result string) {
	if len(listUuid) < 1 {
		return ""
	}

	var listStrUuid []string

	for _, item := range listUuid {
		if item == nil {
			continue
		}

		listStrUuid = append(listStrUuid,
			fmt.Sprintf(
				format,
				item.String(),
			),
		)
	}

	result = strings.Join(listStrUuid, separator)

	return result
}

// ConvertSliceUUIDToStr func
// format: example -> '%s'
func ConvertSliceUUIDToStr(listUuid []guuid.UUID, separator string, format string) (result string) {
	if len(listUuid) < 1 {
		return ""
	}

	var listStrUuid []string

	for _, item := range listUuid {
		listStrUuid = append(listStrUuid,
			fmt.Sprintf(
				format,
				item.String(),
			),
		)
	}

	result = strings.Join(listStrUuid, separator)

	return result
}

// RemoveDuplicatedUUID
func RemoveDuplicatedUUID(listUUID []uuid.UUID) (result []uuid.UUID) {
	if len(listUUID) == 0 {
		return
	}

loopOutside:
	for _, item := range listUUID {
		isMatch := false

	loopInside:
		for _, tempItem := range result {
			if tempItem == item {
				isMatch = true
				break loopInside
			}
		}

		if isMatch {
			continue loopOutside
		}

		result = append(result, item)
	}

	return
}

// FindSliceUUID Find UUID on slice
func FindSliceUUID(slice []uuid.UUID, val uuid.UUID) (index int, isFound bool) {
	for i, item := range slice {
		if item == val {
			index = i
			isFound = true
			return
		}
	}
	return -1, false
}

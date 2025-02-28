package lib

import (
	"reflect"

	"github.com/google/uuid"
)

// RemoveDuplicateString
func RemoveDuplicateString(slice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// RemoveDuplicateString
func RemoveDuplicateUUID(slice []uuid.UUID) []uuid.UUID {
	allKeys := make(map[uuid.UUID]bool)
	list := []uuid.UUID{}
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
func RemoveEmptyString(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// RemoveDuplicateInt
func RemoveDuplicateInt(slice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// RemoveDuplicateFloat64
func RemoveDuplicateFloat64(slice []float64) []float64 {
	allKeys := make(map[float64]bool)
	list := []float64{}
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// FindMatchBetweenString
func FindMatchBetweenString(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// findMinAndMaxFloat64
func FindMinAndMaxFloat64(a []float64) (min float64, max float64) {
	if len(a) > 0 {
		min = a[0]
		max = a[0]
		for _, value := range a {
			if value < min {
				min = value
			}
			if value > max {
				max = value
			}
		}
	}
	return min, max
}

// FindMinAndMaxInt
func FindMinAndMaxInt(a []int) (min int, max int) {
	if len(a) > 0 {
		min = a[0]
		max = a[0]
	}
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

// Slice Contains
func SliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// Slice Int Contains
func SliceIntContains(s []int, num int) bool {
	for _, v := range s {
		if v == num {
			return true
		}
	}

	return false
}

// Slice UUID Contains
func SliceUUIDContains(s []uuid.UUID, compare uuid.UUID) bool {
	for _, v := range s {
		if v == compare {
			return true
		}
	}

	return false
}

// FindMapKeyByValue | find key map with value
func FindMapKeyByValue(m map[string]string, value string) (key string, ok bool) {
	for k, v := range m {
		if v == value {
			key = k
			ok = true
			return
		}
	}
	return
}

func FindMapValueByKey(m map[string]string, key string) (value string, ok bool) {
	for k, v := range m {
		if k == key {
			value = v
			ok = true
		}
	}

	return
}

// ArrStringToCommas
// list := []string{"paid", "due", "partial"}
func ArrStringToCommas(list []string) string {
	var comma string
	for i, id := range list {
		if i == 0 {
			comma = `"` + id + `"`
		} else {
			comma += `, "` + id + `"`
		}
	}
	return comma // expected result "paid", "due", "partial"
}

// SlicesContains
func SlicesContains(slice1, slice2 interface{}) bool {
	sliceValue1 := reflect.ValueOf(slice1)
	sliceValue2 := reflect.ValueOf(slice2)
	if sliceValue1.Kind() != reflect.Slice || sliceValue2.Kind() != reflect.Slice {
		return false
	}
	// Check if the lengths of slices are the same
	if sliceValue1.Len() != sliceValue2.Len() {
		return false
	}
	// Compare individual elements of the slices
	for i := 0; i < sliceValue1.Len(); i++ {
		elem1 := sliceValue1.Index(i).Interface()
		elem2 := sliceValue2.Index(i).Interface()

		if !reflect.DeepEqual(elem1, elem2) {
			return false
		}
	}
	return true
}

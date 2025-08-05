package lib

import (
	"sort"
	"strings"
	"time"
)

func Sorting(sort string) (orderBy string) {
	sortArray := strings.Split(sort, ",")
	orderBy = ""
	for _, val := range sortArray {
		if val != "" {
			val = strings.Trim(val, " ")
			sorting := strings.Split(val, "-")
			sortingLength := len(sorting)
			switch sortingLength {
			case 1:
				orderBy += sorting[0] + " ASC,"
			case 2:
				orderBy += sorting[1] + " DESC,"
			}
		}
	}

	orderBy = TrimSuffix(orderBy, ",")
	return
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func SliceTimeAsc(listTime []time.Time) (result []time.Time) {
	result = listTime

	sort.Slice(result, func(i, j int) bool {
		return result[i].Before(result[j])
	})

	return
}

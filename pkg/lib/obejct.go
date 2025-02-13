package lib

import (
	"regexp"
	"strings"
)

// ObjectToSingleLevel convert nested object to single level object or array
func ObjectToSingleLevel(source map[string]interface{}, fields []string, target interface{}) {
	row := map[string]interface{}{}
	for f := range fields {
		field := fields[f]
		if strings.Contains(field, ".") {
			parent := regexp.MustCompile(`^([^\.]+).+`).ReplaceAllString(field, "$1")
			child := regexp.MustCompile(`^[^\.]+\.(.*)`).ReplaceAllString(field, "$1")
			parentSource, ok := source[parent].(map[string]interface{})
			if ok {
				if !strings.Contains(child, ".") {
					row[field] = parentSource[child]
				} else {
					rowField := map[string]interface{}{}
					ObjectToSingleLevel(parentSource, []string{child}, &rowField)
					if value, ok := rowField[child]; ok {
						row[field] = value
					}
				}
			}
			continue
		}

		if value, ok := source[field]; ok {
			row[field] = value
		}
	}

	Merge(row, target)
}

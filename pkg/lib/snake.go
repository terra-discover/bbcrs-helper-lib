package lib

import (
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

// SnakeCase converts a string into snake case.
func SnakeCase(s string, keepUnderscores ...bool) string {
	snake := strings.ToLower(strcase.ToSnake(s))
	if len(keepUnderscores) == 0 || !keepUnderscores[0] {
		re := regexp.MustCompile(`[_]+`)
		snake = re.ReplaceAllString(snake, "_")
	}
	return snake
}

// UpperSnakeCase converts a string into snake case with capital letters.
func UpperSnakeCase(s string, keepUnderscores ...bool) string {
	return strings.ToUpper(SnakeCase(s, keepUnderscores...))
}

package lib

import (
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

// KebabCase converts a string into kebab case.
func KebabCase(s string, keepDashes ...bool) string {
	kebab := strings.ToLower(strcase.ToKebab(s))
	if len(keepDashes) == 0 || !keepDashes[0] {
		re := regexp.MustCompile(`[-]+`)
		kebab = re.ReplaceAllString(kebab, "-")
	}
	return kebab
}

// UpperKebabCase converts a string into kebab case with capital letters.
func UpperKebabCase(s string) string {
	return strings.ToUpper(KebabCase(s))
}

package names

import (
	"regexp"
)

var fieldFirstCharRegex = regexp.MustCompile("^[_A-Za-z]$")
var fieldBodyCharRegex = regexp.MustCompile("^[_0-9A-Za-z]$")

func FilterNotSupportedFieldNameCharacters(fieldName string) string {
	runes := []rune(fieldName)
	if len(runes) == 0 {
		return fieldName
	}
	for len(runes) > 0 {
		if fieldFirstCharRegex.MatchString(string(runes[0])) {
			break
		}
		runes = runes[1:]
	}
	for i := 1; i < len(runes); i++ {
		if fieldBodyCharRegex.MatchString(string(runes[i])) {
			continue
		}
		runes = append(runes[:i], runes[i+1:]...)
		i--
	}
	return string(runes)
}

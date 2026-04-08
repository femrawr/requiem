package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func ReplaceString(text *string, name string, value string) {
	re := regexp.MustCompile(fmt.Sprintf(
		`%s\s+string\s*=\s*"([^"]*)"`,
		regexp.QuoteMeta(name),
	))

	*text = re.ReplaceAllString(*text, fmt.Sprintf(
		`%s string = "%s"`, name, value,
	))
}

func ReplaceBool(text *string, name string, value bool) {
	re := regexp.MustCompile(fmt.Sprintf(
		`%s\s+bool\s*=\s*(true|false)`,
		regexp.QuoteMeta(name),
	))

	*text = re.ReplaceAllString(*text, fmt.Sprintf(
		`%s bool = %t`, name, value,
	))
}

func ReplaceInt(text *string, name string, value string) {
	re := regexp.MustCompile(fmt.Sprintf(
		`%s\s+int\s*=\s*\d+`,
		regexp.QuoteMeta(name),
	))

	converted, _ := strconv.Atoi(value)

	*text = re.ReplaceAllString(*text, fmt.Sprintf(
		`%s int = %d`, name, converted,
	))
}

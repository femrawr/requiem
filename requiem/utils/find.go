package utils

import (
	"strconv"
	"strings"
)

func FindNumber(text string) (int, bool) {
	for word := range strings.FieldsSeq(text) {
		num, err := strconv.Atoi(word)
		if err == nil {
			return num, true
		}
	}

	return 0, false
}

func FindLastWord(text string) string {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return ""
	}

	return fields[len(fields)-1]
}

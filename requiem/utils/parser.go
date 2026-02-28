package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func UnwrapQuotes(text string) string {
	first := strings.Index(text, "\"")
	last := strings.LastIndex(text, "\"")

	if first == -1 || last == -1 || first == last {
		return ""
	}

	return text[first+1 : last]
}

func ExtractFromQuotes(text string, flag string) string {
	pattern := regexp.MustCompile(fmt.Sprintf(`%s="([^"]*)"`, flag))

	match := pattern.FindStringSubmatch(text)
	if len(match) < 2 {
		return ""
	}

	return match[1]
}

func HasFlag(text string, flag string) bool {
	for part := range strings.SplitSeq(text, " ") {
		if part != fmt.Sprintf("-%s", flag) {
			continue
		}

		return true
	}

	return false
}

package jira

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`\b[a-zA-Z][a-zA-Z0-9]+-\d+\b`)

func ParseKeys(text string) []string {
	return uniq(re.FindAllString(text, -1))
}

func uniq(str []string) []string {
	n := 0
	m := make(map[string]struct{}, len(str))
	for _, s := range str {
		if _, ok := m[s]; !ok {
			str[n] = s
			n++
			m[s] = struct{}{}
		}
	}
	return str[:n]
}

func StripKey(key string) string {
	idx := strings.Index(key, "-")
	if idx != -1 {
		return key[:idx]
	}
	return key
}

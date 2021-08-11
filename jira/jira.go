package jira

import (
	"github.com/pru-mike/rocketchat-jira-webhook/utils"
	"regexp"
	"strings"
)

var findKeysRegexp = regexp.MustCompile(`\b[a-zA-Z][a-zA-Z0-9]+-\d+\b`)

func parseKeys(re *regexp.Regexp, text string) []string {
	return utils.Uniq(re.FindAllString(text, -1))
}

func StripKey(key string) string {
	idx := strings.Index(key, "-")
	if idx != -1 {
		return key[:idx]
	}
	return key
}

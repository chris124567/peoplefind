package htmlhelp

import (
	"regexp"
	"strings"
)

func StringInSlice(slice []string, s string) bool {
	for _, v := range slice {
		if s == v {
			return true
		}
	}
	return false
}

func MatchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
		// (?flags): set flags within current group; non-capturing
		// s: let . match \n (default false)
		// https://github.com/google/re2/wiki/Syntax
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

// returns index to longest string in array
func GetLongestString(array []string) int {
	var currentLength int = 0
	var longestIndex int = 0
	var longestLength int = 0

	for index, value := range array {
		currentLength = len(value)
		if currentLength > longestLength {
			longestIndex = index
			longestLength = currentLength
		}
	}
	return longestIndex
}

func JsonUnescape(str string) string {
	return strings.Replace(str, `\`, ``, -1)
}

func GetStringInBetween(content string, start string, end string) (result string) {
	if content != "" && start != "" && end != "" {
		content := strings.ReplaceAll(content, "\n", " ")
		r := strings.Split(content, start)

		if len(r) < 2 {
			return
		}

		if r[1] != "" {
			r = strings.Split(r[1], end)
		}

		result = strings.TrimSpace(r[0])
		return
	} else {
		return
	}
}

func CleanString(content string) string {
	return strings.TrimSpace(strings.ReplaceAll(content, "\n", " "))
}

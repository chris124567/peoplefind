package httphelp

import (
	"strings"
)

func ResolveRelativeScheme(link string) string {
	if strings.HasPrefix(link, "//") {
		link = strings.Replace(link, "//", "https://", 1)
	}
	return link
}

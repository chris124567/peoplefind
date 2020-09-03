package models

import (
	"strings"
)

func TitleFormat(titles ...string) string {
	var combined string = ""
	var newValue string = ""
	for _, title := range titles {
		newValue = strings.TrimSpace(strings.Title(strings.ToLower(title)))
		combined += newValue
		if newValue != "" {
			combined += " "
		}
	}
	return strings.TrimSpace(combined)
}

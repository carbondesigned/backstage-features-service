package utils

import "strings"

func GenerateSlugFromTitle(title string) string {
	return strings.ToLower(strings.Replace(title, " ", "-", -1))
}

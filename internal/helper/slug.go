package helper

import (
	"regexp"
	"strings"
)

func GenerateSlug(title string) string {
	lowerTitle := strings.ToLower(title)

	// Remove all spaces and replace them with hyphens
	slug := strings.Replace(lowerTitle, " ", "-", -1)

	// Remove all non-alphanumeric characters
	re := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = re.ReplaceAllString(slug, "")

	return slug
}

package models

import (
	"regexp"
	"strings"
	"time"
)

type PostDto struct {
	ID          string    `json:"_id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Content     string    `json:"content"`
	Media       []string  `json:"media"`
	Collections []string  `json:"collections"`
	Tags        []string  `json:"tags"`
	Published   time.Time `json:"published"`
	LastUpdated time.Time `json:"last_updated"`
	Status      string    `json:"status"`
	Featured    bool      `json:"featured"`
}

func GenerateSlug(title string) string {
	// Convert the title to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with dashes
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters using a regular expression
	slug = removeSpecialCharacters(slug)

	// Trim leading and trailing dashes
	slug = strings.Trim(slug, "-")

	return slug
}

// Helper function to remove special characters
func removeSpecialCharacters(s string) string {
	// Define a regular expression to match non-alphanumeric characters except for dashes
	re := regexp.MustCompile(`[^\w-]+`)
	return re.ReplaceAllString(s, "")
}

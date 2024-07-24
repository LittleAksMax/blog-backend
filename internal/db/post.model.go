package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strings"
	"time"
)

type PostStatus uint8

const (
	Draft PostStatus = iota
	Published
	Archived
	Removed
)

const (
	draftString     = "Draft"
	publishedString = "Published"
	archivedString  = "Archived"
	removedString   = "Removed"
)

var postStatusStrings [4]string = [...]string{"Draft", "Published", "Archived", "Removed"}

func (ps PostStatus) String() string {
	return postStatusStrings[ps]
}

func PostStatusFromString(s string) (PostStatus, bool) {
	switch s {
	case draftString:
		return Draft, true
	case publishedString:
		return Published, true
	case archivedString:
		return Archived, true
	case removedString:
		return Removed, true
	default:
		return 255, false
	}
}

type Post struct {
	Id          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	Content     string             `bson:"content"`
	Media       []string           `bson:"media"`
	Collections []string           `bson:"collections"`
	Tags        []string           `bson:"tags"`
	Published   time.Time          `bson:"published"`
	LastUpdated time.Time          `bson:"last_updated"`
	Status      PostStatus         `bson:"status"`
	Featured    bool               `bson:"featured"`
}

func (p *Post) Slug() string {
	// Convert the title to lowercase
	slug := strings.ToLower(p.Title)

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

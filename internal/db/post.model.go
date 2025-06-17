package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostStatus uint8

const (
	None PostStatus = iota
	Published
	Archived
)

const (
	noneString      = "None"
	publishedString = "Published"
	archivedString  = "Archived"
)

var postStatusStrings [2]string = [...]string{"Published", "Archived"}

func (ps PostStatus) String() string {
	return postStatusStrings[ps-1]
}

func PostStatusFromString(s string) (PostStatus, bool) {
	switch s {
	case publishedString:
		return Published, true
	case archivedString:
		return Archived, true
	default:
		return 255, false
	}
}

type Post struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Slug        string             `bson:"slug"`
	Content     string             `bson:"content"`
	Media       []string           `bson:"media"`
	Banner      string             `bson:"banner"`
	Collections []string           `bson:"collections"`
	Tags        []string           `bson:"tags"`
	Published   time.Time          `bson:"published"`
	LastUpdated time.Time          `bson:"last_updated"`
	Status      PostStatus         `bson:"status"`
	Featured    bool               `bson:"featured"`
}

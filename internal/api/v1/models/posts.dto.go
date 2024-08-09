package models

import (
	"time"
)

type PostDto struct {
	ID          string    `json:"_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Media       []string  `json:"media"`
	Collections []string  `json:"collections"`
	Tags        []string  `json:"tags"`
	Published   time.Time `json:"published"`
	LastUpdated time.Time `json:"last_updated"`
	Status      string    `json:"status"`
	Featured    bool      `json:"featured"`
}

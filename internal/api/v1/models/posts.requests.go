package models

const (
	DefaultPageNumber = 1
	DefaultPageSize   = 10
)

type PagedQuery struct {
	PageNum  int
	PageSize int
}

type PostQuery struct {
	Title       string
	Tags        []string
	Collections []string
	Featured    bool
}

type CreatePostRequest struct {
	Title       string   `json:"title" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	Collections []string `json:"collections" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
	Featured    *bool    `json:"featured" binding:"required"` // pointer to make validation work
}

type UpdatePostRequest struct {
	Title       string   `binding:"required" json:"title"`
	Content     string   `binding:"required" json:"content"`
	Collections []string `binding:"required" json:"collections"`
	Tags        []string `binding:"required" json:"tags"`
	Featured    *bool    `binding:"required" json:"featured"`
}

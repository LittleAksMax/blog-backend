package models

const (
	DefaultPageNumber = 1
	DefaultPageSize   = 10
)

type PagedQuery struct {
	PageNum  int `form:"page_num" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1"`
}

type PostFilter struct {
	Title       *string  `form:"title"`
	Tags        []string `form:"tags"`
	Collections []string `form:"collections"`
	Featured    *bool    `form:"featured"`
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
	Status      string   `binding:"required" json:"status"`
	Featured    *bool    `binding:"required" json:"featured"`
}

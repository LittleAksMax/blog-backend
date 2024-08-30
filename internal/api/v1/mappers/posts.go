package mappers

import (
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"github.com/LittleAksMax/blog-backend/internal/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToDto(post *db.Post) models.PostDto {
	return models.PostDto{
		ID:          post.ID.Hex(),
		Title:       post.Title,
		Slug:        post.Slug,
		Content:     post.Content,
		Media:       post.Media,
		Collections: post.Collections,
		Tags:        post.Tags,
		Status:      post.Status.String(),
		Published:   post.Published,
		LastUpdated: post.LastUpdated,
		Featured:    post.Featured,
	}
}

func ToDom(id primitive.ObjectID, postStatus db.PostStatus, dto *models.PostDto) db.Post {
	return db.Post{
		ID:          id,
		Title:       dto.Title,
		Slug:        dto.Slug,
		Content:     dto.Content,
		Media:       dto.Media,
		Collections: dto.Collections,
		Tags:        dto.Tags,
		Status:      postStatus,
		Published:   dto.Published,
		LastUpdated: dto.LastUpdated,
		Featured:    dto.Featured}
}

package services

import (
	"context"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService interface {
	GetPosts(ctx context.Context, pf *models.PagedQuery, pq *models.PostQuery) ([]models.PostDto, error)
	GetPost(ctx context.Context, id primitive.ObjectID) (*models.PostDto, error)
	CreatePost(ctx context.Context, dto *models.PostDto) error
	UpdatePost(ctx context.Context, id primitive.ObjectID, dto *models.PostDto) error
	DeletePost(ctx context.Context, id primitive.ObjectID) error
}

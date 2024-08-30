package services

import (
	"context"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService interface {
	GetPosts(ctx context.Context, pq *models.PagedQuery, pf *models.PostFilter) ([]models.PostDto, int, error)
	GetPostById(ctx context.Context, id primitive.ObjectID) (*models.PostDto, error)
	GetPostBySlug(ctx context.Context, slug string) (*models.PostDto, error)
	CreatePost(ctx context.Context, dto *models.PostDto) error
	UpdatePost(ctx context.Context, id primitive.ObjectID, dto *models.PostDto) error
	DeletePost(ctx context.Context, id primitive.ObjectID) error
}

package services

import (
	"context"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService interface {
	GetPosts(ctx context.Context, pq *models.PagedQuery, pf *models.PostFilter, admin bool) ([]models.PostDto, int, error)
	GetPostById(ctx context.Context, id primitive.ObjectID, admin bool) (*models.PostDto, error)
	GetPostBySlug(ctx context.Context, slug string, admin bool) (*models.PostDto, error)
	CreatePost(ctx context.Context, dto *models.PostDto) error
	UpdatePost(ctx context.Context, id primitive.ObjectID, dto *models.PostDto) error
	ArchivePost(ctx context.Context, id primitive.ObjectID) error
	DeletePost(ctx context.Context, id primitive.ObjectID) error
}

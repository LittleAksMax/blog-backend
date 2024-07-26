package services

import (
	"context"
	"errors"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"github.com/LittleAksMax/blog-backend/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostServiceImpl struct {
	posts *mongo.Collection
}

func NewPostServiceImpl(dbCfg *db.Config) *PostServiceImpl {
	return &PostServiceImpl{
		posts: dbCfg.Posts,
	}
}

func (ps *PostServiceImpl) GetPosts(ctx context.Context, pq *models.PagedQuery, pf *models.PostFilter) ([]models.PostDto, error) {
	cursor, err := ps.posts.Find(ctx, bson.D{}) // TODO: filters from pf and pq

	if err != nil {
		return []models.PostDto{}, err
	}

	defer func() {
		_ = cursor.Close(ctx)
	}()

	// convert result cursor to a slice of DTOs
	dtos := make([]models.PostDto, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var post db.Post
		err := cursor.Decode(&post)
		if err != nil {
			return nil, err
		}
		dtos[cursor.RemainingBatchLength()] = models.PostDto{
			Id:          post.Id.Hex(),
			Title:       post.Title,
			Content:     post.Content,
			Media:       post.Media,
			Collections: post.Collections,
			Tags:        post.Tags,
			Published:   post.Published,
			LastUpdated: post.LastUpdated,
			Status:      post.Status.String(),
			Featured:    post.Featured,
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return dtos, nil
}

func (ps *PostServiceImpl) GetPost(ctx context.Context, id primitive.ObjectID) (*models.PostDto, error) {
	var post db.Post
	err := ps.posts.FindOne(ctx, bson.D{{"_id", id}}).Decode(&post)

	if err != nil {
		return nil, &NotFoundErr{id: id}
	}

	return &models.PostDto{
		Id:          post.Id.Hex(),
		Title:       post.Title,
		Content:     post.Content,
		Media:       post.Media,
		Collections: post.Collections,
		Tags:        post.Tags,
		Status:      post.Status.String(),
		Published:   post.Published,
		LastUpdated: post.LastUpdated,
		Featured:    post.Featured,
	}, nil
}

func (ps *PostServiceImpl) CreatePost(ctx context.Context, dto *models.PostDto) error {
	// add data to mongo database
	postStatus, ok := db.PostStatusFromString(dto.Status)

	if !ok {
		panic("invalid post status")
	}

	post := db.Post{
		Id:          primitive.NewObjectID(),
		Title:       dto.Title,
		Content:     dto.Content,
		Media:       dto.Media,
		Collections: dto.Collections,
		Tags:        dto.Tags,
		Status:      postStatus,
		Published:   dto.Published,
		LastUpdated: dto.LastUpdated,
		Featured:    dto.Featured,
	}
	res, err := ps.posts.InsertOne(ctx, post)

	if err != nil {
		return err
	}

	// set id for dto
	dto.Id = res.InsertedID.(primitive.ObjectID).Hex()

	return nil
}

func (ps *PostServiceImpl) UpdatePost(ctx context.Context, id primitive.ObjectID, dto *models.PostDto) error {
	// add data to mongo database
	postStatus, ok := db.PostStatusFromString(dto.Status)

	if !ok {
		panic("invalid post status field")
	}

	post := db.Post{
		Id:          id,
		Title:       dto.Title,
		Content:     dto.Content,
		Media:       dto.Media,
		Collections: dto.Collections,
		Tags:        dto.Tags,
		Status:      postStatus,
		Published:   dto.Published,
		LastUpdated: dto.LastUpdated,
		Featured:    dto.Featured,
	}
	res, err := ps.posts.ReplaceOne(ctx, primitive.D{{"_id", id}}, post)

	if err != nil {
		// chances are this is a conflict
		return &ConflictErr{msg: err.Error()}
	}

	if res.MatchedCount == 0 {
		return &NotFoundErr{id: id}
	}

	// ensure correctly set DTO value for ID
	dto.Id = id.Hex()

	return nil
}

func (ps *PostServiceImpl) DeletePost(ctx context.Context, id primitive.ObjectID) error {
	var post db.Post
	err := ps.posts.FindOneAndDelete(ctx, bson.D{{"_id", id}}).Decode(&post)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &NotFoundErr{id: id}
		} else {
			return err
		}
	}

	return nil
}

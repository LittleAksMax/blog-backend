package services

import (
	"context"
	"errors"

	"github.com/LittleAksMax/blog-backend/internal/api/v1/mappers"

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

func (ps *PostServiceImpl) GetPosts(ctx context.Context, pq *models.PagedQuery, pf *models.PostFilter, admin bool) ([]models.PostDto, int, error) {
	qfb := createFilterWithPostFilter(pf)

	// non-admins can only see published (i.e., not archived media)
	if !admin {
		qfb.addFilter("status", db.Published)
	}
	filter := qfb.build()

	// count total number of documents for HATEOAS in pagination
	totalCount, err := ps.posts.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if totalCount <= int64(pq.PageNum-1)*int64(pq.PageSize) {
		return []models.PostDto{}, 0, PageNotFoundErr{pageNum: pq.PageNum, pageSize: pq.PageSize}
	}

	// fetch documents in current batch using current filter
	opts := createOptionsWithPagedQuery(pq)
	// sort in descending order of published (newest first)
	opts.SetSort(bson.D{{"published", -1}})
	cursor, err := ps.posts.Find(ctx, filter, opts)

	if err != nil {
		return []models.PostDto{}, 0, err
	}

	defer func() {
		_ = cursor.Close(ctx)
	}()

	// convert result cursor to a slice of DTOs
	batchCount := cursor.RemainingBatchLength()
	dtos := make([]models.PostDto, batchCount)
	for cursor.Next(ctx) {
		var post db.Post
		err := cursor.Decode(&post)
		if err != nil {
			return nil, 0, err
		}
		dtos[cursor.RemainingBatchLength()] = mappers.ToDto(&post)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return dtos, int(totalCount), nil
}

func (ps *PostServiceImpl) GetPostById(ctx context.Context, id primitive.ObjectID, admin bool) (*models.PostDto, error) {
	qfb := newQueryFilterBuilder()
	qfb.addFilter("_id", id)

	// non-admins can only see published (i.e., not archived media)
	if !admin {
		qfb.addFilter("status", db.Published)
	}

	var post db.Post
	err := ps.posts.FindOne(ctx, qfb.build()).Decode(&post)

	if err != nil {
		return nil, NotFoundErr{id: id}
	}

	postDto := mappers.ToDto(&post)
	return &postDto, nil
}

func (ps *PostServiceImpl) GetPostBySlug(ctx context.Context, slug string, admin bool) (*models.PostDto, error) {
	qfb := newQueryFilterBuilder()
	qfb.addFilter("slug", slug)

	// non-admins can only see published (i.e., not archived media)
	if !admin {
		qfb.addFilter("status", db.Published)
	}

	var post db.Post
	err := ps.posts.FindOne(ctx, qfb.build()).Decode(&post)

	if err != nil {
		return nil, SlugNotFoundErr{slug: slug}
	}

	postDto := mappers.ToDto(&post)
	return &postDto, nil
}

func (ps *PostServiceImpl) CreatePost(ctx context.Context, dto *models.PostDto) error {
	// add data to mongo database
	postStatus, ok := db.PostStatusFromString(dto.Status)

	if !ok {
		panic("invalid post status")
	}

	post := mappers.ToDom(primitive.NewObjectID(), postStatus, dto)
	res, err := ps.posts.InsertOne(ctx, post)

	if err != nil {
		return ConflictErr{msg: err.Error()}
	}

	// set id for dto
	dto.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return nil
}

func (ps *PostServiceImpl) UpdatePost(ctx context.Context, id primitive.ObjectID, dto *models.PostDto) error {
	// add data to mongo database
	postStatus, ok := db.PostStatusFromString(dto.Status)

	if !ok {
		panic("invalid post status field")
	}

	post := mappers.ToDom(id, postStatus, dto)
	res, err := ps.posts.ReplaceOne(ctx, primitive.D{{Key: "_id", Value: id}}, post)

	if err != nil {
		// chances are this is a conflict
		return ConflictErr{msg: err.Error()}
	}

	if res.MatchedCount == 0 {
		return NotFoundErr{id: id}
	}

	// ensure correctly set DTO value for ID
	dto.ID = id.Hex()

	return nil
}

func (ps *PostServiceImpl) DeletePost(ctx context.Context, id primitive.ObjectID) error {
	var post db.Post
	err := ps.posts.FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&post)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return NotFoundErr{id: id}
		} else {
			return err
		}
	}

	return nil
}

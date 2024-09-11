package services

import (
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type queryFilterBuilder struct {
	query bson.D
}

func newQueryFilterBuilder() *queryFilterBuilder {
	return &queryFilterBuilder{query: bson.D{}}
}

func (builder *queryFilterBuilder) addFilter(key string, value any) *queryFilterBuilder {
	builder.query = append(builder.query, bson.E{Key: key, Value: value})
	return builder
}

func (builder *queryFilterBuilder) build() bson.D {
	return builder.query
}

func createFilterWithPostFilter(pf *models.PostFilter) *queryFilterBuilder {
	qfb := newQueryFilterBuilder()
	if pf.Title != nil {
		qfb.addFilter("title", *pf.Title)
	}
	if pf.Tags != nil {
		qfb.addFilter("tags", pf.Tags)
	}
	if pf.Collections != nil {
		qfb.addFilter("collections", pf.Collections)
	}
	if pf.Featured != nil {
		qfb.addFilter("featured", *pf.Featured)
	}
	return qfb
}

func createOptionsWithPagedQuery(pq *models.PagedQuery) *options.FindOptions {
	opts := options.Find()
	opts.SetSkip(int64((pq.PageNum - 1) * pq.PageSize)).SetLimit(int64(pq.PageSize))
	return opts
}

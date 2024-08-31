package caching

import (
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strconv"
	"strings"
)

type queryBuilder struct {
	builder *strings.Builder
}

func newQueryBuilder(strBuilder *strings.Builder) queryBuilder {
	return queryBuilder{
		builder: strBuilder,
	}
}

func (b queryBuilder) write(fieldName string, value string) {
	b.builder.WriteString("&" + fieldName + "=" + value)
}

func (b queryBuilder) writeSlice(fieldName string, values []string) {
	for _, value := range values {
		b.write(fieldName, value)
	}
}
func (b queryBuilder) build() string {
	return b.builder.String()
}

func mustGetKey(typeElem reflect.Type, fieldName string) string {
	field, ok := typeElem.FieldByName(fieldName)

	if !ok {
		panic("missing " + fieldName + " field on " + typeElem.Name() + " struct")
	}

	key, ok := field.Tag.Lookup("form")

	if !ok {
		panic("missing form field in tag of " + field.Name)
	}

	return key
}

func addPaginationFields(builder *queryBuilder, pq *models.PagedQuery) {
	typeElem := reflect.TypeOf(pq).Elem()

	key := mustGetKey(typeElem, "PageSize")
	builder.write(key, strconv.Itoa(pq.PageSize))

	key = mustGetKey(typeElem, "PageNum")
	builder.write(key, strconv.Itoa(pq.PageNum))
}

func addFilterFields(builder *queryBuilder, pf *models.PostFilter) {
	typeElem := reflect.TypeOf(pf).Elem()

	if pf.Title != nil {
		key := mustGetKey(typeElem, "Title")
		builder.write(key, *pf.Title)
	}
	if pf.Tags != nil {
		key := mustGetKey(typeElem, "Tags")
		builder.writeSlice(key, pf.Tags)
	}
	if pf.Collections != nil {
		key := mustGetKey(typeElem, "Collections")
		builder.writeSlice(key, pf.Collections)
	}
	if pf.Featured != nil {
		key := mustGetKey(typeElem, "Featured")
		builder.write(key, strconv.FormatBool(*pf.Featured))
	}
}

func HashGetAllPosts(ctx *gin.Context) string {
	pq := ctx.MustGet("pq").(*models.PagedQuery)
	pf := ctx.MustGet("pf").(*models.PostFilter)

	// append prefix for this type of request
	strBuilder := strings.Builder{}
	strBuilder.WriteString("getAllPosts-")

	builder := newQueryBuilder(&strBuilder)
	// append pagination queries first as they have default values
	addPaginationFields(&builder, pq)

	// append required filter queries
	addFilterFields(&builder, pf)

	return builder.build()
}

func HashGetPost(ctx *gin.Context) string {
	// NOTE: must be run after handler is run
	idUsed := ctx.MustGet("idOrSlug").(bool)

	if idUsed {
		id := ctx.MustGet("id").(primitive.ObjectID)
		return "getPost-id-" + id.Hex()
	} else {
		slug := ctx.MustGet("slug").(string)
		return "getPost-slug-" + slug
	}
}

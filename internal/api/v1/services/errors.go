package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

type NotFoundErr struct {
	id primitive.ObjectID
}

func (e NotFoundErr) Error() string {
	return "not found with id: " + e.id.Hex()
}

type SlugNotFoundErr struct {
	slug string
}

func (e SlugNotFoundErr) Error() string {
	return "not found with slug: " + e.slug
}

type PageNotFoundErr struct {
	pageNum  int
	pageSize int
}

func (e PageNotFoundErr) Error() string {
	return "not found page " + strconv.Itoa(e.pageNum) + " with size: " + strconv.Itoa(e.pageSize)
}

type ConflictErr struct {
	msg string
}

func (e ConflictErr) Error() string {
	return e.msg
}

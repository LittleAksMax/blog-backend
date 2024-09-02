package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotFoundErr struct {
	id primitive.ObjectID
}

func (e NotFoundErr) Id() primitive.ObjectID {
	return e.id
}

func (e NotFoundErr) Error() string {
	return "not found with id: " + e.id.Hex()
}

type SlugNotFoundErr struct {
	slug string
}

func (e SlugNotFoundErr) Slug() string {
	return e.slug
}

func (e SlugNotFoundErr) Error() string {
	return "not found with slug: " + e.slug
}

type ConflictErr struct {
	msg string
}

func (e ConflictErr) Error() string {
	return e.msg
}

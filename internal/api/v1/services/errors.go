package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotFoundErr struct {
	id primitive.ObjectID
}

type SlugNotFoundErr struct {
	NotFoundErr
	slug string
}

func (e NotFoundErr) Id() primitive.ObjectID {
	return e.id
}

func (e SlugNotFoundErr) Slug() string {
	return e.slug
}

func (e NotFoundErr) Error() string {
	return "not found with id: " + e.id.Hex()
}

func (e SlugNotFoundErr) Error() string {
	return "not found with slug: " + e.slug
}

type ConflictErr struct {
	msg string
}

func (err ConflictErr) Error() string {
	return err.msg
}

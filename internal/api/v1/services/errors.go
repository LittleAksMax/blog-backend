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

type ConflictErr struct {
	field string
}

func (err ConflictErr) Error() string {
	return "conflict in field '" + err.field + "'. perhaps item with this field already exists"
}

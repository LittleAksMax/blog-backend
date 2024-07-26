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
	msg string
}

func (err ConflictErr) Error() string {
	return err.msg
}

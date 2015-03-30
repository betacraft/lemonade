package models

import (
	"time"
)

//go:generate easytags timestamp.go bson
//go:generate easytags timestamp.go json

type Timestamp struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" bson:"deleted_at"`
}

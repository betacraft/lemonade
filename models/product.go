package models

import (
	"gopkg.in/mgo.v2/bson"
)

//go:generate easytags Product.go bson
//go:generate easytags Product.go json

type Product struct {
	Id bson.ObjectId

	Name         string
	MainCategory string
	SubCategory  string

	PriceValue    int64
	PriceCurrency string

	Description string
	Attributes  map[string]string

	Timestamp
}

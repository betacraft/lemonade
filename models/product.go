package models

import (
	"github.com/rainingclouds/lemonade/db"
	"gopkg.in/mgo.v2/bson"
)

//go:generate easytags Product.go bson
//go:generate easytags Product.go json

const (
	POLL = 1
	DEAL = 2
)

type Product struct {
	Id bson.ObjectId `json:"id" bson:"id"`

	Name         string `json:"name" bson:"name"`
	MainCategory string `json:"main_category" bson:"main_category"`
	SubCategory  string `json:"sub_category" bson:"sub_category"`
	ProductImage string `json:"product_image" bson:"product_image"`
	State        int    `json:"state" bson:"state"`

	PriceValue    int64  `json:"price_value" bson:"price_value"`
	PriceCurrency string `json:"price_currency" bson:"price_currency"`

	Description string            `json:"description" bson:"description"`
	Attributes  map[string]string `json:"attributes" bson:"attributes"`

	Timestamp
}

func (p *Product) Create() error {
	return db.MgCreateStrong(c, data)
}

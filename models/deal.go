package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

//go:generate easytags deal.go bson
//go:generate easytags deal.go json

type Milestone struct {
	GroupSize   int64  `json:"group_size" bson:"group_size"`
	Price       string `json:"price" bson:"price"`
	Description string `json:"description" bson:"description"`
}

type Comment struct {
	By       bson.ObjectId `json:"by" bson:"by"`
	Name     string        `json:"name" bson:"name"`
	Commment string        `json:"commment" bson:"commment"`
	Upvotes  int64         `json:"upvotes" bson:"upvotes"`
}

type Deal struct {
	Id bson.ObjectId `json:"id" bson:"id"`

	Product Product `json:"product" bson:"product"`

	CurrentCount int64         `json:"current_count" bson:"current_count"`
	CreatedBy    bson.ObjectId `json:"created_by" bson:"created_by"`

	ExpiresOn         time.Time `json:"expires_on" bson:"expires_on"`
	DaysLeft          int64     `json:"days_left" bson:"days_left"`
	LowestOnlinePrice string    `json:"lowest_online_price" bson:"lowest_online_price"`

	CurrentMilestone Milestone   `json:"current_milestone" bson:"current_milestone"`
	Milestones       []Milestone `json:"milestones" bson:"milestones"`
	Comments         []Comment   `json:"comments" bson:"comments"`
}

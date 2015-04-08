package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Milestone struct {
	GroupSize   int64
	Description string
}

type Comment struct {
	By       bson.ObjectId
	Name     string
	Commment string
	Upvotes  int64
}

type Deal struct {
	Id bson.ObjectId

	Product Product

	CurrentCount int64
	CreatedBy    bson.ObjectId

	Milestones []Milestone
	Comments   []Comment
}

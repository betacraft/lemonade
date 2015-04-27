package models

import (
	"github.com/rainingclouds/lemonades/db"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//go:generate easytags group.go bson
//go:generate easytags group.go json

type Group struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`

	Product   Product       `json:"product" bson:"product"`
	CreatedBy bson.ObjectId `json:"created_by" bson:"created_by"`

	InterestedUsers      []bson.ObjectId `json:"interested_users" bson:"interested_users"`
	InterestedUsersCount int64           `json:"interested_users_count" bson:"interested_users_count"`

	RequiredUserCount int64 `json:"required_user_count" bson:"required_user_count"`

	ReachedGoalOn time.Time `json:"-" bson:"reached_goal_on"`
	ReachedGoal   int64     `json:"reached_goal" bson:"-"`

	ExpiresOn time.Time `json:"expires_on" bson:"expires_on"`
	IsOn      bool      `json:"is_on" bson:"is_on"`

	ExpiresIn int64 `json:"expires_in" bson:"-"`

	IsJoined bool `json:"is_joined" bson:"-"`

	Timestamp
}

func (g *Group) Create() error {
	g.CreatedAt = time.Now()
	return db.MgCreateStrong(C_GROUP, g)
}

func (g *Group) Update() error {
	g.UpdatedAt = time.Now()
	return db.MgUpdateStrong(C_GROUP, g.Id, g)
}

func GetGroupByProductId(productId bson.ObjectId) (*Group, error) {
	group := new(Group)
	err := db.MgFindOneStrong(C_GROUP, &bson.M{"product._id": productId}, group)
	return group, err
}

func GetGroupById(id bson.ObjectId) (*Group, error) {
	group := new(Group)
	err := db.MgFindOneStrong(C_GROUP, &bson.M{"_id": id}, group)
	return group, err
}

func GetGroups(pageNo int) (*[]Group, error) {
	groups := new([]Group)
	err := db.MgFindPage(C_GROUP, &bson.M{"is_on": true}, pageNo, groups)
	return groups, err
}

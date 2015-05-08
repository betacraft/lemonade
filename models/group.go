package models

import (
	"github.com/rainingclouds/lemonades/db"
	"github.com/rainingclouds/lemonades/logger"
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

	MinDiscount string `json:"min_discount" bson:"min_discount"`

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

func UpdateProductInfo(p *Product) error {
	groups := new([]Group)
	err := db.MgFindAll(C_GROUP, &bson.M{"product._id": p.Id, "is_on": true}, groups)
	if err != nil {
		return err
	}
	for i := 0; i < len(*groups); i++ {
		(*groups)[i].Product = *p
		err = (*groups)[i].Update()
		if err != nil {
			logger.Err("Error while updating group product info", err)
		}
	}
	return nil
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
	err := db.MgFindPageSort(C_GROUP, &bson.M{"is_on": true}, "-interested_users_count", pageNo, groups)
	return groups, err
}

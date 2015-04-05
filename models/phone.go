package models

import (
	"errors"
	"github.com/mholt/binding"
	"github.com/rainingclouds/lemonade/db"
	"github.com/rainingclouds/lemonade/logger"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//go:generate easytags phone.go bson
//go:generate easytags phone.go json

type PhoneDeal struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name         string        `json:"name" bson:"name"`
	Manufacturer string        `json:"manufacturer" bson:"manufacturer"`

	ScreenSize string `json:"screen_size" bson:"screen_size"`
	CPU        string `json:"cpu" bson:"cpu"`
	GPU        string `json:"gpu" bson:"gpu"`
	Memory     string `json:"memory" bson:"memory"`
	Camera     string `json:"camera" bson:"camera"`
	OS         string `json:"os" bson:"os"`
	Battery    string `json:"battery" bson:"battery"`

	FullSpecificationLink string `json:"full_specification_link" bson:"full_specification_link"`
	ImageLink             string `json:"image_link" bson:"image_link"`

	MRP                string    `json:"mrp" bson:"mrp"`
	CurrentDiscount    string    `json:"current_discount" bson:"current_discount"`
	ClosingOn          time.Time `json:"closing_on" bson:"closing_on"`
	CurrentPeopleCount int       `json:"current_people_count" bson:"current_people_count"`
	NextDiscount       string    `json:"next_discount" bson:"next_discount"`
	NextDiscountFor    int       `json:"next_discount_for" bson:"next_discount_for"`
	SellerArea         string    `json:"seller_area" bson:"seller_area"`

	Timestamp
}

func (p PhoneDeal) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&p.Battery:               "battery",
		&p.CPU:                   "cpu",
		&p.Camera:                "camera",
		&p.ClosingOn:             "closing_on",
		&p.CurrentDiscount:       "current_discount",
		&p.CurrentPeopleCount:    "current_people_count",
		&p.FullSpecificationLink: "full_specification_link",
		&p.GPU:             "gpu",
		&p.Id:              "id",
		&p.ImageLink:       "image_link",
		&p.MRP:             "mrp",
		&p.Manufacturer:    "manufacturer",
		&p.Memory:          "memory",
		&p.Name:            "name",
		&p.NextDiscount:    "next_discount",
		&p.NextDiscountFor: "next_discount_for",
		&p.OS:              "os",
		&p.ScreenSize:      "screen_size",
		&p.SellerArea:      "seller_area",
	}
}

func (p *PhoneDeal) Create() error {
	phoneDeal := new(PhoneDeal)
	err := db.MgFindOneStrong(C_USERS, &bson.M{"name": p.Name, "manufacturer": p.Manufacturer}, phoneDeal)
	if err != nil && err.Error() != "not found" {
		logger.Err("While finding phone deal : " + err.Error())
		return err
	}
	if phoneDeal.Id != "" {
		logger.Err("Phone deal already exists")
		return errors.New("Phone deal already exists")
	}
	return db.MgCreateStrong(C_PHONE_DEALS, p)
}

func (p *PhoneDeal) Save() error {
	return db.MgUpdateStrong(C_PHONE_DEALS, p.Id, p)
}

func GetPhoneDealById(id bson.ObjectId) (*PhoneDeal, error) {
	phoneDeal := new(PhoneDeal)
	err := db.MgFindOneStrong(C_PHONE_DEALS, &bson.M{"_id": id}, phoneDeal)
	return phoneDeal, err
}

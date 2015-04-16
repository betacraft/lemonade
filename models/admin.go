package models

import (
	"errors"
	"github.com/mholt/binding"
	"github.com/nu7hatch/gouuid"
	"github.com/rainingclouds/lemonades/db"
	"github.com/rainingclouds/lemonades/framework"
	"github.com/rainingclouds/lemonades/logger"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//go:generate easytags admin.go bson
//go:generate easytags admin.go json

type Admin struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`

	Name     string `json:"name" bson:"name"`
	AuthKey  string `json:"auth_key" bson:"auth_key"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"-" bson:"password"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" bson:"deleted_at"`
}

func (a Admin) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&a.Name:    "name",
		&a.AuthKey: "auth_key",
		&a.Email:   "email",
	}
}

func CreateAdmin(adminMap map[string]interface{}) (*Admin, error) {
	a := Admin{}
	var ok bool
	a.Name, ok = adminMap["name"].(string)
	if !ok {
		return nil, errors.New("Name is not present")
	}
	a.Email, ok = adminMap["email"].(string)
	if !ok {
		return nil, errors.New("Email is not present")
	}
	a.Password, ok = adminMap["password"].(string)
	if !ok {
		return nil, errors.New("Password is not present")
	}
	return &a, nil
}

func GetAdminByAuthKey(authKey string) (*Admin, error) {
	a := new(Admin)
	err := db.MgFindOneStrong(C_ADMINS, &bson.M{"auth_key": authKey}, a)
	return a, err
}

func (a *Admin) Create() error {
	// check if already present
	// we will need to add them into the single query
	admin := new(Admin)
	err := db.MgFindOneStrong(C_ADMINS, &bson.M{"email": a.Email}, admin)
	if err != nil && err.Error() != "not found" {
		return err
	}
	if admin.AuthKey != "" {
		a = admin
		return nil
	}
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	authKey, err := uuid.NewV4()
	if err != nil {
		return err
	}
	a.Password = framework.MD5Hash(a.Password)
	a.AuthKey = authKey.String()
	return db.MgCreateStrong(C_ADMINS, a)
}

func (a *Admin) Update() error {
	a.UpdatedAt = time.Now()
	return db.MgUpdateStrong(C_ADMINS, a.Id, a)
}

func (a *Admin) Authenticate() error {
	logger.Get().Debug(a.Password)
	hashedPassword := framework.MD5Hash(a.Password)
	logger.Get().Debug("Hashed password is ", hashedPassword)
	err := db.MgFindOne(C_ADMINS, &bson.M{"email": a.Email, "password": hashedPassword}, a)
	if err != nil {
		logger.Get().Error(err)
		return errors.New("Admin not found")
	}
	logger.Get().Debug(a)
	if a.Id != "" {
		return nil
	}
	return errors.New("Password did not match")
}

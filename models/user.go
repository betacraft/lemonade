package models

import (
	"errors"
	"github.com/mholt/binding"
	"github.com/nu7hatch/gouuid"
	"github.com/rainingclouds/lemonade/db"
	"github.com/rainingclouds/lemonade/framework"
	"github.com/rainingclouds/lemonade/logger"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//go:generate easytags user.go bson
//go:generate easytags user.go json

type Address struct {
	Address   string  `json:"address" bson:"address"`
	Locality  string  `json:"locality" bson:"locality"`
	City      string  `json:"city" bson:"city"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
	ZipCode   string  `json:"zip_code" bson:"zip_code"`
}

func (a Address) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&a.Address:   "address",
		&a.City:      "city",
		&a.Latitude:  "latitude",
		&a.Locality:  "locality",
		&a.Longitude: "longitude",
		&a.ZipCode:   "zip_code",
	}
}

type User struct {
	Id      bson.ObjectId `bson:"_id,omitempty" json:"id"`
	AuthKey string        `bson:"auth_key" json:"auth_key"`

	Name              string `bson:"name" json:"name"`
	MobileNumber      string `bson:"mobile_number" json:"mobile"`
	MobileNumberValid bool   `bson:"m_valid" json:"m_valid"`
	Email             string `bson:"email" json:"email"`
	Password          string `bson:"password" json:"-"`
	ProfilePicLink    string `bson:"profile_pic" json:"profile_pic"`

	CurrentMobileHandset string `bson:"current_handset" json:"current_handset"`
	LookingForHandset    string `bson:"looking_for_handset" json:"looking_for_handset"`

	OtpCode string `bson:"otp_code" json:"-"`

	Address Address `bson:"address" json:"address"`

	IsConnectedWithFacebook   bool   `bson:"is_fb" json:"is_fb"`
	FacebookAccessToken       string `bson:"fb_token" json:"fb_token"`
	IsConnectedWithGooglePlus bool   `bson:"is_gplus" json:"is_gplus"`

	IsAccessEnabled bool   `bson:"is_access_enabled" json:"is_access_enabled"`
	Reason          string `bson:"reason" json:"reason"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at" json:"deleted_at"`
}

func (u *User) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&u.AuthKey:                   "auth_key",
		&u.Name:                      "name",
		&u.MobileNumber:              "mobile",
		&u.Email:                     "email",
		&u.Password:                  "password",
		&u.ProfilePicLink:            "profile_pic",
		&u.IsConnectedWithFacebook:   "is_fb",
		&u.FacebookAccessToken:       "fb_token",
		&u.IsConnectedWithGooglePlus: "is_gplus",
	}
}

// this function is there because if json="-" for some field then somehow
// bindings doesnt map that field at all even if the FieldMap function
// has a mapping avilable for that field
// we will have to write an alternative for that
func CreateUser(userMap map[string]interface{}) (*User, error) {
	user := new(User)
	var ok bool
	user.Email, ok = userMap["email"].(string)
	if !ok {
		return nil, errors.New("Email is not present")
	}
	user.MobileNumber, ok = userMap["mobile"].(string)
	if !ok {
		return nil, errors.New("Mobile is not present")
	}
	user.Name, ok = userMap["name"].(string)
	if !ok {
		return nil, errors.New("Name is not present")
	}
	user.Password, ok = userMap["password"].(string)
	if !ok {
		return nil, errors.New("Password is not present")
	}
	user.CurrentMobileHandset, ok = userMap["current_handset"].(string)
	if !ok {
		return nil, errors.New("Current Handset is not present")
	}
	user.LookingForHandset, ok = userMap["looking_for_handset"].(string)
	if !ok {
		return nil, errors.New("Looking for Handset is not present")
	}
	user.Address = Address{}
	user.Address.City, ok = userMap["city"].(string)
	if !ok {
		return nil, errors.New("City is not present")
	}
	user.Address.Locality, ok = userMap["locality"].(string)
	if !ok {
		return nil, errors.New("Locality is not present")
	}
	user.Address.Address, ok = userMap["address"].(string)
	if !ok {
		return nil, errors.New("Address is not present")
	}
	user.Address.ZipCode, ok = userMap["zip_code"].(string)
	if !ok {
		return nil, errors.New("Zip Code is not present")
	}
	return user, nil
}

func CreateFacebookUser(userMap map[string]interface{}) (*User, error) {
	user := new(User)
	var ok bool
	user.Email, ok = userMap["email"].(string)
	if !ok {
		return nil, errors.New("Email is not present")
	}
	user.Name, ok = userMap["name"].(string)
	if !ok {
		return nil, errors.New("Name is not present")
	}
	user.ProfilePicLink, ok = userMap["profile_pic"].(string)
	if !ok {
		return nil, errors.New("Profile pic is not present")
	}
	user.IsConnectedWithFacebook, ok = userMap["is_fb"].(bool)
	if !ok {
		return nil, errors.New("Illegal facebook login")
	}
	user.FacebookAccessToken, ok = userMap["fb_token"].(string)
	if !ok && user.IsConnectedWithFacebook {
		return nil, errors.New("Facebook token is not present")
	}
	return user, nil
}

func CreateGplusUser(userMap map[string]interface{}) (*User, error) {
	user := new(User)
	var ok bool
	user.Email, ok = userMap["email"].(string)
	if !ok {
		return nil, errors.New("Email is not present")
	}
	user.Name, ok = userMap["name"].(string)
	if !ok {
		return nil, errors.New("Name is not present")
	}
	user.ProfilePicLink, ok = userMap["profile_pic"].(string)
	if !ok {
		return nil, errors.New("Profile pic is not present")
	}
	user.IsConnectedWithGooglePlus, ok = userMap["is_gplus"].(bool)
	if !ok {
		return nil, errors.New("Illegal Google plus login")
	}
	return user, nil
}

func (u *User) Create() error {
	user := new(User)
	err := db.MgFindOneStrong(C_USERS, &bson.M{"email": u.Email}, user)
	if err != nil && err.Error() != "not found" {
		logger.Get().Error("While finding existing user : " + err.Error())
		return err
	}
	if user.Id != "" {
		if u.IsConnectedWithFacebook {
			user.IsConnectedWithFacebook = true
		}
		if u.IsConnectedWithGooglePlus {
			user.IsConnectedWithGooglePlus = true
		}
		err = user.Save()
		if err != nil {
			return err
		}
		*u = *user
		return nil
	}
	authKey, err := uuid.NewV4()
	if err != nil {
		logger.Get().Error(err)
		return err
	}
	logger.Get().Debug(u.Password)
	u.AuthKey = authKey.String()
	if !u.IsConnectedWithGooglePlus && u.IsConnectedWithFacebook {
		u.Password = framework.MD5Hash(u.Password)
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return db.MgCreateStrong(C_USERS, u)
}

func (u *User) Save() error {
	u.UpdatedAt = time.Now()
	logger.Debug("Saving user with id ", u.Id)
	return db.MgUpdateStrong(C_USERS, u.Id, u)
}

func GetUserByAuthKey(authKey string) (*User, error) {
	u := new(User)
	err := db.MgFindOneStrong(C_USERS, &bson.M{"auth_key": authKey}, u)
	return u, err
}

func (u *User) UpdatePassword(password string) error {
	hashedPassword := framework.MD5Hash(password)
	u.Password = hashedPassword
	return u.Save()
}

func (u *User) AuthenticateUser() error {
	logger.Get().Debug(u.Password)
	hashedPassword := framework.MD5Hash(u.Password)
	logger.Get().Debug("Hashed password is ", hashedPassword)
	err := db.MgFindOneStrong(C_USERS, &bson.M{"email": u.Email, "password": hashedPassword}, u)
	if err != nil {
		logger.Get().Error(err)
		return errors.New("User not found")
	}
	logger.Get().Debug(u)
	if u.Id != "" {
		return nil
	}
	return errors.New("Password did not match")
}

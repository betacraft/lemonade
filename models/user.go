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
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	AuthKey    string        `bson:"auth_key" json:"auth_key"`
	SessionKey string        `bson:"session_key" json:"session_key"`

	Name              string `bson:"name" json:"name"`
	MobileNumber      string `bson:"mobile_number" json:"mobile"`
	MobileNumberValid bool   `bson:"m_valid" json:"m_valid"`
	Email             string `bson:"email" json:"email"`
	IsEmailConfirmed  bool   `bson:"email_confirmed" json:"email_confirmed"`
	Password          string `bson:"password" json:"-"`
	ProfilePicLink    string `bson:"profile_pic" json:"profile_pic"`

	OtpCode string `bson:"otp_code" json:"-"`

	Address Address `bson:"address" json:"address"`

	Gender string `json:"-" bson:"gender"`

	IsConnectedWithFacebook bool   `bson:"is_fb" json:"is_fb"`
	FacebookAccessToken     string `bson:"fb_token" json:"-"`
	FacebookUserId          string `bson:"fb_user_id" json:"-"`

	IsConnectedWithGooglePlus bool   `bson:"is_gplus" json:"is_gplus"`
	GPlusAccessToken          string `bson:"gplus_token" json:"-"`
	GPlusCode                 string `bson:"gplus_code" json:"-"`
	GPlusUserId               string `bson:"gplus_user_id" json:"-"`

	CreatedGroupCount int `bson:"created_group_count" json:"created_group_count"`
	JoinedGroupCount  int `bson:"joined_group_count" json:"joined_group_count"`

	IsAccessEnabled bool   `bson:"is_access_enabled" json:"is_access_enabled"`
	Reason          string `bson:"reason" json:"reason"`

	CreatedGroupIds []bson.ObjectId `bson:"created_group_ids,omitempty" json:"created_group_ids"`
	JoinedGroupIds  []bson.ObjectId `bson:"joined_group_ids,omitempty" json:"joined_group_ids"`

	CreatedGroups *[]Group `bson:"-" json:"created_groups"`
	JoinedGroups  *[]Group `bson:"-" json:"joined_groups"`

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
	user.Password, ok = userMap["password"].(string)
	if !ok {
		return nil, errors.New("Password is not present")
	}
	user.MobileNumber, ok = userMap["mobile_number"].(string)
	if !ok {
		return nil, errors.New("Mobile number is not present")
	}
	user.IsConnectedWithGooglePlus = false
	user.IsConnectedWithFacebook = false
	return user, nil
}

func ParseFacebookUser(userMap map[string]interface{}) (*User, error) {
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
	user.Gender, ok = userMap["gender"].(string)
	if !ok {
		return nil, errors.New("Illegal Facebook login request")
	}
	user.IsConnectedWithFacebook, ok = userMap["is_fb"].(bool)
	if !ok {
		return nil, errors.New("Illegal Facebook login request")
	}
	user.FacebookUserId, ok = userMap["fb_user_id"].(string)
	if !ok {
		return nil, errors.New("Illegal Facebook login request")
	}
	user.FacebookAccessToken, ok = userMap["fb_token"].(string)
	if !ok && user.IsConnectedWithFacebook {
		return nil, errors.New("Illegal Facebook login request")
	}
	return user, nil
}

func ParseGplusUser(userMap map[string]interface{}) (*User, error) {
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
	user.Gender, ok = userMap["gender"].(string)
	if !ok {
		return nil, errors.New("Illegal Google Plus login request")
	}
	user.ProfilePicLink, ok = userMap["profile_pic"].(string)
	if !ok {
		return nil, errors.New("Illegal Google Plus login request")
	}
	user.GPlusUserId, ok = userMap["gplus_user_id"].(string)
	if !ok {
		return nil, errors.New("Illegal Google Plus login request")
	}
	user.GPlusAccessToken, ok = userMap["gplus_token"].(string)
	if !ok && user.IsConnectedWithFacebook {
		return nil, errors.New("Illegal Google Plus login request")
	}
	user.IsConnectedWithGooglePlus, ok = userMap["is_gplus"].(bool)
	if !ok {
		return nil, errors.New("Illegal Google Plus login request")
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
			user.FacebookAccessToken = u.FacebookAccessToken
			user.FacebookUserId = u.FacebookUserId
		}
		if u.IsConnectedWithGooglePlus {
			user.IsConnectedWithGooglePlus = true
			user.GPlusCode = u.GPlusCode
			user.GPlusAccessToken = u.GPlusAccessToken
			user.GPlusUserId = u.GPlusUserId
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
	if !u.IsConnectedWithGooglePlus && !u.IsConnectedWithFacebook {
		u.Password = framework.MD5Hash(u.Password)
	}
	sessionKey, err := uuid.NewV4()
	if err != nil {
		logger.Get().Error(err)
		return err
	}
	u.SessionKey = sessionKey.String()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return db.MgCreateStrong(C_USERS, u)
}

func (u *User) Save() error {
	u.UpdatedAt = time.Now()
	logger.Debug("Saving user with id ", u.Id)
	return db.MgUpdateStrong(C_USERS, u.Id, u)
}

func GetUserByEmailId(email string) (*User, error) {
	user := new(User)
	err := db.MgFindOneStrong(C_USERS, &bson.M{"email": email}, user)
	return user, err
}

func GetUserByFacebookUserId(fbUserId string) (*User, error) {
	user := new(User)
	err := db.MgFindOneStrong(C_USERS, &bson.M{"fb_user_id": fbUserId}, user)
	return user, err
}

func GetUserByGooglePlusUserId(gPlusUserId string) (*User, error) {
	user := new(User)
	err := db.MgFindOneStrong(C_USERS, &bson.M{"gplus_user_id": gPlusUserId}, user)
	return user, err
}

func GetUserByAuthKey(authKey string) (*User, error) {
	u := new(User)
	err := db.MgFindOneStrong(C_USERS, &bson.M{"auth_key": authKey}, u)
	return u, err
}

func GetUserBySessionKey(sessionKey string) (*User, error) {
	u := new(User)
	err := db.MgFindOneStrong(C_USERS, &bson.M{"session_key": sessionKey}, u)
	return u, err
}

func (u *User) UpdatePassword(password string) error {
	hashedPassword := framework.MD5Hash(password)
	u.Password = hashedPassword
	return u.Save()
}

func (u *User) SocialLogin() error {
	sessionKey, err := uuid.NewV4()
	if err != nil {
		logger.Get().Error(err)
		return err
	}
	u.SessionKey = sessionKey.String()
	return u.Save()
}

func (u *User) AuthenticateUser() error {
	logger.Get().Debug(u.Password)
	hashedPassword := framework.MD5Hash(u.Password)
	logger.Get().Debug("Hashed password is ", hashedPassword)
	err := db.MgFindOneStrong(C_USERS, &bson.M{"email": u.Email, "password": hashedPassword}, u)
	if err != nil {
		logger.Get().Error(err)
		return errors.New("Wrong Email/Password")
	}
	logger.Get().Debug(u)
	if u.Id != "" {
		sessionKey, err := uuid.NewV4()
		if err != nil {
			logger.Get().Error(err)
			return err
		}
		u.SessionKey = sessionKey.String()
		return u.Save()
	}
	return errors.New("Password did not match")
}

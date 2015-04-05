package controllers

import (
	"errors"
	"fmt"
	"github.com/go-zoo/bone"
	"github.com/jordan-wright/email"
	"github.com/rainingclouds/lemonade/framework"
	"github.com/rainingclouds/lemonade/logger"
	"github.com/rainingclouds/lemonade/mailer"
	"github.com/rainingclouds/lemonade/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// api to register a new user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	bodyMap, err := framework.ReadBody(r)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	// filling up the user struct
	user, err := models.CreateUser(bodyMap)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	user.OtpCode = framework.GenerateOtp()
	user.IsAccessEnabled = false
	// creating user
	err = user.Create()
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	// notifying to the admins
	go func(user *models.User) {
		u, err := models.GetUserByAuthKey(user.AuthKey)
		if err != nil {
			u = user
		}
		mail := email.NewEmail()
		mail.From = "lemonades@rainingclouds.com"
		mail.Subject = fmt.Sprintf("%v requested access to lemonades", u.Name)
		mail.Text = []byte(fmt.Sprintf("%v is registered, please call him and update the status to akshay@rainingclouds.com", u))
		mailer.SendToMany([]string{"amit@rainingclouds.com", "akshay@rainingclouds.com"}, mail)
	}(user)
	// send otp
	// go user.SendOtp()
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{"success": true,
		"user":    user,
		"message": "User is successfully registered",
	})
}

func RegisterUserWithFacebook(w http.ResponseWriter, r *http.Request) {
	bodyMap, err := framework.ReadBody(r)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	user, err := models.CreateFacebookUser(bodyMap)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	user.OtpCode = framework.GenerateOtp()
	// creating user
	err = user.Create()
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{"success": true,
		"user":    user,
		"message": "User is successfully registered",
	})
}

func RegisterUserWithGplus(w http.ResponseWriter, r *http.Request) {
	bodyMap, err := framework.ReadBody(r)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	user, err := models.CreateGplusUser(bodyMap)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	user.OtpCode = framework.GenerateOtp()
	// creating user
	err = user.Create()
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{"success": true,
		"user":    user,
		"message": "User is successfully registered",
	})
}

func UpdatePhoneNumber(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	bodyMap, err := framework.ReadBody(r.Request)
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusBadRequest, err)
		return
	}
	value, ok := bodyMap["phone_number"].(string)
	if !ok {
		framework.WriteError(w, r.Request, http.StatusBadRequest, errors.New("Phone number is not present in the body"))
		return
	}
	if user.MobileNumber == value && user.MobileNumberValid {
		framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
			"success":  true,
			"is_valid": true,
			"message":  "User is successfully updated",
		})
		return
	}
	user.MobileNumber = value
	err = user.Save()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusBadRequest, err)
		return
	}
	// go user.SendOtp()
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success":  true,
		"is_valid": false,
		"message":  "User is successfully updated",
	})
}

func UpdateUser(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	bodyMap, err := framework.ReadBody(r.Request)
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusBadRequest, err)
		return
	}
	value, ok := bodyMap["password"].(string)
	if ok {
		err = user.UpdatePassword(value)
		if err != nil {
			framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
			return
		}
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"message": "User is successfully updated",
	})

}

func ValidateOtp(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	otp := r.URL.Query().Get("otp")
	logger.Debug("Got ", otp, " from request")
	if otp != user.OtpCode {
		framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
			"success": false,
			"message": "Wrong OTP",
		})
		return
	}
	user.MobileNumberValid = true
	err := user.Save()
	if err != nil {
		logger.Debug("While saving user")
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"message": "Phone number is validated",
	})
}

// api to authenticate a user
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	body, err := framework.ReadBody(r)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	var ok bool
	user := models.User{}
	user.Email, ok = body["email"].(string)
	if !ok {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("Email is not present"))
		return
	}
	user.Password, ok = body["password"].(string)
	if !ok {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("Password is not present"))
		return
	}
	err = user.AuthenticateUser()
	if err != nil {
		framework.WriteError(w, r, http.StatusUnauthorized, err)
		return
	}
	if !user.IsAccessEnabled {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("Your access is still not approved"))
		return
	}
	framework.WriteResponse(w, http.StatusOK,
		framework.JSONResponse{"success": true,
			"user":    user,
			"message": "User is successfully authenticated",
		})
}

func UserLogout(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	user.SessionKey = ""
	err := user.Save()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"message": "User is successfully logged out",
	})
}

func AttachDeal(w http.ResponseWriter, r *framework.Request) {
	authKey := r.Request.Header.Get("User-Auth-Key")
	if authKey == "" {
		framework.WriteError(w, r.Request, http.StatusBadRequest, errors.New("Illegal request"))
		return
	}
	user, err := models.GetUserByAuthKey(authKey)
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusBadRequest, err)
		return
	}
	dealIdHex := bone.GetValue(r.Request, "deal_id")
	dealId := bson.ObjectIdHex(dealIdHex)
	if !dealId.Valid() {
		framework.WriteError(w, r.Request, http.StatusBadRequest, errors.New("Illegal request"))
		return
	}
	deal, err := models.GetPhoneDealById(dealId)
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusBadRequest, err)
		return
	}
	user.DealId = dealId
	deal.CurrentPeopleCount = deal.CurrentPeopleCount + 1
	err = deal.Save()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	err = user.Save()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"message": "User is attached to the deal",
	})
}

func GetUser(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"user":    user,
	})
}

package controllers

import (
	"errors"
	"fmt"
	"github.com/go-zoo/bone"
	"github.com/jordan-wright/email"
	"github.com/rainingclouds/lemonades/captcha"
	"github.com/rainingclouds/lemonades/framework"
	"github.com/rainingclouds/lemonades/logger"
	"github.com/rainingclouds/lemonades/mailer"
	"github.com/rainingclouds/lemonades/models"
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
	response, ok := bodyMap["captcha"].(string)
	if !ok {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("Illegal Request"))
		return
	}
	if !captcha.Validate(response) {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("Illegal Request"))
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
		mail.From = "groupup@rainingclouds.com"
		mail.Subject = "Please confirm your email address"
		mail.Text = []byte("Hello from GroupUP,\nPlease click on the following link to confirm your email address\nhttp://www.groupup.in/#!/user/" + u.AuthKey + "/confirm_email")
		mailer.Send(u.Email, mail)
		mail.Subject = fmt.Sprintf("%v requested access to groupup", u.Name)
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

func LoginWithFacebook(w http.ResponseWriter, r *http.Request) {
	bodyMap, err := framework.ReadBody(r)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	user, err := models.ParseFacebookUser(bodyMap)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	// check if he is already registered with facebook
	fbUser, err := models.GetUserByFacebookUserId(user.FacebookUserId)
	if err != nil && err.Error() != "not found" {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	if fbUser.Id.Hex() != "" {
		fbUser.FacebookAccessToken = user.FacebookAccessToken
		fbUser.IsEmailConfirmed = true
		err = fbUser.Save()
		if err != nil {
			framework.WriteError(w, r, http.StatusInternalServerError, err)
			return
		}
		err = fbUser.SocialLogin()
		if err != nil {
			framework.WriteError(w, r, http.StatusInternalServerError, err)
			return
		}
		framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{"success": true,
			"user":    fbUser,
			"message": "User is successfully logged in",
		})
		return
	}
	// creating user
	user.IsEmailConfirmed = true
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

func LoginWithGooglePlus(w http.ResponseWriter, r *http.Request) {
	bodyMap, err := framework.ReadBody(r)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	user, err := models.ParseGplusUser(bodyMap)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	gPlusUser, err := models.GetUserByGooglePlusUserId(user.GPlusUserId)
	if err != nil && err.Error() != "not found" {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	if gPlusUser.Id.Hex() != "" {
		gPlusUser.GPlusAccessToken = user.GPlusAccessToken
		gPlusUser.IsEmailConfirmed = true
		err = gPlusUser.Save()
		if err != nil {
			framework.WriteError(w, r, http.StatusInternalServerError, err)
			return
		}
		err = gPlusUser.SocialLogin()
		if err != nil {
			framework.WriteError(w, r, http.StatusInternalServerError, err)
			return
		}
		framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{"success": true,
			"user":    gPlusUser,
			"message": "User is successfully logged in",
		})
		return
	}
	user.IsEmailConfirmed = true
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

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	bodyMap, err := framework.ReadBody(r)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	emailId, ok := bodyMap["email"].(string)
	if !ok {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("Email is not present"))
		return
	}
	user, err := models.GetUserByEmailId(emailId)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("This email id is not registered"))
		return
	}
	go func(user *models.User) {
		mail := email.NewEmail()
		mail.From = "groupup@rainingclouds.com"
		mail.Subject = "Password reset instructions"
		mail.Text = []byte("Hello User,\nClick on the link to reset your password:\nhttp://www.groupup.in/#!/user/" + user.AuthKey + "/reset_password")
		mailer.Send(user.Email, mail)
	}(user)
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"message": "Reset password instructions are mailed on the registered email Address",
	})
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	authKey := bone.GetValue(r, "auth_key")
	if authKey == "" {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("Illegal request"))
		return
	}
	user, err := models.GetUserByAuthKey(authKey)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	bodyMap, err := framework.ReadBody(r)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	password, ok := bodyMap["password"].(string)
	if !ok {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("Illegal request"))
		return
	}
	err = user.UpdatePassword(password)
	if err != nil {
		framework.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}
	err = user.RenewSessionKey()
	if err != nil {
		framework.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"user":    user,
		"message": "Password is updated successfully",
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
	framework.WriteResponse(w, http.StatusOK,
		framework.JSONResponse{"success": true,
			"user":    user,
			"message": "User is successfully authenticated",
		})
}

func ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetUserByAuthKey(bone.GetValue(r, "auth_key"))
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("You are not yet registered with us"))
		return
	}
	user.IsEmailConfirmed = true
	err = user.Save()
	if err != nil {
		framework.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"message": "Your email address is confirmed successfully",
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

func GetUser(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"user":    user,
	})
}

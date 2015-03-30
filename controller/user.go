package controllers

import (
	"errors"
	"fmt"
	"github.com/go-zoo/bone"
	"net/http"
	"strings"
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
	// creating user
	err = user.Create()
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	// send otp
	go user.SendOtp()
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
	go user.SendOtp()
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

func GetUploadUrlForProfilePic(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	fileName := bone.GetValue(r.Request, "filename")
	expiryTime := aws.GetExpiryTime()
	components := strings.Split(fileName, ".")
	if len(components) < 2 {
		framework.WriteError(w, r.Request, http.StatusBadRequest, errors.New("Illegal file"))
		return
	}
	extenstion := components[1]
	path := fmt.Sprintf("%s/profile_pic", user.Id.Hex())
	signedUrl := aws.Bucket().UploadSignedURL(path, "PUT",
		fmt.Sprintf("image/%s", extenstion),
		expiryTime)
	signedUrl = strings.Replace(signedUrl, "s3", "s3-ap-southeast-1", 1)
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"url":     signedUrl,
	})
}

func ProfilePicUploaded(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	path := fmt.Sprintf("%s/profile_pic", user.Id.Hex())
	user.ProfilePicLink = aws.Bucket().URL(path)
	err := user.Save()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"message": "Profile pic is updated successfully",
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
	// send otp
	go user.SendOtp()
	framework.WriteResponse(w, http.StatusOK,
		framework.JSONResponse{"success": true,
			"user":    user,
			"message": "User is successfully authenticated",
		})
}

func GetFavoriteRestaurants(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	var restaurant *models.Restaurant
	var err error
	for i := 0; i < len(user.FavoriteRestaurantsIds); i++ {
		restaurant, err = models.GetRestaurantById(user.FavoriteRestaurantsIds[i])
		if err != nil {
			framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
			return
		}
		user.FavoriteRestaurants = append(user.FavoriteRestaurants, *restaurant)
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success":         true,
		"fav_restaurants": user.FavoriteRestaurants,
	})
}

func GetUser(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"user":    user,
	})
}

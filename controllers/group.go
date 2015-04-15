package controllers

import (
	"errors"
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonades/framework"
	"github.com/rainingclouds/lemonades/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

func JoinGroup(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	idString := bone.GetValue(r.Request, "id")
	if idString == "" {
		framework.WriteError(w, r.Request, http.StatusBadRequest, errors.New("Illegal group id"))
		return
	}
	id := bson.ObjectIdHex(idString)
	if !id.Valid() {
		framework.WriteError(w, r.Request, http.StatusBadRequest, errors.New("Illegal group id"))
	}
	group, err := models.GetGroupById(id)
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	group.InterestedUsersCount = group.InterestedUsersCount + 1
	group.InterestedUsers = append(group.InterestedUsers, user.Id)
	err = group.Update()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	user.JoinedGroupCount = user.JoinedGroupCount + 1
	err = user.Save()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	group.IsJoined = true
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"group":   group,
	})
}

func GetGroups(w http.ResponseWriter, r *http.Request) {
	pageNo := 0
	pageNo, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		pageNo = 0
	}
	groups, err := models.GetGroups(pageNo)
	if err != nil {
		framework.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"groups":  groups,
	})
}

func GetGroup(w http.ResponseWriter, r *http.Request) {
	idString := bone.GetValue(r, "id")
	if idString == "" {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("Illegal group id"))
		return
	}
	id := bson.ObjectIdHex(idString)
	if !id.Valid() {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("Illegal group id"))
	}
	group, err := models.GetGroupById(id)
	if err != nil {
		framework.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}
	sessionKey := r.Header.Get("Session-Key")
	if sessionKey != "" {
		user, err := models.GetUserBySessionKey(sessionKey)
		if err != nil {
			framework.WriteError(w, r, http.StatusBadRequest, err)
			return
		}
		for _, id := range group.InterestedUsers {
			if user.Id.Hex() == id.Hex() {
				group.IsJoined = true
			}
		}
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"group":   group,
	})
}

func CreateGroup(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	bodyMap, err := framework.ReadBody(r.Request)
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusBadRequest, err)
		return
	}
	val, ok := bodyMap["product_link"].(string)
	if !ok {
		framework.WriteError(w, r.Request, http.StatusBadRequest, errors.New("Product link is not present"))
		return
	}
	product, err := models.FetchProductInfo(val)
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusBadRequest, err)
		return
	}
	if product.PriceValue < 10000 {
		framework.WriteError(w, r.Request, http.StatusBadRequest, errors.New("Price of the product must be more than 10,000 Rs"))
		return
	}
	has := false
	for _, id := range product.AddedBy {
		if user.Id.Hex() == id.Hex() {
			has = true
		}
	}
	if !has {
		product.AddedBy = append(product.AddedBy, user.Id)
	}
	err = product.CreateOrUpdate()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	group, err := models.GetGroupByProductId(product.Id)
	if err != nil && err.Error() != "not found" {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	if group.Id.Hex() != "" {
		framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
			"success": true,
			"group":   group,
		})
	}
	group = &models.Group{}
	group.Id = bson.NewObjectId()
	group.CreatedBy = user.Id
	group.Product = *product
	group.ExpiresOn = time.Now().Add(time.Hour * 24 * 30)
	group.IsOn = true
	group.InterestedUsers = append(group.InterestedUsers, user.Id)
	group.InterestedUsersCount = 1
	err = group.Create()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	user.CreatedGroupCount = user.CreatedGroupCount + 1
	err = user.Save()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"group":   group,
	})
}

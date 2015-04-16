package controllers

import (
	"errors"
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonades/framework"
	"github.com/rainingclouds/lemonades/logger"
	"github.com/rainingclouds/lemonades/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

func GetUserJoinedGroups(w http.ResponseWriter, r *framework.Request) {
	pageNo := 0
	pageNo, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		pageNo = 0
	}
	user := r.MustGet("user").(*models.User)
	if len(user.JoinedGroupIds) < pageNo*9 {
		framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
			"success":       true,
			"joined_groups": nil,
		})
	}
	var group *models.Group
	user.JoinedGroups = new([]models.Group)
	for _, id := range user.JoinedGroupIds[pageNo*9:] {
		group, err = models.GetGroupById(id)
		if err != nil {
			logger.Err("While getting joined groups by user id", user.Id, err)
			continue
		}
		*user.JoinedGroups = append(*(user.JoinedGroups), *group)
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success":       true,
		"joined_groups": user.JoinedGroups,
	})

}

func GetUserCreatedGroups(w http.ResponseWriter, r *framework.Request) {
	pageNo := 0
	pageNo, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		pageNo = 0
	}
	user := r.MustGet("user").(*models.User)
	if len(user.CreatedGroupIds) < pageNo*9 {
		framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
			"success":        true,
			"created_groups": nil,
		})
	}
	user.CreatedGroups = new([]models.Group)
	var group *models.Group
	for _, id := range user.CreatedGroupIds[pageNo*9:] {
		group, err = models.GetGroupById(id)
		if err != nil {
			logger.Err("While getting created groups by user id", user.Id, err)
			continue
		}
		*user.CreatedGroups = append(*(user.CreatedGroups), *group)
	}

	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success":        true,
		"created_groups": user.CreatedGroups,
	})

}

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
	user.JoinedGroupIds = append(user.JoinedGroupIds, group.Id)
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

	if group.ExpiresOn.After(time.Now()) {
		group.ExpiresIn = int64(group.ExpiresOn.Sub(time.Now()).Hours() / 24)
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
	if product.PriceValue == 0 {
		logger.Err("Could not parse ", val)
		framework.WriteError(w, r.Request, http.StatusBadRequest, errors.New("We could not parse the page, please retry"))
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
	if product.PriceValue < 5000 {
		group.RequiredUserCount = 30
	}
	if product.PriceValue > 5000 && product.PriceValue < 10000 {
		group.RequiredUserCount = 20
	}
	if product.PriceValue > 10000 && product.PriceValue < 25000 {
		group.RequiredUserCount = 10
	}
	if product.PriceValue > 25000 && product.PriceValue < 75000 {
		group.RequiredUserCount = 5
	}
	if product.PriceValue > 75000 {
		group.RequiredUserCount = 3
	}
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
	user.CreatedGroupIds = append(user.CreatedGroupIds, group.Id)
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

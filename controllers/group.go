package controllers

import (
	"errors"
	"fmt"
	"github.com/go-zoo/bone"
	"github.com/jordan-wright/email"
	"github.com/rainingclouds/lemonades/framework"
	"github.com/rainingclouds/lemonades/logger"
	"github.com/rainingclouds/lemonades/mailer"
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

func ShareGroup(w http.ResponseWriter, r *http.Request) {
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
	var groupText string
	switch group.InterestedUsersCount {
	case 0:
		groupText = fmt.Sprintf("Join with a group of people to buy %v with a huge group buying discount on GroupUP.in", group.Product.Name)
	case 1:
		groupText = fmt.Sprintf("1 person is interested in buying %v. Join him on GroupUP.in and get huge group discount.", group.Product.Name)
	default:
		groupText = fmt.Sprintf("%v people are interested in buying %v. Join them on GroupUP.in and get huge group discount.", group.InterestedUsersCount, group.Product.Name)
	}
	framework.WriteText(w, fmt.Sprintf("<!DOCTYPE html><html><head><meta property=\"og:type\" content=\"website\"><link rel=\"canonical\" href=\"http://www.groupup.in/#!/group/%v\"/><meta property=\"og:url\" content=\"http://www.groupup.in/group/%v/share/%v\"><meta property=\"og:url:width\" content=\"300\"><meta property=\"og:url:height\" content=\"300\"><meta property=\"og:title\" content=\"Buy %v with me on GroupUP.in\"><meta property=\"og:image\" content=\"%v\"><meta property=\"og:description\" content=\"%v\"><meta property=\"fb:app_id\" content=\"1608020712745966\"></head><body></body></html>", group.Id, group.Id, group.InterestedUsersCount, group.Product.Name, group.Product.ProductImage, groupText))
}

func LeaveGroup(w http.ResponseWriter, r *framework.Request) {
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
	group.InterestedUsersCount = group.InterestedUsersCount - 1
	// removing that user from interested count
	for i := 0; i < len(group.InterestedUsers); i++ {
		if group.InterestedUsers[i].Hex() == user.Id.Hex() {
			group.InterestedUsers = append(group.InterestedUsers[:i], group.InterestedUsers[i+1:]...)
			break
		}
	}
	err = group.Update()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	user.JoinedGroupCount = user.JoinedGroupCount - 1
	for i := 0; i < len(user.JoinedGroupIds); i++ {
		if user.JoinedGroupIds[i].Hex() == group.Id.Hex() {
			user.JoinedGroupIds = append(user.JoinedGroupIds[:i], user.JoinedGroupIds[i+1:]...)
			break
		}
	}
	err = user.Save()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	if group.ExpiresOn.After(time.Now()) {
		group.ExpiresIn = int64(group.ExpiresOn.Sub(time.Now()).Hours() / 24)
	}
	group.IsJoined = false
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"group":   group,
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
	for _, userId := range group.InterestedUsers {
		if userId.Hex() == user.Id.Hex() {
			group.IsJoined = true
			framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
				"success": true,
				"group":   group,
				"message": "Successfully Joined the group",
			})
			return
		}
	}
	group.InterestedUsersCount = group.InterestedUsersCount + 1
	group.InterestedUsers = append(group.InterestedUsers, user.Id)
	if group.InterestedUsersCount >= group.RequiredUserCount {
		group.ReachedGoalOn = time.Now()
		group.ReachedGoal = int64(time.Now().Sub(group.ReachedGoalOn).Hours() / 24)
		// send email notification to akshay
		go func(group *models.Group) {
			mail := email.NewEmail()
			mail.From = "groupup@rainingclouds.com"
			mail.Subject = "Start getting deal for " + group.Product.Name
			mail.Text = []byte("Group " + group.Id.Hex() + " is done with its expected users\n" + fmt.Sprintf("%v", group))
			mailer.Send("akshay@rainingclouds.com", mail)
		}(group)
	}
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
	if group.ExpiresOn.After(time.Now()) {
		group.ExpiresIn = int64(group.ExpiresOn.Sub(time.Now()).Hours() / 24)
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
	searchTerm := r.URL.Query().Get("search")
	if searchTerm == "" {
		groups, err := models.GetGroups(pageNo)
		if err != nil {
			framework.WriteError(w, r, http.StatusInternalServerError, err)
			return
		}
		framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
			"success": true,
			"groups":  groups,
		})
		return
	}
	groups, err := models.GetGroupsForSearchTerm(pageNo, searchTerm)
	if err != nil {
		framework.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"groups":  groups,
	})
	return
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
		if err == nil {
			for _, id := range group.InterestedUsers {
				if user.Id.Hex() == id.Hex() {
					group.IsJoined = true
				}
			}
		}
	}

	if group.ExpiresOn.After(time.Now()) {
		group.ExpiresIn = int64(group.ExpiresOn.Sub(time.Now()).Hours() / 24)
	}
	if group.InterestedUsersCount >= group.RequiredUserCount {
		group.ReachedGoal = int64(time.Now().Sub(group.ReachedGoalOn).Hours() / 24)
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
	if product.PriceValue < 999 {
		framework.WriteError(w, r.Request, http.StatusBadRequest, errors.New("Minimum cost of the product should be more than 999 Rs"))
		return
	}
	group, err := models.GetGroupByProduct(product)
	if err != nil && err.Error() != "not found" {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	if group.Id.Hex() != "" {
		framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
			"success": true,
			"group":   group,
		})
		return
	}
	logger.Debug("Creating new group")
	group = &models.Group{}
	if product.PriceValue < 5000 {
		group.RequiredUserCount = 20
	}
	if product.PriceValue > 5000 && product.PriceValue < 10000 {
		group.RequiredUserCount = 12
	}
	if product.PriceValue > 10000 && product.PriceValue < 25000 {
		group.RequiredUserCount = 7
	}
	if product.PriceValue > 25000 && product.PriceValue < 75000 {
		group.RequiredUserCount = 5
	}
	if product.PriceValue > 75000 {
		group.RequiredUserCount = 3
	}
	group.Id = bson.NewObjectId()
	group.MinDiscount = "10% Off"
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
	go func(group *models.Group, user *models.User) {
		mail := email.NewEmail()
		mail.From = "groupup@rainingclouds.com"
		mail.Subject = "Group is created for " + group.Product.Name
		mail.Text = []byte("New group is created by " + user.Name + " Email " + user.Email + " Phone number " + user.MobileNumber + " Group " + fmt.Sprintf("%v", group))
		mailer.Send("akshay@rainingclouds.com", mail)
	}(group, user)
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"group":   group,
	})
}

func AdminCreateGroup(w http.ResponseWriter, r *framework.Request) {
	admin := r.MustGet("admin").(*models.Admin)
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
		if admin.Id.Hex() == id.Hex() {
			has = true
		}
	}
	if !has {
		product.AddedBy = append(product.AddedBy, admin.Id)
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
		return
	}
	logger.Debug("Creating new group")
	group = &models.Group{}
	if product.PriceValue < 5000 {
		group.RequiredUserCount = 20
	}
	if product.PriceValue > 5000 && product.PriceValue < 10000 {
		group.RequiredUserCount = 12
	}
	if product.PriceValue > 10000 && product.PriceValue < 25000 {
		group.RequiredUserCount = 7
	}
	if product.PriceValue > 25000 && product.PriceValue < 75000 {
		group.RequiredUserCount = 5
	}
	if product.PriceValue > 75000 {
		group.RequiredUserCount = 3
	}
	group.Id = bson.NewObjectId()
	group.MinDiscount = "10% Off"
	group.CreatedBy = admin.Id
	group.Product = *product
	group.ExpiresOn = time.Now().Add(time.Hour * 24 * 30)
	group.IsOn = true
	group.InterestedUsersCount = 0
	err = group.Create()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"group":   group,
	})
}

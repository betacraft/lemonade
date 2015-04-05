package controllers

import (
	"errors"
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonade/framework"
	"github.com/rainingclouds/lemonade/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func GetDealsForUser(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	if !user.DealId.Valid() {
		framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
			"success": false,
			"message": "Your deal is not approved yet",
		})
		return
	}
	phoneDeal, err := models.GetPhoneDealById(user.DealId)
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"deal":    phoneDeal,
	})
}

func GetDeal(w http.ResponseWriter, r *http.Request) {
	dealId := bone.GetValue(r, "id")
	if dealId == "" {
		framework.WriteError(w, r, http.StatusInternalServerError, errors.New("Illegal request"))
		return
	}
	phoneDeal, err := models.GetPhoneDealById(bson.ObjectIdHex(dealId))
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"deal":    phoneDeal,
	})
}

func CreateDeal(w http.ResponseWriter, r *framework.Request) {
	phoneDeal := new(models.PhoneDeal)
	err := framework.Bind(r.Request, phoneDeal)
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusBadRequest, err)
		return
	}
	err = phoneDeal.Create()
	if err != nil {
		framework.WriteError(w, r.Request, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"message": "Deal is created",
	})
}

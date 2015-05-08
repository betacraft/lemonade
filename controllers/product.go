package controllers

import (
	"errors"
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonades/framework"
	"github.com/rainingclouds/lemonades/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func UpdateProductPrice(w http.ResponseWriter, r *http.Request) {
	productId := bone.GetValue(r, "id")
	if productId == "" {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("Illegal Request"))
		return
	}
	if !bson.IsObjectIdHex(productId) {
		framework.WriteError(w, r, http.StatusBadRequest, errors.New("Illegal Request"))
		return
	}
	product, err := models.GetProductById(bson.ObjectIdHex(productId))
	if err != nil {
		framework.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}
	err = product.UpdatePrice()
	if err != nil {
		framework.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"product": product,
	})
}

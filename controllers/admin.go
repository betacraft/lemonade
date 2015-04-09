package controllers

import (
	"github.com/rainingclouds/lemonades/framework"
	"github.com/rainingclouds/lemonades/models"
	"net/http"
)

func RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	bodyMap, err := framework.ReadBody(r)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	a, err := models.CreateAdmin(bodyMap)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	err = a.Create()
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success":  true,
		"message":  "Admin registered successfully",
		"auth_key": a.AuthKey,
	})
}

func AuthenticateAdmin(w http.ResponseWriter, r *http.Request) {
	bodyMap, err := framework.ReadBody(r)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	a, err := models.CreateAdmin(bodyMap)
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	err = a.Authenticate()
	if err != nil {
		framework.WriteError(w, r, http.StatusBadRequest, err)
		return
	}
	framework.WriteResponse(w, http.StatusOK, framework.JSONResponse{
		"success": true,
		"message": "Admin authenticated successfully",
		"admin":   a,
	})
}

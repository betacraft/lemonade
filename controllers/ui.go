package controllers

import (
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonade/framework"
	"github.com/rainingclouds/lemonade/models"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	framework.WriteHtml(w, "index.html", nil)
}

func Dashboard(w http.ResponseWriter, r *framework.Request) {
	user := r.MustGet("user").(*models.User)
	framework.WriteHtml(w, "deal.html", map[string]interface{}{"dealId": user.DealId.Hex()})
}

func ShareWidget(w http.ResponseWriter, r *http.Request) {
	framework.WriteHtml(w, "shareWidget.html", map[string]interface{}{"dealId": bone.GetValue(r, "id")})
}

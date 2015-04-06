package controllers

import (
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonade/framework"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	framework.WriteHtml(w, "index.html", nil)
}

func ShareWidget(w http.ResponseWriter, r *http.Request) {
	framework.WriteHtml(w, "shareWidget.html", map[string]interface{}{"dealId": bone.GetValue(r, "id")})
}

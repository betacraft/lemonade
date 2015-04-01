package controllers

import (
	"github.com/rainingclouds/lemonade/framework"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	framework.WriteHtml(w, "index.html", nil)
}

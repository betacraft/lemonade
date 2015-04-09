package main

import (
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonades/controllers"
	"net/http"
)

// a function to push all the routes related with the project
// add all the new routes here
// the kernel will take care of adding these routes in the routine
func pushRoutes(mux *bone.Mux) {
	// actual apis
	// admin apis
	mux.Post("/api/v1/admin", http.HandlerFunc(controllers.RegisterAdmin))
	mux.Post("/api/v1/admin/login", http.HandlerFunc(controllers.AuthenticateAdmin))
}

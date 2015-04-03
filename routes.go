package main

import (
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonade/controllers"
	"github.com/rainingclouds/lemonade/interceptors"
	"net/http"
)

// a function to push all the routes related with the project
// add all the new routes here
// the kernel will take care of adding these routes in the routine
func pushRoutes(mux *bone.Mux) {
	mux.Get("/", http.HandlerFunc(controllers.Index))
	// actual apis
	// admin apis
	mux.Post("/api/v1/admin", http.HandlerFunc(controllers.RegisterAdmin))
	mux.Post("/api/v1/admin/login", http.HandlerFunc(controllers.AuthenticateAdmin))
	// user apis
	mux.Post("/api/v1/user", http.HandlerFunc(controllers.RegisterUser))
	mux.Get("/api/v1/user", interceptors.UserAuthenticate(controllers.GetUser))
	mux.Post("/api/v1/user/login", http.HandlerFunc(controllers.AuthenticateUser))
}

package main

import (
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonades/controllers"
	"github.com/rainingclouds/lemonades/framework"
	"github.com/rainingclouds/lemonades/interceptors"
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

	// user apis
	mux.Options("/api/v1/user", framework.OptionsHandler())
	mux.Post("/api/v1/user", http.HandlerFunc(controllers.RegisterUser))
	mux.Options("/api/v1/user/login", framework.OptionsHandler())
	mux.Post("/api/v1/user/login", http.HandlerFunc(controllers.AuthenticateUser))
	mux.Options("/api/v1/user/logout", framework.OptionsHandler())
	mux.Post("/api/v1/user/logout", interceptors.UserAuthenticate(controllers.UserLogout))
	mux.Options("/api/v1/user/confirm_email/:id", framework.OptionsHandler())
	// mux.Post("/api/v1/user/confirm_email/:id", )
	mux.Options("/api/v1/group", framework.OptionsHandler())
	mux.Post("/api/v1/group", interceptors.UserAuthenticate(controllers.CreateGroup))
	mux.Options("/api/v1/group/:id/join", framework.OptionsHandler())
	mux.Post("/api/v1/group/:id/join", interceptors.UserAuthenticate(controllers.JoinGroup))
	mux.Options("/api/v1/group/:id", framework.OptionsHandler())
	mux.Get("/api/v1/group/:id", http.HandlerFunc(controllers.GetGroup))
	mux.Options("/api/v1/groups", framework.OptionsHandler())
	mux.Get("/api/v1/groups", http.HandlerFunc(controllers.GetGroups))
}

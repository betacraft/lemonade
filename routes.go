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
	// global urls
	mux.Get("/api/v1/group/:id/share", http.HandlerFunc(controllers.ShareGroup))
	// admin apis
	mux.Post("/api/v1/admin", http.HandlerFunc(controllers.RegisterAdmin))
	mux.Post("/api/v1/admin/login", http.HandlerFunc(controllers.AuthenticateAdmin))
	mux.Post("/api/v1/admin/group", interceptors.AdminAuthenticate(controllers.AdminCreateGroup))

	// user apis
	mux.Options("/api/v1/user", framework.OptionsHandler())
	mux.Post("/api/v1/user", http.HandlerFunc(controllers.RegisterUser))
	mux.Get("/api/v1/user", interceptors.UserAuthenticate(controllers.GetUser))
	mux.Options("/api/v1/user/login", framework.OptionsHandler())
	mux.Post("/api/v1/user/login", http.HandlerFunc(controllers.AuthenticateUser))

	mux.Options("/api/v1/user/fb_login", framework.OptionsHandler())
	mux.Post("/api/v1/user/fb_login", http.HandlerFunc(controllers.LoginWithFacebook))

	mux.Options("/api/v1/user/gplus_login", framework.OptionsHandler())
	mux.Post("/api/v1/user/gplus_login", http.HandlerFunc(controllers.LoginWithGooglePlus))

	mux.Options("/api/v1/user/forgot_password", framework.OptionsHandler())
	mux.Post("/api/v1/user/forgot_password", http.HandlerFunc(controllers.ForgotPassword))

	mux.Options("/api/v1/user/:auth_key/confirm_email", framework.OptionsHandler())
	mux.Post("/api/v1/user/:auth_key/confirm_email", http.HandlerFunc(controllers.ConfirmEmail))
	mux.Options("/api/v1/user/:auth_key/update_password", framework.OptionsHandler())
	mux.Post("/api/v1/user/:auth_key/update_password", http.HandlerFunc(controllers.UpdatePassword))

	mux.Options("/api/v1/user/logout", framework.OptionsHandler())
	mux.Post("/api/v1/user/logout", interceptors.UserAuthenticate(controllers.UserLogout))
	mux.Options("/api/v1/user/groups/joined", framework.OptionsHandler())
	mux.Get("/api/v1/user/groups/joined", interceptors.UserAuthenticate(controllers.GetUserJoinedGroups))
	mux.Options("/api/v1/user/groups/created", framework.OptionsHandler())
	mux.Get("/api/v1/user/groups/created", interceptors.UserAuthenticate(controllers.GetUserCreatedGroups))

	mux.Options("/api/v1/group", framework.OptionsHandler())
	mux.Post("/api/v1/group", interceptors.UserAuthenticate(controllers.CreateGroup))
	mux.Options("/api/v1/group/:id/join", framework.OptionsHandler())
	mux.Post("/api/v1/group/:id/join", interceptors.UserAuthenticate(controllers.JoinGroup))

	mux.Options("/api/v1/group/:id/leave", framework.OptionsHandler())
	mux.Post("/api/v1/group/:id/leave", interceptors.UserAuthenticate(controllers.LeaveGroup))

	mux.Options("/api/v1/group/:id", framework.OptionsHandler())
	mux.Get("/api/v1/group/:id", http.HandlerFunc(controllers.GetGroup))
	mux.Options("/api/v1/groups", framework.OptionsHandler())
	mux.Get("/api/v1/groups", http.HandlerFunc(controllers.GetGroups))
}

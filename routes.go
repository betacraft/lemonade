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
	mux.Get("/share-widget/:id", http.HandlerFunc(controllers.ShareWidget))
	mux.Get("/share/:id", http.HandlerFunc(controllers.Share))
	mux.Get("/dashboard", interceptors.UserAuthenticate(controllers.Dashboard))
	mux.Get("/deal/:id", http.HandlerFunc(controllers.ShareWidget))
	// actual apis
	// admin apis
	mux.Post("/api/v1/admin", http.HandlerFunc(controllers.RegisterAdmin))
	mux.Post("/api/v1/admin/login", http.HandlerFunc(controllers.AuthenticateAdmin))
	// deals apis
	mux.Post("/api/v1/deals/create", interceptors.AdminAuthenticate(controllers.CreateDeal))
	// user apis
	mux.Get("/api/v1/user/deals", interceptors.UserAuthenticate(controllers.GetDealsForUser))
	mux.Get("/api/v1/deal/:id", http.HandlerFunc(controllers.GetDeal))
	mux.Post("/api/v1/user", http.HandlerFunc(controllers.RegisterUser))
	mux.Get("/api/v1/user", interceptors.UserAuthenticate(controllers.GetUser))
	mux.Post("/api/v1/user/subscribe/:deal_id", interceptors.AdminAuthenticate(controllers.AttachDeal))
	mux.Post("/api/v1/user/login", http.HandlerFunc(controllers.AuthenticateUser))
	mux.Post("/api/v1/user/logout", interceptors.UserAuthenticate(controllers.UserLogout))
}

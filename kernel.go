package main

import (
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonades/db"
	"github.com/rainingclouds/lemonades/interceptors"
	"github.com/rainingclouds/lemonades/logger"
	"github.com/rainingclouds/lemonades/mailer"
	"net/http"
	"os"
)

func getPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "3000"
	}
	return ":" + port
}

func main() {
	logger.Init()
	// initializing databases
	logger.Debug("Initializing mongo")
	err := db.InitMongo()
	if err != nil {
		logger.Panic(err.Error())
	}
	logger.Debug("Initializing mongo : Done")
	defer db.CloseMongo()
	mailer.Init()
	// initializing routes
	logger.Debug("Initializing routes")
	mux := bone.New()
	pushRoutes(mux)
	logger.Debug("Running server on ", getPort())
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	http.ListenAndServe(getPort(), interceptors.NewInterceptor(mux))
	logger.Warn("Closing server")
}

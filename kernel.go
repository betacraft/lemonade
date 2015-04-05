package main

import (
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonade/db"
	"github.com/rainingclouds/lemonade/framework"
	"github.com/rainingclouds/lemonade/interceptors"
	"github.com/rainingclouds/lemonade/logger"
	"github.com/rainingclouds/lemonade/mailer"
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
	err = framework.SetTemplate("views/")
	if err != nil {
		logger.Panic(err.Error())
	}
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	logger.Debug("Running server on ", getPort())
	http.ListenAndServe(getPort(), interceptors.NewInterceptor(mux))
	logger.Warn("Closing server")
}

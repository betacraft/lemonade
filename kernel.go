package main

import (
	"github.com/rainingclouds/lemonade/db"
	"github.com/rainingclouds/lemonade/logger"
	"os"
)

func getPort() string {
	if os.Getenv("ENV") == "prod" {
		return ":80"
	}
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
	// initializing routes
	logger.Debug("Initializing routes")
	mux := bone.New()
}

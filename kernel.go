package main

import (
	"github.com/go-zoo/bone"
	"github.com/rainingclouds/lemonades/aws"
	"github.com/rainingclouds/lemonades/db"
	"github.com/rainingclouds/lemonades/interceptors"
	"github.com/rainingclouds/lemonades/logger"
	"github.com/rainingclouds/lemonades/mailer"
	"github.com/rainingclouds/lemonades/models"
	"github.com/robfig/cron"
	"net/http"
	"os"
)

var c *cron.Cron

func startCronJobs() {
	logger.Debug("Starting the cron job")
	c = cron.New()
	c.AddFunc("@midnight", models.UpdateProductPrices)
	c.Start()
}

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
	aws.Init()
	mailer.Init()
	// initializing routes
	logger.Debug("Initializing routes")
	mux := bone.New()
	pushRoutes(mux)
	logger.Debug("Running server on ", getPort())
	startCronJobs()
	http.ListenAndServe(getPort(), interceptors.NewInterceptor(mux))
	logger.Warn("Closing server")
}

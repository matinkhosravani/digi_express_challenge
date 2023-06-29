package main

import (
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/digi_express_challenge/api/router"
	"github.com/matinkhosravani/digi_express_challenge/app"
	"log"
)

func main() {
	app.Boot()
	r := gin.New()

	api := r.Group("api/v1/")
	router.PartnerRoutes(api)
	// Start server
	err := r.Run(app.GetEnv().ServerAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
}

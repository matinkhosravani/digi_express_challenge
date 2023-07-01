package main

import (
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/digi_express_challenge/api/router"
	"github.com/matinkhosravani/digi_express_challenge/app"
	"github.com/matinkhosravani/digi_express_challenge/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

func main() {
	app.Boot()
	r := gin.New()
	swaggerRoutes(r)
	api := r.Group("api/v1/")
	router.PartnerRoutes(api)
	// Start server
	err := r.Run(app.GetEnv().ServerAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func swaggerRoutes(r *gin.Engine) {
	// Serve the Swagger API documentation
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

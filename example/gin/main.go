package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jailson-dias/swagno/example/gin/handlers"
)

func main() {

	handler := handlers.NewHandler()

	app := gin.Default()

	// set mock routes
	handler.SetRoutes(app)

	// set swagger routes
	handler.SetSwagger(app)

	app.Run()

}

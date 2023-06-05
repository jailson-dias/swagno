package main

import (
	"fmt"

	"github.com/go-swagno/swagno-fiber/swagger"
	"github.com/gofiber/fiber/v2"
	swagno "github.com/jailson-dias/swagno"
	"github.com/jailson-dias/swagno/example/fiber-multi-array/handlers"
)

func main() {

	productHandler := handlers.NewProductHandler()
	merchantHandler := handlers.NewMerchantHandler()

	app := fiber.New()

	// set mock routes
	productHandler.SetProductRoutes(app)
	merchantHandler.SetMerchantRoutes(app)

	// set swagger routes
	sw := swagno.CreateNewSwagger("Swagger API", "1.0")
	swagno.AddEndpoints(handlers.ProductEndpoints)
	swagno.AddEndpoints(handlers.MerchantEndpoints)

	swagger.SwaggerHandler(app, sw.GenerateDocs(), swagger.Config{Prefix: "/swagger"})

	// Listen app
	app.Listen(fmt.Sprintf(
		"%s:%s",
		"localhost",
		"8080"),
	)

}

package handlers

import (
	"github.com/gofiber/fiber/v2"
	. "github.com/jailson-dias/swagno"
	"github.com/jailson-dias/swagno/example/fiber/models"
)

type MerchantHandler struct {
}

func NewMerchantHandler() MerchantHandler {
	return MerchantHandler{}
}

func (h MerchantHandler) SetMerchantRoutes(a *fiber.App) {
	a.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello World!")
	}).Name("index")
}

var MerchantEndpoints = []Endpoint{
	// /merchant/{merchantName}?id={id} -> get product of a merchant
	EndPoint(GET, "/merchant/{merchant}", "merchant", Params(StrParam("merchant", true, ""), IntQuery("id", true, "product id")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
	EndPoint(GET, "/merchant", "merchant", Params(StrQuery("merchant", true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
}

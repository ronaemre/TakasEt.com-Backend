package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	repository := NewRepository()
	service := NewService(repository)
	api := newAPI(&service)
	app := SetupApp(&api)
	app.Listen(":3000")

}

func SetupApp(api *API) *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	app.Get("/products", api.GetProductsHandler)
	app.Post("/product", api.AddProductsHandler)
	/* app.Get("/products/:categoryId2", api.GetProductCategoriesHandler) */
	app.Get("/products/:categoryId", api.GetProductsCategoriesHandler)
	app.Get("/myProducts/:userId", api.GetMyProductsHandler)
	app.Put("/updateProduct/:id", api.UpdateProductHandler)

	app.Post("/users/register", api.CreateUserHandler)
	app.Post("/login", api.LoginHandler)
	app.Get("/user", api.UserHandler)

	app.Post("trade", api.AddTradeOfferHandler)
	app.Get("/tradeOffers", api.GetTradeOffersHandler)
	app.Put("/tradeOffer/:id", api.UpdateTradeOfferHandler)
	app.Delete("/product/:id", api.DeleteProductHandler)

	return app
}

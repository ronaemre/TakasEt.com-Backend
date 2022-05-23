package main

import (
	"example.com/greetings/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type UserCredentialsDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TradeOfferDTO struct {
	Status string `json:"status"`
}

type ProductUpdateDTO struct {
	ProductName    string `json:"productName"`
	ProductPrice   string `json:"productPrice"`
	ProductImage   string `json:"productImage"`
	ProductDetails string `json:"productDetails"`
}

type API struct {
	service *Service
}

func newAPI(service *Service) API {
	return API{
		service: service,
	}
}

func (api *API) GetProductsHandler(c *fiber.Ctx) error {

	products, err := api.service.GetProducts()
	switch err {
	case nil:
		c.JSON(products)
		c.Status(fiber.StatusOK)

	default:
		c.Status(fiber.StatusInternalServerError)
	}

	return nil
}

/* func (api *API) GetProductCategoriesHandler(c *fiber.Ctx) error {

	categoryId := c.Params("categoryId2")
	categories, err := api.service.GetProductCategories(categoryId)
	switch err {
	case nil:
		c.JSON(categories)
		c.Status(fiber.StatusOK)

	default:
		c.Status(fiber.StatusInternalServerError)
	}

	return nil
} */

func (api *API) GetProductsCategoriesHandler(c *fiber.Ctx) error {

	categoryId := c.Params("categoryId")

	comment, err := api.service.GetProductsCategories(categoryId)

	switch err {
	case nil:
		c.JSON(comment)
		c.Status(fiber.StatusOK)
	case CategoryNotFoundError:
		c.Status(fiber.StatusNotFound)
	default:
		c.Status(fiber.StatusInternalServerError)
	}

	return nil
}

func (api *API) GetMyProductsHandler(c *fiber.Ctx) error {

	userId := c.Params("userId")

	comment, err := api.service.GetMyProducts(userId)

	switch err {
	case nil:
		c.JSON(comment)
		c.Status(fiber.StatusOK)
	case CategoryNotFoundError:
		c.Status(fiber.StatusNotFound)
	default:
		c.Status(fiber.StatusInternalServerError)
	}

	return nil
}

func (api *API) AddProductsHandler(c *fiber.Ctx) error {

	createdProduct := models.Product{}
	err := c.BodyParser(&createdProduct)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
	}
	err = api.service.AddProduct(createdProduct)

	switch err {
	case nil:
		c.JSON(createdProduct)
		c.Status(fiber.StatusCreated)
	}

	return nil
}

func (api *API) LoginHandler(c *fiber.Ctx) error {
	userCredentials := UserCredentialsDTO{}
	err := c.BodyParser(&userCredentials)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
	}

	token, cookie, err := api.service.Login(userCredentials)

	switch err {
	case nil:
		c.JSON(token)
		c.Cookie(cookie)
		c.Status(fiber.StatusOK)
	case UserNotFoundError:
		c.Status(fiber.StatusNotFound)
	case WrongPasswordError:
		c.Status(fiber.StatusBadRequest)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) CreateUserHandler(c *fiber.Ctx) error {

	userDTO := UserDTO{}
	err := c.BodyParser(&userDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
	}

	user, err := api.service.CreateUser(userDTO)

	switch err {
	case nil:
		c.JSON(user)
		c.Status(fiber.StatusCreated)
	case UserAlreadyRegisteredError:
		c.Status(fiber.StatusBadRequest)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) UserHandler(c *fiber.Ctx) error {

	cookie := c.Cookies("user-token")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return nil
	}

	claims := token.Claims.(*jwt.StandardClaims)

	user, err := api.service.repository.GetUser(claims.Issuer)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return nil
	}

	return c.JSON(user)
}

func (api *API) AddTradeOfferHandler(c *fiber.Ctx) error {
	createdTrade := models.TradeOffer{}

	err := c.BodyParser(&createdTrade)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
	}

	err = api.service.AddTrade(createdTrade)

	switch err {
	case nil:
		c.JSON(createdTrade)
		c.Status(fiber.StatusCreated)
	}

	return nil

}

func (api *API) GetTradeOffersHandler(c *fiber.Ctx) error {

	offers, err := api.service.GetTradeOffers()
	switch err {
	case nil:
		c.JSON(offers)
		c.Status(fiber.StatusOK)

	default:
		c.Status(fiber.StatusInternalServerError)
	}

	return nil
}

func (api *API) UpdateTradeOfferHandler(c *fiber.Ctx) error {

	tradeOfferID := c.Params("id")
	tradeOfferDTO := TradeOfferDTO{}
	err := c.BodyParser(&tradeOfferDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
	}

	trade, err := api.service.UpdateTradeOffer(tradeOfferID, tradeOfferDTO)

	switch err {
	case nil:
		c.JSON(trade)
		c.Status(fiber.StatusOK)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) UpdateProductHandler(c *fiber.Ctx) error {

	productID := c.Params("id")
	productUpdateDTO := ProductUpdateDTO{}
	err := c.BodyParser(&productUpdateDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
	}

	trade, err := api.service.UpdateProduct(productID, productUpdateDTO)

	switch err {
	case nil:
		c.JSON(trade)
		c.Status(fiber.StatusOK)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) DeleteProductHandler(c *fiber.Ctx) error {

	id := c.Params("id")
	err := api.service.DeleteProduct(id)

	switch err {
	case nil:
		c.Status(fiber.StatusNoContent)
	default:
		c.Status(fiber.StatusInternalServerError)
	}

	return nil
}

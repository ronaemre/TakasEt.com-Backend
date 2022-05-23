package main

import (
	"errors"
	"strings"
	"time"

	"example.com/greetings/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	Token string `json:"token"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type Service struct {
	repository *Repository
}

var CategoryNotFoundError error = errors.New("Category not found!")
var UserNotFoundError error = errors.New("User not found!")
var WrongPasswordError error = errors.New("Wrong password!")
var UserAlreadyRegisteredError error = errors.New("User already registered!")
var TradeOfferNotFound error = errors.New("Trade Offer Not found")

const SecretKey = "14465375-b4a8-47fa-9692-c986d4a825ee"

func NewService(repository *Repository) Service {
	return Service{
		repository: repository,
	}
}

func (service *Service) GetProducts() ([]models.Product, error) {
	products, err := service.repository.GetProducts()

	if err != nil {
		return nil, err
	}

	return products, nil
}

/* func (service *Service) GetProductCategories(categoryId string) ([]models.Product, error) {
	categories, err := service.repository.GetProductCategories(categoryId)

	if err != nil {
		return nil, err
	}

	return categories, nil
} */

func (service *Service) GetProductsCategories(categoryId string) ([]models.Product, error) {
	category, err := service.repository.GetProductsCategories(categoryId)

	if err != nil {
		return nil, CategoryNotFoundError
	}

	return category, nil
}

func (service *Service) GetMyProducts(userId string) ([]models.Product, error) {
	products, err := service.repository.GetMyProducts(userId)

	if err != nil {
		return nil, CategoryNotFoundError
	}

	return products, nil
}

func (service *Service) AddProduct(product models.Product) error {
	product.ID = GenerateUUID(11)
	err := service.repository.CreateProduct(product)

	if err != nil {
		return err
	}
	return nil

}

func GenerateUUID(length int) string {
	uuid := uuid.New().String()

	uuid = strings.ReplaceAll(uuid, "-", "")

	if length < 1 {
		return uuid
	}
	if length > len(uuid) {
		length = len(uuid)
	}

	return uuid[0:length]
}

func (service *Service) Login(userCredentials UserCredentialsDTO) (*Token, *fiber.Cookie, error) {

	user, err := service.repository.GetUserByUsername(userCredentials.Username)

	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, UserNotFoundError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userCredentials.Password)); err != nil {

		return nil, nil, WrongPasswordError
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		return nil, nil, err
	}

	cookie := fiber.Cookie{
		Name:    "user-token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
	}

	return &Token{
		Token: token,
	}, &cookie, nil
}

func (service *Service) CreateUser(userDTO UserDTO) (*User, error) {

	existingUser, _ := service.repository.GetUserByUsername(userDTO.Username)

	if existingUser != nil {
		return nil, UserAlreadyRegisteredError
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := User{
		ID:       GenerateUUID(8),
		Username: userDTO.Username,
		Password: string(hashedPassword),
	}

	newUser, err := service.repository.CreateUser(user)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (service *Service) AddTrade(tradeOffer models.TradeOffer) error {
	tradeOffer.ID = GenerateUUID(11)
	err := service.repository.CreateTradeOffer(tradeOffer)

	if err != nil {
		return err
	}
	return nil

}

func (service *Service) GetTradeOffers() ([]models.TradeOffer, error) {
	offers, err := service.repository.GetTradeOffers()

	if err != nil {
		return nil, err
	}

	return offers, nil
}

func (service *Service) UpdateTradeOffer(tradeOfferID string, tradeOfferDTO TradeOfferDTO) (*models.TradeOffer, error) {
	existingTradeOffer, err := service.repository.GetTradeOffer(tradeOfferID)

	if err != nil {
		return nil, TradeOfferNotFound
	}

	existingTradeOffer.Status = tradeOfferDTO.Status

	_, err = service.repository.UpdateTradeOffer(*existingTradeOffer)

	if err != nil {
		return nil, err
	}

	return service.GetTradeOffer(tradeOfferID)
}

func (service *Service) UpdateProduct(productID string, productUpdateDTO ProductUpdateDTO) (*models.UpdateProduct, error) {
	existingProduct, err := service.repository.GetProduct(productID)

	if err != nil {
		return nil, TradeOfferNotFound
	}

	existingProduct.ProductName = productUpdateDTO.ProductName
	existingProduct.ProductPrice = productUpdateDTO.ProductPrice
	existingProduct.ProductImage = productUpdateDTO.ProductImage
	existingProduct.ProductDetails = productUpdateDTO.ProductDetails

	_, err = service.repository.UpdateProduct(*existingProduct)

	if err != nil {
		return nil, err
	}

	return service.GetUpdateProduct(productID)
}

func (service *Service) GetTradeOffer(tradeOfferID string) (*models.TradeOffer, error) {
	tradeOffer, err := service.repository.GetTradeOffer(tradeOfferID)

	if err != nil {
		return nil, TradeOfferNotFound
	}

	return tradeOffer, nil
}

func (service *Service) GetUpdateProduct(productID string) (*models.UpdateProduct, error) {
	updateProduct, err := service.repository.GetProduct(productID)

	if err != nil {
		return nil, TradeOfferNotFound
	}

	return updateProduct, nil
}

func (service *Service) DeleteProduct(id string) error {
	err := service.repository.DeleteProduct(id)

	if err != nil {
		return err
	}

	return nil
}

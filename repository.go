package main

import (
	"context"
	"log"
	"time"

	"example.com/greetings/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductEntity struct {
	ID             string `bson:"id"`
	ProductName    string `bson:"productName"`
	ProductPrice   string `bson:"productPrice"`
	ProductImage   string `bson:"productImage"`
	ProductDetails string `bson:"productDetails"`
	CategoryID     string `bson:"categoryId"`
	UserId         string `bson:"userId"`
}

type TradeOfferEntity struct {
	ID                string `bson:"id"`
	ProductName       string `bson:"productName"`
	ProductPrice      string `bson:"productPrice"`
	ProductImage      string `bson:"productImage"`
	ProductDetails    string `bson:"productDetails"`
	CategoryID        string `bson:"categoryId"`
	UserId            string `bson:"userId"`
	SenderUserId      string `bson:"senderUserId"`
	SenderProdName    string `bson:"senderProdName"`
	SenderProdPrice   string `bson:"senderProdPrice"`
	SenderProdImage   string `bson:"senderProdImage"`
	SenderProdDetails string `bson:"senderProdDetails"`
	Status            string `bson:"status"`
}

type ProductUpdateEntity struct {
	ID                string `bson:"id"`
	ProductName       string `bson:"productName"`
	ProductPrice      string `bson:"productPrice"`
	ProductImage      string `bson:"productImage"`
	ProductDetails    string `bson:"productDetails"`
	CategoryID        string `bson:"categoryId"`
	UserId            string `bson:"userId"`
	SenderUserId      string `bson:"senderUserId"`
	SenderProdName    string `bson:"senderProdName"`
	SenderProdPrice   string `bson:"senderProdPrice"`
	SenderProdImage   string `bson:"senderProdImage"`
	SenderProdDetails string `bson:"senderProdDetails"`
	Status            string `bson:"status"`
}

type UserEntity struct {
	ID       string `bson:"id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type Repository struct {
	client *mongo.Client
}

func NewRepository() *Repository {
	uri := "mongodb+srv://ronaemre:123@cluster0.4uapv.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return &Repository{client}
}

func (repository *Repository) GetProducts() ([]models.Product, error) {
	collection := repository.client.Database("ProductTestDatabase").Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{}

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	products := []models.Product{}
	for cur.Next(ctx) {
		var product models.Product
		err := cur.Decode(&product)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, product)
	}
	return products, nil
}

/* func (repository *Repository) GetProductCategories(categoryId string) ([]models.Product, error) {
	collection := repository.client.Database("ProductTestDatabase").Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"categoryId": categoryId}

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	products := []models.Product{}
	for cur.Next(ctx) {
		 var product models.Product
		ProductEntity := ProductEntity{}
		err := cur.Decode(&ProductEntity)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, convertProductEntityToModel(ProductEntity))
	}
	return products, nil
} */

func (repository *Repository) GetProductsCategories(categoryId string) ([]models.Product, error) {
	collection := repository.client.Database("ProductTestDatabase").Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"categoryId": categoryId}
	cur, err := collection.Find(ctx, filter)

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, CategoryNotFoundError
	}

	productEntity := []ProductEntity{}
	for cur.Next(ctx) {
		var product ProductEntity
		err := cur.Decode(&product)

		if err != nil {
			log.Fatal(err)
		}
		productEntity = append(productEntity, product)
	}

	if err != nil {
		return nil, err
	}

	products := convertProductEntitiesToProductModels(productEntity)

	return products, nil
}

func (repository *Repository) GetMyProducts(userId string) ([]models.Product, error) {
	collection := repository.client.Database("ProductTestDatabase").Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"userId": userId}
	cur, err := collection.Find(ctx, filter)

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, CategoryNotFoundError
	}

	productEntity := []ProductEntity{}
	for cur.Next(ctx) {
		var product ProductEntity
		err := cur.Decode(&product)

		if err != nil {
			log.Fatal(err)
		}
		productEntity = append(productEntity, product)
	}

	if err != nil {
		return nil, err
	}

	products := convertProductEntitiesToProductModels(productEntity)
	return products, nil
}

func (repository *Repository) CreateProduct(product models.Product) error {

	collection := repository.client.Database("ProductTestDatabase").Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	productEntity := convertProductModelToEntity(product)
	_, err := collection.InsertOne(ctx, productEntity)
	if err != nil {
		return err
	}
	return nil

}

func (repository *Repository) GetUserByUsername(username string) (*User, error) {
	collection := repository.client.Database("UsersDB").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	cur := collection.FindOne(ctx, bson.M{"username": username})

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, UserNotFoundError
	}

	userEntity := UserEntity{}
	err := cur.Decode(&userEntity)

	if err != nil {
		return nil, err
	}

	user := convertUserEntityToUserModel(userEntity)

	return &user, nil
}

func (repository *Repository) CreateUser(user User) (*User, error) {
	collection := repository.client.Database("UsersDB").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	userEntity := convertUserModelToUserEntity(user)

	_, err := collection.InsertOne(ctx, userEntity)

	if err != nil {
		return nil, err
	}

	return repository.GetUser(userEntity.ID)
}

func (repository *Repository) GetUser(userID string) (*User, error) {
	collection := repository.client.Database("UsersDB").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	cur := collection.FindOne(ctx, bson.M{"id": userID})

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, UserNotFoundError
	}

	userEntity := UserEntity{}
	err := cur.Decode(&userEntity)

	if err != nil {
		return nil, err
	}

	user := convertUserEntityToUserModel(userEntity)

	return &user, nil
}

func (repository *Repository) CreateTradeOffer(tradeOffer models.TradeOffer) error {
	collection := repository.client.Database("TradeOfferDatabase").Collection("TradeOffer")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tradeOfferEntity := convertTradeOfferModelToEntity(tradeOffer)

	_, err := collection.InsertOne(ctx, tradeOfferEntity)

	if err != nil {
		return err
	}
	return nil
}

func (repository *Repository) GetTradeOffers() ([]models.TradeOffer, error) {
	collection := repository.client.Database("TradeOfferDatabase").Collection("TradeOffer")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{}

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	offers := []models.TradeOffer{}
	for cur.Next(ctx) {
		var offer models.TradeOffer
		err := cur.Decode(&offer)
		if err != nil {
			log.Fatal(err)
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

func (repository *Repository) GetTradeOffer(tradeOfferID string) (*models.TradeOffer, error) {
	collection := repository.client.Database("TradeOfferDatabase").Collection("TradeOffer")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur := collection.FindOne(ctx, bson.M{"id": tradeOfferID})

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, TradeOfferNotFound
	}

	tradeOfferEntity := TradeOfferEntity{}
	err := cur.Decode(&tradeOfferEntity)

	if err != nil {
		return nil, err
	}

	trade := convertTradeOfferEntityToModel(tradeOfferEntity)

	return &trade, nil
}

func (repository *Repository) GetProduct(productID string) (*models.UpdateProduct, error) {
	collection := repository.client.Database("ProductTestDatabase").Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur := collection.FindOne(ctx, bson.M{"id": productID})

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, TradeOfferNotFound
	}

	productUpdateEntity := ProductUpdateEntity{}
	err := cur.Decode(&productUpdateEntity)

	if err != nil {
		return nil, err
	}

	UpdateProduct := convertProductUpdateEntityToModel(productUpdateEntity)

	return &UpdateProduct, nil
}

func (repository *Repository) UpdateTradeOffer(tradeOffer models.TradeOffer) (*models.TradeOffer, error) {
	collection := repository.client.Database("TradeOfferDatabase").Collection("TradeOffer")
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	tradeOfferEntity := convertTradeOfferModelToEntity(tradeOffer)

	result := collection.FindOneAndReplace(ctx, bson.M{"id": tradeOffer.ID}, tradeOfferEntity)

	if result.Err() != nil {
		return nil, result.Err()
	}

	return repository.GetTradeOffer(tradeOffer.ID)
}

func (repository *Repository) UpdateProduct(updateProduct models.UpdateProduct) (*models.UpdateProduct, error) {
	collection := repository.client.Database("ProductTestDatabase").Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	updateProductEntity := convertProductUpdateModelToEntity(updateProduct)

	result := collection.FindOneAndReplace(ctx, bson.M{"id": updateProduct.ID}, updateProductEntity)

	if result.Err() != nil {
		return nil, result.Err()
	}

	return repository.GetProduct(updateProduct.ID)
}

func (repository *Repository) DeleteProduct(id string) error {
	collection := repository.client.Database("ProductTestDatabase").Collection("Products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": id}

	_, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	return nil

}

func convertProductEntityToModel(productEntity ProductEntity) models.Product {
	return models.Product{
		ID:             productEntity.ID,
		ProductName:    productEntity.ProductName,
		ProductPrice:   productEntity.ProductPrice,
		ProductImage:   productEntity.ProductImage,
		ProductDetails: productEntity.ProductDetails,
		CategoryID:     productEntity.CategoryID,
		UserId:         productEntity.UserId,
	}
}

func convertProductModelToEntity(product models.Product) ProductEntity {
	return ProductEntity{ //CommentEntitiy
		ID:             product.ID,
		ProductName:    product.ProductName,
		ProductPrice:   product.ProductPrice,
		ProductImage:   product.ProductImage,
		ProductDetails: product.ProductDetails,
		CategoryID:     product.CategoryID,
		UserId:         product.UserId,
	}
}

//TradeOffer
func convertTradeOfferEntityToModel(tradeOfferEntity TradeOfferEntity) models.TradeOffer {
	return models.TradeOffer{
		ID:                tradeOfferEntity.ID,
		ProductName:       tradeOfferEntity.ProductName,
		ProductPrice:      tradeOfferEntity.ProductPrice,
		ProductImage:      tradeOfferEntity.ProductImage,
		ProductDetails:    tradeOfferEntity.ProductDetails,
		CategoryID:        tradeOfferEntity.CategoryID,
		UserId:            tradeOfferEntity.UserId,
		SenderUserId:      tradeOfferEntity.SenderUserId,
		SenderProdName:    tradeOfferEntity.SenderProdName,
		SenderProdPrice:   tradeOfferEntity.SenderProdPrice,
		SenderProdImage:   tradeOfferEntity.SenderProdImage,
		SenderProdDetails: tradeOfferEntity.SenderProdDetails,
		Status:            tradeOfferEntity.Status,
	}
}

func convertTradeOfferModelToEntity(tradeOffer models.TradeOffer) TradeOfferEntity {
	return TradeOfferEntity{
		ID:                tradeOffer.ID,
		ProductName:       tradeOffer.ProductName,
		ProductPrice:      tradeOffer.ProductPrice,
		ProductImage:      tradeOffer.ProductImage,
		ProductDetails:    tradeOffer.ProductDetails,
		CategoryID:        tradeOffer.CategoryID,
		UserId:            tradeOffer.UserId,
		SenderUserId:      tradeOffer.SenderUserId,
		SenderProdName:    tradeOffer.SenderProdName,
		SenderProdPrice:   tradeOffer.SenderProdPrice,
		SenderProdImage:   tradeOffer.SenderProdImage,
		SenderProdDetails: tradeOffer.SenderProdDetails,
		Status:            tradeOffer.Status,
	}
}

func convertProductEntitiesToProductModels(productEntity []ProductEntity) []models.Product {
	products := []models.Product{}

	for _, entity := range productEntity {
		product := convertProductEntityToModel(entity)
		products = append(products, product)

	}
	return products
}

//TradeOffer
func convertTradeOfferEntitiesToTradeModels(tradeOfferEntity []TradeOfferEntity) []models.TradeOffer {
	tradeOffers := []models.TradeOffer{}

	for _, entity := range tradeOfferEntity {
		tradeOffer := convertTradeOfferEntityToModel(entity)
		tradeOffers = append(tradeOffers, tradeOffer)

	}
	return tradeOffers
}

//UPDATE PRODUCT
func convertProductUpdateEntityToModel(productUpdateEntity ProductUpdateEntity) models.UpdateProduct {
	return models.UpdateProduct{
		ID:                productUpdateEntity.ID,
		ProductName:       productUpdateEntity.ProductName,
		ProductPrice:      productUpdateEntity.ProductPrice,
		ProductImage:      productUpdateEntity.ProductImage,
		ProductDetails:    productUpdateEntity.ProductDetails,
		CategoryID:        productUpdateEntity.CategoryID,
		UserId:            productUpdateEntity.UserId,
		SenderUserId:      productUpdateEntity.SenderUserId,
		SenderProdName:    productUpdateEntity.SenderProdName,
		SenderProdPrice:   productUpdateEntity.SenderProdPrice,
		SenderProdImage:   productUpdateEntity.SenderProdImage,
		SenderProdDetails: productUpdateEntity.SenderProdDetails,
		Status:            productUpdateEntity.Status,
	}
}

func convertProductUpdateModelToEntity(updateProduct models.UpdateProduct) ProductUpdateEntity {
	return ProductUpdateEntity{
		ID:                updateProduct.ID,
		ProductName:       updateProduct.ProductName,
		ProductPrice:      updateProduct.ProductPrice,
		ProductImage:      updateProduct.ProductImage,
		ProductDetails:    updateProduct.ProductDetails,
		CategoryID:        updateProduct.CategoryID,
		UserId:            updateProduct.UserId,
		SenderUserId:      updateProduct.SenderUserId,
		SenderProdName:    updateProduct.SenderProdName,
		SenderProdPrice:   updateProduct.SenderProdPrice,
		SenderProdImage:   updateProduct.SenderProdImage,
		SenderProdDetails: updateProduct.SenderProdDetails,
		Status:            updateProduct.Status,
	}
}

func convertUserEntityToUserModel(userEntity UserEntity) User {
	return User{
		ID:       userEntity.ID,
		Username: userEntity.Username,
		Password: userEntity.Password,
	}
}

func convertUserModelToUserEntity(user User) UserEntity {
	return UserEntity{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}
}

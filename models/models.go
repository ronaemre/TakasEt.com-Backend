package models

type Product struct {
	ID             string `json:"id,omitempty"`
	ProductName    string `json:"productName"`
	ProductPrice   string `json:"productPrice"`
	ProductImage   string `json:"ProductImage"`
	ProductDetails string `json:"productDetails"`
	CategoryID     string `json:"categoryId"`
	UserId         string `json:"userId"`
}

type TradeOffer struct {
	ID                string `json:"id,omitempty"`
	ProductName       string `json:"productName"`
	ProductPrice      string `json:"productPrice"`
	ProductImage      string `json:"productImage"`
	ProductDetails    string `json:"productDetails"`
	CategoryID        string `json:"categoryId"`
	UserId            string `json:"userId"`
	SenderUserId      string `json:"senderUserId"`
	SenderProdName    string `json:"senderProdName"`
	SenderProdPrice   string `json:"senderProdPrice"`
	SenderProdImage   string `json:"senderProdImage"`
	SenderProdDetails string `json:"senderProdDetails"`
	Status            string `json:"status"`
}

type UpdateProduct struct {
	ID                string `json:"id,omitempty"`
	ProductName       string `json:"productName"`
	ProductPrice      string `json:"productPrice"`
	ProductImage      string `json:"productImage"`
	ProductDetails    string `json:"productDetails"`
	CategoryID        string `json:"categoryId"`
	UserId            string `json:"userId"`
	SenderUserId      string `json:"senderUserId"`
	SenderProdName    string `json:"senderProdName"`
	SenderProdPrice   string `json:"senderProdPrice"`
	SenderProdImage   string `json:"senderProdImage"`
	SenderProdDetails string `json:"senderProdDetails"`
	Status            string `json:"status"`
}

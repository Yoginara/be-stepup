package models

import "time"

// Cart represents the structure of a shopping cart
type Cart struct {
	CartID     string     `bson:"cart_id" json:"cart_id"`         // Unique identifier for the cart
	UserID     string     `bson:"user_id" json:"user_id"`         // Reference to the User
	UserName   string     `bson:"user_name" json:"user_name"`     // Name of the user
	Items      []CartItem `bson:"items" json:"items"`             // List of cart items
	CreatedAt  time.Time  `bson:"created_at" json:"created_at"`   // Timestamp when the cart was created
	ModifiedAt time.Time  `bson:"modified_at" json:"modified_at"` // Timestamp when the cart was last updated
}

// CartItem defines the structure of an item in the cart
type CartItem struct {
	ProductID   string  `bson:"product_id" json:"product_id"`     // Unique identifier for the product
	ProductCode string  `bson:"product_code" json:"product_code"` // Product code for reference
	ProductName string  `bson:"product_name" json:"product_name"` // Name of the product
	Quantity    int     `bson:"quantity" json:"quantity"`         // Quantity of the product in the cart
	Price       float64 `bson:"price" json:"price"`               // Price of the product
	ImageURL    string  `bson:"image_url" json:"image_url"`       // URL of the product image
}

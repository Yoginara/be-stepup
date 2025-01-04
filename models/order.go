package models

import "time"

// Order represents the structure of a user's order
type Order struct {
	OrderID         string      `bson:"order_id" json:"order_id"`                 // Unique identifier for the order
	UserID          string      `bson:"user_id" json:"user_id"`                   // Reference to the User
	UserName        string      `bson:"user_name" json:"user_name"`               // Name of the user
	TotalPrice      float64     `bson:"total_price" json:"total_price"`           // Total price of the order
	Items           []OrderItem `bson:"items" json:"items"`                       // List of items in the order
	OrderStatus     string      `bson:"order_status" json:"order_status"`         // Current status of the order (e.g., pending, shipped, delivered)
	CreatedAt       time.Time   `bson:"created_at" json:"created_at"`             // Timestamp when the order was created
	ModifiedAt      time.Time   `bson:"modified_at" json:"modified_at"`           // Timestamp when the order was last updated
	ShippingAddress string      `bson:"shipping_address" json:"shipping_address"` // Address to deliver the order
}

// OrderItem defines the structure of an item in the order
type OrderItem struct {
	ProductID   string  `bson:"product_id" json:"product_id"`     // Unique identifier for the product
	ProductCode string  `bson:"product_code" json:"product_code"` // Product code for reference
	ProductName string  `bson:"product_name" json:"product_name"` // Name of the product
	Quantity    int     `bson:"quantity" json:"quantity"`         // Quantity of the product in the order
	Price       float64 `bson:"price" json:"price"`               // Price of the product
	ImageURL    string  `bson:"image_url" json:"image_url"`       // URL of the product image
}

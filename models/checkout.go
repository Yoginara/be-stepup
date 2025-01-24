package models

import "time"

// Checkout represents the structure of the checkout process
type Checkout struct {
	CheckoutID  string     `bson:"checkout_id" json:"checkout_id"`   // Unique identifier for the checkout
	UserID      string     `bson:"user_id" json:"user_id"`           // Reference to the User
	UserName    string     `bson:"user_name" json:"user_name"`       // Name of the user
	Items       []CartItem `bson:"items" json:"items"`               // List of items to be purchased
	TotalPrice  float64    `bson:"total_price" json:"total_price"`   // Total price of the checkout
	Address     string     `bson:"address" json:"address"`           // Address for delivery
	PhoneNumber string     `bson:"phone_number" json:"phone_number"` // Phone number of the user
	Status      string     `bson:"status" json:"status"`             // Status of the checkout (e.g., "Pending", "Completed")
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`     // Timestamp when the checkout was created
	ModifiedAt  time.Time  `bson:"modified_at" json:"modified_at"`   // Timestamp when the checkout was last updated
}

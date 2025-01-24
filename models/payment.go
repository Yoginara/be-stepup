package models

import "time"

// Payment represents the structure of a payment
type Payment struct {
	PaymentID     string    `bson:"payment_id" json:"payment_id"`         // Unique identifier for the payment
	CheckoutID    string    `bson:"checkout_id" json:"checkout_id"`       // Reference to the checkout
	UserID        string    `bson:"user_id" json:"user_id"`               // Reference to the User
	PaymentImage  string    `bson:"payment_image" json:"payment_image"`   // Path to the payment image (e.g., PNG, JPEG)
	PaymentStatus string    `bson:"payment_status" json:"payment_status"` // Status of the payment (e.g., "Pending", "Verified")
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`         // Timestamp when the payment was created
	ModifiedAt    time.Time `bson:"modified_at" json:"modified_at"`       // Timestamp when the payment was last updated
}

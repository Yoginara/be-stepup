package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID     primitive.ObjectID `bson:"product_id" json:"product_id"`
	ProductName   string             `bson:"product_name" json:"product_name"`
	ProductCode   string             `bson:"product_code" json:"product_code"`
	Quantity      int                `bson:"quantity" json:"quantity"`
	TotalPrice    float64            `bson:"total_price" json:"total_price"`
	CustomerName  string             `bson:"customer_name" json:"customer_name"`
	CustomerEmail string             `bson:"customer_email" json:"customer_email"`
	OrderDate     primitive.DateTime `bson:"order_date" json:"order_date"`
	Status        string             `bson:"status" json:"status"`
}

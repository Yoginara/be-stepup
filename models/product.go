package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Product defines the structure of product data
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID   string             `bson:"product_id" json:"product_id"`
	Code        string             `bson:"code" json:"code"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Brand       string             `bson:"brand" json:"brand"`
	Category    string             `bson:"category" json:"category"`
	Color       string             `bson:"color" json:"color"`
	Price       float64            `bson:"price" json:"price"`
	Stock       int                `bson:"stock" json:"stock"`
	ImageURL    string             `bson:"image_url" json:"image_url"`
}

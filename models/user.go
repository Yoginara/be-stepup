package models

import "time"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents the structure of a user
type User struct {
	UserID    string     `bson:"userid"`
	Email     string     `bson:"email"`
	Password  string     `bson:"password"`
	Role      string     `bson:"role"`
	Name      string     `bson:"name"`       // Menambahkan field nama
	CreatedAt time.Time  `bson:"created_at"` // Menambahkan field tanggal pembuatan akun
	Cart      []CartItem `bson:"cart"`       // Field untuk menyimpan keranjang belanja
}

// CartItem defines the structure of an item in the scart
type CartItem struct {
	ProductID   string  `bson:"product_id" json:"product_id"`
	ProductCode string  `bson:"product_code" json:"product_code"`
	Quantity    int     `bson:"quantity" json:"quantity"`
	Price       float64 `bson:"price" json:"price"`
}

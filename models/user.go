package models

import "time"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents the structure of a user
type User struct {
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	Role      string    `bson:"role"`
	Name      string    `bson:"name"`       // Menambahkan field nama
	CreatedAt time.Time `bson:"created_at"` // Menambahkan field tanggal pembuatan akun
}

package controllers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"

	"be-stepup/config"
	"be-stepup/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAllUsers retrieves all users from the database
func GetAllUsers(c *fiber.Ctx) error {
	// Membuat konteks dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mengambil koleksi user dari database
	userCollection := config.GetCollection("users")

	// Query semua user
	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	defer cursor.Close(ctx)

	// Menyimpan hasil query ke dalam slice
	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode users",
		})
	}

	// Mengembalikan data user sebagai JSON
	return c.JSON(users)
}

// GetUserByID retrieves a user by their ID from the database
func GetUserByID(c *fiber.Ctx) error {
	// Ambil ID dari parameter URL
	id := c.Params("id")

	// Validasi ID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Membuat konteks dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mengambil koleksi user dari database
	userCollection := config.GetCollection("users")

	// Mencari user berdasarkan ID
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Mengembalikan data user sebagai JSON
	return c.JSON(user)
}

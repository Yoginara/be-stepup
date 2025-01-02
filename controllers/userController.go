package controllers

import (
	"context"
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

func GetUserByID(c *fiber.Ctx) error {
	// Ambil ID dari parameter URL
	id := c.Params("id")

	// Membuat konteks dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mengambil koleksi user dari database
	userCollection := config.GetCollection("users")

	// Mencari user berdasarkan `userid`
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"userid": id}).Decode(&user)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	// Mengembalikan data user sebagai JSON
	return c.JSON(user)
}

package controllers

import (
	"be-stepup/config"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// Fungsi untuk menghitung jumlah produk
func GetProductCount(c *fiber.Ctx) error {
	// Menggunakan GetCollection untuk mengambil koleksi produk
	collection := config.GetCollection("products")

	// Menghitung jumlah produk
	count, err := collection.CountDocuments(c.Context(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghitung jumlah produk",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"product_count": count,
	})
}

// Fungsi untuk menghitung jumlah pengguna
func GetUserCount(c *fiber.Ctx) error {
	// Menggunakan GetCollection untuk mengambil koleksi pengguna
	collection := config.GetCollection("users")

	// Menghitung jumlah pengguna
	count, err := collection.CountDocuments(c.Context(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghitung jumlah pengguna",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user_count": count,
	})
}

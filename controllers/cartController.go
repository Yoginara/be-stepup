package controllers

import (
	"be-stepup/config"
	"be-stepup/models"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AddToCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var item models.CartItem
	if err := c.BodyParser(&item); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	collection := config.GetCollection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := map[string]string{"userId": userID}
	update := map[string]interface{}{
		"$push": map[string]interface{}{"cart": item},
		"$set":  map[string]interface{}{"updatedAt": time.Now()},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add item to cart"})
	}

	return c.JSON(fiber.Map{"message": "Item added to cart"})
}

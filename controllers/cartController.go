package controllers

import (
	"be-stepup/config"
	"be-stepup/models"
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

// AddToCart handles adding items to a user's cart
func AddToCart(c *fiber.Ctx) error {
	// Parse request body
	var cartItem models.CartItem
	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get the user ID from the token (from JWT middleware)
	userID := c.Locals("userID").(string)

	// Get the user from the database
	collection := config.GetCollection("users")
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, map[string]string{"userid": userID}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Add the item to the user's cart
	user.Cart = append(user.Cart, cartItem)

	// Update the user document in the database
	_, err = collection.UpdateOne(
		ctx,
		map[string]string{"userid": userID},
		map[string]interface{}{"$set": map[string]interface{}{"cart": user.Cart}},
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update cart",
		})
	}

	// Respond with the updated cart
	return c.JSON(fiber.Map{
		"message": "Item added to cart successfully",
		"cart":    user.Cart,
	})
}

// RemoveSingleCartItem handles removing a single item from the user's cart
func RemoveSingleCartItem(c *fiber.Ctx) error {
	// Parse request body
	var request struct {
		ProductID string `json:"product_id"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get the user ID from the token (from JWT middleware)
	userID := c.Locals("userID").(string)

	// Get the user from the database
	collection := config.GetCollection("users")
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, map[string]string{"userid": userID}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Find and remove the item from the cart
	updatedCart := []models.CartItem{}
	itemFound := false

	for _, item := range user.Cart {
		if item.ProductID == request.ProductID {
			// Skip this item (remove it)
			itemFound = true
			continue
		}
		// Keep the other items
		updatedCart = append(updatedCart, item)
	}

	// If no matching item is found
	if !itemFound {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Item not found in cart",
		})
	}

	// Update the user's cart in the database
	_, err = collection.UpdateOne(
		ctx,
		map[string]string{"userid": userID},
		map[string]interface{}{"$set": map[string]interface{}{"cart": updatedCart}},
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update cart",
		})
	}

	// Respond with the updated cart
	return c.JSON(fiber.Map{
		"message": "Item removed from cart successfully",
		"cart":    updatedCart,
	})
}

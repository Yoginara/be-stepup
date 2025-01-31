package controllers

import (
	"be-stepup/config"
	"be-stepup/models"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

var validate = validator.New()

// Fungsi untuk membuat ID unik khusus checkout
func generateCheckoutID() string {
	return uuid.New().String()
}

// CreateCheckout handles the creation of a checkout
func CreateCheckout(c *fiber.Ctx) error {
	// Parse request body
	var checkoutRequest struct {
		Address     string `json:"address" validate:"required,min=5,max=100"`
		PhoneNumber string `json:"phone_number" validate:"required,min=10,max=15"`
	}

	if err := c.BodyParser(&checkoutRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Body request tidak valid",
		})
	}

	// Validasi data request menggunakan validator
	if err := validate.Struct(checkoutRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Data tidak valid: " + err.Error(),
		})
	}

	// Mendapatkan userID dari token (dari middleware JWT)
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User tidak terautentikasi",
		})
	}
	log.Println("UserID:", userID)

	// Mengambil koleksi `cart`
	cartCollection := config.GetCollection("cart")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mencari keranjang berdasarkan userID
	var cart models.Cart
	err := cartCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Keranjang tidak ditemukan",
		})
	}

	// Mengambil koleksi `products`
	productCollection := config.GetCollection("products")

	// Menghitung total harga dan mengurangi stok produk
	var totalPrice float64
	for _, item := range cart.Items {
		var product models.Product
		err := productCollection.FindOne(ctx, bson.M{"product_id": item.ProductID}).Decode(&product)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{
					"success": false,
					"error":   "Produk tidak ditemukan",
				})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Gagal mendapatkan produk",
			})
		}

		// Mengecek apakah stok cukup
		if product.Stock < item.Quantity {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Stok produk tidak cukup untuk " + product.Name,
			})
		}

		// Mengurangi stok produk
		newStock := product.Stock - item.Quantity
		_, err = productCollection.UpdateOne(
			ctx,
			bson.M{"product_id": item.ProductID},
			bson.M{"$set": bson.M{"stock": newStock}},
		)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Gagal memperbarui stok produk",
			})
		}

		totalPrice += item.Price * float64(item.Quantity)
	}

	// Mengambil data pengguna untuk mendapatkan user_name
	userCollection := config.GetCollection("users")
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"userid": userID}).Decode(&user)
	if err != nil {
		log.Printf("Error finding user for userID %s: %v\n", userID, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mendapatkan nama pengguna",
		})
	}

	// Membuat checkout baru
	checkout := models.Checkout{
		CheckoutID:  generateCheckoutID(),
		UserID:      userID,
		UserName:    user.Name,
		Items:       cart.Items,
		TotalPrice:  totalPrice,
		Address:     checkoutRequest.Address,
		PhoneNumber: checkoutRequest.PhoneNumber,
		Status:      "Pending", // Status awal adalah Pending
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
	}

	// Mengambil koleksi `checkout`
	checkoutCollection := config.GetCollection("checkout")
	_, err = checkoutCollection.InsertOne(ctx, checkout)
	if err != nil {
		log.Printf("Error creating checkout for userID %s: %v\n", userID, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal membuat checkout",
		})
	}

	// Hapus keranjang setelah checkout berhasil
	_, err = cartCollection.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		log.Printf("Error deleting cart for userID %s: %v\n", userID, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal menghapus keranjang",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Checkout berhasil dibuat",
		"data":    checkout,
	})
}

// GetCheckoutByID handles getting the checkout by ID
func GetCheckoutByID(c *fiber.Ctx) error {
	// Mendapatkan checkoutID dari parameter URL
	checkoutID := c.Params("checkout_id")

	// Mengambil koleksi `checkout`
	collection := config.GetCollection("checkout")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mencari checkout berdasarkan checkoutID
	var checkout models.Checkout
	err := collection.FindOne(ctx, bson.M{"checkout_id": checkoutID}).Decode(&checkout)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Checkout tidak ditemukan",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mendapatkan status checkout",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    checkout,
	})
}

func UpdateCheckout(c *fiber.Ctx) error {
	var updateRequest struct {
		Status string `json:"status" validate:"required"`
	}

	if err := c.BodyParser(&updateRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Body request tidak valid",
		})
	}

	if err := validate.Struct(updateRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Data tidak valid: " + err.Error(),
		})
	}

	checkoutID := c.Params("checkout_id")

	collection := config.GetCollection("checkout")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"checkout_id": checkoutID},
		bson.M{"$set": bson.M{
			"status":      updateRequest.Status,
			"modified_at": time.Now(),
		}},
	)
	if err != nil {
		log.Printf("Error updating checkout status for checkoutID %s: %v\n", checkoutID, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal memperbarui status checkout",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Status checkout berhasil diperbarui",
		"data": fiber.Map{
			"checkout_id": checkoutID,
			"status":      updateRequest.Status,
		},
	})
}

func GetAllCheckout(c *fiber.Ctx) error {
	collection := config.GetCollection("checkout")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding checkouts: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mendapatkan data checkout",
		})
	}
	defer cursor.Close(ctx)

	var checkouts []models.Checkout
	for cursor.Next(ctx) {
		var checkout models.Checkout
		if err := cursor.Decode(&checkout); err != nil {
			log.Printf("Error decoding checkout: %v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Gagal mendekode data checkout",
			})
		}
		checkouts = append(checkouts, checkout)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mengambil data checkout",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    checkouts,
	})
}

// DeleteCheckout handles the deletion of a checkout by checkoutID
func DeleteCheckout(c *fiber.Ctx) error {
	// Mendapatkan checkoutID dari parameter URL
	checkoutID := c.Params("checkout_id")

	// Mengambil koleksi checkout
	collection := config.GetCollection("checkout")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Menghapus checkout berdasarkan checkoutID
	_, err := collection.DeleteOne(ctx, bson.M{"checkout_id": checkoutID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Checkout tidak ditemukan",
			})
		}
		log.Printf("Error deleting checkout for checkoutID %s: %v\n", checkoutID, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal menghapus checkout",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Checkout berhasil dihapus",
	})
}

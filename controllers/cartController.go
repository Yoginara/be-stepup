package controllers

import (
	"be-stepup/config"
	"be-stepup/models"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

// GetAllCart menangani pengambilan semua item dalam keranjang
func GetAllCart(c *fiber.Ctx) error {
	// Mendapatkan userID dari token (dari middleware JWT)
	userID := c.Locals("userID").(string)
	log.Println("UserID:", userID) // Menambahkan log untuk melihat nilai userID

	// Mengambil koleksi `cart`
	collection := config.GetCollection("cart")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mencari keranjang berdasarkan userID
	var cart models.Cart
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		log.Println("Error finding cart:", err) // Menambahkan log error untuk mencari cart
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Keranjang tidak ditemukan",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mendapatkan keranjang",
		})
	}

	// Mengambil nama pengguna
	userCollection := config.GetCollection("users")
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"userid": userID}).Decode(&user)
	if err != nil {
		log.Println("Error finding user:", err) // Menambahkan log error untuk mencari user
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mendapatkan nama pengguna",
		})
	}

	// Mengembalikan item-item di dalam keranjang
	return c.JSON(fiber.Map{
		"cart_id":     cart.CartID,
		"user_id":     cart.UserID,
		"user_name":   user.Name, // Menambahkan nama pengguna
		"items":       cart.Items,
		"created_at":  cart.CreatedAt,
		"modified_at": cart.ModifiedAt,
	})
}

// generateUniqueID menghasilkan ID unik baru
func generateUniqueID() string {
	return uuid.New().String()
}

// AddToCart menangani penambahan item ke keranjang
func AddToCart(c *fiber.Ctx) error {
	// Mengambil data dari body request
	var cartItem models.CartItem
	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Body request tidak valid",
		})
	}

	// Mendapatkan userID dari token (dari middleware JWT)
	userID := c.Locals("userID").(string)
	log.Println("UserID:", userID) // Menambahkan log untuk melihat nilai userID

	// Mengambil koleksi `cart`
	collection := config.GetCollection("cart")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mengambil data pengguna untuk mendapatkan user_name
	userCollection := config.GetCollection("users")
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"userid": userID}).Decode(&user)
	if err != nil {
		log.Println("Error finding user:", err) // Menambahkan log error untuk mencari user
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mendapatkan nama pengguna",
		})
	}

	// Mencari produk berdasarkan ProductID
	productCollection := config.GetCollection("products")
	var product models.Product
	err = productCollection.FindOne(ctx, bson.M{"product_id": cartItem.ProductID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Produk tidak ditemukan",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mendapatkan produk",
		})
	}

	// Mengecek apakah stok cukup
	if product.Stock < cartItem.Quantity {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Stok produk tidak cukup",
		})
	}

	// Mencari keranjang berdasarkan userID
	var cart models.Cart
	err = collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Jika keranjang belum ada, buat keranjang baru dengan user_name
			cartItem.ProductCode = product.Code
			cartItem.ProductName = product.Name
			cartItem.ImageURL = product.ImageURL

			cart = models.Cart{
				CartID:     generateUniqueID(),
				UserID:     userID,
				UserName:   user.Name, // Menyimpan user_name
				Items:      []models.CartItem{cartItem},
				CreatedAt:  time.Now(),
				ModifiedAt: time.Now(),
			}
			_, err := collection.InsertOne(ctx, cart)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"error": "Gagal membuat keranjang baru",
				})
			}
			return c.JSON(fiber.Map{
				"message": "Item berhasil ditambahkan ke keranjang",
				"cart":    cart.Items,
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mendapatkan keranjang",
		})
	}

	// Menambahkan atau memperbarui item dalam keranjang
	itemFound := false
	for i, item := range cart.Items {
		if item.ProductID == cartItem.ProductID {
			// Jika item ditemukan, tambahkan kuantitas
			cart.Items[i].Quantity += cartItem.Quantity
			itemFound = true
			break
		}
	}

	if !itemFound {
		// Tambahkan item baru dengan informasi produk yang diperlukan
		cartItem.ProductCode = product.Code
		cartItem.ProductName = product.Name
		cartItem.ImageURL = product.ImageURL

		cart.Items = append(cart.Items, cartItem)
	}

	// Perbarui keranjang di database dengan user_name
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{
			"items":       cart.Items,
			"user_name":   user.Name, // Update user_name di cart
			"modified_at": time.Now(),
		}})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memperbarui keranjang",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Item berhasil ditambahkan ke keranjang",
		"cart":    cart.Items,
	})
}

// RemoveSingleCartItem menangani penghapusan satu item dari keranjang
func RemoveSingleCartItem(c *fiber.Ctx) error {
	// Parse request body
	var request struct {
		ProductID string `json:"product_id"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Body request tidak valid",
		})
	}

	// Mendapatkan userID dari token (dari middleware JWT)
	userID := c.Locals("userID").(string)

	// Mengambil koleksi `cart`
	collection := config.GetCollection("cart")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mencari keranjang berdasarkan userID
	var cart models.Cart
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Keranjang tidak ditemukan",
		})
	}

	// Menghapus item dari keranjang
	itemFound := false
	updatedItems := []models.CartItem{}

	for _, item := range cart.Items {
		if item.ProductID == request.ProductID {
			itemFound = true
			continue // Hapus item dari keranjang
		}
		updatedItems = append(updatedItems, item)
	}

	if !itemFound {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Item tidak ditemukan di keranjang",
		})
	}

	// Jika Items kosong setelah penghapusan, hapus dokumen keranjang dari koleksi
	if len(updatedItems) == 0 {
		_, err := collection.DeleteOne(ctx, bson.M{"user_id": userID})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menghapus keranjang kosong",
			})
		}
	} else {
		// Perbarui keranjang dengan item yang sudah diperbarui
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"user_id": userID},
			bson.M{"$set": bson.M{
				"items":       updatedItems,
				"modified_at": time.Now(),
			}})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal memperbarui keranjang",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Item berhasil dihapus dari keranjang",
		"cart":    updatedItems,
	})
}

// UpdateCartItem menangani pembaruan item di keranjang
func UpdateCartItem(c *fiber.Ctx) error {
	// Parse request body untuk mendapatkan data pembaruan
	var updateItemRequest struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := c.BodyParser(&updateItemRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Body request tidak valid",
		})
	}

	// Mendapatkan userID dari token (dari middleware JWT)
	userID := c.Locals("userID").(string)

	// Mengambil koleksi `cart`
	collection := config.GetCollection("cart")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mencari keranjang berdasarkan userID
	var cart models.Cart
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Keranjang tidak ditemukan",
		})
	}

	// Mencari item dalam keranjang berdasarkan ProductID dan memperbarui kuantitasnya
	itemFound := false
	for i, item := range cart.Items {
		if item.ProductID == updateItemRequest.ProductID {
			// Jika item ditemukan, perbarui kuantitasnya
			cart.Items[i].Quantity = updateItemRequest.Quantity
			itemFound = true
			break
		}
	}

	// Jika item tidak ditemukan dalam keranjang
	if !itemFound {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Item tidak ditemukan di keranjang",
		})
	}

	// Perbarui keranjang di database dengan kuantitas baru
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{
			"items":       cart.Items,
			"modified_at": time.Now(),
		}})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memperbarui keranjang",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Item berhasil diperbarui",
		"cart":    cart.Items,
	})
}

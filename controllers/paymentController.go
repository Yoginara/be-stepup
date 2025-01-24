package controllers

import (
	"be-stepup/config"
	"be-stepup/models"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Fungsi untuk membuat ID unik khusus pembayaran
func generatePaymentID() string {
	return uuid.New().String()
}

// SavePayment handles saving the payment proof image
func SavePayment(c *fiber.Ctx) error {
	// Mendapatkan checkoutID dari parameter URL
	checkoutID := c.Params("checkout_id")

	// Periksa apakah checkoutID tersedia
	if checkoutID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "checkout_id tidak ditemukan dalam URL",
		})
	}

	// Mengambil file gambar bukti pembayaran
	file, err := c.FormFile("payment_image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Gambar pembayaran tidak ditemukan atau tidak valid",
		})
	}

	// Validasi ekstensi file (hanya PNG, JPEG, JPG yang diperbolehkan)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" && ext != ".jpeg" && ext != ".jpg" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Tipe file tidak valid. Hanya file PNG, JPEG, atau JPG yang diperbolehkan",
		})
	}

	// Membuat direktori "payment" jika belum ada
	paymentDir := "./payment"
	if _, err := os.Stat(paymentDir); os.IsNotExist(err) {
		if err := os.Mkdir(paymentDir, os.ModePerm); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Gagal membuat direktori payment",
			})
		}
	}

	// Menyimpan file gambar pembayaran dengan nama unik
	fileName := fmt.Sprintf("%s%s", generatePaymentID(), ext)
	filePath := filepath.Join(paymentDir, fileName) // Path lokal untuk menyimpan file
	if err := c.SaveFile(file, filePath); err != nil {
		log.Println("Error saving payment image:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal menyimpan gambar pembayaran",
		})
	}

	// URL untuk gambar (gunakan `path.Join` untuk URL-friendly path)
	imageURL := fmt.Sprintf("http://localhost:3000/%s", path.Join("payment", fileName))

	// Membuat data pembayaran baru
	payment := models.Payment{
		PaymentID:     generatePaymentID(),
		CheckoutID:    checkoutID,
		UserID:        c.Locals("userID").(string), // Mendapatkan userID dari token
		PaymentImage:  imageURL,                    // URL publik
		PaymentStatus: "Pending",                   // Status awal adalah Pending
		CreatedAt:     time.Now(),
		ModifiedAt:    time.Now(),
	}

	// Menyimpan data pembayaran ke dalam koleksi `payment`
	collection := config.GetCollection("payment")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, payment)
	if err != nil {
		log.Println("Error saving payment to database:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal menyimpan data pembayaran",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Bukti pembayaran berhasil disimpan",
		"data":    payment,
	})
}

// GetPaymentByCheckoutID handles getting the payment by checkout ID
func GetPaymentByCheckoutID(c *fiber.Ctx) error {
	// Mendapatkan checkoutID dari parameter URL
	checkoutID := c.Params("checkout_id")

	// Mengambil koleksi `payment`
	collection := config.GetCollection("payment")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mencari pembayaran berdasarkan checkoutID
	var payment models.Payment
	err := collection.FindOne(ctx, bson.M{"checkout_id": checkoutID}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Pembayaran tidak ditemukan",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mendapatkan pembayaran",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    payment,
	})
}

// UpdatePaymentStatus handles updating the payment status
func UpdatePaymentStatus(c *fiber.Ctx) error {
	// Parse request body untuk mendapatkan status pembayaran
	var updateRequest struct {
		Status string `json:"status" validate:"required"`
	}

	if err := c.BodyParser(&updateRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Body request tidak valid",
		})
	}

	// Mendapatkan paymentID dari parameter URL
	paymentID := c.Params("payment_id")

	// Mengambil koleksi `payment`
	collection := config.GetCollection("payment")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Memperbarui status pembayaran
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"payment_id": paymentID},
		bson.M{
			"$set": bson.M{
				"payment_status": updateRequest.Status,
				"modified_at":    time.Now(),
			},
		},
	)
	if err != nil {
		log.Printf("Error updating payment status for paymentID %s: %v\n", paymentID, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal memperbarui status pembayaran",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Status pembayaran berhasil diperbarui",
		"data": fiber.Map{
			"payment_id": paymentID,
			"status":     updateRequest.Status,
		},
	})
}

// GetAllPayments handles getting all payments
func GetAllPayments(c *fiber.Ctx) error {
	// Mengambil koleksi `payment`
	collection := config.GetCollection("payment")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Mencari semua pembayaran
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error fetching payments:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mendapatkan pembayaran",
		})
	}
	defer cursor.Close(ctx)

	// Menyimpan semua data pembayaran dalam slice
	var payments []models.Payment
	for cursor.Next(ctx) {
		var payment models.Payment
		if err := cursor.Decode(&payment); err != nil {
			log.Println("Error decoding payment:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Gagal mendecode data pembayaran",
			})
		}
		payments = append(payments, payment)
	}

	// Periksa jika cursor ada kesalahan
	if err := cursor.Err(); err != nil {
		log.Println("Error iterating payments:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal iterasi pembayaran",
		})
	}

	// Mengembalikan semua data pembayaran
	return c.JSON(fiber.Map{
		"success": true,
		"data":    payments,
	})
}

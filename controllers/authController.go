package controllers

import (
	"be-stepup/config"
	"be-stepup/models"
	"context"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Register handles user registration
func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if the email already exists
	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := collection.FindOne(ctx, map[string]string{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Email already in use",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not hash password",
		})
	}
	user.Password = string(hashedPassword)

	// Set Name and CreatedAt
	user.CreatedAt = time.Now() // Menambahkan tanggal pembuatan akun

	// Insert user into database
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user, // Kembalikan informasi pengguna yang baru dibuat
	})
}

// Login handles user authentication
func Login(c *fiber.Ctx) error {
	var loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Fetch user from database
	collection := config.GetCollection("users")
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, map[string]string{"email": loginReq.Email}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Respond with user role
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"role":    user.Role,
	})
}

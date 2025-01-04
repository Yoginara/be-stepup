package main

import (
	"be-stepup/config"
	"be-stepup/controllers"
	"be-stepup/middleware"
	"be-stepup/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Inisialisasi aplikasi Fiber
	app := fiber.New()

	// Middleware untuk logging
	app.Use(logger.New())

	// Middleware untuk mengatasi CORS (Didefinisikan sebelum rute)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1:5500", // Domain frontend Anda
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Menghubungkan ke database MongoDB
	err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.DisconnectDB() // Menutup koneksi ke MongoDB ketika server dimatikan

	// Atur semua rute
	routes.SetupRoutes(app)

	// Tambahkan rute untuk register
	app.Post("/api/register", controllers.Register)

	// Route untuk login
	app.Post("/api/login", controllers.Login)

	// Proteksi endpoint AddToCart dengan middleware JWT
	app.Post("/cart/add", middleware.JWTAuthMiddleware, controllers.AddToCart)

	// Middleware untuk melayani file statis dari folder uploads
	app.Static("/uploads", "./uploads")

	// Menyiapkan channel untuk menangkap sinyal shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Menjalankan server di port tertentu
	port := ":3000" // Sesuaikan dengan port yang diinginkan
	go func() {
		log.Printf("Server is running on http://localhost%s", port)
		if err := app.Listen(port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Menunggu sinyal shutdown
	<-c
	log.Println("Gracefully shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

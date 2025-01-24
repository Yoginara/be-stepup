package routes

import (
	"be-stepup/controllers"
	"be-stepup/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes mengatur semua rute yang digunakan dalam aplikasi
func SetupRoutes(app *fiber.App) {
	// Grup rute untuk produk
	productGroup := app.Group("/api/products")
	productGroup.Get("/", controllers.GetAllProducts)             // Mengambil semua produk
	productGroup.Get("/:id", controllers.GetProductByID)          // Mengambil produk berdasarkan ID
	productGroup.Get("/code/:code", controllers.GetProductByCode) // Mengambil produk berdasarkan kode unik
	productGroup.Post("/", controllers.CreateProduct)             // Membuat produk baru
	productGroup.Put("/:id", controllers.UpdateProduct)           // Memperbarui produk berdasarkan ID
	productGroup.Delete("/:id", controllers.DeleteProduct)        // Menghapus produk berdasarkan ID

	// Rute untuk user
	app.Get("/api/users", controllers.GetAllUsers)
	app.Get("/api/users/:id", controllers.GetUserByID)

	// Rute untuk mengunggah gambar
	app.Post("/api/upload", controllers.UploadImage) // Mengunggah gambar produk

	// Grup rute untuk autentikasi
	authGroup := app.Group("/api/auth")
	authGroup.Post("/login", controllers.Login)       // Login user/admin
	authGroup.Post("/register", controllers.Register) // Registrasi user/admin

	// Rute untuk menghitung jumlah produk dan pengguna
	app.Get("/api/count/products", controllers.GetProductCount)
	app.Get("/api/count/users", controllers.GetUserCount)

	// Rute keranjang belanja (dengan autentikasi)
	app.Post("/api/cart/add", middleware.JWTAuthMiddleware, controllers.AddToCart)                        // Menambahkan item ke keranjang
	app.Delete("/api/cart/remove-single", middleware.JWTAuthMiddleware, controllers.RemoveSingleCartItem) // Menghapus item dari keranjang
	app.Get("/api/cart", middleware.JWTAuthMiddleware, controllers.GetAllCart)                            // Mendapatkan semua item dalam keranjang
	app.Put("/api/cart", middleware.JWTAuthMiddleware, controllers.UpdateCartItem)

	// Rute checkout (dengan autentikasi)
	app.Post("/api/checkout", middleware.JWTAuthMiddleware, controllers.CreateCheckout)              // Membuat checkout baru
	app.Get("/api/checkout/:checkout_id", middleware.JWTAuthMiddleware, controllers.GetCheckoutByID) // Mendapatkan checkout berdasarkan ID
	app.Put("/api/checkout/:checkout_id", controllers.UpdateCheckout)                                // Memperbarui checkout berdasarkan ID
	app.Get("/checkouts", controllers.GetAllCheckout)
	app.Delete("/checkout/:checkout_id", controllers.DeleteCheckout)

	// Rute pembayaran (dengan autentikasi)
	paymentGroup := app.Group("/api/payment", middleware.JWTAuthMiddleware)
	paymentGroup.Post("/:checkout_id", controllers.SavePayment)           // Menyimpan bukti pembayaran
	paymentGroup.Get("/:checkout_id", controllers.GetPaymentByCheckoutID) // Mendapatkan bukti pembayaran berdasarkan checkoutID

	app.Put("/:payment_id/status", controllers.UpdatePaymentStatus) // Memperbarui status pembayaran
	app.Get("/payments", controllers.GetAllPayments)
}

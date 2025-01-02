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
	app.Post("/api/cart/add", middleware.JWTAuthMiddleware, controllers.AddToCart)
	// Proteksi endpoint RemoveSingleCartItem dengan middleware JWT
	app.Delete("/cart/remove-single", middleware.JWTAuthMiddleware, controllers.RemoveSingleCartItem)

	app.Get("/products/:id/sizes", controllers.GetStockBySize)

}

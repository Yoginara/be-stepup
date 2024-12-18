package controllers

import (
	"be-stepup/config"
	"be-stepup/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// Definisikan folder untuk menyimpan file upload
const uploadDir = "uploads"

// GetAllProducts fetches all products from the database
func GetAllProducts(c *fiber.Ctx) error {
	var products []models.Product
	collection := config.GetCollection("products")
	ctx := c.Context()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching products"})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error decoding product"})
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Cursor error"})
	}

	return c.JSON(products)
}

// GetProductByID fetches a product by its ID
func GetProductByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	productID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	collection := config.GetCollection("products")
	var product models.Product
	err = collection.FindOne(c.Context(), bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching product"})
	}

	return c.JSON(product)
}

// CreateProduct creates a new product
func CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	product.Code = "SKU-" + uuid.New().String()[:8]

	// Parsing data produk
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request"})
	}

	// Validasi harga dan stok
	if product.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Price must be greater than zero"})
	}
	if product.Stock < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Stock cannot be negative"})
	}

	// Penanganan gambar
	file, err := c.FormFile("image")
	if err == nil { // Gambar berhasil diterima
		savePath := fmt.Sprintf("./uploads/%s", file.Filename)
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save image"})
		}
		product.ImageURL = fmt.Sprintf("http://localhost:3000/uploads/%s", file.Filename)
	} else {
		product.ImageURL = "" // Kosongkan jika gambar tidak diunggah
	}

	// Simpan produk ke database
	collection := config.GetCollection("products")
	_, err = collection.InsertOne(c.Context(), product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Product created successfully",
		"productID":   product.ID.Hex(),
		"productCode": product.Code,
		"imageURL":    product.ImageURL,
	})
}

// UpdateProduct updates an existing product
// UpdateProduct updates an existing product
func UpdateProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	productID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	// Struktur untuk data yang akan di-update
	var productData struct {
		Name        string  `json:"name"`
		Brand       string  `json:"brand"`
		Category    string  `json:"category"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		Description string  `json:"description"`
		ImageURL    string  `json:"image_url"`
	}

	// Parsing body request
	if err := c.BodyParser(&productData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request"})
	}

	collection := config.GetCollection("products")
	filter := bson.M{"_id": productID}

	// Mendapatkan data produk yang ada
	var existingProduct models.Product
	err = collection.FindOne(c.Context(), filter).Decode(&existingProduct)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	// Jika ImageURL kosong, gunakan URL gambar lama
	if productData.ImageURL == "" {
		productData.ImageURL = existingProduct.ImageURL
	}

	// Validasi input
	if productData.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Price must be greater than zero"})
	}
	if productData.Stock < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Stock cannot be negative"})
	}

	// Update produk di database
	update := bson.M{
		"$set": bson.M{
			"name":        productData.Name,
			"brand":       productData.Brand,
			"category":    productData.Category,
			"price":       productData.Price,
			"stock":       productData.Stock,
			"description": productData.Description,
			"image_url":   productData.ImageURL, // Update image URL
		},
	}

	_, err = collection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
	}

	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

// DeleteProduct deletes a product from the database
func DeleteProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	productID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	collection := config.GetCollection("products")
	filter := bson.M{"_id": productID}

	_, err = collection.DeleteOne(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product"})
	}

	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

// GetProductByCode fetches a product by its unique code
func GetProductByCode(c *fiber.Ctx) error {
	code := c.Params("code")

	collection := config.GetCollection("products")
	var product models.Product
	err := collection.FindOne(c.Context(), bson.M{"code": code}).Decode(&product)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching product"})
	}

	return c.JSON(product)
}

// UploadImage handles uploading an image and returning its URL
func UploadImage(c *fiber.Ctx) error {
	// Ambil file dari request dengan key "image"
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	// Tentukan lokasi penyimpanan file
	savePath := fmt.Sprintf("./uploads/%s", file.Filename)

	// Simpan file ke folder uploads
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// URL file yang bisa diakses
	fileURL := fmt.Sprintf("http://localhost:3000/uploads/%s", file.Filename)

	// Kembalikan URL sebagai respons
	return c.JSON(fiber.Map{"image_url": fileURL})
}

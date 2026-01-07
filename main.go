package main

import (
	"log"
	"os"
	
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	godotenv.Load()
	
	// Koneksi database
	ConnectDB()
	
	// Setup router
	r := gin.Default()
	
	// ===== SERVE STATIC FILES (untuk akses gambar) =====
	r.Static("/uploads", "./uploads")
	
	// Routes PUBLIC (tidak perlu login)
	public := r.Group("/api")
	{
		// Auth routes
		public.POST("/register", Register)
		public.POST("/login", Login)
		
		// Produk routes (bisa diakses tanpa login)
		public.GET("/produk", GetAllProduk)
		public.GET("/produk/:id", GetProduk)
	}
	
	// Routes PROTECTED (harus login)
	protected := r.Group("/api")
	protected.Use(AuthMiddleware())
	{
		// Profile
		protected.GET("/profile", GetProfile)
		
		// Produk management (harus login)
		protected.POST("/produk", CreateProduk)
		protected.PUT("/produk/:id", UpdateProduk)
		protected.DELETE("/produk/:id", DeleteProduk)
		
		// ===== UPLOAD GAMBAR PRODUK =====
		protected.POST("/produk/:id/upload", UploadProdukImage)
		protected.DELETE("/produk/:id/image", DeleteProdukImage)
	}
	
	// Routes ADMIN ONLY
	admin := r.Group("/api/admin")
	admin.Use(AuthMiddleware(), AdminOnly())
	{
		// Route khusus admin bisa ditambah di sini
		admin.GET("/users", GetAllUsers)
	}
	
	// Baca port dari .env
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("🚀 Server running on http://localhost:%s", port)
	r.Run(":" + port)
}

// Function untuk admin get all users
func GetAllUsers(c *gin.Context) {
	var users []User
	DB.Find(&users)
	
	c.JSON(200, gin.H{
		"status": "success",
		"data":   users,
	})
}
package main

import (
	"log"
	"os"
	
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ConnectDB()
	
	r := gin.Default()
	
	r.Static("/uploads", "./uploads")
	
	// Routes PUBLIC
	public := r.Group("/api")
	{
		public.POST("/register", Register)
		public.POST("/login", Login)
		
		public.GET("/produk", GetAllProduk)
		public.GET("/produk/:id", GetProduk)
	}
	
	// Routes PROTECTED
	protected := r.Group("/api")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/profile", GetProfile)
		
		// Produk
		protected.POST("/produk", CreateProduk)
		protected.PUT("/produk/:id", UpdateProduk)
		protected.DELETE("/produk/:id", DeleteProduk)
		
		protected.POST("/produk/:id/upload", UploadProdukImage)
		protected.DELETE("/produk/:id/image", DeleteProdukImage)
		
		protected.GET("/produk/trash/list", GetDeletedProduk)
		protected.POST("/produk/:id/restore", RestoreProduk)
		protected.DELETE("/produk/:id/permanent", PermanentDeleteProduk)
		
		// ===== TRANSAKSI =====
		protected.POST("/transaksi", CreateTransaksi)              // Buat order
		protected.GET("/transaksi/my", GetMyTransaksi)             // Order saya
		protected.GET("/transaksi/:id", GetTransaksi)              // Detail order
		protected.POST("/transaksi/:id/cancel", CancelTransaksi)   // Cancel order
	}
	
	// Routes ADMIN
	admin := r.Group("/api/admin")
	admin.Use(AuthMiddleware(), AdminOnly())
	{
		admin.GET("/users", GetAllUsers)
		
		// Transaksi management (admin)
		admin.GET("/transaksi", GetAllTransaksi)                      // Semua transaksi
		admin.PUT("/transaksi/:id/status", UpdateStatusTransaksi)     // Update status
	}
	
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("🚀 Server running on http://localhost:%s", port)
	r.Run(":" + port)
}

func GetAllUsers(c *gin.Context) {
	var users []User
	DB.Find(&users)
	c.JSON(200, gin.H{"status": "success", "data": users})
}
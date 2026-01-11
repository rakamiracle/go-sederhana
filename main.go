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
	
	// Serve static files
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
		
		// Produk management
		protected.POST("/produk", CreateProduk)
		protected.PUT("/produk/:id", UpdateProduk)
		protected.DELETE("/produk/:id", DeleteProduk) // Soft delete
		
		// Upload
		protected.POST("/produk/:id/upload", UploadProdukImage)
		protected.DELETE("/produk/:id/image", DeleteProdukImage)
		
		// ===== TRASH MANAGEMENT =====
		protected.GET("/produk/trash/list", GetDeletedProduk)      // Lihat trash
		protected.POST("/produk/:id/restore", RestoreProduk)       // Restore
		protected.DELETE("/produk/:id/permanent", PermanentDeleteProduk) // Hapus permanen
	}
	
	// Routes ADMIN
	admin := r.Group("/api/admin")
	admin.Use(AuthMiddleware(), AdminOnly())
	{
		admin.GET("/users", GetAllUsers)
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
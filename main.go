package main

import (
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	// Koneksi database
	ConnectDB()
	
	// Setup router
	r := gin.Default()
	
	// Routes
	api := r.Group("/api")
	{
		api.GET("/produk", GetAllProduk)
		api.GET("/produk/:id", GetProduk)
		api.POST("/produk", CreateProduk)
		api.PUT("/produk/:id", UpdateProduk)
		api.DELETE("/produk/:id", DeleteProduk)
	}
	
	log.Println("🚀 Server running on http://localhost:8080")
	r.Run(":8080")
}
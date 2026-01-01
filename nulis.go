package belajargo

package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Struct untuk tabel produk
type Produk struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var db *gorm.DB

func main() {
	// 1. KONEKSI KE DATABASE
	// Format: username:password@tcp(host:port)/nama_database
	dsn := "root:root123@tcp(localhost:3306)/toko?charset=utf8mb4&parseTime=True"
	
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	
	if err != nil {
		panic("Gagal konek ke database!")
	}
	
	println("✓ Berhasil konek ke database")

	// 2. AUTO CREATE TABLE (kalau belum ada)
	db.AutoMigrate(&Produk{})
	println("✓ Tabel produk siap")

	// 3. SETUP ROUTER (jalur API)
	router := gin.Default()

	// 4. DAFTAR ENDPOINT API
	router.GET("/produk", ambilSemuaProduk)       // Lihat semua
	router.GET("/produk/:id", ambilSatuProduk)    // Lihat 1 produk
	router.POST("/produk", buatProduk)            // Tambah produk
	router.PUT("/produk/:id", updateProduk)       // Edit produk
	router.DELETE("/produk/:id", hapusProduk)     // Hapus produk

	// 5. JALANKAN SERVER
	println("✓ Server jalan di http://localhost:8080")
	router.Run(":8080")
}

// LIHAT SEMUA PRODUK
func ambilSemuaProduk(c *gin.Context) {
	var produk []Produk
	
	// Ambil semua data dari database
	db.Find(&produk)
	
	// Kirim response JSON
	c.JSON(200, gin.H{
		"status": "success",
		"data":   produk,
	})
}

// LIHAT 1 PRODUK
func ambilSatuProduk(c *gin.Context) {
	var produk Produk
	id := c.Param("id")
	
	// Cari produk berdasarkan ID
	if err := db.First(&produk, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	c.JSON(200, gin.H{
		"status": "success",
		"data":   produk,
	})
}

// TAMBAH PRODUK BARU
func buatProduk(c *gin.Context) {
	var produk Produk
	
	// Ambil data JSON dari request
	if err := c.ShouldBindJSON(&produk); err != nil {
		c.JSON(400, gin.H{"error": "Data tidak valid"})
		return
	}
	
	// Simpan ke database
	db.Create(&produk)
	
	c.JSON(201, gin.H{
		"status":  "success",
		"message": "Produk berhasil ditambahkan",
		"data":    produk,
	})
}

// UPDATE PRODUK
func updateProduk(c *gin.Context) {
	var produk Produk
	id := c.Param("id")
	
	// Cek apakah produk ada
	if err := db.First(&produk, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	// Ambil data baru dari request
	var input Produk
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Data tidak valid"})
		return
	}
	
	// Update ke database
	db.Model(&produk).Updates(input)
	
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Produk berhasil diupdate",
		"data":    produk,
	})
}

// HAPUS PRODUK
func hapusProduk(c *gin.Context) {
	var produk Produk
	id := c.Param("id")
	
	// Cek apakah produk ada
	if err := db.First(&produk, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	// Hapus dari database
	db.Delete(&produk)
	
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Produk berhasil dihapus",
	})
}
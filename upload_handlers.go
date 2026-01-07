package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/gin-gonic/gin"
)

// Upload gambar produk
func UploadProdukImage(c *gin.Context) {
	// Ambil ID produk dari URL
	produkID := c.Param("id")
	
	// Cek apakah produk ada
	var produk Produk
	if err := DB.First(&produk, produkID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	// Ambil file dari form
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
		return
	}
	
	// Validasi ukuran file (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file maksimal 5MB"})
		return
	}
	
	// Validasi ekstensi file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format file harus jpg, jpeg, png, atau gif"})
		return
	}
	
	// Generate nama file unik
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("produk_%s_%d%s", produkID, timestamp, ext)
	
	// Path untuk simpan file
	uploadPath := filepath.Join("uploads", "produk", filename)
	
	// Simpan file
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload file"})
		return
	}
	
	// Hapus gambar lama kalau ada
	if produk.ImageURL != "" {
		oldPath := filepath.Join("uploads", "produk", filepath.Base(produk.ImageURL))
		os.Remove(oldPath)
	}
	
	// Update database dengan URL gambar
	imageURL := fmt.Sprintf("/uploads/produk/%s", filename)
	DB.Model(&produk).Update("image_url", imageURL)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Gambar berhasil diupload",
		"data": gin.H{
			"image_url": imageURL,
			"produk":    produk,
		},
	})
}

// Delete gambar produk
func DeleteProdukImage(c *gin.Context) {
	produkID := c.Param("id")
	
	var produk Produk
	if err := DB.First(&produk, produkID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	// Cek apakah ada gambar
	if produk.ImageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Produk tidak punya gambar"})
		return
	}
	
	// Hapus file gambar
	imagePath := filepath.Join("uploads", "produk", filepath.Base(produk.ImageURL))
	if err := os.Remove(imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus file"})
		return
	}
	
	// Update database
	DB.Model(&produk).Update("image_url", "")
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Gambar berhasil dihapus",
	})
}
package main

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type PaginationResponse struct {
	Status     string      `json:"status"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// GET semua produk (hanya yang tidak di-soft delete)
func GetAllProduk(c *gin.Context) {
	var produk []Produk
	var total int64
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	
	offset := (page - 1) * limit
	
	// Query builder - GORM otomatis filter data yang deleted_at = NULL
	query := DB.Model(&Produk{})
	
	// Search
	if search := c.Query("search"); search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}
	
	// Filter by kategori
	if kategori := c.Query("kategori"); kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}
	
	// Filter by harga
	if minHarga := c.Query("min_harga"); minHarga != "" {
		query = query.Where("harga >= ?", minHarga)
	}
	if maxHarga := c.Query("max_harga"); maxHarga != "" {
		query = query.Where("harga <= ?", maxHarga)
	}
	
	// Filter by stok
	if minStok := c.Query("min_stok"); minStok != "" {
		query = query.Where("stok >= ?", minStok)
	}
	
	// Sorting
	sortBy := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")
	
	allowedSorts := map[string]bool{
		"id": true, "nama": true, "harga": true, 
		"stok": true, "created_at": true,
	}
	if !allowedSorts[sortBy] {
		sortBy = "id"
	}
	
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	
	query = query.Order(sortBy + " " + order)
	
	query.Count(&total)
	query.Limit(limit).Offset(offset).Find(&produk)
	
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	
	c.JSON(http.StatusOK, PaginationResponse{
		Status: "success",
		Data:   produk,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// GET 1 produk
func GetProduk(c *gin.Context) {
	var produk Produk
	
	if err := DB.First(&produk, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   produk,
	})
}

// POST produk baru
func CreateProduk(c *gin.Context) {
	var produk Produk
	
	if err := c.ShouldBindJSON(&produk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	DB.Create(&produk)
	
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Produk berhasil ditambahkan",
		"data":    produk,
	})
}

// PUT update produk
func UpdateProduk(c *gin.Context) {
	var produk Produk
	
	if err := DB.First(&produk, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	if err := c.ShouldBindJSON(&produk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	DB.Save(&produk)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Produk berhasil diupdate",
		"data":    produk,
	})
}

// DELETE produk (SOFT DELETE)
func DeleteProduk(c *gin.Context) {
	var produk Produk
	
	// Cari produk
	if err := DB.First(&produk, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	// Soft delete - GORM otomatis set deleted_at
	DB.Delete(&produk)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Produk berhasil dihapus (soft delete)",
	})
}

// GET produk yang sudah di-delete (Trash)
func GetDeletedProduk(c *gin.Context) {
	var produk []Produk
	
	// Unscoped() = tampilkan data yang sudah di-delete
	DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&produk)
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   produk,
	})
}

// RESTORE produk yang di-delete
func RestoreProduk(c *gin.Context) {
	var produk Produk
	
	// Cari di data yang sudah di-delete
	if err := DB.Unscoped().First(&produk, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	// Cek apakah memang sudah di-delete
	if produk.DeletedAt.Time.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Produk tidak dalam status deleted"})
		return
	}
	
	// Restore = set deleted_at = NULL
	DB.Unscoped().Model(&produk).Update("deleted_at", nil)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Produk berhasil di-restore",
		"data":    produk,
	})
}

// PERMANENT DELETE (hapus permanen)
func PermanentDeleteProduk(c *gin.Context) {
	var produk Produk
	
	// Cari di semua data (termasuk yang sudah di-delete)
	if err := DB.Unscoped().First(&produk, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	// Hapus gambar kalau ada
	if produk.ImageURL != "" {
		// Import filepath dan os di atas
		// filepath.Join dan os.Remove
	}
	
	// Permanent delete
	DB.Unscoped().Delete(&produk)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Produk berhasil dihapus permanen",
	})
}
package main

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

// Struct untuk response pagination
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

// GET semua produk dengan pagination, search, filter & sorting
func GetAllProduk(c *gin.Context) {
	var produk []Produk
	var total int64
	
	// ===== PAGINATION =====
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	
	offset := (page - 1) * limit
	
	// ===== QUERY BUILDER =====
	query := DB.Model(&Produk{})
	
	// ===== SEARCH (cari di nama produk) =====
	if search := c.Query("search"); search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}
	
	// ===== FILTER BY KATEGORI =====
	if kategori := c.Query("kategori"); kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}
	
	// ===== FILTER BY HARGA (min & max) =====
	if minHarga := c.Query("min_harga"); minHarga != "" {
		query = query.Where("harga >= ?", minHarga)
	}
	if maxHarga := c.Query("max_harga"); maxHarga != "" {
		query = query.Where("harga <= ?", maxHarga)
	}
	
	// ===== FILTER BY STOK =====
	if minStok := c.Query("min_stok"); minStok != "" {
		query = query.Where("stok >= ?", minStok)
	}
	
	// ===== SORTING =====
	sortBy := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")
	
	// Validasi sort field
	allowedSorts := map[string]bool{
		"id": true, "nama": true, "harga": true, 
		"stok": true, "created_at": true,
	}
	if !allowedSorts[sortBy] {
		sortBy = "id"
	}
	
	// Validasi order
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	
	query = query.Order(sortBy + " " + order)
	
	// ===== COUNT TOTAL =====
	query.Count(&total)
	
	// ===== GET DATA =====
	query.Limit(limit).Offset(offset).Find(&produk)
	
	// ===== HITUNG TOTAL PAGES =====
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	
	// ===== RESPONSE =====
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

// GET 1 produk (tidak berubah)
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

// POST produk baru (tidak berubah)
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

// PUT update produk (tidak berubah)
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

// DELETE produk (tidak berubah)
func DeleteProduk(c *gin.Context) {
	if err := DB.Delete(&Produk{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Produk berhasil dihapus",
	})
}
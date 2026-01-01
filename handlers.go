package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// GET semua produk
func GetAllProduk(c *gin.Context) {
	var produk []Produk
	DB.Find(&produk)
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   produk,
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

// DELETE produk
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
package main

import (
	"fmt"
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
)

// Buat transaksi/order baru
func CreateTransaksi(c *gin.Context) {
	var req CreateTransaksiRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Ambil user ID dari token
	userID, _ := c.Get("userID")
	
	// Mulai transaction database
	tx := DB.Begin()
	
	// Buat kode transaksi unik
	kodeTransaksi := fmt.Sprintf("TRX-%d-%d", userID, time.Now().Unix())
	
	// Hitung total harga
	var totalHarga int
	var transaksiItems []TransaksiItem
	
	for _, item := range req.Items {
		// Cek produk ada dan stok cukup
		var produk Produk
		if err := tx.First(&produk, item.ProdukID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Produk ID %d tidak ditemukan", item.ProdukID)})
			return
		}
		
		if produk.Stok < item.Jumlah {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Stok %s tidak cukup. Tersedia: %d", produk.Nama, produk.Stok),
			})
			return
		}
		
		// Hitung subtotal
		subtotal := produk.Harga * item.Jumlah
		totalHarga += subtotal
		
		// Kurangi stok produk
		produk.Stok -= item.Jumlah
		if err := tx.Save(&produk).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update stok"})
			return
		}
		
		// Tambahkan ke items
		transaksiItems = append(transaksiItems, TransaksiItem{
			ProdukID:    item.ProdukID,
			Jumlah:      item.Jumlah,
			HargaSatuan: produk.Harga,
			Subtotal:    subtotal,
		})
	}
	
	// Buat transaksi
	transaksi := Transaksi{
		UserID:          userID.(uint),
		KodeTransaksi:   kodeTransaksi,
		TotalHarga:      totalHarga,
		Status:          "pending",
		PaymentMethod:   req.PaymentMethod,
		ShippingAddress: req.ShippingAddress,
		Notes:           req.Notes,
		Items:           transaksiItems,
	}
	
	if err := tx.Create(&transaksi).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat transaksi"})
		return
	}
	
	// Commit transaction
	tx.Commit()
	
	// Load relasi untuk response
	DB.Preload("User").Preload("Items.Produk").First(&transaksi, transaksi.ID)
	
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Transaksi berhasil dibuat",
		"data":    transaksi,
	})
}

// Get transaksi by ID
func GetTransaksi(c *gin.Context) {
	var transaksi Transaksi
	
	if err := DB.Preload("User").Preload("Items.Produk").First(&transaksi, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}
	
	// Cek apakah transaksi milik user yang login (kecuali admin)
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	
	if role != "admin" && transaksi.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transaksi,
	})
}

// Get semua transaksi user yang login
func GetMyTransaksi(c *gin.Context) {
	var transaksi []Transaksi
	userID, _ := c.Get("userID")
	
	// Filter berdasarkan status (opsional)
	query := DB.Where("user_id = ?", userID)
	
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	
	query.Preload("Items.Produk").Order("created_at DESC").Find(&transaksi)
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transaksi,
	})
}

// Get semua transaksi (Admin only)
func GetAllTransaksi(c *gin.Context) {
	var transaksi []Transaksi
	
	query := DB.Preload("User").Preload("Items.Produk")
	
	// Filter berdasarkan status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	
	// Filter berdasarkan user
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	
	query.Order("created_at DESC").Find(&transaksi)
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transaksi,
	})
}

// Update status transaksi
func UpdateStatusTransaksi(c *gin.Context) {
	var req UpdateStatusRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validasi status
	validStatus := map[string]bool{
		"pending": true, "paid": true, "processing": true,
		"shipped": true, "completed": true, "cancelled": true,
	}
	
	if !validStatus[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak valid"})
		return
	}
	
	var transaksi Transaksi
	if err := DB.First(&transaksi, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}
	
	// Update status
	transaksi.Status = req.Status
	DB.Save(&transaksi)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Status transaksi berhasil diupdate",
		"data":    transaksi,
	})
}

// Cancel transaksi (kembalikan stok)
func CancelTransaksi(c *gin.Context) {
	var transaksi Transaksi
	
	if err := DB.Preload("Items").First(&transaksi, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}
	
	// Cek apakah transaksi milik user yang login
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	
	if role != "admin" && transaksi.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
		return
	}
	
	// Cek status, hanya bisa cancel kalau pending atau paid
	if transaksi.Status != "pending" && transaksi.Status != "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaksi tidak bisa dibatalkan"})
		return
	}
	
	// Kembalikan stok produk
	tx := DB.Begin()
	
	for _, item := range transaksi.Items {
		var produk Produk
		if err := tx.First(&produk, item.ProdukID).Error; err == nil {
			produk.Stok += item.Jumlah
			tx.Save(&produk)
		}
	}
	
	// Update status jadi cancelled
	transaksi.Status = "cancelled"
	tx.Save(&transaksi)
	
	tx.Commit()
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transaksi berhasil dibatalkan",
		"data":    transaksi,
	})
}
package main

import (
	"time"
	"gorm.io/gorm"
)

// Transaksi/Order
type Transaksi struct {
	ID              uint              `json:"id" gorm:"primaryKey"`
	UserID          uint              `json:"user_id"`
	User            User              `json:"user" gorm:"foreignKey:UserID"`
	KodeTransaksi   string            `json:"kode_transaksi" gorm:"unique"`
	TotalHarga      int               `json:"total_harga"`
	Status          string            `json:"status" gorm:"default:'pending'"` // pending, paid, processing, shipped, completed, cancelled
	PaymentMethod   string            `json:"payment_method"`
	ShippingAddress string            `json:"shipping_address"`
	Notes           string            `json:"notes"`
	Items           []TransaksiItem   `json:"items" gorm:"foreignKey:TransaksiID"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `json:"-" gorm:"index"`
}

// Item dalam transaksi
type TransaksiItem struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	TransaksiID uint    `json:"transaksi_id"`
	ProdukID    uint    `json:"produk_id"`
	Produk      Produk  `json:"produk" gorm:"foreignKey:ProdukID"`
	Jumlah      int     `json:"jumlah"`
	HargaSatuan int     `json:"harga_satuan"`
	Subtotal    int     `json:"subtotal"`
	CreatedAt   time.Time `json:"created_at"`
}

// Request untuk buat transaksi
type CreateTransaksiRequest struct {
	Items           []TransaksiItemRequest `json:"items" binding:"required,min=1"`
	PaymentMethod   string                 `json:"payment_method" binding:"required"`
	ShippingAddress string                 `json:"shipping_address" binding:"required"`
	Notes           string                 `json:"notes"`
}

type TransaksiItemRequest struct {
	ProdukID uint `json:"produk_id" binding:"required"`
	Jumlah   int  `json:"jumlah" binding:"required,min=1"`
}

// Request untuk update status
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
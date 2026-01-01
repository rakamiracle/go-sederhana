package main

import "time"

type Produk struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Nama      string    `json:"nama" binding:"required"`
	Harga     int       `json:"harga" binding:"required"`
	Stok      int       `json:"stok" binding:"required"`
	Kategori  string    `json:"kategori"`
	CreatedAt time.Time `json:"created_at"`
}
package main

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique" binding:"required"`
	Email     string         `json:"email" gorm:"unique" binding:"required,email"`
	Password  string         `json:"-"` // Tidak ditampilkan di JSON
	Nama      string         `json:"nama"`
	Role      string         `json:"role" gorm:"default:'customer'"` // customer atau admin
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Struct untuk request register
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nama     string `json:"nama" binding:"required"`
}

// Struct untuk request login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Struct untuk response login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
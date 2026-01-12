package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  File .env tidak ditemukan, pakai default")
	}
	
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "toko")
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal konek database:", err)
	}
	
	// Auto migrate - TAMBAHKAN Transaksi & TransaksiItem
	db.AutoMigrate(&Produk{}, &User{}, &Transaksi{}, &TransaksiItem{})
	
	DB = db
	log.Println("✅ Database connected")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
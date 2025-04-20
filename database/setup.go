package database

import (
	"go-auth-api/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // Driver SQLite
)

// DB adalah instance koneksi database GORM
var DB *gorm.DB

// ConnectDatabase menginisialisasi koneksi ke database
func ConnectDatabase() {
	// Menggunakan SQLite untuk contoh ini. File test.db akan dibuat.
	// Untuk produksi, gunakan PostgreSQL atau MySQL.
	database, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Gagal terhubung ke database!")
	}

	// Migrasi schema - membuat tabel User jika belum ada
	database.AutoMigrate(&models.User{})

	DB = database
}

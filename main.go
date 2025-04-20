package main

import (
	"go-auth-api/controllers"
	"go-auth-api/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi koneksi database
	database.ConnectDatabase()

	// Inisialisasi Gin router
	r := gin.Default()

	// Grouping route untuk API (contoh: /api/v1)
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.POST("/refresh", controllers.RefreshToken)
		}
	}

	// Jalankan server
	log.Println("Menjalankan server di port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

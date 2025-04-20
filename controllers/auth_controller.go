package controllers

import (
	"go-auth-api/database"
	"go-auth-api/models"
	"go-auth-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// RegisterInput struct untuk binding data register
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginInput struct untuk binding data login
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenInput struct untuk binding data refresh token
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Register menangani pendaftaran user baru
func Register(c *gin.Context) {
	var input RegisterInput

	// Bind JSON ke struct & validasi basic
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hash password"})
		return
	}

	// Buat user baru
	user := models.User{Username: input.Username, Password: hashedPassword}

	// Simpan ke database
	if err := database.DB.Create(&user).Error; err != nil {
		// Cek jika error karena username sudah ada
		// Ini contoh sederhana, penanganan error DB bisa lebih kompleks
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan atau terjadi kesalahan database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil"})
}

// Login menangani login user dan pembuatan token
func Login(c *gin.Context) {
	var input LoginInput
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari user berdasarkan username
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan database"})
		}
		return
	}

	// Cek password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// Generate Access Token
	accessToken, err := utils.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat access token"})
		return
	}

	// Generate Refresh Token
	refreshToken, err := utils.GenerateRefreshToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshToken menangani pembuatan access token baru menggunakan refresh token
func RefreshToken(c *gin.Context) {
	var input RefreshTokenInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Dapatkan secret key untuk refresh token
	_, refreshSecret := utils.GetSecrets()

	// Validasi refresh token
	token, err := utils.ValidateToken(input.RefreshToken, refreshSecret)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token tidak valid atau expired"})
		return
	}

	// Ekstrak username dari refresh token
	username, err := utils.ExtractUsernameFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal membaca data dari refresh token"})
		return
	}

	// Cek apakah user masih ada di DB (opsional tapi bagus)
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User terkait token tidak ditemukan"})
		return
	}

	// Generate access token baru
	newAccessToken, err := utils.GenerateToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat access token baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}

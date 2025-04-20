package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key untuk JWT. Sebaiknya disimpan di environment variable!
// Ambil dari env var atau gunakan default jika tidak ada (TIDAK DISARANKAN UNTUK PRODUKSI)
var jwtSecretKey = []byte(getEnv("JWT_SECRET_KEY", "ini_rahasia_banget_lho"))
var jwtRefreshSecretKey = []byte(getEnv("JWT_REFRESH_SECRET_KEY", "ini_rahasia_refresh_juga"))

// Helper untuk mendapatkan environment variable atau default value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GenerateToken membuat access token JWT baru
func GenerateToken(username string) (string, error) {
	// Durasi token akses (contoh: 15 menit)
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "go-auth-api", // Opsional: identitas penerbit token
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)

	return tokenString, err
}

// GenerateRefreshToken membuat refresh token JWT baru
func GenerateRefreshToken(username string) (string, error) {
	// Durasi refresh token (contoh: 7 hari)
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &jwt.RegisteredClaims{
		Subject:   username, // Biasanya cukup subject saja untuk refresh token
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "go-auth-api-refresh", // Bedakan issuer jika perlu
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtRefreshSecretKey) // Gunakan secret key berbeda

	return tokenString, err
}

// ValidateToken memvalidasi token JWT (bisa untuk access atau refresh, tergantung secret)
func ValidateToken(tokenStr string, secretKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method tidak valid: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err // Termasuk error jika token expired atau signature tidak valid
	}

	return token, nil
}

// ExtractUsernameFromToken mengekstrak username (Subject) dari token yang valid
func ExtractUsernameFromToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("token tidak valid atau claims tidak bisa dibaca")
	}

	// Coba ambil 'sub' (Subject) dari claims
	username, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("claim 'sub' (username) tidak ditemukan atau bukan string")
	}

	return username, nil
}

// GetSecrets digunakan oleh controller untuk akses mudah
func GetSecrets() ([]byte, []byte) {
	return jwtSecretKey, jwtRefreshSecretKey
}

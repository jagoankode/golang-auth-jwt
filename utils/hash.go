package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword mengenkripsi password menggunakan bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // Cost factor 14
	return string(bytes), err
}

// CheckPasswordHash membandingkan password plain text dengan hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // err == nil berarti password cocok
}

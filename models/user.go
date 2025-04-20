package models

import "github.com/jinzhu/gorm"

// User struct merepresentasikan tabel users di database
type User struct {
	gorm.Model        // Includes fields like ID, CreatedAt, UpdatedAt, DeletedAt
	Username   string `gorm:"size:255;not null;unique" json:"username"`
	Password   string `gorm:"size:255;not null;" json:"-"` // Jangan expose password hash di JSON
}

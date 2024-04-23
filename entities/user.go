package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        int            `gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

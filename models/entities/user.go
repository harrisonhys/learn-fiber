package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id               int            `gorm:"primaryKey;autoIncrement"`
	Name             string         `json:"name"`
	Email            string         `json:"email" gorm:"uniqueIndex:idx_email_deleted_at,where:deleted_at IS NULL;size:191"`
	Password         string         `json:"password"`
	Active           bool           `json:"active"`
	EmailConfirmedAt *time.Time     `json:"email_confirmed_at" gorm:"default:null"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

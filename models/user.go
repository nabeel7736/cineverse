package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName     string         `gorm:"type:varchar(255);not null" json:"full_name"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Role         string         `gorm:"type:varchar(50);default:'user';not null" json:"role"`
	IsBlocked    bool           `gorm:"default:false;not null" json:"is_blocked"`
	Address      string         `gorm:"type:text" json:"address"`
	IsVerified   bool           `gorm:"default:false;not null" json:"is_verified"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"update_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

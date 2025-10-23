package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null" json:"user_id"`
	ShowtimeID uint           `gorm:"not null" json:"showtime_id"`
	Seats      string         `gorm:"type:text;not null" json:"seats"` // JSON string of seat numbers
	TotalPrice float64        `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status     string         `gorm:"type:varchar(50);default:confirmed;not null" json:"status"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

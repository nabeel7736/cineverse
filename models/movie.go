package models

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Genre       string         `gorm:"type:varchar(100)" json:"genre"`
	Language    string         `gorm:"type:varchar(100)" json:"language"`
	Duration    int            `gorm:"not null" json:"duration"` // in minutes
	ReleaseDate time.Time      `json:"release_date"`
	PosterURL   string         `gorm:"type:text" json:"poster_url"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

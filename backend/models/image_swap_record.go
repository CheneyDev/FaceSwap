package models

import (
	"time"
)

type ImageSwapRecord struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"not null"`
	OriginalImage string    `gorm:"not null"`
	ImageToSwap   string    `gorm:"not null"`
	ResultImage   string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

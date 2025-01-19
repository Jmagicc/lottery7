package models

import (
	"time"
)

type LotteryResult struct {
	ID        uint      `gorm:"primaryKey"`
	DrawNo    string    `gorm:"type:varchar(20);not null;uniqueIndex;index:idx_draw_no"`
	DrawDate  time.Time `gorm:"not null;index:idx_draw_date"`
	Num1      uint8     `gorm:"not null"`
	Num2      uint8     `gorm:"not null"`
	Num3      uint8     `gorm:"not null"`
	Num4      uint8     `gorm:"not null"`
	Num5      uint8     `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

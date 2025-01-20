package models

import (
	"time"
)

type LicenseKey struct {
	ID        uint      `gorm:"primaryKey"`
	Key       string    `gorm:"column:license_key;type:varchar(8);not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (LicenseKey) TableName() string {
	return "license_keys"
}

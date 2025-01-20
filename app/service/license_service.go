package service

import (
	"crypto/rand"
	"encoding/hex"
	"lottery7/models"
	"time"

	"gorm.io/gorm"
)

type LicenseService struct {
	db *gorm.DB
}

func NewLicenseService(db *gorm.DB) *LicenseService {
	return &LicenseService{db: db}
}

// GenerateKey 生成新的密钥
func (s *LicenseService) GenerateKey() (string, error) {
	// 生成4字节的随机数（生成8位十六进制字符）
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// 转换为16进制字符串
	key := hex.EncodeToString(bytes)

	// 保存到数据库
	licenseKey := models.LicenseKey{
		Key: key,
	}

	if err := s.db.Create(&licenseKey).Error; err != nil {
		return "", err
	}

	return key, nil
}

// ValidateKey 验证密钥并返回创建时间
func (s *LicenseService) ValidateKey(key string) (*time.Time, error) {
	var license models.LicenseKey
	if err := s.db.Where("license_key = ?", key).First(&license).Error; err != nil {
		return nil, err
	}

	return &license.CreatedAt, nil
}

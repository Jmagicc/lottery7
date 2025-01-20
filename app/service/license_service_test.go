package service

import (
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGenerateKey(t *testing.T) {
	// 连接数据库
	dsn := "lottery7:EniQXpY6x8rMjMsz@tcp(192.168.0.200:3306)/lottery7?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	service := NewLicenseService(db)

	// 生成密钥
	key, err := service.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	// 验证密钥
	createdAt, err := service.ValidateKey(key)
	if err != nil {
		t.Fatal(err)
	}

	// 检查创建时间是否在最近1分钟内
	if time.Since(*createdAt) > time.Minute {
		t.Error("Created time is too old")
	}

	t.Logf("Generated key: %s", key)
}

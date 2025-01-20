package main

import (
	"flag"
	"fmt"
	"lottery7/handler"
	"lottery7/models"
	"lottery7/router"
	"lottery7/service"

	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initConfig() {
	// 定义命令行参数
	env := flag.String("env", "dev", "environment: dev or prod")
	flag.Parse()

	// 设置配置文件路径
	viper.SetConfigName(fmt.Sprintf("config.%s", *env))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func main() {
	// 初始化配置
	initConfig()
	// 数据库连接
	dsn := viper.GetString("database.dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移
	db.AutoMigrate(&models.LotteryResult{}, &models.LicenseKey{})

	// 初始化服务和处理器
	lotteryService := service.NewLotteryService(db)
	licenseService := service.NewLicenseService(db)

	lotteryHandler := handler.NewLotteryHandler(lotteryService)
	licenseHandler := handler.NewLicenseHandler(licenseService)

	// 设置路由
	r := router.SetupRouter(lotteryHandler, licenseHandler)

	// 启动服务器
	fmt.Println("Server is running on port 10025")
	r.Run(":10025")
}

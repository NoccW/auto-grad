package main

import (
	"auto-grad-backend/internal/api"
	"auto-grad-backend/internal/config"
	"auto-grad-backend/internal/db"
	"auto-grad-backend/internal/models"
	"auto-grad-backend/internal/services"
	"auto-grad-backend/internal/storage"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	if err := db.InitDB(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 自动迁移数据库表
	if err := autoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化服务
	authService := services.NewAuthService(cfg.JWTSecret)
	ocrService := services.NewBaiduOCRService(cfg.BaiduAPIKey, cfg.BaiduSecretKey)
	aiService := services.NewDeepSeekService(cfg.DeepSeekAPIKey)
	fileStorage := storage.NewFileStorage(cfg.UploadPath)

	// 创建Fiber应用
	app := fiber.New()

	// 设置路由
	api.SetupRoutes(app, authService, ocrService, aiService, fileStorage)

	// 静态文件服务
	app.Static("/uploads", cfg.UploadPath)

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Fatal(app.Listen(":" + cfg.ServerPort))
}

// 自动迁移数据库表
func autoMigrate() error {
	return db.GetDB().AutoMigrate(
		&models.User{},
		&models.GradingRecord{},
		&models.TeacherTask{},
	)
}

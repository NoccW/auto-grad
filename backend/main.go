package main

import (
	"auto-grad-backend/internal/api"
	"auto-grad-backend/internal/config"
	"auto-grad-backend/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	// 加载本地 .env（环境变量优先）
	_ = godotenv.Load(".env")

	cfg := config.LoadConfig()
	pool, err := db.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to init postgres: %v", err)
	}

	// 创建Fiber应用
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// 设置统一系统路由（包含家长端和教师端）
	api.SetupUnifiedRoutes(app, pool)

	// 静态文件服务（模拟）
	app.Static("/uploads", "./uploads")

	// 启动服务器
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 智能改卷统一系统启动成功!")
	log.Printf("📍 服务地址: http://localhost:%s", port)
	log.Printf("🔗 API文档: http://localhost:%s/api", port)
	log.Printf("👨‍🏫 教师端: http://localhost:%s/api/teacher/dashboard", port)
	log.Printf("👨‍👩‍👧‍👦 家长端: http://localhost:%s/api/parent/dashboard", port)
	log.Printf("❤️ 健康检查: http://localhost:%s/health", port)
	log.Printf("🔐 用户登录: http://localhost:%s/api/auth/login", port)

	log.Fatal(app.Listen(":" + port))
}

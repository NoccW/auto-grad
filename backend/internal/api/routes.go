package api

import (
	"auto-grad-backend/internal/services"
	"auto-grad-backend/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App, authService *services.AuthService, ocrService *services.BaiduOCRService, aiService *services.DeepSeekService, fileStorage *storage.FileStorage) {
	// 中间件
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	handler := NewHandler(authService, ocrService, aiService, fileStorage)

	// 健康检查
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API路由组
	api := app.Group("/api")

	// 认证路由（无需JWT验证）
	auth := api.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	// 受保护的路由（需要JWT验证）
	protected := api.Group("/")
	protected.Use(jwtMiddleware)

	// 用户相关
	protected.Get("/auth/me", handler.GetMe)

	// 文件上传
	protected.Post("/upload", handler.UploadFile)

	// 批改相关
	protected.Post("/grading", handler.CreateGrading)
	protected.Get("/grading", handler.GetGradingList)
	protected.Get("/grading/:id", handler.GetGradingDetail)
	protected.Post("/grading/:id/process", handler.ProcessGrading)
	protected.Delete("/grading/:id", handler.DeleteGrading)
}

// JWT中间件
func jwtMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Missing authorization token"})
	}

	// 移除 "Bearer " 前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// 这里需要访问authService来验证token
	// 暂时简化处理
	// TODO: 实现完整的JWT验证

	c.Locals("user", &services.Claims{
		UserID: 1, // 临时硬编码
	})

	return c.Next()
}

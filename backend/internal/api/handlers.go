package api

import (
	"auto-grad-backend/internal/db"
	"auto-grad-backend/internal/models"
	"auto-grad-backend/internal/services"
	"auto-grad-backend/internal/storage"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	authService *services.AuthService
	ocrService  *services.BaiduOCRService
	aiService   *services.DeepSeekService
	fileStorage *storage.FileStorage
}

func NewHandler(authService *services.AuthService, ocrService *services.BaiduOCRService, aiService *services.DeepSeekService, fileStorage *storage.FileStorage) *Handler {
	return &Handler{
		authService: authService,
		ocrService:  ocrService,
		aiService:   aiService,
		fileStorage: fileStorage,
	}
}

// 用户注册
func (h *Handler) Register(c *fiber.Ctx) error {
	type RegisterRequest struct {
		OpenID   string `json:"openId"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserRole string `json:"userRole"` // parent, teacher
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// 检查用户是否已存在
	var existingUser models.User
	if err := db.GetDB().Where("open_id = ?", req.OpenID).First(&existingUser).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "User already exists"})
	}

	// 哈希密码
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// 创建用户
	user := models.User{
		OpenID:       req.OpenID,
		Name:         &req.Name,
		Email:        &req.Email,
		UserRole:     &req.UserRole,
		PasswordHash: &hashedPassword,
		Role:         "user",
		LastSignedIn: time.Now(),
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// 生成JWT token
	token, err := h.authService.GenerateToken(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

// 用户登录
func (h *Handler) Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		OpenID   string `json:"openId"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// 查找用户
	var user models.User
	if err := db.GetDB().Where("open_id = ?", req.OpenID).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// 验证密码
	if user.PasswordHash == nil || !h.authService.CheckPassword(req.Password, *user.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// 生成JWT token
	token, err := h.authService.GenerateToken(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

// 获取用户信息
func (h *Handler) GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*services.Claims)

	var dbUser models.User
	if err := db.GetDB().First(&dbUser, user.UserID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(dbUser)
}

// 文件上传
func (h *Handler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
	}

	// 保存文件
	relativePath, err := h.fileStorage.SaveFile(file, "papers")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save file"})
	}

	return c.JSON(fiber.Map{
		"filename":     file.Filename,
		"relativePath": relativePath,
		"size":         file.Size,
	})
}

// 创建批改记录
func (h *Handler) CreateGrading(c *fiber.Ctx) error {
	user := c.Locals("user").(*services.Claims)

	type CreateGradingRequest struct {
		PaperImageUrl  string `json:"paperImageUrl"`
		AnswerImageUrl string `json:"answerImageUrl,omitempty"`
	}

	var req CreateGradingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// 创建批改记录
	record := models.GradingRecord{
		UserID:        user.UserID,
		PaperImageUrl: req.PaperImageUrl,
		Status:        "pending",
	}

	if req.AnswerImageUrl != "" {
		record.AnswerImageUrl = &req.AnswerImageUrl
	}

	if err := db.GetDB().Create(&record).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create grading record"})
	}

	return c.JSON(record)
}

// 获取批改记录列表
func (h *Handler) GetGradingList(c *fiber.Ctx) error {
	user := c.Locals("user").(*services.Claims)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	var records []models.GradingRecord
	var total int64

	db.GetDB().Model(&models.GradingRecord{}).Where("user_id = ?", user.UserID).Count(&total)

	if err := db.GetDB().Where("user_id = ?", user.UserID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch records"})
	}

	return c.JSON(fiber.Map{
		"records": records,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// 获取批改记录详情
func (h *Handler) GetGradingDetail(c *fiber.Ctx) error {
	user := c.Locals("user").(*services.Claims)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid record ID"})
	}

	var record models.GradingRecord
	if err := db.GetDB().Where("id = ? AND user_id = ?", id, user.UserID).First(&record).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
	}

	return c.JSON(record)
}

// 处理批改（OCR + AI评分）
func (h *Handler) ProcessGrading(c *fiber.Ctx) error {
	user := c.Locals("user").(*services.Claims)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid record ID"})
	}

	var record models.GradingRecord
	if err := db.GetDB().Where("id = ? AND user_id = ?", id, user.UserID).First(&record).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
	}

	// 更新状态为处理中
	record.Status = "processing"
	db.GetDB().Save(&record)

	// OCR 试卷
	paperPath := h.fileStorage.GetFilePath(record.PaperImageUrl)
	paperBytes, err := os.ReadFile(paperPath)
	if err != nil {
		record.Status = "failed"
		db.GetDB().Save(&record)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to read paper image: %v", err)})
	}

	ocrText, err := h.ocrService.RecognizeText(paperBytes)
	if err != nil {
		record.Status = "failed"
		db.GetDB().Save(&record)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("OCR failed: %v", err)})
	}

	// OCR 参考答案（如有上传）
	var referenceAnswer string
	if record.AnswerImageUrl != nil && *record.AnswerImageUrl != "" {
		answerPath := h.fileStorage.GetFilePath(*record.AnswerImageUrl)
		answerBytes, err := os.ReadFile(answerPath)
		if err != nil {
			record.Status = "failed"
			db.GetDB().Save(&record)
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to read answer image: %v", err)})
		}

		refText, err := h.ocrService.RecognizeText(answerBytes)
		if err != nil {
			record.Status = "failed"
			db.GetDB().Save(&record)
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Answer OCR failed: %v", err)})
		}
		referenceAnswer = refText
	}

	// AI 评分
	gradeResult, err := h.aiService.GradePaper(ocrText, referenceAnswer)
	if err != nil {
		record.Status = "failed"
		db.GetDB().Save(&record)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("AI grading failed: %v", err)})
	}

	record.Status = "completed"
	record.AIScore = &gradeResult.Score
	ocrClean := strings.TrimSpace(ocrText)
	record.OCRResult = &ocrClean

	if len(gradeResult.WrongQuestions) > 0 {
		if b, err := json.Marshal(map[string][]string{"questions": gradeResult.WrongQuestions}); err == nil {
			str := string(b)
			record.WrongQuestions = &str
		}
	}
	if len(gradeResult.CorrectAnswers) > 0 {
		if b, err := json.Marshal(map[string][]string{"answers": gradeResult.CorrectAnswers}); err == nil {
			str := string(b)
			record.CorrectAnswers = &str
		}
	}

	db.GetDB().Save(&record)

	return c.JSON(record)
}

// 删除批改记录
func (h *Handler) DeleteGrading(c *fiber.Ctx) error {
	user := c.Locals("user").(*services.Claims)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid record ID"})
	}

	var record models.GradingRecord
	if err := db.GetDB().Where("id = ? AND user_id = ?", id, user.UserID).First(&record).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Record not found"})
	}

	// 删除相关文件
	if record.PaperImageUrl != "" {
		h.fileStorage.DeleteFile(record.PaperImageUrl)
	}
	if record.AnswerImageUrl != nil && *record.AnswerImageUrl != "" {
		h.fileStorage.DeleteFile(*record.AnswerImageUrl)
	}

	// 删除数据库记录
	if err := db.GetDB().Delete(&record).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete record"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Record deleted successfully"})
}

package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 统一系统路由设置 - 支持家长端和教师端
var teacherHandler = &TeacherTaskHandler{}
var gradingStore *GradingStore
var userStore *UserStore
var pgPool *pgxpool.Pool

type GradingRequest struct {
	ID            string   `json:"id"`
	Subject       string   `json:"subject"`
	Images        []string `json:"images"`
	PaperImage    string   `json:"paperImageUrl"`
	AnswerImage   string   `json:"answerImageUrl,omitempty"`
	Description   string   `json:"description"`
	Status        string   `json:"status"`
	Score         int      `json:"score"`
	AiScore       int      `json:"aiScore"`
	TotalScore    int      `json:"totalScore"`
	SubmitTime    string   `json:"submitTime"`
	CreatedAt     string   `json:"createdAt"`
	CompleteTime  string   `json:"completeTime,omitempty"`
	Feedback      string   `json:"feedback,omitempty"`
	OcrResult     string   `json:"ocrResult,omitempty"`
	OwnerUsername string   `json:"ownerUsername,omitempty"`
	OwnerRole     string   `json:"ownerRole,omitempty"`
}

type GradingStore struct {
	pool *pgxpool.Pool
}

type User struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	StudentName string `json:"studentName,omitempty"`
	Class       string `json:"class,omitempty"`
	School      string `json:"school,omitempty"`
}

type UserStore struct {
	pool *pgxpool.Pool
}

func NewGradingStore(pool *pgxpool.Pool) *GradingStore {
	return &GradingStore{pool: pool}
}

func NewUserStore(pool *pgxpool.Pool) *UserStore {
	return &UserStore{pool: pool}
}

func (u *UserStore) key(username, role string) string {
	return fmt.Sprintf("%s:%s", role, username)
}

func (u *UserStore) Get(username, role string) (User, bool) {
	row := u.pool.QueryRow(context.Background(), `SELECT username, password, role, name, email, student_name, class, school FROM users WHERE username=$1 AND role=$2`, username, role)
	var usr User
	if err := row.Scan(&usr.Username, &usr.Password, &usr.Role, &usr.Name, &usr.Email, &usr.StudentName, &usr.Class, &usr.School); err != nil {
		return User{}, false
	}
	return usr, true
}

func (u *UserStore) Create(user User) error {
	_, err := u.pool.Exec(context.Background(), `
INSERT INTO users (username, role, password, name, email, student_name, class, school)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
ON CONFLICT (username, role) DO NOTHING
`, user.Username, user.Role, user.Password, user.Name, user.Email, user.StudentName, user.Class, user.School)
	if err != nil {
		return err
	}
	// check if inserted
	_, ok := u.Get(user.Username, user.Role)
	if ok {
		return nil
	}
	return fmt.Errorf("用户已存在")
}

func (u *UserStore) Update(user User) {
	_, _ = u.pool.Exec(context.Background(), `
INSERT INTO users (username, role, password, name, email, student_name, class, school)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
ON CONFLICT (username, role) DO UPDATE SET
 password=EXCLUDED.password,
 name=EXCLUDED.name,
 email=EXCLUDED.email,
 student_name=EXCLUDED.student_name,
 class=EXCLUDED.class,
 school=EXCLUDED.school;
`, user.Username, user.Role, user.Password, user.Name, user.Email, user.StudentName, user.Class, user.School)
}

func (s *GradingStore) add(req GradingRequest) {
	_, _ = s.pool.Exec(context.Background(), `
INSERT INTO gradings 
  (id, subject, paper_image, answer_image, description, status, score, ai_score, total_score, submit_time, created_at, complete_time, feedback, ocr_result, owner_username, owner_role)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
ON CONFLICT (id) DO UPDATE SET
  subject=excluded.subject,
  paper_image=excluded.paper_image,
  answer_image=excluded.answer_image,
  description=excluded.description,
  status=excluded.status,
  score=excluded.score,
  ai_score=excluded.ai_score,
  total_score=excluded.total_score,
  submit_time=excluded.submit_time,
  created_at=excluded.created_at,
  complete_time=excluded.complete_time,
  feedback=excluded.feedback,
  ocr_result=excluded.ocr_result,
  owner_username=excluded.owner_username,
  owner_role=excluded.owner_role;
`, req.ID, req.Subject, req.PaperImage, req.AnswerImage, req.Description, req.Status, req.Score, req.AiScore, req.TotalScore, parseTime(req.SubmitTime), parseTime(req.CreatedAt), parseTime(req.CompleteTime), req.Feedback, req.OcrResult, req.OwnerUsername, req.OwnerRole)
}

func (s *GradingStore) update(id string, fn func(*GradingRequest)) *GradingRequest {
	item := s.get(id)
	if item == nil {
		return nil
	}
	fn(item)
	_, _ = s.pool.Exec(context.Background(), `
UPDATE gradings SET 
  subject=$2, paper_image=$3, answer_image=$4, description=$5, status=$6, score=$7, ai_score=$8, total_score=$9, submit_time=$10, created_at=$11, complete_time=$12, feedback=$13, ocr_result=$14, owner_username=$15, owner_role=$16
WHERE id=$1
`, item.ID, item.Subject, item.PaperImage, item.AnswerImage, item.Description, item.Status, item.Score, item.AiScore, item.TotalScore, parseTime(item.SubmitTime), parseTime(item.CreatedAt), parseTime(item.CompleteTime), item.Feedback, item.OcrResult, item.OwnerUsername, item.OwnerRole)
	return item
}

func (s *GradingStore) getAll() []GradingRequest {
	rows, err := s.pool.Query(context.Background(), `SELECT id, subject, paper_image, answer_image, description, status, score, ai_score, total_score, submit_time, created_at, complete_time, feedback, ocr_result, owner_username, owner_role FROM gradings ORDER BY submit_time DESC`)
	if err != nil {
		return []GradingRequest{}
	}
	defer rows.Close()
	var res []GradingRequest
	for rows.Next() {
		var g GradingRequest
		var submit, created, complete *time.Time
		_ = rows.Scan(&g.ID, &g.Subject, &g.PaperImage, &g.AnswerImage, &g.Description, &g.Status, &g.Score, &g.AiScore, &g.TotalScore, &submit, &created, &complete, &g.Feedback, &g.OcrResult, &g.OwnerUsername, &g.OwnerRole)
		g.SubmitTime = formatTime(submit)
		g.CreatedAt = formatTime(created)
		g.CompleteTime = formatTime(complete)
		res = append(res, g)
	}
	return res
}

func (s *GradingStore) get(id string) *GradingRequest {
	row := s.pool.QueryRow(context.Background(), `SELECT id, subject, paper_image, answer_image, description, status, score, ai_score, total_score, submit_time, created_at, complete_time, feedback, ocr_result, owner_username, owner_role FROM gradings WHERE id=$1`, id)
	var g GradingRequest
	var submit, created, complete *time.Time
	if err := row.Scan(&g.ID, &g.Subject, &g.PaperImage, &g.AnswerImage, &g.Description, &g.Status, &g.Score, &g.AiScore, &g.TotalScore, &submit, &created, &complete, &g.Feedback, &g.OcrResult, &g.OwnerUsername, &g.OwnerRole); err != nil {
		return nil
	}
	g.SubmitTime = formatTime(submit)
	g.CreatedAt = formatTime(created)
	g.CompleteTime = formatTime(complete)
	return &g
}

func SetupUnifiedRoutes(app *fiber.App, pool *pgxpool.Pool) {
	pgPool = pool
	gradingStore = NewGradingStore(pool)
	userStore = NewUserStore(pool)
	ensureDefaultUsers()
	ensureDefaultUsers()
	// 中间件
	app.Use(cors.New(cors.Config{
		// 允许本地调试来源避免开发时的跨域限制
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
	}))

	// 健康检查
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "智能改卷系统运行正常",
			"version": "1.0.0",
		})
	})

	// API路由组
	api := app.Group("/api")

	// 用户认证相关路由
	auth := api.Group("/auth")
	auth.Get("/me", getUserInfo)
	auth.Post("/login", userLogin)
	auth.Post("/register", userRegister)
	auth.Post("/logout", userLogout)

	// 文件上传
	api.Post("/upload", handleFileUpload)

	// 改卷流程
	grading := api.Group("/grading")
	grading.Get("/", listGradingRequests)
	grading.Post("/", createGradingRequest)
	grading.Get("/:id", getGradingDetail)
	grading.Post("/:id/process", processGradingRequest)

	// 家长端路由
	parent := api.Group("/parent")
	parent.Get("/dashboard", getParentDashboard)
	parent.Post("/submit", submitGradingRequest)
	parent.Get("/results", getParentResults)
	parent.Get("/result/:id", getParentResultDetail)
	parent.Get("/history", getParentHistory)

	// 教师端路由
	teacher := api.Group("/teacher")
	teacher.Get("/dashboard", getTeacherDashboard)
	teacher.Post("/tasks", teacherHandler.CreateTeacherTask)
	teacher.Get("/tasks", teacherHandler.GetTeacherTasks)
	teacher.Get("/tasks/:id", teacherHandler.GetTeacherTask)
	teacher.Post("/tasks/:id/execute", teacherHandler.ExecuteTeacherTask)
	teacher.Post("/tasks/:id/cancel", teacherHandler.CancelTeacherTask)
	teacher.Get("/tasks/:id/status", teacherHandler.GetTaskStatus)
	teacher.Get("/tasks/:id/statistics", teacherHandler.GetTaskStatistics)
	teacher.Get("/tasks/:id/analytics", teacherHandler.GetTaskAnalytics)
	teacher.Delete("/tasks/:id", teacherHandler.DeleteTeacherTask)
	teacher.Get("/history", getTeacherHistory)

	// 管理员路由
	admin := api.Group("/admin")
	admin.Get("/users", getAllUsers)
	admin.Get("/tasks", getAllTasks)
	admin.Get("/statistics", getSystemStatistics)

	// 用户资料
	auth.Put("/profile", updateProfile)

	// 学生信息
	parent.Get("/student", getStudentInfo)
	parent.Put("/student", updateStudentInfo)
}

// 用户认证相关函数
func getUserInfo(c *fiber.Ctx) error {
	user := currentUser(c)
	return c.JSON(fiber.Map{
		"userId":      1,
		"openId":      user.Username,
		"role":        user.Role,
		"userRole":    user.Role,
		"username":    user.Name,
		"email":       user.Email,
		"studentName": user.StudentName,
		"class":       user.Class,
		"school":      user.School,
	})
}

func userLogin(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	role := req.Role
	if role == "" {
		role = "parent"
	}

	user, ok := userStore.Get(req.Username, role)
	if !ok || user.Password != req.Password {
		return c.Status(401).JSON(fiber.Map{"error": "用户名或密码错误"})
	}

	token := "mock_token_" + user.Username + "_" + user.Role
	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"userId":      1,
			"username":    user.Name,
			"role":        user.Role,
			"userRole":    user.Role,
			"studentName": user.StudentName,
			"class":       user.Class,
			"school":      user.School,
			"email":       user.Email,
		},
	})
}

func userRegister(c *fiber.Ctx) error {
	type RegisterRequest struct {
		OpenId          string `json:"openId"`
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
		UserRole        string `json:"userRole"`
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	if req.OpenId == "" || req.Password == "" || req.ConfirmPassword == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
	}

	if req.Password != req.ConfirmPassword {
		return c.Status(400).JSON(fiber.Map{"error": "Passwords do not match"})
	}

	role := req.UserRole
	if role == "" {
		role = "parent"
	}

	newUser := User{
		Username:    req.OpenId,
		Password:    req.Password,
		Role:        role,
		Name:        req.Name,
		Email:       req.Email,
		StudentName: req.Name + "的孩子",
	}
	if err := userStore.Create(newUser); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	token := "mock_token_" + newUser.Username + "_" + newUser.Role
	user := fiber.Map{
		"userId":      time.Now().Unix(),
		"openId":      newUser.Username,
		"username":    newUser.Name,
		"email":       newUser.Email,
		"role":        newUser.Role,
		"userRole":    newUser.Role,
		"studentName": newUser.StudentName,
	}

	return c.JSON(fiber.Map{
		"message": "注册成功",
		"token":   token,
		"user":    user,
	})
}

func userLogout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "退出登录成功",
	})
}

func handleFileUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "文件上传失败"})
	}

	filename := fmt.Sprintf("%s-%s%s",
		strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename)),
		time.Now().Format("20060102150405"),
		filepath.Ext(file.Filename),
	)
	saveDir := filepath.Join("uploads", "papers")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建上传目录失败"})
	}
	if err := c.SaveFile(file, filepath.Join(saveDir, filename)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "保存文件失败"})
	}

	return c.JSON(fiber.Map{
		"message":      "文件上传成功",
		"filename":     filename,
		"size":         file.Size,
		"url":          "/uploads/papers/" + filename,
		"relativePath": "papers/" + filename,
	})
}

// 家长端功能
func getParentDashboard(c *fiber.Ctx) error {
	user := currentUser(c)
	results := gradingStore.getAll()
	recent := []fiber.Map{}
	totalSubmissions := len(results)
	completed := 0
	totalScore := 0
	for _, r := range results {
		if r.Status == "completed" {
			completed++
			totalScore += r.Score
		}
		if len(recent) < 5 {
			recent = append(recent, fiber.Map{
				"id":         r.ID,
				"subject":    r.Subject,
				"score":      r.Score,
				"totalScore": r.TotalScore,
				"submitTime": r.SubmitTime,
				"status":     r.Status,
			})
		}
	}

	avg := 0.0
	if completed > 0 {
		avg = float64(totalScore) / float64(completed)
	}

	return c.JSON(fiber.Map{
		"studentInfo": fiber.Map{
			"name":   firstNonEmpty(user.StudentName, "李小明"),
			"class":  firstNonEmpty(user.Class, "三年级一班"),
			"school": firstNonEmpty(user.School, "示例小学"),
		},
		"recentResults": recent,
		"statistics": fiber.Map{
			"totalSubmissions": totalSubmissions,
			"averageScore":     avg,
			"completedTasks":   completed,
			"pendingTasks":     totalSubmissions - completed,
		},
	})
}

func submitGradingRequest(c *fiber.Ctx) error {
	type SubmitRequest struct {
		Subject     string   `json:"subject"`
		Images      []string `json:"images"`
		Description string   `json:"description"`
	}

	var req SubmitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	return createGradingInternal(c, req.Subject, req.Images, req.Description)
}

func getParentResults(c *fiber.Ctx) error {
	results := gradingStore.getAll()
	resp := []fiber.Map{}
	for _, r := range results {
		resp = append(resp, fiber.Map{
			"id":           r.ID,
			"subject":      r.Subject,
			"score":        r.Score,
			"totalScore":   r.TotalScore,
			"submitTime":   r.SubmitTime,
			"completeTime": r.CompleteTime,
			"status":       r.Status,
			"feedback":     r.Feedback,
			"ocrResult":    r.OcrResult,
		})
	}

	return c.JSON(fiber.Map{
		"results": resp,
		"total":   len(resp),
		"page":    1,
		"limit":   10,
	})
}

func getParentResultDetail(c *fiber.Ctx) error {
	resultId := c.Params("id")

	item := gradingStore.get(resultId)
	if item == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Result not found"})
	}

	return c.JSON(fiber.Map{
		"id":           item.ID,
		"subject":      item.Subject,
		"score":        item.Score,
		"totalScore":   item.TotalScore,
		"submitTime":   item.SubmitTime,
		"completeTime": item.CompleteTime,
		"status":       item.Status,
		"feedback":     item.Feedback,
		"ocrResult":    item.OcrResult,
		"images":       item.Images,
		"details":      []fiber.Map{},
	})
}

func getParentHistory(c *fiber.Ctx) error {
	items := gradingStore.getAll()
	history := []fiber.Map{}
	total := 0
	scoreSum := 0
	for _, r := range items {
		history = append(history, fiber.Map{
			"id":        r.ID,
			"subject":   r.Subject,
			"score":     r.Score,
			"date":      r.SubmitTime,
			"status":    r.Status,
			"feedback":  r.Feedback,
			"ocrResult": r.OcrResult,
		})
		if r.Status == "completed" {
			scoreSum += r.Score
			total++
		}
	}

	avg := 0.0
	if total > 0 {
		avg = float64(scoreSum) / float64(total)
	}

	return c.JSON(fiber.Map{
		"history": history,
		"statistics": fiber.Map{
			"totalTasks":   len(items),
			"averageScore": avg,
			"bestSubject":  "",
			"improvement":  "",
		},
	})
}

// 教师端功能
func getTeacherDashboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"teacherInfo": fiber.Map{
			"name":    "张老师",
			"school":  "示例小学",
			"subject": "数学",
			"class":   "三年级一班",
		},
		"taskSummary": fiber.Map{
			"totalTasks":      5,
			"runningTasks":    1,
			"completedTasks":  3,
			"failedTasks":     1,
			"totalPapers":     450,
			"completedPapers": 380,
		},
		"recentTasks": []fiber.Map{
			{"id": "task_1640996805", "targetUrl": "https://www.7net.cc", "account": "123123", "status": "completed", "totalPapers": 100, "completedPapers": 95, "averageScore": 78.5, "createdAt": "2024-01-17T10:30:00Z"},
			{"id": "task_1640996806", "targetUrl": "https://www.7net.cc", "account": "test123", "status": "running", "totalPapers": 200, "completedPapers": 120, "averageScore": 82.3, "createdAt": "2024-01-17T09:00:00Z"},
		},
	})
}

func getTeacherHistory(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"tasks": []fiber.Map{
			{"id": "task_1640996805", "targetUrl": "https://www.7net.cc", "account": "123123", "status": "completed", "totalPapers": 100, "completedPapers": 95, "failedPapers": 5, "averageScore": 78.5, "createdAt": "2024-01-17T10:30:00Z", "updatedAt": "2024-01-17T11:30:00Z"},
			{"id": "task_1640996806", "targetUrl": "https://www.7net.cc", "account": "test123", "status": "running", "totalPapers": 200, "completedPapers": 120, "failedPapers": 5, "averageScore": 82.3, "createdAt": "2024-01-17T09:00:00Z", "updatedAt": "2024-01-17T09:00:00Z"},
		},
		"total": 2,
		"page":  1,
		"limit": 10,
	})
}

// 管理员功能
func getAllUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"users": []fiber.Map{
			{"id": 1, "username": "张老师", "role": "teacher", "email": "teacher@example.com", "status": "active"},
			{"id": 2, "username": "李家长", "role": "parent", "email": "parent@example.com", "status": "active"},
		},
		"total": 2,
	})
}

func getAllTasks(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"tasks": []fiber.Map{
			{"id": "task_001", "creator": "张老师", "type": "batch_grading", "status": "completed", "progress": 100},
		},
		"total": 1,
	})
}

func getSystemStatistics(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"users": fiber.Map{
			"total":    150,
			"teachers": 25,
			"parents":  125,
		},
		"tasks": fiber.Map{
			"total":     1250,
			"completed": 1100,
			"running":   50,
			"failed":    100,
		},
		"performance": fiber.Map{
			"avgProcessingTime": "3.5分钟",
			"successRate":       "88%",
			"dailyTasks":        45,
		},
	})
}

// 兼容性函数
func NewTeacherTask(c *fiber.Ctx) error {
	return teacherHandler.CreateTeacherTask(c)
}

func GetTeacherTasks(c *fiber.Ctx) error {
	return teacherHandler.GetTeacherTasks(c)
}

func GetTeacherTask(c *fiber.Ctx) error {
	return teacherHandler.GetTeacherTask(c)
}

func ExecuteTeacherTask(c *fiber.Ctx) error {
	return teacherHandler.ExecuteTeacherTask(c)
}

func CancelTeacherTask(c *fiber.Ctx) error {
	return teacherHandler.CancelTeacherTask(c)
}

func GetTaskStatus(c *fiber.Ctx) error {
	return teacherHandler.GetTaskStatus(c)
}

func GetTaskStatistics(c *fiber.Ctx) error {
	return teacherHandler.GetTaskStatistics(c)
}

func GetTaskAnalytics(c *fiber.Ctx) error {
	return teacherHandler.GetTaskAnalytics(c)
}

func DeleteTeacherTask(c *fiber.Ctx) error {
	return teacherHandler.DeleteTeacherTask(c)
}

// 改卷统一处理
func createGradingInternal(c *fiber.Ctx, subject string, images []string, desc string) error {
	user := currentUser(c)
	if subject == "" {
		subject = "未指定科目"
	}
	id := fmt.Sprintf("grading_%d", time.Now().UnixNano())
	now := time.Now().Format(time.RFC3339)
	paperImg := ""
	if len(images) > 0 {
		paperImg = images[0]
	}
	item := GradingRequest{
		ID:            id,
		Subject:       subject,
		Images:        images,
		PaperImage:    paperImg,
		Description:   desc,
		Status:        "processing",
		Score:         0,
		AiScore:       0,
		TotalScore:    100,
		SubmitTime:    now,
		CreatedAt:     now,
		OwnerUsername: user.Username,
		OwnerRole:     user.Role,
	}
	gradingStore.add(item)

	go runGradingPipeline(id)

	return c.JSON(fiber.Map{
		"id":            id,
		"status":        item.Status,
		"message":       "改卷请求已提交，正在处理",
		"submitTime":    now,
		"estimatedTime": "5-10分钟",
	})
}

func listGradingRequests(c *fiber.Ctx) error {
	items := gradingStore.getAll()
	records := []fiber.Map{}
	for _, it := range items {
		records = append(records, fiber.Map{
			"id":             it.ID,
			"status":         it.Status,
			"aiScore":        it.Score,
			"createdAt":      firstNonEmpty(it.CreatedAt, it.SubmitTime),
			"paperImageUrl":  it.PaperImage,
			"answerImageUrl": it.AnswerImage,
		})
	}
	return c.JSON(fiber.Map{
		"items":   items,
		"records": records,
		"total":   len(items),
		"page":    1,
		"limit":   10,
	})
}

func createGradingRequest(c *fiber.Ctx) error {
	type Req struct {
		Subject     string   `json:"subject"`
		Images      []string `json:"images"`
		Description string   `json:"description"`
		PaperImage  string   `json:"paperImageUrl"`
		AnswerImage string   `json:"answerImageUrl"`
	}
	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if req.PaperImage != "" {
		req.Images = append(req.Images, req.PaperImage)
	}
	user := currentUser(c)
	item := GradingRequest{
		ID:            fmt.Sprintf("grading_%d", time.Now().UnixNano()),
		Subject:       firstNonEmpty(req.Subject, "未指定科目"),
		Images:        req.Images,
		PaperImage:    req.PaperImage,
		AnswerImage:   req.AnswerImage,
		Description:   req.Description,
		Status:        "processing",
		Score:         0,
		AiScore:       0,
		TotalScore:    100,
		SubmitTime:    time.Now().Format(time.RFC3339),
		CreatedAt:     time.Now().Format(time.RFC3339),
		OwnerUsername: user.Username,
		OwnerRole:     user.Role,
	}
	if user.Role == "parent" && user.StudentName != "" {
		item.Description = firstNonEmpty(item.Description, user.StudentName)
	}
	gradingStore.add(item)
	go runGradingPipeline(item.ID)
	return c.JSON(item)
}

func getGradingDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	item := gradingStore.get(id)
	if item == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	return c.JSON(item)
}

func processGradingRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	updated := gradingStore.update(id, func(r *GradingRequest) {
		r.Status = "processing"
		r.AiScore = 0
		r.CompleteTime = ""
		r.Feedback = "已提交，等待真实评分处理"
	})
	if updated == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	go runGradingPipeline(id)
	return c.JSON(updated)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func runGradingPipeline(id string) {
	req := gradingStore.get(id)
	if req == nil {
		return
	}

	// locate paper image path
	imagePath := ""
	if req.PaperImage != "" {
		imagePath = filepath.Join("uploads", req.PaperImage)
	} else if len(req.Images) > 0 {
		imagePath = filepath.Join("uploads", req.Images[0])
	}
	updateWithError := func(msg string) {
		gradingStore.update(id, func(r *GradingRequest) {
			r.Status = "failed"
			r.Feedback = msg
		})
		log.Printf("[grading:%s] failed: %s", id, msg)
	}

	if imagePath == "" {
		updateWithError("未找到试卷图片路径")
		return
	}

	imgBytes, err := os.ReadFile(imagePath)
	if err != nil {
		updateWithError("读取试卷图片失败: " + err.Error())
		return
	}

	imageBase64 := base64.StdEncoding.EncodeToString(imgBytes)

	ocrText, err := callBaiduOCR(imageBase64)
	if err != nil {
		updateWithError("OCR 识别失败: " + err.Error())
		return
	}

	score, feedback, err := callDeepSeekScore(ocrText, req.Subject)
	if err != nil {
		updateWithError("DeepSeek 评分失败: " + err.Error())
		return
	}

	gradingStore.update(id, func(r *GradingRequest) {
		r.Status = "completed"
		r.AiScore = score
		r.Score = score
		r.Feedback = feedback
		r.OcrResult = ocrText
		r.CompleteTime = time.Now().Format(time.RFC3339)
	})
	log.Printf("[grading:%s] completed. score=%d", id, score)
}

// 用户工具
func currentUser(c *fiber.Ctx) User {
	token := c.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	token = strings.TrimPrefix(token, "mock_token_")
	parts := strings.Split(token, "_")
	if len(parts) >= 2 {
		username := strings.Join(parts[:len(parts)-1], "_")
		role := parts[len(parts)-1]
		if user, ok := userStore.Get(username, role); ok {
			return user
		}
	}
	// 默认家长
	user, _ := userStore.Get("123123", "parent")
	return user
}

func callBaiduOCR(imageBase64 string) (string, error) {
	apiKey := os.Getenv("BAIDU_API_KEY")
	secretKey := os.Getenv("BAIDU_SECRET_KEY")
	if apiKey == "" || secretKey == "" {
		return "", fmt.Errorf("缺少 BAIDU_API_KEY/BAIDU_SECRET_KEY")
	}

	tokenURL := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", apiKey, secretKey)
	resp, err := http.Post(tokenURL, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil || tokenResp.AccessToken == "" {
		return "", fmt.Errorf("获取百度 token 失败: %s", string(body))
	}

	ocrURL := "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic?access_token=" + tokenResp.AccessToken
	form := "image=" + urlEncode(imageBase64) + "&language_type=CHN_ENG"
	ocrResp, err := http.Post(ocrURL, "application/x-www-form-urlencoded", strings.NewReader(form))
	if err != nil {
		return "", err
	}
	defer ocrResp.Body.Close()
	ocrBody, _ := io.ReadAll(ocrResp.Body)

	var ocrData struct {
		WordsResult []struct {
			Words string `json:"words"`
		} `json:"words_result"`
		ErrorMsg string `json:"error_msg"`
	}
	if err := json.Unmarshal(ocrBody, &ocrData); err != nil {
		return "", fmt.Errorf("解析 OCR 响应失败: %s", err.Error())
	}
	if len(ocrData.WordsResult) == 0 {
		if ocrData.ErrorMsg != "" {
			return "", fmt.Errorf("%s", ocrData.ErrorMsg)
		}
		return "", fmt.Errorf("未识别到文本")
	}

	var lines []string
	for _, w := range ocrData.WordsResult {
		lines = append(lines, w.Words)
	}
	return strings.Join(lines, "\n"), nil
}

func callDeepSeekScore(text, subject string) (int, string, error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return 0, "", fmt.Errorf("缺少 DEEPSEEK_API_KEY")
	}

	prompt := fmt.Sprintf("你是一名阅卷老师，请根据学生答案给出0-100的分数并简要反馈。\n【科目】%s\n【学生答案】%s\n输出格式：分数（数字）+简短中文反馈。", subject, text)
	payload := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
	}
	bodyBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return 0, "", fmt.Errorf("DeepSeek API 返回错误: %s", string(respBody))
	}

	var dsResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &dsResp); err != nil {
		return 0, "", fmt.Errorf("解析 DeepSeek 响应失败: %v", err)
	}
	if len(dsResp.Choices) == 0 {
		return 0, "", fmt.Errorf("DeepSeek 无返回内容")
	}
	content := dsResp.Choices[0].Message.Content
	score := extractFirstNumber(content)
	feedback := content
	return score, feedback, nil
}

func extractFirstNumber(s string) int {
	num := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			num = num*10 + int(s[i]-'0')
			// break after first contiguous number
			j := i + 1
			for j < len(s) && s[j] >= '0' && s[j] <= '9' {
				num = num*10 + int(s[j]-'0')
				j++
			}
			return num
		}
	}
	return 0
}

func urlEncode(s string) string {
	return url.QueryEscape(s)
}

func parseTime(t string) *time.Time {
	if t == "" {
		return nil
	}
	tt, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return nil
	}
	return &tt
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func ensureDefaultUsers() {
	_ = userStore.Create(User{
		Username:    "123123",
		Password:    "123123",
		Role:        "parent",
		Name:        "李家长",
		Email:       "parent@example.com",
		StudentName: "李小明",
		Class:       "三年级一班",
		School:      "示例小学",
	})
	_ = userStore.Create(User{
		Username: "123123",
		Password: "123123",
		Role:     "teacher",
		Name:     "张老师",
		Email:    "teacher@example.com",
		School:   "示例小学",
	})
}

// 资料与学生信息
func updateProfile(c *fiber.Ctx) error {
	user := currentUser(c)
	type Req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	userStore.Update(user)
	return c.JSON(fiber.Map{"message": "资料已更新", "user": user})
}

func getStudentInfo(c *fiber.Ctx) error {
	user := currentUser(c)
	return c.JSON(fiber.Map{
		"name":   user.StudentName,
		"class":  user.Class,
		"school": user.School,
	})
}

func updateStudentInfo(c *fiber.Ctx) error {
	user := currentUser(c)
	type Req struct {
		Name   string `json:"name"`
		Class  string `json:"class"`
		School string `json:"school"`
	}
	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if req.Name != "" {
		user.StudentName = req.Name
	}
	if req.Class != "" {
		user.Class = req.Class
	}
	if req.School != "" {
		user.School = req.School
	}
	userStore.Update(user)
	return c.JSON(fiber.Map{"message": "学生信息已更新", "studentInfo": req})
}

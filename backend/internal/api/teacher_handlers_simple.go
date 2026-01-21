package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"sync"
	"time"
)

type TeacherTask struct {
	ID              string  `json:"id"`
	TargetURL       string  `json:"targetUrl"`
	Account         string  `json:"account"`
	Password        string  `json:"password"`
	Status          string  `json:"status"`
	TotalPapers     int     `json:"totalPapers"`
	CompletedPapers int     `json:"completedPapers"`
	FailedPapers    int     `json:"failedPapers"`
	AverageScore    float64 `json:"averageScore"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

type TeacherTaskHandler struct {
	mu    sync.Mutex
	tasks []TeacherTask
}

func (h *TeacherTaskHandler) CreateTeacherTask(c *fiber.Ctx) error {
	type CreateTaskRequest struct {
		TargetURL  string `json:"targetUrl"`
		Account    string `json:"account"`
		Password   string `json:"password"`
		PaperLimit int    `json:"paperLimit"`
	}

	var req CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "Invalid request format"})
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	task := TeacherTask{
		ID:              fmt.Sprintf("task_%d", time.Now().UnixNano()),
		TargetURL:       req.TargetURL,
		Account:         req.Account,
		Password:        req.Password,
		Status:          "pending",
		TotalPapers:     0,
		CompletedPapers: 0,
		FailedPapers:    0,
		AverageScore:    0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	h.mu.Lock()
	h.tasks = append(h.tasks, task)
	h.mu.Unlock()

	return c.JSON(fiber.Map{
		"id":      task.ID,
		"message": "Task created successfully",
		"task":    task,
	})
}

func (h *TeacherTaskHandler) GetTeacherTasks(c *fiber.Ctx) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return c.JSON(fiber.Map{
		"tasks": h.tasks,
		"total": len(h.tasks),
		"page":  1,
		"limit": 10,
	})
}

func (h *TeacherTaskHandler) GetTeacherTask(c *fiber.Ctx) error {
	taskID := c.Params("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	task := h.findTask(taskID)
	if task == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	statistics := map[string]interface{}{
		"taskId":          task.ID,
		"totalPapers":     task.TotalPapers,
		"completedPapers": task.CompletedPapers,
		"failedPapers":    task.FailedPapers,
		"averageScore":    task.AverageScore,
	}

	return c.JSON(fiber.Map{
		"task":       task,
		"executions": []map[string]interface{}{},
		"statistics": statistics,
	})
}

func (h *TeacherTaskHandler) ExecuteTeacherTask(c *fiber.Ctx) error {
	taskID := c.Params("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	if task := h.findTask(taskID); task != nil {
		task.Status = "running"
		task.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	return c.JSON(fiber.Map{
		"message": "Task execution started",
		"taskId":  taskID,
		"status":  "running",
	})
}

func (h *TeacherTaskHandler) CancelTeacherTask(c *fiber.Ctx) error {
	taskID := c.Params("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	if task := h.findTask(taskID); task != nil {
		task.Status = "cancelled"
		task.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	return c.JSON(fiber.Map{
		"message": "Task cancelled",
		"taskId":  taskID,
		"status":  "cancelled",
	})
}

func (h *TeacherTaskHandler) GetTaskStatus(c *fiber.Ctx) error {
	taskID := c.Params("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	task := h.findTask(taskID)

	return c.JSON(fiber.Map{
		"taskId":          taskID,
		"status":          statusOrPending(task),
		"totalPapers":     intValue(task, func(t *TeacherTask) int { return t.TotalPapers }),
		"completedPapers": intValue(task, func(t *TeacherTask) int { return t.CompletedPapers }),
		"failedPapers":    intValue(task, func(t *TeacherTask) int { return t.FailedPapers }),
		"averageScore":    floatValue(task, func(t *TeacherTask) float64 { return t.AverageScore }),
		"progress":        "0%",
		"message":         "任务执行中...",
		"timestamp":       time.Now().Format("2006-01-02 15:04:05"),
	})
}

func (h *TeacherTaskHandler) GetTaskStatistics(c *fiber.Ctx) error {
	taskID := c.Params("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	task := h.findTask(taskID)

	return c.JSON(fiber.Map{
		"taskId":          taskID,
		"totalPapers":     intValue(task, func(t *TeacherTask) int { return t.TotalPapers }),
		"completedPapers": intValue(task, func(t *TeacherTask) int { return t.CompletedPapers }),
		"failedPapers":    intValue(task, func(t *TeacherTask) int { return t.FailedPapers }),
		"averageScore":    floatValue(task, func(t *TeacherTask) float64 { return t.AverageScore }),
		"passRate":        0.0,
	})
}

func (h *TeacherTaskHandler) GetTaskAnalytics(c *fiber.Ctx) error {
	taskID := c.Params("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	task := h.findTask(taskID)

	return c.JSON(fiber.Map{
		"taskId": taskID,
		"trend":  []map[string]interface{}{},
		"ranking": []map[string]interface{}{
			{"taskId": taskID, "averageScore": floatValue(task, func(t *TeacherTask) float64 { return t.AverageScore }), "completionRate": 0},
		},
	})
}

func (h *TeacherTaskHandler) DeleteTeacherTask(c *fiber.Ctx) error {
	taskID := c.Params("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	for i := range h.tasks {
		if h.tasks[i].ID == taskID {
			h.tasks = append(h.tasks[:i], h.tasks[i+1:]...)
			break
		}
	}

	return c.JSON(fiber.Map{
		"message": "Task deleted successfully",
		"taskId":  taskID,
		"status":  "deleted",
	})
}

func (h *TeacherTaskHandler) findTask(taskID string) *TeacherTask {
	for i := range h.tasks {
		if h.tasks[i].ID == taskID {
			return &h.tasks[i]
		}
	}
	return nil
}

func statusOrPending(t *TeacherTask) string {
	if t == nil || t.Status == "" {
		return "pending"
	}
	return t.Status
}

func intValue(t *TeacherTask, f func(*TeacherTask) int) int {
	if t == nil {
		return 0
	}
	return f(t)
}

func floatValue(t *TeacherTask, f func(*TeacherTask) float64) float64 {
	if t == nil {
		return 0
	}
	return f(t)
}

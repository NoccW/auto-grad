package services

import (
	"fmt"
	"sync"
	"time"
)

type AutomationService struct {
	taskManager *TaskManager
	mutex       sync.RWMutex
}

type TaskManager struct {
	tasks map[string]*TaskStatus
	mutex sync.RWMutex
}

type TaskStatus struct {
	TaskID          string    `json:"taskId"`
	Status          string    `json:"status"`
	TotalPapers     int       `json:"totalPapers"`
	CompletedPapers int       `json:"completedPapers"`
	FailedPapers    int       `json:"failedPapers"`
	AverageScore    float64   `json:"averageScore"`
	CurrentPaper    int       `json:"currentPaper"`
	Message         string    `json:"message"`
	StartTime       time.Time `json:"startTime"`
	LastUpdateTime  time.Time `json:"lastUpdateTime"`
}

// 创建新的自动化服务
func NewAutomationService() *AutomationService {
	return &AutomationService{
		taskManager: NewTaskManager(),
	}
}

// 创建新的任务管理器
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*TaskStatus),
	}
}

// 开始执行任务
func (s *AutomationService) StartTask(taskID string, totalPapers int, targetURL, account, password string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 创建任务状态
	status := &TaskStatus{
		TaskID:          taskID,
		Status:          "running",
		TotalPapers:     totalPapers,
		CompletedPapers: 0,
		FailedPapers:    0,
		AverageScore:    0,
		CurrentPaper:    0,
		Message:         "任务开始执行...",
		StartTime:       time.Now(),
		LastUpdateTime:  time.Now(),
	}

	s.taskManager.tasks[taskID] = status

	// 启动模拟执行
	go s.simulateTaskExecution(taskID, totalPapers)

	return nil
}

// 模拟任务执行
func (s *AutomationService) simulateTaskExecution(taskID string, totalPapers int) {
	for i := 1; i <= totalPapers; i++ {
		// 检查任务是否还存在
		status, exists := s.taskManager.GetTask(taskID)
		if !exists || status.Status == "cancelled" {
			break
		}

		// 模拟处理一张试卷
		time.Sleep(2 * time.Second) // 模拟处理时间

		// 更新进度
		score := 60 + (i % 40) // 模拟分数 60-99
		s.taskManager.UpdateProgress(taskID, i, score >= 60, float64(score))

		// 每10张试卷更新一次状态
		if i%10 == 0 || i == totalPapers {
			message := fmt.Sprintf("已处理 %d/%d 张试卷，平均分: %.1f", i, totalPapers, status.AverageScore)
			s.taskManager.UpdateMessage(taskID, message)
		}
	}

	// 任务完成
	s.taskManager.CompleteTask(taskID)
}

// 停止任务
func (s *AutomationService) StopTask(taskID string) error {
	s.taskManager.CancelTask(taskID)
	return nil
}

// 获取任务状态
func (s *AutomationService) GetTaskStatus(taskID string) (*TaskStatus, bool) {
	return s.taskManager.GetTask(taskID)
}

// 获取所有任务状态
func (s *AutomationService) GetAllTasks() map[string]*TaskStatus {
	return s.taskManager.GetAllTasks()
}

// TaskManager 方法
func (tm *TaskManager) GetTask(taskID string) (*TaskStatus, bool) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	task, exists := tm.tasks[taskID]
	return task, exists
}

func (tm *TaskManager) UpdateProgress(taskID string, paperNum int, success bool, score float64) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if task, exists := tm.tasks[taskID]; exists {
		task.CurrentPaper = paperNum
		task.LastUpdateTime = time.Now()

		if success {
			task.CompletedPapers++
		} else {
			task.FailedPapers++
		}

		// 计算平均分
		if task.CompletedPapers > 0 {
			totalScore := task.AverageScore*float64(task.CompletedPapers-1) + score
			task.AverageScore = totalScore / float64(task.CompletedPapers)
		}
	}
}

func (tm *TaskManager) UpdateMessage(taskID string, message string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if task, exists := tm.tasks[taskID]; exists {
		task.Message = message
		task.LastUpdateTime = time.Now()
	}
}

func (tm *TaskManager) CompleteTask(taskID string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if task, exists := tm.tasks[taskID]; exists {
		task.Status = "completed"
		task.Message = "任务执行完成"
		task.LastUpdateTime = time.Now()
	}
}

func (tm *TaskManager) CancelTask(taskID string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if task, exists := tm.tasks[taskID]; exists {
		task.Status = "cancelled"
		task.Message = "任务已取消"
		task.LastUpdateTime = time.Now()
	}
}

func (tm *TaskManager) GetAllTasks() map[string]*TaskStatus {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	// 创建副本
	result := make(map[string]*TaskStatus)
	for id, task := range tm.tasks {
		taskCopy := *task
		result[id] = &taskCopy
	}
	return result
}

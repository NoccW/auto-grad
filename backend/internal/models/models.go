package models

import (
	"time"
)

type User struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	OpenID       string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"openId"`
	Name         *string   `json:"name"`
	Email        *string   `json:"email"`
	LoginMethod  *string   `json:"loginMethod"`
	Role         string    `gorm:"default:user" json:"role"`
	UserRole     *string   `json:"userRole"` // parent, teacher
	PasswordHash *string   `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	LastSignedIn time.Time `json:"lastSignedIn"`
}

type GradingRecord struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	UserID         uint64    `gorm:"not null;index" json:"userId"`
	PaperImageUrl  string    `gorm:"not null" json:"paperImageUrl"`
	AnswerImageUrl *string   `json:"answerImageUrl"`
	OCRResult      *string   `gorm:"type:text" json:"ocrResult"`
	AIScore        *int      `json:"aiScore"`
	WrongQuestions *string   `gorm:"type:text" json:"wrongQuestions"`
	CorrectAnswers *string   `gorm:"type:text" json:"correctAnswers"`
	Status         string    `gorm:"default:pending;index" json:"status"` // pending, processing, completed, failed
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type TeacherTask struct {
	ID              uint64     `gorm:"primaryKey" json:"id"`
	UserID          uint64     `gorm:"not null;index" json:"userId"`
	TargetURL       string     `gorm:"not null" json:"targetUrl"`
	Account         string     `gorm:"not null" json:"account"`
	Password        string     `gorm:"not null" json:"-"`
	TaskID          string     `gorm:"not null" json:"taskId"`
	Status          string     `gorm:"default:pending;index" json:"status"` // pending, running, completed, failed, cancelled
	Progress        *string    `json:"progress"`
	TotalPapers     int        `gorm:"default:0" json:"totalPapers"`
	CompletedPapers int        `gorm:"default:0" json:"completedPapers"`
	FailedPapers    int        `gorm:"default:0" json:"failedPapers"`
	AverageScore    float64    `gorm:"default:0" json:"averageScore"`
	ErrorMessage    *string    `gorm:"type:text" json:"errorMessage"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	CompletedAt     *time.Time `json:"completedAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type TeacherTaskExecution struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	TaskID       uint64    `gorm:"not null;index" json:"taskId"`
	PaperID      string    `json:"paperId"`
	StudentName  string    `json:"studentName"`
	Score        int       `gorm:"default:0" json:"score"`
	OCRResult    string    `gorm:"type:text" json:"ocrResult"`
	AIFeedback   string    `gorm:"type:text" json:"aiFeedback"`
	Status       string    `gorm:"default:pending;index" json:"status"` // pending, processing, completed, failed
	ErrorMessage *string   `gorm:"type:text" json:"errorMessage"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`

	Task TeacherTask `gorm:"foreignKey:TaskID" json:"task,omitempty"`
}

type TaskStatistics struct {
	ID                uint64    `gorm:"primaryKey" json:"id"`
	TaskID            uint64    `gorm:"uniqueIndex;not null" json:"taskId"`
	TotalPapers       int       `gorm:"default:0" json:"totalPapers"`
	CompletedPapers   int       `gorm:"default:0" json:"completedPapers"`
	FailedPapers      int       `gorm:"default:0" json:"failedPapers"`
	AverageScore      float64   `gorm:"default:0" json:"averageScore"`
	MaxScore          int       `gorm:"default:0" json:"maxScore"`
	MinScore          int       `gorm:"default:100" json:"minScore"`
	PassRate          float64   `gorm:"default:0" json:"passRate"`
	ExcellenceRate    float64   `gorm:"default:0" json:"excellenceRate"`
	ScoreDistribution string    `gorm:"type:json" json:"scoreDistribution"`
	ErrorDistribution string    `gorm:"type:json" json:"errorDistribution"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`

	Task TeacherTask `gorm:"foreignKey:TaskID" json:"task,omitempty"`
}

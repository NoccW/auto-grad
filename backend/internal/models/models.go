package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
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
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null;index" json:"userId"`
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
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `gorm:"not null;index" json:"userId"`
	TargetURL       string     `gorm:"not null" json:"targetUrl"`
	Account         string     `gorm:"not null" json:"account"`
	Password        string     `gorm:"not null" json:"-"`
	TaskID          string     `gorm:"not null" json:"taskId"`
	Status          string     `gorm:"default:pending;index" json:"status"` // pending, running, completed, failed
	Progress        *string    `json:"progress"`
	TotalPapers     *int       `json:"totalPapers"`
	CompletedPapers int        `gorm:"default:0" json:"completedPapers"`
	ErrorMessage    *string    `gorm:"type:text" json:"errorMessage"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	CompletedAt     *time.Time `json:"completedAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

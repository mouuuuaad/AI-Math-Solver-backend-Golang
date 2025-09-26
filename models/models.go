package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"not null"`
	FullName  string         `json:"full_name" gorm:"not null"`
	Password  string         `json:"-" gorm:"not null"` // Hidden from JSON
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Solution struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Expression  string         `json:"expression" gorm:"not null"`
	StepsJSON   string         `json:"steps_json" gorm:"type:text"`
	FinalAnswer string         `json:"final_answer" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type UsageLimit struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_date"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Date      string    `json:"date" gorm:"type:date;uniqueIndex:idx_user_date"` // YYYY-MM-DD format
	Count     int       `json:"count" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Request/Response DTOs
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required,min=2"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type SolveMathRequest struct {
	Expression string `json:"expression" binding:"required"`
}

type SolveMathResponse struct {
	Steps []SolutionStep `json:"steps"`
	Final string         `json:"final"`
}

type SolutionStep struct {
	Index int    `json:"index"`
	Latex string `json:"latex"`
}

type HistoryResponse struct {
	Solutions []Solution `json:"solutions"`
	Total     int64      `json:"total"`
}

type UsageLimitResponse struct {
	Count     int    `json:"count"`
	Limit     int    `json:"limit"`
	Exceeded  bool   `json:"exceeded"`
	ResetTime string `json:"reset_time"` // ISO string for next reset
}

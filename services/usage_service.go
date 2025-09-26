package services

import (
	"time"

	"maths-solution-backend/models"

	"gorm.io/gorm"
)

const DAILY_LIMIT = 10

type UsageService struct {
	db *gorm.DB
}

func NewUsageService(db *gorm.DB) *UsageService {
	return &UsageService{db: db}
}

func (s *UsageService) CheckUsageLimit(userID uint) (*models.UsageLimitResponse, error) {
	today := time.Now().Format("2006-01-02")

	var usage models.UsageLimit
	result := s.db.Where("user_id = ? AND date = ?", userID, today).First(&usage)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// No usage today, create new record
			usage = models.UsageLimit{
				UserID: userID,
				Date:   today,
				Count:  0,
			}
			s.db.Create(&usage)
		} else {
			return nil, result.Error
		}
	}

	// Calculate next reset time (tomorrow at 00:00)
	tomorrow := time.Now().AddDate(0, 0, 1)
	nextReset := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())

	return &models.UsageLimitResponse{
		Count:     usage.Count,
		Limit:     DAILY_LIMIT,
		Exceeded:  usage.Count >= DAILY_LIMIT,
		ResetTime: nextReset.Format(time.RFC3339),
	}, nil
}

func (s *UsageService) IncrementUsage(userID uint) (*models.UsageLimitResponse, error) {
	today := time.Now().Format("2006-01-02")

	var usage models.UsageLimit
	result := s.db.Where("user_id = ? AND date = ?", userID, today).First(&usage)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new usage record
			usage = models.UsageLimit{
				UserID: userID,
				Date:   today,
				Count:  1,
			}
			if err := s.db.Create(&usage).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, result.Error
		}
	} else {
		// Increment existing usage
		usage.Count++
		if err := s.db.Save(&usage).Error; err != nil {
			return nil, err
		}
	}

	// Calculate next reset time
	tomorrow := time.Now().AddDate(0, 0, 1)
	nextReset := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())

	return &models.UsageLimitResponse{
		Count:     usage.Count,
		Limit:     DAILY_LIMIT,
		Exceeded:  usage.Count >= DAILY_LIMIT,
		ResetTime: nextReset.Format(time.RFC3339),
	}, nil
}

func (s *UsageService) GetUsageStats(userID uint) (*models.UsageLimitResponse, error) {
	return s.CheckUsageLimit(userID)
}

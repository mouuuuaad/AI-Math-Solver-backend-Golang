package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"maths-solution-backend/config"
	"maths-solution-backend/database"
	"maths-solution-backend/models"
	"maths-solution-backend/services"

	"github.com/gin-gonic/gin"
)

type MathHandler struct {
	config       *config.Config
	aiService    *services.AIService
	usageService *services.UsageService
}

func NewMathHandler(cfg *config.Config) *MathHandler {
	return &MathHandler{
		config:       cfg,
		aiService:    services.NewAIService(cfg),
		usageService: services.NewUsageService(database.DB),
	}
}

func (h *MathHandler) SolveMath(c *gin.Context) {
	var req models.SolveMathRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check usage limit before processing
	usage, err := h.usageService.CheckUsageLimit(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check usage limit"})
		return
	}

	if usage.Exceeded {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Daily usage limit exceeded",
			"usage": usage,
		})
		return
	}

	// Call AI service
	aiResp, err := h.aiService.SolveMath(req.Expression)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI service unavailable: " + err.Error()})
		return
	}

	// Convert steps to JSON
	stepsJSON, err := json.Marshal(aiResp.Steps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process solution"})
		return
	}

	// Save to database (best-effort). If it fails in dev, still return the AI result.
	_ = database.DB.Create(&models.Solution{
		UserID:      userIDUint,
		Expression:  req.Expression,
		StepsJSON:   string(stepsJSON),
		FinalAnswer: aiResp.Final,
	}).Error

	// Increment usage count after successful solve
	_, _ = h.usageService.IncrementUsage(userIDUint)

	// Return response
	c.JSON(http.StatusOK, models.SolveMathResponse{
		Steps: aiResp.Steps,
		Final: aiResp.Final,
	})
}

func (h *MathHandler) GetHistory(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Query solutions
	var solutions []models.Solution
	var total int64

	if err := database.DB.Model(&models.Solution{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count solutions"})
		return
	}

	if err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&solutions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch solutions"})
		return
	}

	c.JSON(http.StatusOK, models.HistoryResponse{
		Solutions: solutions,
		Total:     total,
	})
}

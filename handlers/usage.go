package handlers

import (
	"net/http"

	"maths-solution-backend/services"

	"github.com/gin-gonic/gin"
)

type UsageHandler struct {
	usageService *services.UsageService
}

func NewUsageHandler(usageService *services.UsageService) *UsageHandler {
	return &UsageHandler{usageService: usageService}
}

// GetUsageStats returns current usage statistics for the authenticated user
func (h *UsageHandler) GetUsageStats(c *gin.Context) {
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

	usage, err := h.usageService.GetUsageStats(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get usage stats"})
		return
	}

	c.JSON(http.StatusOK, usage)
}

// CheckUsageLimit checks if user can make another request
func (h *UsageHandler) CheckUsageLimit(c *gin.Context) {
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

	usage, err := h.usageService.CheckUsageLimit(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check usage limit"})
		return
	}

	c.JSON(http.StatusOK, usage)
}

// IncrementUsage increments the usage count for the user
func (h *UsageHandler) IncrementUsage(c *gin.Context) {
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

	usage, err := h.usageService.IncrementUsage(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment usage"})
		return
	}

	c.JSON(http.StatusOK, usage)
}

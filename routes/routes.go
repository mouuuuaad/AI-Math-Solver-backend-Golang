package routes

import (
	"maths-solution-backend/config"
	"maths-solution-backend/database"
	"maths-solution-backend/handlers"
	"maths-solution-backend/middleware"
	"maths-solution-backend/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(cfg *config.Config) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	r := gin.New()

	// Middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.CORSMiddleware(cfg))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg)
	mathHandler := handlers.NewMathHandler(cfg)
	usageHandler := handlers.NewUsageHandler(services.NewUsageService(database.DB))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes (public)
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		api.POST("/solve-math", mathHandler.SolveMath)
		api.GET("/history", mathHandler.GetHistory)
		api.GET("/usage", usageHandler.GetUsageStats)
		api.GET("/usage/check", usageHandler.CheckUsageLimit)
	}

	// Legacy route for backward compatibility
	r.POST("/solve-math", middleware.AuthMiddleware(cfg), mathHandler.SolveMath)

	return r
}

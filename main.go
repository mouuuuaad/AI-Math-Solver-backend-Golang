package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"maths-solution-backend/config"
	"maths-solution-backend/database"
	"maths-solution-backend/routes"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database (best-effort for local dev)
	if err := database.Connect(cfg); err != nil {
		log.Println("[warn] Database not available, continuing without DB:", err)
	} else {
		defer database.Close()
		if err := database.Migrate(); err != nil {
			log.Println("[warn] Database migration failed, continuing:", err)
		}
	}

	// Setup routes
	router := routes.SetupRoutes(cfg)

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)

	// Graceful shutdown
	go func() {
		if err := router.Run(":" + cfg.Server.Port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

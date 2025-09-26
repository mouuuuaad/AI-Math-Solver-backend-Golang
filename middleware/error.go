package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Error: %v", err.Error())

			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			case gin.ErrorTypePublic:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
		}
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}

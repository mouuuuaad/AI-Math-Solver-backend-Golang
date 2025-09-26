package middleware

import (
	"strings"

	"maths-solution-backend/config"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	origins := strings.Split(cfg.CORS.AllowedOrigins, ",")
	
	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})

	return func(ctx *gin.Context) {
		c.HandlerFunc(ctx.Writer, ctx.Request)
		ctx.Next()
	}
}

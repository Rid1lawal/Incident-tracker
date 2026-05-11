package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	router.Run(":" + port)
}

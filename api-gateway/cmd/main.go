package main

import (
	"log"

	"api-gateway/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize router
	router := gin.Default()

	// Initialize handlers
	otpHandler := handlers.NewOTPHandler("http://localhost:8081")

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// OTP endpoints
	otpGroup := router.Group("/api/v1/otp")
	{
		otpGroup.POST("/request", otpHandler.RequestOTP)
	}

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

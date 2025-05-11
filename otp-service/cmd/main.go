package main

import (
	"log"

	"otp-service/internal/handlers"
	"otp-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize OTP service
	otpService := service.NewOTPService()

	// Initialize router
	router := gin.Default()

	// Initialize handlers
	otpHandler := handlers.NewOTPHandler(otpService)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// OTP endpoints
	otpGroup := router.Group("/api/v1/otp")
	{
		otpGroup.POST("/generate", otpHandler.GenerateOTP)
		otpGroup.GET("/validate/:request_id", otpHandler.ValidateOTP)
	}

	// Start server
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

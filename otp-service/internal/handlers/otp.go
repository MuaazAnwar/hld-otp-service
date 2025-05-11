package handlers

import (
	"net/http"

	"common/models"
	"otp-service/internal/service"

	"github.com/gin-gonic/gin"
)

type OTPHandler struct {
	otpService *service.OTPService
}

func NewOTPHandler(otpService *service.OTPService) *OTPHandler {
	return &OTPHandler{
		otpService: otpService,
	}
}

func (h *OTPHandler) GenerateOTP(c *gin.Context) {
	var req models.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	otp, err := h.otpService.GenerateOTP(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := models.OTPResponse{
		RequestID:  otp.ID,
		Message:    "OTP generated successfully",
		ExpiresIn:  300, // 5 minutes
		RetryAfter: 60,  // 1 minute
	}

	c.JSON(http.StatusOK, response)
}

func (h *OTPHandler) ValidateOTP(c *gin.Context) {
	requestID := c.Param("request_id")
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "OTP code is required",
		})
		return
	}

	valid, err := h.otpService.ValidateOTP(requestID, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": valid,
	})
}

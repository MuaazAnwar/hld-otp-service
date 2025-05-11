package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"common/models"

	"github.com/gin-gonic/gin"
)

type OTPHandler struct {
	otpServiceURL string
	httpClient    *http.Client
}

func NewOTPHandler(otpServiceURL string) *OTPHandler {
	return &OTPHandler{
		otpServiceURL: otpServiceURL,
		httpClient:    &http.Client{},
	}
}

func (h *OTPHandler) RequestOTP(c *gin.Context) {
	var req models.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Forward request to OTP service
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process request",
		})
		return
	}

	resp, err := h.httpClient.Post(
		fmt.Sprintf("%s/api/v1/otp/generate", h.otpServiceURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to OTP service",
		})
		return
	}
	defer resp.Body.Close()

	// Forward the response from OTP service
	var otpResp models.OTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&otpResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process OTP service response",
		})
		return
	}

	c.JSON(resp.StatusCode, otpResp)
}

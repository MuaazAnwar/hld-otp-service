package models

import "time"

type OTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Purpose     string `json:"purpose" binding:"required"` // e.g., "login", "payment"
}

type OTPResponse struct {
	RequestID  string `json:"request_id"`
	Message    string `json:"message"`
	ExpiresIn  int    `json:"expires_in"`            // in seconds
	RetryAfter int    `json:"retry_after,omitempty"` // in seconds
}

type OTP struct {
	ID          string    `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Code        string    `json:"code"`
	Purpose     string    `json:"purpose"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	IsUsed      bool      `json:"is_used"`
}

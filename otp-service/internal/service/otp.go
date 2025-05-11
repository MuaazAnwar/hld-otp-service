package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"common/models"
)

const (
	OTPLength     = 6
	OTPExpiration = 5 * time.Minute
)

type OTPService struct {
	// In-memory storage for now, we'll replace with Redis later
	otps map[string]*models.OTP
}

func NewOTPService() *OTPService {
	return &OTPService{
		otps: make(map[string]*models.OTP),
	}
}

func (s *OTPService) GenerateOTP(req *models.OTPRequest) (*models.OTP, error) {
	// Generate random OTP code
	code, err := generateRandomCode(OTPLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Create OTP record
	now := time.Now()
	otp := &models.OTP{
		ID:          generateRequestID(),
		PhoneNumber: req.PhoneNumber,
		Code:        code,
		Purpose:     req.Purpose,
		CreatedAt:   now,
		ExpiresAt:   now.Add(OTPExpiration),
		IsUsed:      false,
	}

	// Store OTP
	s.otps[otp.ID] = otp

	return otp, nil
}

func (s *OTPService) ValidateOTP(requestID, code string) (bool, error) {
	otp, exists := s.otps[requestID]
	if !exists {
		return false, fmt.Errorf("OTP not found")
	}

	if otp.IsUsed {
		return false, fmt.Errorf("OTP already used")
	}

	if time.Now().After(otp.ExpiresAt) {
		return false, fmt.Errorf("OTP expired")
	}

	if otp.Code != code {
		return false, fmt.Errorf("invalid OTP code")
	}

	otp.IsUsed = true
	return true, nil
}

// Helper functions
func generateRandomCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to numeric string
	code := ""
	for _, b := range bytes {
		code += fmt.Sprintf("%d", b%10)
	}
	return code, nil
}

func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

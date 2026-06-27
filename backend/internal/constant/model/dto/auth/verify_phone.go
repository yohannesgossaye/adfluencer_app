package auth

import "time"

type VerifyPhoneRequest struct {
	OTP string `json:"otp"`
}

type VerifyPhoneResponse struct {
	User VerifyPhoneUserResponse `json:"user"`
}

type VerifyPhoneUserResponse struct {
	ID            string `json:"id"`
	PhoneVerified bool   `json:"phone_verified"`
}

type PhoneVerificationToken struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

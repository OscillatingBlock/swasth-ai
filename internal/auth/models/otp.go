package models

import "time"

type OTP struct {
	ID        int       `db:"id"`
	Phone     string    `db:"phone"`
	OTP       string    `db:"otp"`
	Attempts  int       `db:"attempts"`
	ExpiresAt time.Time `db:"expires_at"`
	Verified  bool      `db:"verified"`
}

type VerifyOTPInput struct {
	Phone string `json:"phone" validate:"required"`
	OTP   string `json:"otp" validate:"required"`
}

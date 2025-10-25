package models

import (
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	ID        uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()" json:"id" validate:"omitempty,numeric"`
	Phone     string    `bun:",notnull" json:"phone" validate:"required,e164"`            // e.g. +919876543210
	OTP       string    `bun:",notnull" json:"otp" validate:"required,len=6,numeric"`     // must be 6 digits
	Attempts  int       `bun:",notnull,default:0" json:"attempts" validate:"gte=0,lte=3"` // 0–3 attempts
	ExpiresAt time.Time `bun:",notnull" json:"expires_at" validate:"required"`            // must be valid time
	Verified  bool      `bun:",notnull,default:false" json:"verified" validate:"boolean"` // true if verified
}

type VerifyOTPInput struct {
	Phone string `json:"phone" validate:"required"`
	OTP   string `json:"otp" validate:"required"`
}

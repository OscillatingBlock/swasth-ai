package models

type OTP struct {
	ID        int    `db:"id"`
	Phone     string `db:"phone"`
	OTP       string `db:"otp"`
	Attempts  int    `db:"attempts"`
	ExpiresAt string `db:"expires_at"`
}

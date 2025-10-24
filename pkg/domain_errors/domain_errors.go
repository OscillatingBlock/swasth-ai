package domain_errors

import (
	"net/http"
	"strings"
	"swasthAI/pkg/errors"
)

// User Domain Errors
var (
	ErrInvalidPhoneFormat = errors.New("USER_INVALID_PHONE", "Invalid phone number format. Use +91xxxxxxxxxx", http.StatusBadRequest, nil)
	ErrUserAlreadyExists  = errors.New("USER_AlREADY_EXISTS", "Phone number already registered", http.StatusConflict, nil)
	ErrUserNotFound       = errors.New("USER_NOT_FOUND", "User not found", http.StatusNotFound, nil)
	ErrInvalidLanguage    = errors.New("USER_INVALID_LANGUAGE", "Unsupported language. Use: hi,en,ta,te,bn,mr", http.StatusUnprocessableEntity, nil)

	// Validation errors for User struct
	ErrInvalidFirstName = errors.New("USER_INVALID_FIRST_NAME", "First name must be 2-50 alphabetic characters", http.StatusBadRequest, nil)
	ErrInvalidLastName  = errors.New("USER_INVALID_LAST_NAME", "Last name must be 2-50 alphabetic characters", http.StatusBadRequest, nil)
)

// Auth Domain Errors
var (
	ErrInvalidOTP          = errors.New("AUTH_INVALID_OTP", "Invalid or expired OTP", http.StatusBadRequest, nil)
	ErrOTPAttemptsExceeded = errors.New("AUTH_OTP_ATTEMPTS_EXCEEDED", "Too many failed OTP attempts", http.StatusTooManyRequests, nil)
	ErrOTPExpired          = errors.New("AUTH_OTP_EXPIRED", "OTP expired. Please request new OTP", http.StatusBadRequest, nil)
	ErrResendCooldown      = errors.New("AUTH_RESEND_COOLDOWN", "Please wait 60 seconds before requesting new OTP", http.StatusTooManyRequests, nil)
	ErrFailedToSendOTP     = errors.New("AUTH_OTP_FAILED", "Failed to send OTP", http.StatusInternalServerError, nil)
)

// Voice Domain Errors
var (
	ErrInvalidAudioFormat  = errors.New("VOICE_INVALID_FORMAT", "Only WAV/MP3 audio supported", http.StatusBadRequest, nil)
	ErrAudioTooLarge       = errors.New("VOICE_AUDIO_TOO_LARGE", "Audio file must be less than 10MB", http.StatusRequestEntityTooLarge, nil)
	ErrUnsupportedLanguage = errors.New("VOICE_UNSUPPORTED_LANGUAGE", "Language not supported for voice analysis", http.StatusUnprocessableEntity, nil)
	ErrTranscriptionFailed = errors.New("VOICE_TRANSCRIPTION_FAILED", "Failed to transcribe audio", http.StatusInternalServerError, nil)
)

// Vision Domain Errors
var (
	ErrInvalidImageFormat = errors.New("VISION_INVALID_IMAGE", "Only JPEG/PNG images supported", http.StatusBadRequest, nil)
	ErrImageTooLarge      = errors.New("VISION_IMAGE_TOO_LARGE", "Image must be less than 5MB", http.StatusRequestEntityTooLarge, nil)
	ErrImageTooBlurry     = errors.New("VISION_IMAGE_BLURRY", "Image too blurry for analysis", http.StatusUnprocessableEntity, nil)
	ErrInvalidPDFFormat   = errors.New("VISION_INVALID_PDF", "Invalid PDF format", http.StatusBadRequest, nil)
	ErrPDFTooLarge        = errors.New("VISION_PDF_TOO_LARGE", "PDF must be less than 10MB", http.StatusRequestEntityTooLarge, nil)
	ErrOCRFailed          = errors.New("VISION_OCR_FAILED", "Unable to read text from report", http.StatusUnprocessableEntity, nil)
)

// Video Domain Errors
var (
	ErrInvalidCategory   = errors.New("VIDEO_INVALID_CATEGORY", "Invalid video category", http.StatusBadRequest, nil)
	ErrVideoNotFound     = errors.New("VIDEO_NOT_FOUND", "Video not found", http.StatusNotFound, nil)
	ErrAlreadyDownloaded = errors.New("VIDEO_ALREADY_DOWNLOADED", "Video already downloaded", http.StatusConflict, nil)
)

// Doctor Domain Errors
var (
	ErrInvalidCoordinates = errors.New("DOCTOR_INVALID_COORDS", "Invalid latitude/longitude", http.StatusBadRequest, nil)
	ErrInvalidRadius      = errors.New("DOCTOR_INVALID_RADIUS", "Radius must be between 1-50 km", http.StatusUnprocessableEntity, nil)
	ErrDoctorNotAvailable = errors.New("DOCTOR_NOT_AVAILABLE", "Doctor not available", http.StatusNotFound, nil)
	ErrSlotAlreadyBooked  = errors.New("CONSULTATION_SLOT_BOOKED", "Time slot already booked", http.StatusConflict, nil)
	ErrInvalidTimeFormat  = errors.New("CONSULTATION_INVALID_TIME", "Invalid datetime format", http.StatusUnprocessableEntity, nil)
)

// Helper functions for User validation
func ValidateUserPhone(phone string) error {
	if !strings.HasPrefix(phone, "+91") || len(phone) != 13 {
		return ErrInvalidPhoneFormat
	}
	return nil
}

func ValidateUserLanguage(lang string) error {
	supported := map[string]bool{"hi": true, "en": true, "ta": true, "te": true, "bn": true, "mr": true}
	if !supported[lang] {
		return ErrInvalidLanguage
	}
	return nil
}

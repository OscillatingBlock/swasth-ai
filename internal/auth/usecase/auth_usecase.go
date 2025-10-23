package usecase

import (
	"context"
	"fmt"
	"time"

	"swasthAI/config"
	"swasthAI/internal/auth"
	"swasthAI/internal/auth/models"
	"swasthAI/pkg/domain_errors"
	appErrors "swasthAI/pkg/errors"
	"swasthAI/pkg/logger"
	"swasthAI/pkg/utils"
)

type AuthUsecase struct {
	userRepo auth.UserRepository
	otpRepo  auth.OTPRepository
	cfg      config.Config
	logger   logger.Logger
}

func NewAuthUsecase(repo auth.UserRepository, otpRepo auth.OTPRepository, cfg config.Config, logger logger.Logger) *AuthUsecase {
	return &AuthUsecase{userRepo: repo, otpRepo: otpRepo, cfg: cfg, logger: logger}
}

func (uc *AuthUsecase) Register(ctx context.Context, user *models.User) error {
	domain_errors.ValidateUserPhone(user.Phone)

	// Check if user already exists
	existingUser, err := uc.userRepo.FindByPhone(ctx, user.Phone)
	if err != nil {
		uc.logger.Error("Error while finding user (authUC.Register.userRepo.FindByPhone)", "error", err)
		return appErrors.ErrDatabase
	}
	if existingUser != nil {
		uc.logger.Error("user already exists ", "phone", user.Phone)
		return appErrors.ErrAlreadyExists
	}

	// Create temporary user record (or just store in OTP table)
	tempUser := &models.User{
		Phone:     user.Phone,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	tempUser.PrepareCreate()

	// Generate & store OTP
	otp := utils.GenerateOTP()
	expiresAt := time.Now().Add(5 * time.Minute)

	err = uc.otpRepo.Create(ctx, &models.OTP{
		Phone:     user.Phone,
		OTP:       otp,
		ExpiresAt: expiresAt.String(),
		Attempts:  0,
	})
	if err != nil {
		uc.logger.Error("Failed to store OTP (authUC.Register.otpRepo.Create)", "error", err)
		return domain_errors.ErrInvalidOTP
	}

	// Send SMS (integrate with your SMS provider)
	//NOTE: right now , not sending any real sms just mocking
	err = uc.sendSMS(user.Phone, fmt.Sprintf("Your OTP is %s. Valid for 5 minutes.", otp))
	if err != nil {
		uc.logger.Error("SMS send failed", "error", err)
		return domain_errors.ErrFailedToSendOTP
	}
	uc.logger.Info("OTP sent successfully", "phone", user.Phone)
	return nil
}

func (uc *AuthUsecase) sendSMS(phone, message string) error {
	// TODO: Implement SMS sending logic
	return nil
}

func (uc *AuthUsecase) VerifyOTP(ctx context.Context, phone, otp string) (*models.UserWithToken, error) {

	//TODO: Implement OTP verification logic

	// Create or get user
	user, err := uc.userRepo.FindOrCreate(ctx, &models.User{
		Phone:     phone,
		FirstName: "", // Will be updated later or from OTP record
		LastName:  "",
	})
	if err != nil {
		uc.logger.Error("Failed to create/find user", "error", err)
		return nil, appErrors.ErrDatabase
	}

	// Generate JWT tokens
	token, refreshToken, err := utils.GenerateJWTToken(user, uc.cfg.JWT)
	if err != nil {
		uc.logger.Error("Token generation failed", "error", err)
		return nil, appErrors.ErrJWTGeneration
	}

	// Delete used OTP
	err = uc.otpRepo.Delete(ctx, phone)
	if err != nil {
		uc.logger.Error("Failed to delete OTP", "error", err)
		// Non-critical - continue
	}

	return &models.UserWithToken{
		User:         user,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*models.Tokens, error) {

	claims, err := utils.ValidateRefreshToken(refreshToken, uc.cfg.JWT)
	if err != nil {
		uc.logger.Error("Invalid refresh token")
		return nil, appErrors.ErrUnauthorized
	}

	user, err := uc.userRepo.FindByID(ctx, claims.ID)
	if err != nil {
		uc.logger.Error("Failed to find user (authUC.RefreshToken.FindByID)", "error", err)
		return nil, appErrors.ErrDatabase
	}

	newToken, refreshToken, err := utils.GenerateJWTToken(user, uc.cfg.JWT)
	if err != nil {
		uc.logger.Error("Failed to generate new token (authUC.RefreshToken.GenerateJWTToken)", "error", err)
		return nil, appErrors.ErrJWTGeneration
	}

	return &models.Tokens{
		Token:        newToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400,
	}, nil
}

func (uc *AuthUsecase) ResendOTP(ctx context.Context, phone string) error {
	//TODO: Implement OTP resend logic
	return nil
}

func (uc *AuthUsecase) GetUserByID(ctx context.Context, token string) (*models.User, error) {
	claims, ok := ctx.Value("claims").(*utils.JWTClaims)
	if !ok {
		uc.logger.Error("Invalid claims in context")
		return nil, appErrors.ErrUnauthorized
	}

	user, err := uc.userRepo.FindByID(ctx, claims.ID)
	if err != nil {
		uc.logger.Error("Failed to find user (authUC.GetUserByID.FindByID)", "error", err)
		return nil, appErrors.ErrDatabase
	}
	if user == nil {
		return nil, appErrors.ErrUnauthorized
	}
	return user, nil
}

func (uc *AuthUsecase) UpdateProfile(ctx context.Context, input *models.UpdateProfileInput) (*models.User, error) {
	claims, ok := ctx.Value("claims").(*utils.JWTClaims)
	if !ok {
		uc.logger.Error("Invalid claims in context")
		return nil, appErrors.ErrUnauthorized
	}

	user, err := uc.userRepo.FindByID(ctx, claims.ID)
	if err != nil {
		uc.logger.Error("Failed to find user (authUC.UpdateProfile.FindByID)", "error", err)
		return nil, appErrors.ErrDatabase
	}
	if user == nil {
		return nil, domain_errors.ErrUserNotFound
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.Language = input.LastName
	}
	if input.Language != "" {
		if err := domain_errors.ValidateUserLanguage(input.Language); err != nil {
			return nil, domain_errors.ErrInvalidLanguage
		}
		user.Language = input.Language
	}
	user.FullName = user.FirstName + " " + user.LastName

	updatedUser, err := uc.userRepo.Update(ctx, user)
	if err != nil {
		uc.logger.Error("Failed to update user (authUC.UpdateProfile.Update)", "error", err)
		return nil, appErrors.ErrDatabase
	}
	return updatedUser, nil
}

package usecase

import (
	"context"
	"errors"
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

func (uc *AuthUsecase) SendOTP(ctx context.Context, phone string) error {
	// TODO: Implement SMS sending logic
	err := domain_errors.ValidateUserPhone(phone)
	if err != nil {
		uc.logger.Error("error", err)
		return domain_errors.ErrInvalidPhoneFormat
	}

	//Rate limiting , max 3 otp / hour
	count, err := uc.otpRepo.CountRecent(ctx, phone, time.Now().UTC().Truncate(time.Second))
	if err != nil {
		uc.logger.Error("error while getting recent otp count (otpRepo.CountRecent)", "error", err)
		return appErrors.ErrInternal
	}
	if count >= 3 {
		return domain_errors.ErrOTPAttemptsExceeded
	}

	//create and store otp
	otp := utils.GenerateOTP()
	expiresAt := time.Now().Add(5 * time.Minute)

	err = uc.otpRepo.Create(ctx, &models.OTP{
		Phone:     phone,
		OTP:       otp,
		ExpiresAt: expiresAt,
		Attempts:  0,
		Verified:  false,
	})
	if err != nil {
		uc.logger.Error("failed to store OTP (authUC.sendSMS.otpRepo.Create)", "error", err)
		return appErrors.ErrDatabase
	}

	//send sms , mocking for now
	if err = uc.sendSMS(ctx, fmt.Sprintf("Your OTP is %s. Valid for 5 minutes.", otp), phone); err != nil {
		uc.logger.Error("failed to send SMS", "error", err)
		return domain_errors.ErrFailedToSendOTP
	}
	return nil
}

func (uc *AuthUsecase) VerifyOTP(ctx context.Context, phone, otp string) (*models.UserWithToken, bool, error) {

	// otpRecord, err := uc.otpRepo.FindByPhone(ctx, phone)
	// if err !=

	// //TODO: Implement OTP verification logic
	// otpRecord.Verified = true

	//check if user already exists
	existingUser, err := uc.userRepo.FindByPhone(ctx, phone)

	if err != nil {
		if errors.Is(err, domain_errors.ErrUserNotFound) {
			// user does not exist → continue to signup
			existingUser = nil
		} else {
			uc.logger.Error("Error while searching user (authUC.VerifyOTP.userRepo.FindByPhone)", "error", err)
			return nil, false, appErrors.ErrDatabase
		}
	}

	if existingUser != nil {
		// user already exists, generate new token
		token, refreshToken, err := utils.GenerateJWTToken(existingUser, uc.cfg.JWT)
		if err != nil {
			uc.logger.Error("Failed to generate token (authUC.VerifyOTP.GenerateJWTToken)", "error", err)
			return nil, false, appErrors.ErrJWTGeneration
		}

		uc.otpRepo.Delete(ctx, otp)

		return &models.UserWithToken{
			User:         existingUser,
			Token:        token,
			RefreshToken: refreshToken,
		}, true, nil
	}
	//signup
	return nil, false, nil
}

func (uc *AuthUsecase) RegisterUser(ctx context.Context, input *models.RegisterUserInput) (*models.UserWithToken, error) {
	existingUser, err := uc.userRepo.FindByPhone(ctx, input.Phone)

	if err != nil && !errors.Is(err, domain_errors.ErrUserNotFound) {
		uc.logger.Error("failed to search user (authUC.RegisterUser.userRepo.FindByPhone)", "error", err)
		return nil, appErrors.ErrDatabase
	}
	if existingUser != nil {
		uc.logger.Error("user already exists", "phone", input.Phone)
		return nil, domain_errors.ErrUserAlreadyExists
	}

	err = domain_errors.ValidateUserLanguage(input.Language)
	if err != nil {
		uc.logger.Error("invalid language", "error", err)
		return nil, domain_errors.ErrInvalidLanguage
	}

	//create user
	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		Language:  input.Language,
	}
	user.PrepareCreate()
	createdUser, err := uc.userRepo.Create(ctx, &user)
	if err != nil {
		uc.logger.Error("failed to create user (authUC.RegisterUser.userRepo.Create)", "error", err)
		return nil, appErrors.ErrDatabase
	}

	//generate token
	token, refreshToken, err := utils.GenerateJWTToken(createdUser, uc.cfg.JWT)
	if err != nil {
		uc.logger.Error("failed to generate token (authUC.RegisterUser.GenerateJWTToken)", "error", err)
		return nil, appErrors.ErrJWTGeneration
	}
	return &models.UserWithToken{
		User:         createdUser,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*models.Tokens, error) {
	user, err := utils.ValidateRefreshToken(refreshToken, uc.cfg.JWT)
	if err != nil {
		uc.logger.Error("Invalid refresh token", "error", err)
		return nil, appErrors.ErrUnauthorized
	}

	newToken, newRefreshToken, err := utils.GenerateJWTToken(user, uc.cfg.JWT)
	if err != nil {
		uc.logger.Error("Failed to generate new tokens", "error", err)
		return nil, appErrors.ErrJWTGeneration
	}

	return &models.Tokens{
		Token:        newToken,
		RefreshToken: newRefreshToken,
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

func (uc *AuthUsecase) sendSMS(ctx context.Context, message, phone string) error {
	//TODO: Implement OTP resend logic

	return nil
}

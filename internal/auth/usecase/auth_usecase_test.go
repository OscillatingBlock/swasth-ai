package usecase

import (
	"context"
	"testing"
	"time"

	"swasthAI/config"
	mocks "swasthAI/internal/auth/mocks"
	"swasthAI/internal/auth/models"
	"swasthAI/pkg/domain_errors"
	"swasthAI/pkg/logger"
	"swasthAI/pkg/utils"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (AuthUsecase, *mocks.MockUserRepository, *mocks.MockOTPRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockOTPRepo := mocks.NewMockOTPRepository(ctrl)

	cfg := config.Config{
		JWT: config.JWT{
			Secret: "test-secret",
		},
	}
	log, _ := logger.NewLogger(&config.Config{LoggerMode: config.LoggerMode{Development: true}})

	uc := NewAuthUsecase(mockUserRepo, mockOTPRepo, cfg, *log)
	return *uc, mockUserRepo, mockOTPRepo, ctrl
}

func TestAuthUsecase_SendOTP_Success(t *testing.T) {
	uc, _, mockOTPRepo, ctrl := setupTest(t)
	defer ctrl.Finish()

	ctx := context.Background()
	phone := "+919876543210"

	mockOTPRepo.EXPECT().CountRecent(ctx, phone, time.Now().UTC().Truncate(time.Second)).Return(0, nil)

	// Mock: OTP create
	mockOTPRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	err := uc.SendOTP(ctx, phone)
	assert.NoError(t, err)
}

func TestAuthUsecase_SendOTP_InvalidPhone(t *testing.T) {
	uc, _, _, ctrl := setupTest(t)
	defer ctrl.Finish()

	err := uc.SendOTP(context.Background(), "invalid")
	assert.Error(t, err)
	assert.Equal(t, domain_errors.ErrInvalidPhoneFormat, err)
}

func TestAuthUsecase_SendOTP_RateLimit(t *testing.T) {
	uc, _, mockOTPRepo, ctrl := setupTest(t)
	defer ctrl.Finish()

	ctx := context.Background()
	phone := "+919876543210"

	mockOTPRepo.EXPECT().CountRecent(ctx, phone, time.Now().UTC().Truncate(time.Second)).Return(3, nil)

	err := uc.SendOTP(ctx, phone)
	assert.Error(t, err)
	assert.Equal(t, domain_errors.ErrOTPAttemptsExceeded, err)
}

func TestAuthUsecase_VerifyOTP_Login_Success(t *testing.T) {
	uc, mockUserRepo, mockOTPRepo, ctrl := setupTest(t)
	defer ctrl.Finish()

	ctx := context.Background()
	phone := "+919876543210"
	otp := "123456"

	user := &models.User{
		ID:        uuid.New(),
		Phone:     phone,
		FirstName: "रमेश",
		LastName:  "कुमार",
	}

	mockUserRepo.EXPECT().FindByPhone(ctx, phone).Return(user, nil)
	mockOTPRepo.EXPECT().Delete(ctx, otp).Return(nil)

	// Mock token generation
	// GenerateTokens := func(u *models.User, cfg config.JWT) (string, string, error) {
	// 	return "access-token", "refresh-token", nil
	// }

	result, registered, err := uc.VerifyOTP(ctx, phone, otp)
	assert.NoError(t, err)
	assert.True(t, registered)
	assert.Equal(t, user.ID, result.User.ID)
	assert.Equal(t, "access-token", result.Token)
}

func TestAuthUsecase_VerifyOTP_Signup_Flow(t *testing.T) {
	uc, mockUserRepo, _, ctrl := setupTest(t)
	defer ctrl.Finish()

	ctx := context.Background()
	phone := "+919876543210"
	otp := "123456"

	// otpRec := &models.OTP{
	// 	Phone:     phone,
	// 	OTP:       otp,
	// 	ExpiresAt: time.Now().Add(5 * time.Minute),
	// 	Attempts:  0,
	// }

	// mockOTPRepo.EXPECT().FindByPhone(ctx, phone).Return(otpRec, nil)
	mockUserRepo.EXPECT().FindByPhone(ctx, phone).Return((*models.User)(nil), nil)

	result, registered, err := uc.VerifyOTP(ctx, phone, otp)
	assert.NoError(t, err)
	assert.False(t, registered)
	assert.Nil(t, result)
}

func TestAuthUsecase_RegisterUser_Success(t *testing.T) {
	uc, mockUserRepo, _, ctrl := setupTest(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := &models.RegisterUserInput{
		Phone:     "+919876543210",
		FirstName: "रमेश",
		LastName:  "कुमार",
		Language:  "hi",
	}

	// otpRec := &models.OTP{
	// 	Phone:    input.Phone,
	// 	Verified: true,
	// }

	mockUserRepo.EXPECT().FindByPhone(ctx, input.Phone).Return((*models.User)(nil), nil)
	mockUserRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, u *models.User) (*models.User, error) {
		u.ID = uuid.New()
		return u, nil
	})

	// utils.GenerateTokens = func(u *models.User, cfg config.JWTConfig) (string, string, error) {
	// 	return "access-token", "refresh-token", nil
	// }

	result, err := uc.RegisterUser(ctx, input)
	assert.NoError(t, err)
	assert.Equal(t, input.FirstName, result.User.FirstName)
	assert.Equal(t, "access-token", result.Token)
}

// func TestAuthUsecase_RegisterUser_OTPNotVerified(t *testing.T) {
// 	uc, _, _, ctrl := setupTest(t)
// 	defer ctrl.Finish()

// 	input := &models.RegisterUserInput{Phone: "+919876543210"}

// 	_, err := uc.RegisterUser(context.Background(), input)
// 	assert.Error(t, err)
// 	assert.Equal(t, domain_errors.ErrInvalidOTP, err)
// }

func TestAuthUsecase_UpdateProfile_Success(t *testing.T) {
	uc, mockUserRepo, _, ctrl := setupTest(t)
	defer ctrl.Finish()

	ctx := context.WithValue(context.Background(), "claims", &utils.JWTClaims{ID: uuid.New()})
	userID := ctx.Value("claims").(*utils.JWTClaims).ID

	existing := &models.User{
		ID:        uuid.New(),
		FirstName: "रमेश",
		LastName:  "कुमार",
		Language:  "hi",
	}

	input := &models.UpdateProfileInput{
		FirstName: "राम",
		Language:  "ta",
	}

	mockUserRepo.EXPECT().FindByID(ctx, userID).Return(existing, nil)
	mockUserRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, u *models.User) (*models.User, error) {
		assert.Equal(t, "राम", u.FirstName)
		assert.Equal(t, "ta", u.Language)
		assert.Equal(t, "राम कुमार", u.FullName)
		return u, nil
	})

	updated, err := uc.UpdateProfile(ctx, input)
	assert.NoError(t, err)
	assert.Equal(t, "राम", updated.FirstName)
	assert.Equal(t, "ta", updated.Language)
}

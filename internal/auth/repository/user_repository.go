package repository

import (
	"context"
	"database/sql"
	"swasthAI/internal/auth/models"
	"swasthAI/pkg/domain_errors"
	"swasthAI/pkg/logger"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB, logger logger.Logger) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := r.db.NewInsert().Model(user).Returning("*").Exec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.CreateUser.InsertUser: ")
	}
	return user, nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	user := new(models.User)
	err := r.db.NewSelect().Model(user).Where("phone = ?", phone).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain_errors.ErrUserNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByPhone.Select")
	}
	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := new(models.User)
	err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByID.Select")
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := r.db.NewUpdate().Model(user).Where("id = ?", user.ID).Returning("*").Exec(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain_errors.ErrUserNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.UpdateUser.Update")
	}
	return user, nil
}

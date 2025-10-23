package repository

import (
	"context"
	"swasthAI/internal/auth/models"

	"github.com/pkg/errors"

	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
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
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByPhone.Select")
	}
	return user, nil
}

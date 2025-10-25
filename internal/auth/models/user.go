package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// User represents a user in the system.
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID        uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()" json:"id" validate:"required,uuid4"`
	Phone     string    `bun:",unique,notnull" json:"phone" validate:"required,e164"` // e164 = +919876543210 format
	FirstName string    `bun:",notnull" json:"first_name" validate:"required,alpha,min=2,max=50"`
	LastName  string    `bun:",notnull" json:"last_name" validate:"required,alpha,min=2,max=50"`
	FullName  string    `bun:",notnull" json:"full_name" validate:"required"`
	Language  string    `bun:",notnull" json:"language" validate:"required,alpha,len=2"` // e.g. 'en', 'hi'
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
}

func (u *User) PrepareCreate() {
	u.FullName = u.FirstName + " " + u.LastName
	u.Phone = strings.TrimSpace(u.Phone)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
}

type UserWithToken struct {
	User         *User  `json:"user" validate:"required"`
	Token        string `json:"token" validate:"required,jwt"`
	RefreshToken string `json:"refresh_token" validate:"required,jwt"`
}

type Tokens struct {
	Token        string `json:"token" validate:"required,jwt"`
	RefreshToken string `json:"refresh_token" validate:"required,jwt"`
	ExpiresIn    int    `json:"expires_in" validate:"required,gt=0"`
}

type UpdateProfileInput struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Language  string `json:"language" validate:"required,alpha,len=2"`
}

type RegisterUserInput struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Language  string `json:"language" validate:"required,alpha,len=2"`
	Phone     string `json:"phone" validate:"required,e164"`
}

type SendOTPInput struct {
	Phone string `json:"phone" validate:"required,e164"`
}

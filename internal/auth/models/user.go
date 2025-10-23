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

	ID        uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id" validate:"required,uuid4"`
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
	User         *User
	Token        string
	RefreshToken string
}

type Tokens struct {
	Token        string
	RefreshToken string
	ExpiresIn    int
}

type UpdateProfileInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Language  string `json:"language"`
}

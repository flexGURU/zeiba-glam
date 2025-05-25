package repository

import (
	"context"
	"time"

	"github.com/flexGURU/zeiba-glam/backend/pkg"
)

type User struct {
	ID           uint32    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number"`
	Password     *string   `json:"password,omitempty"`
	IsAdmin      bool      `json:"is_admin"`
	RefreshToken *string   `json:"refresh_token,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type UpdateUser struct {
	ID           uint32  `json:"id"`
	Name         *string `json:"name"`
	Email        *string `json:"email"`
	PhoneNumber  *string `json:"phone_number"`
	Password     *string `json:"password"`
	IsAdmin      *bool   `json:"is_admin"`
	RefreshToken *string `json:"refresh_token"`
}

type UserFilter struct {
	Pagination *pkg.Pagination
	Search     *string
	IsAdmin    *bool
}

type UserRepository interface {
	GetUserInternal(ctx context.Context, id uint32, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, id uint32, email string) (*User, error)
	ListUsers(ctx context.Context, filter *UserFilter) ([]*User, *pkg.Pagination, error)
	UpdateUser(ctx context.Context, user *UpdateUser) (*User, error)
}

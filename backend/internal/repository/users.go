package repository

import (
	"context"
	"time"
)

type User struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"password"`
	IsAdmin     bool      `json:"is_admin"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	UpdateUser(ctx context.Context, user User) (User, error)
	DeleteUser(ctx context.Context, id int64) error
}

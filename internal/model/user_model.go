package model

import (
	"time"
)

type UserResponse struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	Username  string     `json:"username,omitempty"`
	Token     string     `json:"token,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Username string `json:"username" validate:"required,min=5,max=30"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=100"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=100"`
}

type UpdateUserRequest struct {
	ID       string `json:"-" validate:"uuid4"`
	Name     string `json:"name,omitempty" validate:"max=100"`
	Email    string `json:"email,omitempty" validate:"omitempty,validateEmail,max=255"`
	Password string `json:"password,omitempty" validate:"max=100"`
}

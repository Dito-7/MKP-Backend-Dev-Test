package entity

import (
	"time"
)

type UserRole string

const (
	RoleAdmin    UserRole = "ADMIN"
	RoleStaff    UserRole = "STAFF"
	RoleCustomer UserRole = "CUSTOMER"
)

type User struct {
	ID           string    `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	PhoneNumber  string    `json:"phone_number"`
	Role         UserRole  `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	FullName    string   `json:"full_name" binding:"required,min=2,max=100"`
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=6"`
	PhoneNumber string   `json:"phone_number" binding:"required,min=8,max=20"`
	Role        UserRole `json:"role"`
}

type LoginResponse struct {
	User      UserResponse `json:"user"`
	Token     string       `json:"token"`
	TokenType string       `json:"token_type"`
	ExpiresIn int64        `json:"expires_in_seconds"`
}

type UserResponse struct {
	ID          string   `json:"id"`
	FullName    string   `json:"full_name"`
	Email       string   `json:"email"`
	PhoneNumber string   `json:"phone_number"`
	Role        UserRole `json:"role"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		FullName:    u.FullName,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Role:        u.Role,
	}
}

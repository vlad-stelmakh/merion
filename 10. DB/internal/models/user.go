package models

import (
	"database/sql"
	"time"
)

// User представляет пользователя в системе
type User struct {
	ID        int            `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	Email     string         `json:"email" db:"email"`
	Age       sql.NullInt64  `json:"age" db:"age"`
	Bio       sql.NullString `json:"bio" db:"bio"`
	IsActive  bool           `json:"is_active" db:"is_active"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest представляет запрос на создание пользователя
type CreateUserRequest struct {
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Age      sql.NullInt64  `json:"age"`
	Bio      sql.NullString `json:"bio"`
	IsActive bool           `json:"is_active"`
}

// UpdateUserRequest представляет запрос на обновление пользователя
type UpdateUserRequest struct {
	Name     *string        `json:"name,omitempty"`
	Email    *string        `json:"email,omitempty"`
	Age      *sql.NullInt64 `json:"age,omitempty"`
	Bio      *sql.NullString `json:"bio,omitempty"`
	IsActive *bool          `json:"is_active,omitempty"`
}

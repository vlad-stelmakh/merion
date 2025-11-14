package repository

import (
	"context"
	"database-examples/internal/models"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.User, error)
	Update(ctx context.Context, id int, req *models.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}

// Repository объединяет все репозитории
type Repository struct {
	User UserRepository
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"database-examples/internal/models"
)

// UserRepository реализует интерфейс repository.UserRepository для PostgreSQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create создает нового пользователя
func (r *UserRepository) Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	query := `
		INSERT INTO users (name, email, age, bio, is_active) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, name, email, age, bio, is_active, created_at, updated_at`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, req.Name, req.Email, req.Age, req.Bio, req.IsActive).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.Bio,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return user, nil
}

// GetByID получает пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, name, email, age, bio, is_active, created_at, updated_at 
		FROM users 
		WHERE id = $1`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.Bio,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь с ID %d не найден", id)
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return user, nil
}

// GetByEmail получает пользователя по email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, email, age, bio, is_active, created_at, updated_at 
		FROM users 
		WHERE email = $1`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.Bio,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь с email %s не найден", email)
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return user, nil
}

// GetAll получает всех пользователей с пагинацией
func (r *UserRepository) GetAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	query := `
		SELECT id, name, email, age, bio, is_active, created_at, updated_at 
		FROM users 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователей: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Age,
			&user.Bio,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования пользователя: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по строкам: %w", err)
	}

	return users, nil
}

// Update обновляет данные пользователя
func (r *UserRepository) Update(ctx context.Context, id int, req *models.UpdateUserRequest) (*models.User, error) {
	// Сначала получаем текущие данные пользователя
	currentUser, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Обновляем только переданные поля
	if req.Name != nil {
		currentUser.Name = *req.Name
	}
	if req.Email != nil {
		currentUser.Email = *req.Email
	}
	if req.Age != nil {
		currentUser.Age = *req.Age
	}
	if req.Bio != nil {
		currentUser.Bio = *req.Bio
	}
	if req.IsActive != nil {
		currentUser.IsActive = *req.IsActive
	}

	query := `
		UPDATE users 
		SET name = $2, email = $3, age = $4, bio = $5, is_active = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 
		RETURNING id, name, email, age, bio, is_active, created_at, updated_at`

	user := &models.User{}
	err = r.db.QueryRowContext(ctx, query, id, currentUser.Name, currentUser.Email, 
		currentUser.Age, currentUser.Bio, currentUser.IsActive).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.Bio,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка обновления пользователя: %w", err)
	}

	return user, nil
}

// Delete удаляет пользователя
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("пользователь с ID %d не найден", id)
	}

	return nil
}

// Count возвращает общее количество пользователей
func (r *UserRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка подсчета пользователей: %w", err)
	}

	return count, nil
}

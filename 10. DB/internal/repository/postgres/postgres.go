package postgres

import (
	"database/sql"
	"fmt"

	"database-examples/internal/config"
	"database-examples/internal/repository"

	_ "github.com/lib/pq"
)

// NewConnection создает новое подключение к PostgreSQL
func NewConnection(cfg *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения с БД: %w", err)
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	return db, nil
}

// NewRepository создает новый экземпляр Repository с PostgreSQL реализациями
func NewRepository(db *sql.DB) *repository.Repository {
	return &repository.Repository{
		User: NewUserRepository(db),
	}
}

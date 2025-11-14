package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"database-examples/internal/config"
	"database-examples/internal/models"
	"database-examples/internal/repository"
	"database-examples/internal/repository/postgres"

	_ "github.com/lib/pq"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключаемся к базе данных
	db, err := postgres.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	fmt.Println("✅ Успешно подключились к PostgreSQL!")

	// Создаем репозитории
	userRepo := postgres.NewRepository(db)

	// Демонстрация работы с пользователями
	demonstrateUserOperations(userRepo)
}

// demonstrateUserOperations демонстрирует операции с пользователями
func demonstrateUserOperations(repo *repository.Repository) {
	ctx := context.Background()
	fmt.Println("\n🔄 Демонстрация операций с пользователями:")

	// CREATE - Создание пользователей
	fmt.Println("\n📝 Создание пользователей...")

	user1Req := &models.CreateUserRequest{
		Name:     "Алексей Иванов",
		Email:    "alexey@example.com",
		Age:      sql.NullInt64{Int64: 28, Valid: true},
		Bio:      sql.NullString{String: "Разработчик Go", Valid: true},
		IsActive: true,
	}

	user1, err := repo.User.Create(ctx, user1Req)
	if err != nil {
		log.Printf("Ошибка создания пользователя: %v", err)
	} else {
		fmt.Printf("✅ Создан пользователь с ID: %d\n", user1.ID)
	}

	user2Req := &models.CreateUserRequest{
		Name:     "Мария Петрова",
		Email:    "maria@example.com",
		Age:      sql.NullInt64{Int64: 25, Valid: true},
		Bio:      sql.NullString{String: "Дизайнер", Valid: true},
		IsActive: true,
	}

	user2, err := repo.User.Create(ctx, user2Req)
	if err != nil {
		log.Printf("Ошибка создания пользователя: %v", err)
	} else {
		fmt.Printf("✅ Создан пользователь с ID: %d\n", user2.ID)
	}

	// READ - Чтение пользователя по ID
	fmt.Println("\n📖 Чтение пользователя по ID...")
	if user1 != nil {
		foundUser, err := repo.User.GetByID(ctx, user1.ID)
		if err != nil {
			log.Printf("Ошибка получения пользователя: %v", err)
		} else {
			fmt.Printf("✅ Найден пользователь: %s (%s)\n", foundUser.Name, foundUser.Email)
			if foundUser.Age.Valid {
				fmt.Printf("   Возраст: %d\n", foundUser.Age.Int64)
			}
			if foundUser.Bio.Valid {
				fmt.Printf("   Биография: %s\n", foundUser.Bio.String)
			}
		}
	}

	// READ - Получение всех пользователей
	fmt.Println("\n📖 Получение всех пользователей...")
	users, err := repo.User.GetAll(ctx, 10, 0)
	if err != nil {
		log.Printf("Ошибка получения пользователей: %v", err)
	} else {
		fmt.Printf("✅ Найдено пользователей: %d\n", len(users))
		for _, u := range users {
			ageStr := "не указан"
			if u.Age.Valid {
				ageStr = fmt.Sprintf("%d", u.Age.Int64)
			}
			fmt.Printf("   - %s (%s), возраст: %s\n", u.Name, u.Email, ageStr)
		}
	}

	// UPDATE - Обновление пользователя
	fmt.Println("\n✏️  Обновление пользователя...")
	if user1 != nil {
		newName := "Алексей Иванович Иванов"
		newAge := sql.NullInt64{Int64: 29, Valid: true}

		updateReq := &models.UpdateUserRequest{
			Name: &newName,
			Age:  &newAge,
		}

		updatedUser, err := repo.User.Update(ctx, user1.ID, updateReq)
		if err != nil {
			log.Printf("Ошибка обновления пользователя: %v", err)
		} else {
			fmt.Printf("✅ Пользователь обновлен: %s, возраст: %d\n",
				updatedUser.Name, updatedUser.Age.Int64)
		}
	}

	// DELETE - Удаление пользователя
	fmt.Println("\n🗑️  Удаление пользователя...")
	if user2 != nil {
		err = repo.User.Delete(ctx, user2.ID)
		if err != nil {
			log.Printf("Ошибка удаления пользователя: %v", err)
		} else {
			fmt.Println("✅ Пользователь удален")
		}
	}

	// Подсчет пользователей
	count, err := repo.User.Count(ctx)
	if err != nil {
		log.Printf("Ошибка подсчета пользователей: %v", err)
	} else {
		fmt.Printf("📊 Всего пользователей в системе: %d\n", count)
	}
}

package main

import (
	"errors"
	"testing"

	"advanced-patterns/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotificationService_WithMocks(t *testing.T) {
	t.Run("успешная отправка уведомления", func(t *testing.T) {
		// Создаем мок
		mockEmailSender := mocks.NewEmailSender(t)

		// Настраиваем ожидания
		mockEmailSender.On("SendEmail",
			"test@example.com",
			"Уведомление",
			"Привет, Тест!\n\nВаш заказ готов").Return(nil)

		// Создаем сервис с моком
		service := NewNotificationService(mockEmailSender)

		// Создаем пользователя
		user := &User{
			Name:  "Тест",
			Email: "test@example.com",
			Age:   25,
		}

		// Вызываем тестируемый метод
		err := service.NotifyUser(user, "Ваш заказ готов")

		// Проверяем результат
		assert.NoError(t, err)

		// Проверяем, что мок был вызван с правильными параметрами
		mockEmailSender.AssertExpectations(t)
	})

	t.Run("ошибка при отправке email", func(t *testing.T) {
		mockEmailSender := mocks.NewEmailSender(t)

		// Мок возвращает ошибку
		mockEmailSender.On("SendEmail",
			mock.Anything, mock.Anything, mock.Anything).Return(errors.New("сервер недоступен"))

		service := NewNotificationService(mockEmailSender)
		user := &User{
			Name:  "Тест",
			Email: "test@example.com",
			Age:   25,
		}

		err := service.NotifyUser(user, "Тест")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "сервер недоступен")
		mockEmailSender.AssertExpectations(t)
	})

	t.Run("некорректный email пользователя", func(t *testing.T) {
		mockEmailSender := mocks.NewEmailSender(t)
		service := NewNotificationService(mockEmailSender)

		user := &User{
			Name:  "Тест",
			Email: "invalid-email",
			Age:   25,
		}

		err := service.NotifyUser(user, "Тест")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "некорректный email")

		// Мок не должен был вызываться
		mockEmailSender.AssertNotCalled(t, "SendEmail")
	})
}

func TestSkippedTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем интеграционный тест")
	}

	t.Log("Выполняем интеграционный тест")
}

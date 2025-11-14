package main

import (
	"errors"
	"testing"

	"table-tests-with-mock/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotificationService_table_tests(t *testing.T) {
	tests := []struct {
		name                  string
		user                  *User
		message               string
		mockSetup             func(*mocks.EmailSender)
		expectedError         bool
		expectedErrorContains string
		shouldCallSendEmail   bool
	}{
		{
			name: "успешная_отправка_уведомления",
			user: &User{
				Name:  "Тест",
				Email: "test@example.com",
				Age:   25,
			},
			message: "Ваш заказ готов",
			mockSetup: func(mockEmailSender *mocks.EmailSender) {
				mockEmailSender.On("SendEmail",
					"test@example.com",
					"Уведомление",
					"Привет, Тест!\n\nВаш заказ готов").Return(nil)
			},
			expectedError:       false,
			shouldCallSendEmail: true,
		},
		{
			name: "ошибка_при_отправке_email",
			user: &User{
				Name:  "Тест",
				Email: "test@example.com",
				Age:   25,
			},
			message: "Тест",
			mockSetup: func(mockEmailSender *mocks.EmailSender) {
				mockEmailSender.On("SendEmail",
					mock.Anything, mock.Anything, mock.Anything).Return(errors.New("сервер недоступен"))
			},
			expectedError:         true,
			expectedErrorContains: "сервер недоступен",
			shouldCallSendEmail:   true,
		},
		{
			name: "некорректный_email_пользователя",
			user: &User{
				Name:  "Тест",
				Email: "invalid-email",
				Age:   25,
			},
			message: "Тест",
			mockSetup: func(mockEmailSender *mocks.EmailSender) {
				// Мок не должен вызываться для некорректного email
			},
			expectedError:         true,
			expectedErrorContains: "некорректный email",
			shouldCallSendEmail:   false,
		},
		{
			name: "пустой_email",
			user: &User{
				Name:  "Тест",
				Email: "",
				Age:   25,
			},
			message: "Тест сообщение",
			mockSetup: func(mockEmailSender *mocks.EmailSender) {
				// Мок не должен вызываться для пустого email
			},
			expectedError:         true,
			expectedErrorContains: "некорректный email",
			shouldCallSendEmail:   false,
		},
		{
			name: "длинное_сообщение",
			user: &User{
				Name:  "Александр",
				Email: "aleksandr@example.com",
				Age:   30,
			},
			message: "Это очень длинное сообщение для проверки того, как система обрабатывает большие объемы текста в уведомлениях",
			mockSetup: func(mockEmailSender *mocks.EmailSender) {
				mockEmailSender.On("SendEmail",
					"aleksandr@example.com",
					"Уведомление",
					"Привет, Александр!\n\nЭто очень длинное сообщение для проверки того, как система обрабатывает большие объемы текста в уведомлениях").Return(nil)
			},
			expectedError:       false,
			shouldCallSendEmail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок для каждого теста
			mockEmailSender := mocks.NewEmailSender(t)

			// Настраиваем ожидания мока
			tt.mockSetup(mockEmailSender)

			// Создаем сервис с моком
			service := NewNotificationService(mockEmailSender)

			// Выполняем тестируемый метод
			err := service.NotifyUser(tt.user, tt.message)

			// Проверяем результат
			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrorContains != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			// Проверяем вызовы мока
			if !tt.shouldCallSendEmail {
				mockEmailSender.AssertNotCalled(t, "SendEmail")
			}
			// AssertExpectations автоматически вызывается через cleanup функцию в моке
		})
	}
}

func TestSkippedTest_table_style(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем интеграционный тест")
	}

	t.Log("Выполняем интеграционный тест в table style")
}

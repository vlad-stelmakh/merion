package main

import (
	"testing"
)

func TestAdd_Simple(t *testing.T) {
	result := Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2, 3) = %d; ожидалось 5", result)
	}
}

func TestAdd_Multiple(t *testing.T) {
	if Add(2, 3) != 5 {
		t.Error("2 + 3 должно быть 5")
	}

	if Add(-1, 1) != 0 {
		t.Error("-1 + 1 должно быть 0")
	}

	if Add(0, 0) != 0 {
		t.Error("0 + 0 должно быть 0")
	}
}

func Test_WithLogging(t *testing.T) {
	t.Log("Начинаем тест с логированием")

	result := Add(2, 3)
	t.Logf("Результат сложения: %d", result)

	if result != 5 {
		t.Error("ошибка, ожидалось 5 в результате сложения")
		return
	}

	t.Log("Тест завершен успешно")
}

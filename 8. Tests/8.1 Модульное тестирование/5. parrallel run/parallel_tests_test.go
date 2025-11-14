package main

import (
	"testing"
	"time"
)

// 880 ms - 320ms
// Тесты БЕЗ параллельного выполнения - для сравнения
func TestSequential_Calculator(t *testing.T) {
	start := time.Now()

	tests := []struct {
		name     string
		function func() int
		expected int
	}{
		{"медленные_вычисления_5", func() int { return SlowCalculation(5) }, 25},
		{"медленные_вычисления_3", func() int { return SlowCalculation(3) }, 9},
		{"медленные_вычисления_4", func() int { return SlowCalculation(4) }, 16},
		{"медленные_вычисления_2", func() int { return SlowCalculation(2) }, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function()
			if result != tt.expected {
				t.Errorf("получено %d, ожидалось %d", result, tt.expected)
			}
		})
	}

	duration := time.Since(start)

	t.Logf("Последовательное выполнение заняло: %v", duration)
}

// Тесты С параллельным выполнением
func TestParallel_Calculator(t *testing.T) {
	start := time.Now()

	t.Parallel() // включаем параллельное выполнение

	tests := []struct {
		name     string
		function func() int
		expected int
	}{
		{"медленные_вычисления_5", func() int { return SlowCalculation(5) }, 25},
		{"медленные_вычисления_3", func() int { return SlowCalculation(3) }, 9},
		{"медленные_вычисления_4", func() int { return SlowCalculation(4) }, 16},
		{"медленные_вычисления_2", func() int { return SlowCalculation(2) }, 4},
	}

	for _, tt := range tests {
		tt := tt // важно для замыкания в горутинах
		t.Run(tt.name, func(t *testing.T) {

			result := tt.function()
			if result != tt.expected {
				t.Errorf("получено %d, ожидалось %d", result, tt.expected)
			}
		})
	}

	duration := time.Since(start)
	t.Logf("Параллельное выполнение заняло: %v", duration)
}

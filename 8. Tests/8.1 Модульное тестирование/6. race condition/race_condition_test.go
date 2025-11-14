package main

import (
	"sync"
	"testing"
)

// Демонстрация race condition
func TestRace_Condition_Demo(t *testing.T) {
	t.Run("небезопасный_счетчик", func(t *testing.T) {
		ResetCounter()
		var wg sync.WaitGroup

		// Запускаем 100 горутин для инкремента
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				UnsafeIncrement()
			}()
		}

		wg.Wait()
		result := GetCounter()
		t.Logf("Небезопасный счетчик: %d (ожидалось 100, но может быть меньше из-за race condition)", result)

		if result == 100 {
			t.Log("Повезло! Получили правильный результат, но это случайность")
		} else {
			t.Logf("Race condition обнаружен: потеряно %d инкрементов", 100-result)
		}
	})

	t.Run("безопасный_счетчик", func(t *testing.T) {
		ResetCounter()
		var wg sync.WaitGroup

		// Запускаем 100 горутин для безопасного инкремента
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				SafeIncrement()
			}()
		}

		wg.Wait()
		result := GetCounter()

		if result != 100 {
			t.Errorf("Безопасный счетчик = %d; ожидалось 100", result)
		}
		t.Logf("Безопасный счетчик: %d (всегда правильный результат)", result)
	})
}

// Более сложная демонстрация с несколькими операциями
func TestComplex_Race_Condition(t *testing.T) {
	t.Run("несколько_операций_небезопасно", func(t *testing.T) {
		ResetCounter()
		var wg sync.WaitGroup

		// Одни горутины увеличивают
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					UnsafeIncrement()
				}
			}()
		}

		// Другие тоже увеличивают
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					UnsafeIncrement()
				}
			}()
		}

		wg.Wait()
		result := GetCounter()
		expected := 1000 // 50 * 10 + 50 * 10

		t.Logf("Сложный небезопасный тест: %d из %d", result, expected)
		if result < expected {
			t.Logf("Потеряно %d инкрементов из-за race condition", expected-result)
		}
	})

	t.Run("несколько_операций_безопасно", func(t *testing.T) {
		ResetCounter()
		var wg sync.WaitGroup

		// Одни горутины увеличивают
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					SafeIncrement()
				}
			}()
		}

		// Другие тоже увеличивают
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					SafeIncrement()
				}
			}()
		}

		wg.Wait()
		result := GetCounter()
		expected := 1000

		if result != expected {
			t.Errorf("Безопасный сложный тест = %d; ожидалось %d", result, expected)
		}
		t.Logf("Сложный безопасный тест: %d (всегда правильно)", result)
	})
}

// Бенчмарк для сравнения производительности
func BenchmarkUnsafeIncrement(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ResetCounter()
		UnsafeIncrement()
	}
}

func BenchmarkSafeIncrement(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ResetCounter()
		SafeIncrement()
	}
}

package main

import (
	"log"
	"time"
)

func FormattingExamples() {
	log.Println("\n=== Форматирование сообщений ===")

	name := "Алексей"
	email := "alexey@example.com"
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("Пользователь: %s (%s) зарегистрирован в %s", name, email, timestamp)

	user := struct {
		ID    int
		Name  string
		Email string
	}{
		ID:    1,
		Name:  "Иван",
		Email: "ivan@example.com",
	}

	log.Printf("Данные пользователя: %+v", user)
}

func PerformanceLogging() {
	log.Println("\n=== Логирование производительности ===")

	logDuration := func(operation string, startTime time.Time) {
		duration := time.Since(startTime)
		log.Printf("[PERF] %s выполнено за %v",
			operation, duration)
	}

	start := time.Now()
	time.Sleep(10 * time.Millisecond) // Имитация работы
	logDuration("Быстрая операция", start)

	start = time.Now()
	time.Sleep(100 * time.Millisecond) // Имитация работы
	logDuration("Медленная операция", start)
}

func main() {
	log.Println("=== ДЕМОНСТРАЦИЯ ФОРМАТИРОВАНИЯ ЛОГОВ ===")

	FormattingExamples()
	PerformanceLogging()

	log.Println("\n=== Техники форматирования ===")
	log.Println("• Printf() - форматирование с placeholder'ами")
	log.Println("• %v, %+v - вывод структур и значений")
	log.Println("• %T - тип данных")
	log.Println("• Выравнивание текста и чисел")
	log.Println("• JSON-подобное структурированное логирование")
}

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	demonstrateConcurrentDebugging()
}

func demonstrateConcurrentDebugging() {
	fmt.Println("=== Отладка конкурентного кода ===")

	simpleGoroutines()
	channelExample()
	raceConditionExample()
}

func simpleGoroutines() {
	fmt.Println("\n--- Простые горутины ---")

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			for j := 0; j < 3; j++ {
				fmt.Printf("Горутина %d: шаг %d\n", id, j)
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Все горутины завершены")
}

func channelExample() {
	fmt.Println("\n--- Каналы ---")

	messages := make(chan string, 3)
	results := make(chan string, 3)

	go func() {
		for i := 0; i < 3; i++ {
			msg := fmt.Sprintf("Сообщение %d", i)
			messages <- msg
			fmt.Printf("Отправлено: %s\n", msg)
		}
		close(messages)
	}()

	go func() {
		for msg := range messages {
			processed := fmt.Sprintf("Обработано: %s", msg)
			results <- processed
			time.Sleep(50 * time.Millisecond)
		}
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}

func raceConditionExample() {
	fmt.Println("\n--- Race Condition ---")

	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	numGoroutines := 3
	incrementsPerGoroutine := 5

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			for j := 0; j < incrementsPerGoroutine; j++ {
				mu.Lock()
				oldValue := counter
				counter = oldValue + 1
				newValue := counter
				mu.Unlock()

				fmt.Printf("Горутина %d: инкремент %d, %d -> %d\n", id, j, oldValue, newValue)
			}
		}(i)
	}

	wg.Wait()

	expected := numGoroutines * incrementsPerGoroutine
	fmt.Printf("Финальный counter: %d (ожидалось: %d)\n", counter, expected)

	if counter == expected {
		fmt.Println("✅ Race condition предотвращен!")
	} else {
		fmt.Println("❌ Race condition произошел!")
	}
}

func deadlockExample() {
	fmt.Println("\n--- Deadlock (НЕ ЗАПУСКАЙТЕ) ---")

	/*
		mu1 := &sync.Mutex{}
		mu2 := &sync.Mutex{}

		go func() {
			mu1.Lock()
			fmt.Println("Горутина 1: заблокировала mu1")
			time.Sleep(100 * time.Millisecond)

			mu2.Lock()
			fmt.Println("Горутина 1: заблокировала mu2")
			mu2.Unlock()
			mu1.Unlock()
		}()

		go func() {
			mu2.Lock()
			fmt.Println("Горутина 2: заблокировала mu2")
			time.Sleep(100 * time.Millisecond)

			mu1.Lock()
			fmt.Println("Горутина 2: заблокировала mu1")
			mu1.Unlock()
			mu2.Unlock()
		}()

		time.Sleep(1 * time.Second)
		fmt.Println("Deadlock произошел!")
	*/
}

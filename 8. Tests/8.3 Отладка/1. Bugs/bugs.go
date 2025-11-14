package main

import (
	"fmt"
	"strings"
)

func main() {
	indexOutOfBoundsBug()
	logicBug()
	concurrencyBug()
}

func indexOutOfBoundsBug() {
	fmt.Println("\n--- Index Out of Bounds Bug ---")

	numbers := []int{1, 2, 3, 4, 5}

	for i := 0; i <= 10; i++ {
		if i < len(numbers) {
			fmt.Printf("numbers[%d] = %d\n", i, numbers[i])
		} else {
			fmt.Printf("Индекс %d вне границ массива (размер: %d)\n", i, len(numbers))
		}
	}
}

func logicBug() {
	fmt.Println("\n--- Logic Bug ---")

	text := "Hello, World! This is a test."
	words := strings.Fields(text)

	var longWords []string

	for i, word := range words {
		cleanWord := strings.Trim(word, ".,!?")

		if len(cleanWord) > 5 {
			longWords = append(longWords, cleanWord)
		}

		fmt.Printf("Слово %d: '%s' (длина: %d)\n", i, cleanWord, len(cleanWord))
	}

	fmt.Printf("Длинные слова: %v\n", longWords)
}

func concurrencyBug() {
	fmt.Println("\n--- Concurrency Bug ---")

	counter := 0
	done := make(chan bool)

	for i := 0; i < 3; i++ {
		go func(id int) {
			for j := 0; j < 5; j++ {
				oldValue := counter
				counter = oldValue + 1

				fmt.Printf("Горутина %d: итерация %d, counter=%d\n", id, j, counter)
			}
			done <- true
		}(i)
	}

	// Ждем завершения всех горутин
	for i := 0; i < 3; i++ {
		<-done
	}

	fmt.Printf("Финальное значение counter: %d (ожидалось: 15)\n", counter)
}

func complexCalculation(x, y int) int {
	result := multiply(x, y)
	return result
}

func multiply(a, b int) int {
	return a * b
}

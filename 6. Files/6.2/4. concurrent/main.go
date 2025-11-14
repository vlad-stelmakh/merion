package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	var mu sync.Mutex
	file, err := os.Create("concurrent.txt")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer file.Close()

	mu.Lock()
	file.WriteString("Запись из горутины")
	mu.Unlock()
}

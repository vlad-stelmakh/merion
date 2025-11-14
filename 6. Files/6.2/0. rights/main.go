package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.Chmod("img.png", 0642)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Права доступа изменены на 0644")
}

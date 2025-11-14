package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	data, err := os.ReadFile("file.txt")
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)

		return
	}

	fmt.Println(string(data))

	time.Sleep(time.Minute)
}

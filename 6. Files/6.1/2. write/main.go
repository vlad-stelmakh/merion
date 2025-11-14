package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("Привет, Go!")
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
		return
	}
}

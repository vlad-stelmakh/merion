package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("output.csv")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	records := [][]string{
		{"Имя", "Фамилия", "Возраст"},
		{"Иван", "Иванов", "30"},
		{"Мария", "Петрова", "25"},
	}

	if err = writer.WriteAll(records); err != nil {
		fmt.Println("Ошибка записи:", err)
		return
	}
}

package main

import (
	"bufio"
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

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	_, err = writer.WriteString("Пример буферизованной записи")
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
	}
}

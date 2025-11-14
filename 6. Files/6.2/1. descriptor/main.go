package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("file.txt")
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)

		return
	}

	file.Write([]byte("Hello Go!"))
	file.Write([]byte("Bye Go!"))

	file.Close()

	fileInfo, err := os.Stat("file.txt")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Имя:", fileInfo.Name())
	fmt.Println("Размер:", fileInfo.Size(), "байт")
	fmt.Println("Время изменения:", fileInfo.ModTime())
	fmt.Println("Права доступа:", fileInfo.Mode())
	fmt.Println("Это директория?", fileInfo.IsDir())
}

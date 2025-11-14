package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	tmpFile, err := os.CreateTemp("", "example_*.txt")
	if err != nil {
		fmt.Printf("Ошибка создания временного файла: %v\n", err)
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	fmt.Printf("Создан временный файл: %s\n", tmpFile.Name())

	content := []byte("Временные данные для теста\nВторая строка")
	if _, err = tmpFile.Write(content); err != nil {
		fmt.Printf("Ошибка записи во временный файл: %v\n", err)
		return
	}

	readData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		fmt.Printf("Ошибка чтения временного файла: %v\n", err)
		return
	}
	fmt.Printf("\nСодержимое файла:\n%s\n", string(readData))

	tmpDir, err := os.MkdirTemp("", "example_dir_*")
	if err != nil {
		fmt.Printf("Ошибка создания временной директории: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	fmt.Printf("\nСоздана временная директория: %s\n", tmpDir)

	for i := 1; i <= 3; i++ {
		filePath := filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i))
		if err = os.WriteFile(filePath, []byte(fmt.Sprintf("Файл %d", i)), 0644); err != nil {
			fmt.Printf("Ошибка создания файла %s: %v\n", filePath, err)
			return
		}
	}

	files, err := os.ReadDir(tmpDir)
	if err != nil {
		fmt.Printf("Ошибка чтения директории: %v\n", err)
		return
	}

	fmt.Println("\nСодержимое временной директории:")
	for _, file := range files {
		content, _ = os.ReadFile(tmpFile.Name())
		fmt.Printf("- %s (%d bytes)\n", file.Name(), len(content))
	}
}

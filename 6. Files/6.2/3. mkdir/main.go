package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}

func main() {
	dirName := "myapp_config"
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		fmt.Printf("Ошибка при создании директории: %v\n", err)
		return
	}
	fmt.Printf("Директория '%s' создана\n", dirName)

	config := Config{
		AppName: "MyAwesomeApp",
		Version: "1.0.0",
	}

	filePath := fmt.Sprintf("%s/config.json", dirName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Ошибка при создании файла: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(config); err != nil {
		fmt.Printf("Ошибка при записи JSON: %v\n", err)
		return
	}
	fmt.Printf("Файл конфигурации создан: %s\n", filePath)

	fmt.Println("\nСодержимое директории:")
	files, err := os.ReadDir(dirName)
	if err != nil {
		fmt.Printf("Ошибка при чтении директории: %v\n", err)
		return
	}

	for _, file := range files {
		info, _ := file.Info()
		fmt.Printf("- %s (%.1f KB)\n", file.Name(), float64(info.Size())/1024)
	}
}

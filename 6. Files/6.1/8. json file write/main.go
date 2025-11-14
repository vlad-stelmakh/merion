package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	type Person struct {
		Name     string `json:"name"`
		LastName string `json:"lastName,omitempty"`
		Age      int    `json:"age"`
	}

	people := []Person{
		{Name: "Иван", Age: 30, LastName: "Иванов"},
		{Name: "Мария", Age: 25},
	}

	file, err := os.Create("people.json")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(people); err != nil {
		fmt.Println("Ошибка записи JSON:", err)
	}
}

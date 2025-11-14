package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Person struct {
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Email string   `json:"email,omitempty"`
	Cars  []string `json:"cars,omitempty"`
}

func main() {
	data, err := os.ReadFile("cfg.json")
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	var person Person
	err = json.Unmarshal(data, &person)
	if err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		return
	}

	fmt.Printf("Имя: %s\nВозраст: %d\nEmail: %sn",
		person.Name, person.Age, person.Email)
}

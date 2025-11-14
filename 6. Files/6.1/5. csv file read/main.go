package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

//
//func main() {
//	file, err := os.Open("data.csv")
//	if err != nil {
//		fmt.Println("Ошибка открытия файла:", err)
//		return
//	}
//	defer file.Close()
//
//	reader := csv.NewReader(file)
//
//	records, err := reader.ReadAll()
//	if err != nil {
//		fmt.Println("Ошибка чтения CSV:", err)
//		return
//	}
//
//	for i, record := range records {
//		fmt.Println("Запись", i, ":", record[1])
//	}
//}

func main() {
	file, err := os.Open("data.csv")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("Ошибка чтения CSV:", err)
			return
		}

		fmt.Println("Запись:", record)
	}
}

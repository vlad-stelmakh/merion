package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	path := "."

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Ошибка доступа к пути %q: %v\n", path, err)

			return err
		}

		if info.IsDir() {
			return nil
		}

		fmt.Println(path)

		return nil
	})

	if err != nil {
		fmt.Printf("Ошибка обхода директории: %v\n", err)
	}
}

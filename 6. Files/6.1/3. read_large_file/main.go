package main

import (
	"fmt"
	"io"
	"os"
)

//
//func main() {
//	file, err := os.Open("largefile.txt")
//	if err != nil {
//		fmt.Println("Ошибка открытия файла:", err)
//		return
//	}
//	defer file.Close()
//
//	scanner := bufio.NewScanner(file)
//	for scanner.Scan() {
//		fmt.Println("----------------------------------------")
//		fmt.Println(scanner.Text(), len(scanner.Bytes()))
//	}
//
//	if err = scanner.Err(); err != nil {
//		fmt.Println("Ошибка чтения файла:", err)
//	}
//}

func main() {
	file, err := os.Open("largefile.txt")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 256)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Ошибка чтения файла:", err)
			}

			break
		}

		fmt.Println("----------------------------------------")
		fmt.Printf("Прочитано %d байт: %s\n", n, string(buffer[:n]))
	}
}

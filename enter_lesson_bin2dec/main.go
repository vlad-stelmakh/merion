package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	var input string

	for {
		fmt.Println("Введите число в двоичном виде:")
		fmt.Scanln(&input)

		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		if input == "exit" {
			break
		}

		res, err := binToDecimal(input)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Result:", res)
	}

	fmt.Println("Program is finished!")
}

func binToDecimal(input string) (int, error) {
	if len(input) > 8 {
		return 0, fmt.Errorf("input should be less than 8 characters")
	}

	for i := 0; i < len(input); i++ {
		if input[i] != '0' && input[i] != '1' {
			return 0, fmt.Errorf("input should be 0 or 1")
		}
	}

	var res int
	for i := 0; i < len(input); i++ {
		if input[i] == '1' {
			res += int(math.Pow(2, float64(i)))
		}
	}

	return res, nil
}

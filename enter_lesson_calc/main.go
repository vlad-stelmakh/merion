package main

import "fmt"

func main() {
	var operator string
	var firstDigit, secondDigit, res int

	for {
		fmt.Println("Enter operator number: ")
		fmt.Scan(&operator)

		fmt.Println("Enter first number: ")
		fmt.Scan(&firstDigit)

		fmt.Println("Enter second number: ")
		fmt.Scan(&secondDigit)

		if operator == "^" {
			break
		}

		switch operator {
		case "+":
			res = firstDigit + secondDigit
			fmt.Printf("Result: %d\n", res)
		case "-":
			res = firstDigit - secondDigit
			fmt.Printf("Result: %d\n", res)
		case "*":
			res = firstDigit * secondDigit
			fmt.Printf("Result: %d\n", res)
		case "/":
			res = firstDigit / secondDigit
			fmt.Printf("Result: %d\n", res)
		default:
			fmt.Println("Invalid operator")
		}
	}

	fmt.Println("App closed")
}

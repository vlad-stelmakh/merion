package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "grpc-client/generated"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculatorServiceClient(conn)

	// Контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Примеры вызовов различных операций
	operations := []struct {
		name string
		a    float64
		b    float64
		op   string
	}{
		{"Сложение", 10.5, 5.3, "add"},
		{"Вычитание", 20.0, 8.5, "subtract"},
		{"Умножение", 4.0, 7.0, "multiply"},
		{"Деление", 15.0, 3.0, "divide"},
		{"Неизвестная операция", 5.0, 2.0, "unknown"},
	}

	fmt.Println("=== gRPC Калькулятор Клиент ===\n")

	for _, op := range operations {
		fmt.Printf("Операция: %s (%.2f %s %.2f)\n", op.name, op.a, op.op, op.b)

		// Создаём запрос
		req := &pb.CalculateRequest{
			A:         op.a,
			B:         op.b,
			Operation: op.op,
		}

		// Вызываем метод сервера
		resp, err := client.Calculate(ctx, req)
		if err != nil {
			fmt.Printf("Ошибка вызова: %v\n", err)
		} else {
			if resp.Error != "" {
				fmt.Printf("Результат: Ошибка - %s\n", resp.Error)
			} else {
				fmt.Printf("Результат: %.2f\n", resp.Result)
			}
		}
		fmt.Println("---")
	}

	fmt.Println("До свидания!")
}

//// Интерактивный режим
//fmt.Println("\n=== Интерактивный режим ===")
//fmt.Println("Введите операцию в формате: число операция число")
//fmt.Println("Доступные операции: add, subtract, multiply, divide")
//fmt.Println("Пример: 10.5 add 5.3")
//fmt.Println("Для выхода введите 'exit'")
//
//var a, b float64
//var operation, input string
//
//for {
//fmt.Print("\n> ")
//_, err := fmt.Scanln(&input)
//if err != nil {
//continue
//}
//
//if input == "exit" {
//break
//}
//
//// Простой парсинг ввода
//n, err := fmt.Sscanf(input+" ", "%f %s %f", &a, &operation, &b)
//if err != nil || n != 3 {
//fmt.Println("Неверный формат. Используйте: число операция число")
//continue
//}
//
//req := &pb.CalculateRequest{
//A:         a,
//B:         b,
//Operation: operation,
//}
//
//resp, err := client.Calculate(ctx, req)
//if err != nil {
//fmt.Printf("Ошибка: %v\n", err)
//} else {
//if resp.Error != "" {
//fmt.Printf("Ошибка: %s\n", resp.Error)
//} else {
//fmt.Printf("Результат: %.2f\n", resp.Result)
//}
//}
//}

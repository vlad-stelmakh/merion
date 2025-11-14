package main

import (
	"log"
)

func FlagsExamples() {
	log.Println("\n--- Различные флаги ---")

	log.SetFlags(log.Ldate)
	log.SetPrefix("[DATE] ")
	log.Println("Только дата")

	log.SetFlags(log.Ltime)
	log.SetPrefix("[TIME] ")
	log.Println("Только время")

	log.SetFlags(log.LstdFlags)
	log.SetPrefix("[DATETIME] ")
	log.Println("Дата и время")

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[MICRO] ")
	log.Println("С микросекундами")

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[SHORT] ")
	log.Println("С коротким именем файла")

	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetPrefix("[LONG] ")
	log.Println("С полным путем к файлу")

	log.SetFlags(log.LstdFlags | log.LUTC)
	log.SetPrefix("[UTC] ")
	log.Println("UTC время")
}

func LoggerConfiguration() {
	log.Println("\n=== Конфигурации логгера ===")

	// Сохраняем текущие настройки
	originalFlags := log.Flags()
	originalPrefix := log.Prefix()

	defer func() {
		log.SetFlags(originalFlags)
		log.SetPrefix(originalPrefix)
	}()

	// Конфигурация для разработки
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[DEV] ")
	log.Println("Конфигурация для разработки")

	// Конфигурация для продакшена
	log.SetFlags(log.LstdFlags | log.LUTC)
	log.SetPrefix("[PROD] ")
	log.Println("Конфигурация для продакшена")

	// Конфигурация для тестирования
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.SetPrefix("[TEST] ")
	log.Println("Конфигурация для тестирования")
}

func main() {
	log.Println("=== ДЕМОНСТРАЦИЯ ФЛАГОВ И ПРЕФИКСОВ ===")

	FlagsExamples()
	LoggerConfiguration()

	log.Println("\n=== Флаги логирования ===")
	log.Println("• log.Ldate - дата (2009/01/23)")
	log.Println("• log.Ltime - время (01:23:23)")
	log.Println("• log.Lmicroseconds - микросекунды (01:23:23.123123)")
	log.Println("• log.Llongfile - полный путь к файлу (/a/b/c/d.go:23)")
	log.Println("• log.Lshortfile - имя файла (d.go:23)")
	log.Println("• log.LUTC - UTC время")
	log.Println("• log.LstdFlags - стандартные флаги (Ldate | Ltime)")
}

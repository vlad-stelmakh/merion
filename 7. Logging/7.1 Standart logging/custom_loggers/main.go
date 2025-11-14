package main

import (
	"io"
	"log"
	"os"
)

func CustomLoggerExamples() {
	log.Println("\n=== Пользовательские логгеры ===")

	appLogger := log.New(os.Stdout, "[APP] ", 0)
	appLogger.Println("Это сообщение от логгера приложения")

	errorLogger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	errorLogger.Println("Это сообщение об ошибке")

	debugLogger := log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lmicroseconds)
	debugLogger.Println("Отладочная информация с микросекундами")
}

func FileLogger() {
	file, err := os.OpenFile("file.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Не удалось создать файл лога: %v", err)
		return
	}
	defer file.Close()

	fileWriter := io.Writer(file)
	fileLogger := log.New(fileWriter, "[FILE] ", log.LstdFlags)

	fileLogger.Println("Это сообщение записывается и в консоль, и в файл")
	log.Println("Проверьте файл app.log")
}

func MultiWriterLogger() {
	log.Println("\n=== Логирование в несколько мест ===")

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Не удалось создать файл лога: %v", err)
		return
	}
	defer file.Close()

	multiWriter := io.MultiWriter(os.Stdout, file)
	multiLogger := log.New(multiWriter, "[MULTI] ", log.LstdFlags)

	multiLogger.Println("Это сообщение записывается и в консоль, и в файл")
	log.Println("Проверьте файл app.log")
}

func main() {
	log.Println("=== ДЕМОНСТРАЦИЯ ПОЛЬЗОВАТЕЛЬСКИХ ЛОГГЕРОВ ===")

	CustomLoggerExamples()
	FileLogger()
	MultiWriterLogger()

	log.Println("\n=== Основные возможности ===")
	log.Println("• log.New() - создание пользовательского логгера")
	log.Println("• io.MultiWriter() - запись в несколько мест одновременно")
	log.Println("• Различные префиксы для разных типов сообщений")
	log.Println("• Настройка флагов для разных логгеров")
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<h1>Добро пожаловать на мой HTTP сервер!!!</h1>`)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Гость"
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<h1>Привет, %s!</h1>`, name)
}

func apiHandler(w http.ResponseWriter, _ *http.Request) {
	data := map[string]interface{}{
		"status":    "success",
		"message":   "Данные от сервера",
		"timestamp": time.Now().Unix(),
		"server":    "Go HTTP Server",
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(data)
}

func main() {
	mux := http.NewServeMux()

	// Настройка маршрутов
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/api/data", apiHandler)

	port := ":8080"
	log.Printf("Сервер запущен на http://localhost%s", port)
	log.Printf("Доступные endpoints:")
	log.Printf("  - http://localhost%s/ (Главная страница)", port)
	log.Printf("  - http://localhost%s/hello?name=Имя (Приветствие)", port)
	log.Printf("  - http://localhost%s/api/data (JSON API)", port)
	log.Printf("  - http://localhost%s/health (Статус сервера)", port)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type DataStructure struct {
	ID      int
	Name    string
	Values  []float64
	SubData map[string]interface{}
}

func main() {
	start := time.Now()

	createManySlices()
	createLargeStructures()
	simulateMemoryLeaks()

	runtime.GC()
	runtime.GC()
	
	memFile, err := os.Create("mem.prof")
	if err != nil {
		panic(err)
	}
	defer memFile.Close()

	if err := pprof.WriteHeapProfile(memFile); err != nil {
		panic(err)
	}

	fmt.Printf("Memory профилирование завершено за: %v\n", time.Since(start))
	printMemStats()
	fmt.Println("Профиль сохранен в mem.prof")
	fmt.Println("Анализ: go tool pprof -http=:8080 mem.prof")
}

func createManySlices() {
	const numSlices = 10000
	slices := make([][]int, numSlices)

	for i := 0; i < numSlices; i++ {
		size := 100 + (i % 900)
		slices[i] = make([]int, size)

		for j := 0; j < size; j++ {
			slices[i][j] = i * j
		}
	}

	fmt.Printf("Создано %d slice'ов\n", numSlices)
	for i := 0; i < numSlices/2; i++ {
		slices[i] = nil
	}
}

func createLargeStructures() {
	const numStructs = 50000
	structures := make([]*DataStructure, numStructs)

	for i := 0; i < numStructs; i++ {
		structures[i] = &DataStructure{
			ID:   i,
			Name: fmt.Sprintf("Structure_%d", i),
			Values: func() []float64 {
				values := make([]float64, 50)
				for j := range values {
					values[j] = float64(i*j) * 3.14
				}
				return values
			}(),
			SubData: map[string]interface{}{
				"timestamp": time.Now(),
				"counter":   i,
				"metadata":  fmt.Sprintf("meta_%d", i),
			},
		}
	}

	fmt.Printf("Создано %d структур\n", numStructs)

	_ = structures
}

var globalData [][]byte

func simulateMemoryLeaks() {
	for i := 0; i < 1000; i++ {
		data := make([]byte, 1024*10)

		for j := range data {
			data[j] = byte(i % 256)
		}

		globalData = append(globalData, data)
	}

	fmt.Printf("Накоплено %d блоков данных\n", len(globalData))
}

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Println("\n=== Статистика памяти ===")
	fmt.Printf("Аллоцировано: %d KB\n", bToKb(m.Alloc))
	fmt.Printf("Всего аллокаций: %d KB\n", bToKb(m.TotalAlloc))
	fmt.Printf("Получено от системы: %d KB\n", bToKb(m.Sys))
	fmt.Printf("Количество GC: %d\n", m.NumGC)
	fmt.Printf("Количество горутин: %d\n", runtime.NumGoroutine())
}

func bToKb(b uint64) uint64 {
	return b / 1024
}

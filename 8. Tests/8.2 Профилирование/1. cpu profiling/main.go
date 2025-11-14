package main

import (
	"fmt"
	"math"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	cpuProfile := "cpu.prof"
	cpuFile, err := os.Create(cpuProfile)
	if err != nil {
		fmt.Printf("Could not create CPU profile: %v\n", err)
		return
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		fmt.Printf("Could not start CPU profile: %v\n", err)
		return
	}
	defer pprof.StopCPUProfile()

	start := time.Now()

	// Run the intensive computation multiple times to get more samples
	var primes []int
	for i := 0; i < 100; i++ {
		primes = calculatePrimes(50000)
		// Add some additional computational work
		performAdditionalWork()
	}

	duration := time.Since(start)

	fmt.Printf("Found %d primes up to %d\n", len(primes), 50000)
	fmt.Printf("Execution time: %v\n", duration)

	fmt.Println("\nProfiling completed:")
	fmt.Printf("- CPU profile: %s\n", cpuProfile)
	fmt.Println("\nTo analyze profiles:")
	fmt.Println("  go tool pprof -http=:8080", cpuProfile)
}

func calculatePrimes(limit int) []int {
	primes := make([]int, 0, limit/2) // Pre-allocate capacity for better performance

	for num := 2; num <= limit; num++ {
		if isPrime(num) {
			primes = append(primes, num)
		}
	}

	return primes
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	sqrt := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrt; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func performAdditionalWork() {
	size := 100
	matrix1 := make([][]int, size)
	matrix2 := make([][]int, size)
	result := make([][]int, size)

	// Initialize matrices
	for i := 0; i < size; i++ {
		matrix1[i] = make([]int, size)
		matrix2[i] = make([]int, size)
		result[i] = make([]int, size)
		for j := 0; j < size; j++ {
			matrix1[i][j] = i + j
			matrix2[i][j] = i * j
		}
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			for k := 0; k < size; k++ {
				result[i][j] += matrix1[i][k] * matrix2[k][j]
			}
		}
	}
}

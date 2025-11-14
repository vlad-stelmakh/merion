package main

import "time"

func SlowCalculation(n int) int {
	time.Sleep(100 * time.Millisecond)

	return n * n
}

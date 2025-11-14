package main

import "sync"

var counter int
var mutex sync.Mutex

func UnsafeIncrement() {
	counter++
}

func SafeIncrement() {
	mutex.Lock()
	counter++
	mutex.Unlock()
}

func ResetCounter() {
	counter = 0
}

func GetCounter() int {
	return counter
}

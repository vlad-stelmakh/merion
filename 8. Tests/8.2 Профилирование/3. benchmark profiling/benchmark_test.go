package main

import (
	"testing"
)

func BenchmarkStringConcat(b *testing.B) {
	strs := []string{"hello", "world", "benchmark", "test", "performance"}

	b.Run("Slow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = SlowStringConcat(strs)
		}
	})

	b.Run("Fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FastStringConcat(strs)
		}
	})
}
